package middleware

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/streadway/amqp"
	"net/http"
)

type RabbitMQAccessor struct {
	//rx_conn amqp.Connection
	//tx_conn amqp.Connection
	conn amqp.Connection

	//rx_url string
	//tx_url string
	url string
}

type RMQCH struct {
	//TX amqp.Channel
	//RX amqp.Channel
	Ch amqp.Channel
}

func NewRabbitMQAccesor(url string) (*RabbitMQAccessor, error) {
	conn, err := amqp.Dial(url)
	//rx_conn, err := amqp.Dial(rx_url)

	common.Check(err)

	logrus.Info("RabbitMQ Connected")

	return &RabbitMQAccessor{*conn, url}, nil
}

func (rmqa *RabbitMQAccessor) Set(request *http.Request, ch amqp.Channel) context.Context {
	channels := RMQCH{ch}

	return context.WithValue(request.Context(), "rabbitmq", channels)
}

type RabbitMQConnection struct {
	rmqa RabbitMQAccessor
}

func NewRabbitMQConnection(RabbitMQAccessor RabbitMQAccessor) *RabbitMQConnection {
	return &RabbitMQConnection{RabbitMQAccessor}
}

func (r *RabbitMQConnection) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		//tx_ch, err := r.rmqa.tx_conn.Channel()
		//rx_ch, err := r.rmqa.rx_conn.Channel()
		ch, err := r.rmqa.conn.Channel()
		common.Check(err)
		//defer tx_ch.Close()
		//defer rx_ch.Close()
		defer ch.Close()
		ctx := r.rmqa.Set(request, *ch)
		next(writer, request.WithContext(ctx))
	}
}
