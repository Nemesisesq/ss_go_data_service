package streamsavvy

import (
	"net/http"

	"github.com/nemesisesq/ss_data_service/common"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/streadway/amqp"
)

func HandleRecomendations(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	rmqc := r.Context().Value("rabbitmq").(middleware.RMQCH)

	common.Check(err)

	for {
		messageType, p, err := conn.ReadMessage()
		common.Check(err)

		PublishShowInfo(p, rmqc)

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

			msgs, err := rmqc.RX.Consume(
				rx_q.Name, // queue
				"",        // consumer
				true,      // auto-ack
				false,     // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)

			for {
				select {
				case m := <-msgs:
					conn.WriteMessage(1, m.Body)

				}
			}
		}()

		p = append([]byte("Hello World"), p...)

		if err = conn.WriteMessage(messageType, p); err != nil {
			common.Check(err)
		}
	}
}

func PublishShowInfo(p []byte, rmqc middleware.RMQCH) {
	tx_q, err := rmqc.TX.QueueDeclare(
		"reco_engine",
		false,
		false,
		false,
		false,
		nil,
	)
	common.Check(err)

	err = rmqc.TX.Publish(
		"",        // exchange
		tx_q.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        p,
		})

	common.Check(err)
}
