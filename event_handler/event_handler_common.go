package event_handler

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/pdas/conf"
	"github.com/giskook/pdas/conn"
	"github.com/giskook/pdas/pkg"
	"github.com/giskook/pdas/protocol"
	"log"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	checkinterval := conf.GetConf().Server.ConnCheckInterval
	readlimit := conf.GetConf().Server.ReadLimit
	writelimit := conf.GetConf().Server.WriteLimit
	config := &conn.ConnConfig{
		ConnCheckInterval: uint16(checkinterval),
		ReadLimit:         uint16(readlimit),
		WriteLimit:        uint16(writelimit),
	}
	connection := conn.NewConn(c, config)

	c.PutExtraData(connection)

	connection.Do()
	conn.NewConns().Add(connection)

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*conn.Conn)
	log.Printf("close %d\n", connection.ID)
	connection.Close()
	conn.NewConns().Remove(connection)
	log.Println(conn.NewConns())
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	prison_pkg := p.(*pkg.Prison_Packet)
	switch prison_pkg.Type {
	case protocol.PROTOCOL_BLUETOOTH_SAMPLE:
		log.Println("PROTOCOL_BLUETOOTH_SAMPLE")
		event_handler_blue_tooth_sample(c, prison_pkg)
	case protocol.PROTOCOL_BLUETOOTH_LOCATE:
		log.Println("PROTOCOL_BLUETOOTH_LOCATE")
		event_handler_blue_tooth_locate(c, prison_pkg)
	}

	return true
}
