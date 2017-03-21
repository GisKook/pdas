package zmq_socket

import (
	"github.com/giskook/pdas/pb"
	"github.com/golang/protobuf/proto"
	"time"
)

func (z *ZmqWorker) ProccessSendMsg() {
	for i, m := range z.ZmqLocateQueue {
		time_send := time.Now().Unix()
		time_recv := time_send
		serial_number := int32(0)
		tid := ""
		uuid := "das"

		location_type := Report.Location_EIndoor
		location_from := Report.Location_ETerminal
		indoor := &Report.Indoor{
			RingMac:    &m.RingMac,
			DegreeX:    &m.DegreeX,
			DegreeY:    &m.DegreeY,
			DegreeZ:    &m.DegreeZ,
			Batt:       &m.Bett,
			Alarm:      &m.WarnInfo,
			TagMacRssi: make([]*Report.TagMacRssi, len(m.Rssis)),
		}

		for j, _ := range m.Rssis {
			indoor.TagMacRssi[j].TagMac = &m.Rssis[j].TagMac
			indoor.TagMacRssi[j].Rssi = &m.Rssis[j].Rssi
		}

		locations := []*Report.Location{
			&Report.Location{
				Locationtype: &location_type,
				From:         &location_from,
				Indoor:       indoor,
			},
		}
		location_report := &Report.LocationReport{
			TimeSend:     &time_send,
			TimeRecv:     &time_recv,
			SerialNumber: &serial_number,
			Tid:          &tid,
			Uuid:         &uuid,
			Locations:    locations,
		}
		data, _ := proto.Marshal(location_report)

		GetZmqSocket().Send(string(data))
		delete(z.ZmqLocateQueue, i)
	}
}
