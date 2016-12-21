package streamsavvy

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/streadway/amqp"
	"gopkg.in/redis.v5"
	"net/http"
	"sync"
)

type Reco struct {
	mu         sync.Mutex
	GuideboxId string
	sock       *websocket.Conn
}

func HandleRecomendations(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	rmqc := r.Context().Value("rabbitmq").(middleware.RMQCH)
	r_client := r.Context().Value("redis_client").(*redis.Client)
	cleanup := r.Context().Value("cleanup").(chan string)

	SimKey := "ss_reco:%v:%v"
	categories := []string{"genres", "tags", "cast"}

	common.Check(err)
	if err != nil {
		conn.Close()
	}

	reco := Reco{}

	reco.sock = conn

	for {
		messageType, p, err := conn.ReadMessage()
		common.Check(err)
		if err != nil {
			//common.Check(err)
			logrus.Error(err)
			conn.Close()
			return
		}

		reco.GuideboxId = string(p[:])

		reco_ids := []string{}

		for _, cat := range categories {

			q := fmt.Sprintf(SimKey, cat, reco.GuideboxId)

			res := r_client.ZRange(q, 0, 3)

			reco_ids = append(reco_ids, res.Val()...)
			//common.Check(err)

		}

		RemoveDuplicates(&reco_ids)

		for _, val := range reco_ids {

			reco.PublishShowInfo(val, rmqc)
		}

		//TODO use this to update show recomendations later.
		//reco.PublishShowInfo(p, rmqc)

		go func() {
			rx_q, err := rmqc.RX.QueueDeclare(
				"reco_engine_results",
				false,
				false,
				false,
				false,
				nil,
			)
			common.Check(err)
			if err != nil {
				logrus.Error(err)
				//rmqc.RX.Close()
				conn.Close()
			}
			logrus.Info("Listening to recomendation results channel")
			msgs, err := rmqc.RX.Consume(
				rx_q.Name, // queue
				"", // consumer
				true, // auto-ack
				false, // exclusive
				false, // no-local
				false, // no-wait
				nil, // args
			)
			common.Check(err)
			if err != nil {
				common.Check(err)
				logrus.Error(err)
				//rmqc.RX.Close()
				conn.Close()
			}

			for {

				select {
				case m := <-msgs:

					if len(m.Body) > 2 {

						err = reco.send(messageType, m.Body)
						common.Check(err)
						if err != nil {
							logrus.Error(err)
							reco.sock.Close()
						}
					}

				case <-cleanup:
					rmqc.RX.Close()
					break
				}
			}

			return

		}()

		select {
		case <-cleanup:
			rmqc.RX.Close()
			rmqc.TX.Close()
			return
		}
	}
}

//TODO ionert ping pong listeners

func (r Reco) send(msg int, m []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.sock.WriteMessage(msg, m)
}

func (r Reco) sendJson(msg int, m interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.sock.WriteJSON(m)
}

//TODO it would be nice to access database straight from  go
func (r Reco) PublishShowInfo(show_id string, rmqc middleware.RMQCH) {

	logrus.WithField("connection address", r.sock).Info("publishing", show_id)
	tx_q, err := rmqc.TX.QueueDeclare(
		"reco_engine",
		false,
		false,
		false,
		false,
		nil,
	)

	common.Check(err)
	if err != nil {
		logrus.Error(err)
		rmqc.TX.Close()
	}

	err = rmqc.TX.Publish(
		"", // exchange
		tx_q.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(show_id),
		})
	common.Check(err)
	if err != nil {
		common.Check(err)
		rmqc.TX.Close()
	}
}

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}