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
		DegreeX: float64(sample_pkg.DegreeX) / 32,
		DegreeY: float64(sample_pkg.DegreeY) / 32,
		DegreeZ: float64(sample_pkg.DegreeZ) / 32,
		Bett:    sample_pkg.Bett,
	})
}
