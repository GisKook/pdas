package mqtt_srv

import (
	"encoding/json"
	"github.com/giskook/pdas/base"
	"log"
)

func (mqtt_socket *Mqtt_srv_socket) send(p *base.BluetoothRing) {
	payload, _ := json.Marshal(p)
	log.Println(payload)
	log.Println(mqtt_socket.Options.TopicPub)
	if token := mqtt_socket.Mqtt_socket.Publish(mqtt_socket.Options.TopicPub, 0, false, payload); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

}
