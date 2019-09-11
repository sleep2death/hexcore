package hexcore

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	"github.com/sleep2death/hexcore/router"
)

func getCoreHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		log.Print("received a request")
		if err != nil {
			log.Print("upgrade failed:", err)
			http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
			return
		}

		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read error:", err)
				break
			}

			switch mt {
			case websocket.BinaryMessage:
				writer, _ := c.NextWriter(websocket.BinaryMessage)
				handle(message, writer)
			case websocket.TextMessage:
				_ = c.WriteMessage(websocket.TextMessage, []byte("protocol not supported"))
				return
			case websocket.CloseMessage:
				// log.Printf("client closed the connection")
				return
			}

			log.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func handle(msg []byte, w io.WriteCloser) error {
	anyMsg := &any.Any{}
	if err := proto.Unmarshal(msg, anyMsg); err != nil {
		return err
	}

	path := anyMsg.GetTypeUrl()
	r.Serve(path, anyMsg.GetValue(), w)
	return nil
}

func initHandlers(r *router.Engine) {
	// r.Handle("/pb/echo", echoHandler)
}

var addr = flag.String("addr", "localhost:9090", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var r = router.Default()

// Serve -
func Serve(addr string) {
	initHandlers(r)

	http.HandleFunc("/core", getCoreHandler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
