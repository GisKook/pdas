package protocol

import (
	"github.com/giskook/pdas/base"
)

type LocatePacket struct {
	Location *base.LocateMessage
}

func (p *LocatePacket) Serialize() []byte {
	return nil
}

func ParseLocate(buffer []byte) *LocatePacket {
	reader, _, _ := ParseHeader(buffer)
	ring_mac := base.ReadMac(reader)
	_degree_x, _ := reader.ReadByte()
	degree_x := int8(_degree_x)
	_degree_y, _ := reader.ReadByte()
	degree_y := int8(_degree_y)
	_degree_z, _ := reader.ReadByte()
	degree_z := int8(_degree_z)
	bett, _ := reader.ReadByte()
	warn_info, _ := reader.ReadByte()
	rssi_count, _ := reader.ReadByte()

	var tag_mac int64
	var _rssi uint8
	rssis := make([]*base.TagMacRssi, rssi_count)

	for i := uint8(0); i < rssi_count; i++ {
		tag_mac = base.ReadMacInt(reader)
		_rssi, _ = reader.ReadByte()
		rssis[i].TagMac = tag_mac
		rssis[i].Rssi = float64(_rssi)
	}

	return &LocatePacket{
		Location: &base.LocateMessage{
			RingMac:  ring_mac,
			DegreeX:  float64(degree_x),
			DegreeY:  float64(degree_y),
			DegreeZ:  float64(degree_z),
			Bett:     int32(bett),
			WarnInfo: int32(warn_info),
			Rssis:    rssis,
		},
	}
}
