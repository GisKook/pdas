package zmq_socket

import (
	"github.com/giskook/pdas/conf"
	"github.com/giskook/pdas/pb"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

func (z *ZmqWorker) ProccessSendMsg() {
	for i, m := range z.ZmqLocateQueue {
		time_send := time.Now().Unix() * 1000
		time_recv := time_send
		serial_number := int32(0)
		uuid := "das"

		tag_mac_rssi_count := len(m.Rssis)
		if tag_mac_rssi_count > int(conf.GetConf().Zmq.MaxReportCount) {
			tag_mac_rssi_count = int(conf.GetConf().Zmq.MaxReportCount)
		}
		location_type := Report.Location_EIndoor
		location_from := Report.Location_ETerminal
		ring_mac := strings.ToLower(m.RingMac)
		indoor := &Report.Indoor{
			RingMac:    &ring_mac,
			DegreeX:    &m.DegreeX,
			DegreeY:    &m.DegreeY,
			DegreeZ:    &m.DegreeZ,
			Batt:       &m.Bett,
			Alarm:      &m.WarnInfo,
			TagMacRssi: make([]*Report.TagMacRssi, tag_mac_rssi_count),
		}

		for j := 0; j < tag_mac_rssi_count; j++ {
			indoor.TagMacRssi[j] = &Report.TagMacRssi{
				TagMac: &m.Rssis[j].TagMac,
				Rssi:   &m.Rssis[j].Rssi,
			}
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
			Tid:          &m.RingMac,
			Uuid:         &uuid,
			Locations:    locations,
		}
		data, _ := proto.Marshal(location_report)

		GetZmqSocket().Send(string(data))
		delete(z.ZmqLocateQueue, i)
	}
}
