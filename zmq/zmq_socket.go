package zmq_socket

import (
	"github.com/giskook/pdas/conf"
	zmq "github.com/pebbe/zmq3"
	"log"
)

type Zmq_Socket struct {
	Socket_PosUp *zmq.Socket
}

var G_Zmq_Socket *Zmq_Socket = nil

func GetZmqSocket() *Zmq_Socket {
	if G_Zmq_Socket == nil {
		socket, _ := zmq.NewSocket(zmq.PUSH)
		G_Zmq_Socket = &Zmq_Socket{
			Socket_PosUp: socket,
		}
	}

	return G_Zmq_Socket
}

func (z *Zmq_Socket) ConnectToServer() {
	err := z.Socket_PosUp.Connect(conf.GetConf().Zmq.PosUpAddr)
	if err != nil {
		log.Println(err)
	}
}

func (z *Zmq_Socket) Send(value string) {
	z.Socket_PosUp.Send(value, 0)
}
