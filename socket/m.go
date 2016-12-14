package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nemesisesq/ss_data_service/common"
	"math/rand"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func print_binary(s []byte) {
	fmt.Printf("recieved b:")
	for _, n := range s {
		fmt.Printf("%d,", n)

	}
	fmt.Printf("\n")
}

func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	common.Check(err)

	for {
		messageType, p, err := conn.ReadMessage()
		common.Check(err)

		print_binary(p)

		var x string = string(p[:])

		x = Reverse(x)

		go func() {

			for i := 1; i <= 20; i++ {
				v := rand.Intn(5)

				time.Sleep(time.Duration(v) * time.Second)
				err = conn.WriteMessage(messageType, []byte(x))
				common.Check(err)

			}
			return
		}()

	}
}
