package mqtt_srv

import (
	"encoding/json"
	"github.com/giskook/pdas/base"
	"log"
)

func (mqtt_socket *Mqtt_srv_socket) Recv() {
	for {
		select {
		case p := <-mqtt_socket.RecvStringChan:
			mqtt_socket.MutexRecvChan.Lock()

			if mqtt_socket.RecvChan == nil {
				mqtt_socket.RecvChan = make(chan *base.BluetoothRing, 1024)
			}
			mqtt_socket.MutexRecvChan.Unlock()
			var btr base.BluetoothRing
			if err := json.Unmarshal([]byte(p), &btr); err != nil {
				log.Println(err)
			}
			mqtt_socket.RecvChan <- &btr
		}
	}
}
