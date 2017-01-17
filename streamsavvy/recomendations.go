package streamsavvy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/streadway/amqp"
	"gopkg.in/redis.v5"
)

type Reco struct {
	mu         sync.Mutex
	GuideboxId string
	sock       *websocket.Conn
}

func HandleRecomendations(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	rmqc := r.Context().Value("rabbitmq").(middleware.RMQCH)
	r_client := r.Context().Value("redis_client").(*redis.Client)
	//cleanup := r.Context().Value("cleanup").(chan string)

	SimKey := "ss_reco:%v:%v"
	categories := []string{"genres", "cast"}

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

			res := r_client.ZRange(q, 0, 5)

			reco_ids = append(reco_ids, res.Val()...)
			//common.Check(err)

		}

		RemoveDuplicates(&reco_ids)

		logrus.Info(reco_ids)

		//reco.PublishShowInfo(reco_ids, rmqc, corrId)

		//TODO use this to update show recomendations later.
		//reco.PublishShowInfo(p, rmqc)

		q, err := rmqc.Ch.QueueDeclare(
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		common.Check(err)
		logrus.Info("Listening to recomendation results channel")
		msgs, err := rmqc.Ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		common.Check(err)

		corrId := common.RandomString(32)

		the_json, err := json.Marshal(reco_ids)

		err = rmqc.Ch.Publish(
			"",
			"reco_rpc_queue",
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          the_json,
			})

		go func() {
			for {

				select {
				case m := <-msgs:

					if corrId == m.CorrelationId {
						//logrus.Info(m.Body)

						temp := []interface{}{}

						json.Unmarshal(m.Body, &temp)

						for _, val := range temp {
							res, err := json.Marshal(&val)
							common.Check(err)
							err = reco.send(messageType, res)
							common.Check(err)
						}

						//err = reco.send(messageType, m.Body)

						return
					}

				}
			}

			return

		}()

	}
}

//TODO ionert ping pong listeners

func (r Reco) send(msg int, m []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	defer logrus.Info("Unlocking mutex")
	return r.sock.WriteMessage(msg, m)
}

func (r Reco) sendJson(msg int, m interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.sock.WriteJSON(m)
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
