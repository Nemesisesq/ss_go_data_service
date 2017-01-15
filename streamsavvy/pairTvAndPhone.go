package streamsavvy

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/streadway/amqp"
	"github.com/Sirupsen/logrus"
)

func pairHander(w http.ResponseWriter, r *http.Request) {
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

	var pairId string = r.Form.Get("id")

	for {
		messageType, p, err := reco.sock.ReadMessage()

		common.Check(err)

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

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				reco.send(messageType, d)
			}
		}()

		logrus.Printf(" [*] Waiting for logs. To exit press CTRL+C")
		<-forever




	}

}
