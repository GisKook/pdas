package event_handler

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/pdas/pkg"
	"github.com/giskook/pdas/protocol"
	"github.com/giskook/pdas/zmq"
)

func event_handler_blue_tooth_locate(c *gotcp.Conn, p *pkg.Prison_Packet) {
	//connection := c.GetExtraData().(*conn.Conn)
	locate_pkg := p.Packet.(*protocol.LocatePacket)
	zmq_socket.GetZmqWorker().PushLocate(locate_pkg.Location.RingMac, locate_pkg.Location)
}
