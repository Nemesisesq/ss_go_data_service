package streamsavvy

//import (
//	"net/http"
//	"github.com/gorilla/websocket"
//	log "github.com/Sirupsen/logrus"
//)
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}
//
//func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	ws, err := upgrader.Upgrade(w, r, nil)
//
//	if err != nil {
//		m := "Unable to upgrade to websockets"
//		log.WithField("err", err).Println(m)
//		http.Error(w, m, http.StatusBadRequest)
//		return
//	}
//
//	id := r.register(ws)
//

//}
