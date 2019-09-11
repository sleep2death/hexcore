package hexcore

import (
	"log"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/sleep2death/hexcore/pb"
)

func TestProtobuf(t *testing.T) {
	// go Serve("127.0.0.1:8080")
	msg := &pb.Echo{Message: "Hello"}
	anyMsg, err := ptypes.MarshalAny(msg)

	if err != nil {
		t.Error(err)
	}
	anyMsg.TypeUrl = "hex/Echo"
	buf, err := proto.Marshal(anyMsg)

	if err != nil {
		t.Error(err)
	}

	anyMsg = &any.Any{}

	if proto.Unmarshal(buf, anyMsg); err != nil {
		t.Error(err)
	}

	assert.Equal(t, "hex/Echo", anyMsg.GetTypeUrl())
}

func TestServe(t *testing.T) {
	go Serve("localhost:9090")

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:9090/core", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.WriteMessage(websocket.CloseMessage, []byte{})
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	_, message, err := c.ReadMessage()
	assert.Equal(t, "Hello", string(message))
}
