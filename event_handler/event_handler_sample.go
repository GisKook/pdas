package event_handler

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/pdas/base"
	"github.com/giskook/pdas/mqtt_srv"
	"github.com/giskook/pdas/pkg"
	"github.com/giskook/pdas/protocol"
)

func event_handler_blue_tooth_sample(c *gotcp.Conn, p *pkg.Prison_Packet) {
	//connection := c.GetExtraData().(*conn.Conn)
	sample_pkg := p.Packet.(*protocol.SamplePacket)
	mqtt_srv.GetMqttSocket().Send(&base.BluetoothRing{
		TagMac:  sample_pkg.TagMac,
		Rssi:    sample_pkg.Rssi,
		RingMac: sample_pkg.RingMac,
		DegreeX: sample_pkg.DegreeX,
		DegreeY: sample_pkg.DegreeY,
		DegreeZ: sample_pkg.DegreeZ,
		Bett:    sample_pkg.Bett,
	})
}
