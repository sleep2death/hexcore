package hexcore

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"github.com/sleep2death/hexcore/pb"
	"github.com/sleep2death/hexcore/router"
)

func getCoreHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, r.Header)
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
				writer, err := c.NextWriter(websocket.BinaryMessage)
				if err != nil {
					log.Println("get writer error:", err)
					return
				}
				rt.Serve(message, writer)
				// r.Serve(message, writer)

			case websocket.TextMessage:
				_ = c.WriteMessage(websocket.TextMessage, []byte("protocol not supported"))
				return
			}
		}
	}
}

func initHandlers(r *router.Engine) {
	r.Handle("/hex/echo", echoHandler)
}

func echoHandler(ctx *router.Context) {
	defer ctx.Writer.Close()

	echo := &pb.Echo{}
	if err := proto.Unmarshal(ctx.Value, echo); err != nil {
		ctx.Error(err)
		return
	}

	echo.Message = "echo: " + echo.GetMessage()
	any, err := ptypes.MarshalAny(echo)
	if err != nil {
		ctx.Error(err)
		return
	}
	any.TypeUrl = "hex/echo"
	buf, err := proto.Marshal(any)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Writer.Write(buf)
}

var upgrader = websocket.Upgrader{} // use default options

var rt = router.Default()
var srv *http.Server

// Serve -
func Serve(addr string, done <-chan struct{}) {
	initHandlers(rt)

	m := http.NewServeMux()
	s := http.Server{Addr: addr, Handler: m}

	m.HandleFunc("/core", getCoreHandler())

	go func() {
		log.Print("server started...")
		if err := s.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("server has shutdown...")
}
