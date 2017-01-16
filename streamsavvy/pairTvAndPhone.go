package streamsavvy

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/streadway/amqp"
	"github.com/Sirupsen/logrus"
	"fmt"
)

func PairHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	common.Check(err)
	rmqc := r.Context().Value("rabbitmq").(middleware.RMQCH)

	reco := Reco{sock: conn}

	var pairId string = r.URL.Query().Get("id")
	logrus.Info(pairId, "pair id")
	if pairId == "" {
		fmt.Fprint(w, "no pair id found, please provide pair id")
		return
	}

	err = rmqc.Ch.ExchangeDeclare(
		pairId,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	common.Check(err)



	q, err := rmqc.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = rmqc.Ch.QueueBind(
		q.Name,
		"",
		pairId,
		false,
		nil,
	)

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

	for {
		messageType, p, err := reco.sock.ReadMessage()

		common.Check(err)
		err = rmqc.Ch.Publish(
			pairId,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        p,
			},
		)



		//forever := make(chan bool)

		go func() {
			for d := range msgs {
				logrus.Info("getting broadcast")
				reco.send(messageType, append([]byte(q.Name),d.Body...))
			}
		}()

		logrus.Printf(" [*] Waiting for logs. To exit press CTRL+C")





	}

}
