package zmq_socket

import (
	"github.com/giskook/pdas/base"
	"github.com/giskook/pdas/conf"
	"log"
	"time"
)

type ZmqLocatePair struct {
	Key   string
	Value *base.LocateMessage
}

type ZmqWorker struct {
	ZmqLocateQueue map[string]*base.LocateMessage
	ZmqLocateAdd   chan *base.LocateMessage
	ZmqLocateDel   chan string

	ticker *time.Ticker
}

var G_ZmqWorker *ZmqWorker = nil

func GetZmqWorker() *ZmqWorker {
	if G_ZmqWorker == nil {
		G_ZmqWorker = &ZmqWorker{
			ZmqLocateQueue: make(map[string]*base.LocateMessage),
			ZmqLocateAdd:   make(chan *base.LocateMessage),
			ZmqLocateDel:   make(chan string),
			ticker:         time.NewTicker(time.Duration(conf.GetConf().Zmq.ReportInterval) * time.Second),
		}
	}

	return G_ZmqWorker
}

func (z *ZmqWorker) Close() {
	close(z.ZmqLocateAdd)
	close(z.ZmqLocateDel)
	z.ticker.Stop()
}

func (z *ZmqWorker) PushLocate(key string, locate *base.LocateMessage) {
	z.ZmqLocateAdd <- locate
}

func (z *ZmqWorker) RemoveLocate(key string) {
	z.ZmqLocateDel <- key
}

func (z *ZmqWorker) Run() {
	go func() {
		for {
			select {
			case <-z.ticker.C:
				log.Println("ticker")
				z.ProccessSendMsg()
			case locate := <-z.ZmqLocateAdd:
				z.insert_locate(locate.RingMac, locate)
			}
		}
	}()
}

func (z *ZmqWorker) insert_locate(key string, value *base.LocateMessage) {
	log.Println("insert_locate")

	_, ok := z.ZmqLocateQueue[key]
	if ok {
		z.ZmqLocateQueue[key].DegreeX = value.DegreeX
		z.ZmqLocateQueue[key].DegreeY = value.DegreeY
		z.ZmqLocateQueue[key].DegreeZ = value.DegreeZ
		z.ZmqLocateQueue[key].Bett = value.Bett
		z.ZmqLocateQueue[key].WarnInfo = value.WarnInfo

		bHave := false
		for i, v := range value.Rssis {
			bHave = false
			for j, zv := range z.ZmqLocateQueue[key].Rssis {
				if v.TagMac == zv.TagMac {
					z.ZmqLocateQueue[key].Rssis[j].Rssi = v.Rssi
					bHave = true

					break
				}
			}

			if !bHave {
				z.ZmqLocateQueue[key].Rssis = append(z.ZmqLocateQueue[key].Rssis, value.Rssis[i])
			}
		}
		//	z.ZmqLocateQueue[key].Rssis = append(z.ZmqLocateQueue[key].Rssis, value.Rssis...)
	} else {
		z.ZmqLocateQueue[key] = value
	}
}
