package zmq_socket

import (
	"github.com/giskook/pdas/base"
	"github.com/giskook/pdas/conf"
	"github.com/giskook/pdas/pb"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

type _tag_mac_rssi struct {
	TagMac int64
	Rssi   float64
	Count  int
}

func (z *ZmqWorker) GetAvgRssi(in []*base.TagMacRssi) []*base.TagMacRssi {
	if in == nil {
		return nil
	}
	tag_mac_rssi := []*_tag_mac_rssi{
		&_tag_mac_rssi{
			TagMac: in[0].TagMac,
			Rssi:   in[0].Rssi,
			Count:  1,
		},
	}
	bhave := false
	for i := 1; i < len(in); i++ {
		bhave = false
		for j := 0; j < len(tag_mac_rssi); j++ {
			if tag_mac_rssi[j].TagMac == in[i].TagMac {
				bhave = true
				tag_mac_rssi[j].Rssi += in[i].Rssi
				tag_mac_rssi[j].Count++
			}
		}

		if !bhave {
			tag_mac_rssi = append(tag_mac_rssi, &_tag_mac_rssi{
				TagMac: in[i].TagMac,
				Rssi:   in[i].Rssi,
				Count:  1,
			})
		}
	}

	length := len(tag_mac_rssi)

	mac_rssis := make([]*base.TagMacRssi, length)

	for k, value := range tag_mac_rssi {
		mac_rssis[k].TagMac = value.TagMac
		mac_rssis[k].Rssi = value.Rssi / float64(value.Count)
	}

	return mac_rssis
}

func (z *ZmqWorker) PreProccessMsg() {
	for i, m := range z.ZmqLocateQueue {
		z.ZmqLocateQueue[i].Rssis = z.GetAvgRssi(m.Rssis)
	}
}

func (z *ZmqWorker) ProccessSendMsg() {
	z.PreProccessMsg()
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
