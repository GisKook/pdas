package mqtt_srv

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/giskook/pdas/base"
	"github.com/giskook/pdas/conf"
	"github.com/giskook/pdas/db"
	"log"
	"sync"
	"time"
)

var G_Mqtt_Socket *Mqtt_srv_socket

type Mqtt_srv_socket struct {
	Mqtt_socket    *MQTT.Client
	Options        *conf.MqttConfiguration
	SendChan       chan *base.BluetoothRing
	MutexRecvChan  sync.Mutex
	RecvChan       chan *base.BluetoothRing
	RecvStringChan chan string
	ticker         *time.Ticker
}

func MsgRecvfun(client *MQTT.Client, msg MQTT.Message) {
	GetMqttSocket().RecvStringChan <- string(msg.Payload())
}

func NewMqttSocket(mqtt_opts *conf.MqttConfiguration) *Mqtt_srv_socket {
	opts := MQTT.NewClientOptions().AddBroker(mqtt_opts.Addr).SetClientID(mqtt_opts.ClientID)
	opts.SetCleanSession(true)

	mqtt_socket := MQTT.NewClient(opts)
	if token := mqtt_socket.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Connected to MqttServer")

	if token := mqtt_socket.Subscribe(mqtt_opts.TopicSub, 0, MsgRecvfun); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
	G_Mqtt_Socket = &Mqtt_srv_socket{
		Mqtt_socket:    mqtt_socket,
		Options:        mqtt_opts,
		SendChan:       make(chan *base.BluetoothRing, 256),
		RecvChan:       nil,
		RecvStringChan: make(chan string, 256),
		ticker:         time.NewTicker(time.Duration(5) * time.Second),
	}

	return G_Mqtt_Socket
}

func (mqtt_socket *Mqtt_srv_socket) Send(payload *base.BluetoothRing) {
	log.Println("Mqtt add send list")
	mqtt_socket.SendChan <- payload
}

func (mqtt_socket *Mqtt_srv_socket) Proccess() {
	for {
		select {
		case send_pkg := <-mqtt_socket.SendChan:
			mqtt_socket.send(send_pkg)
			//		case <-mqtt_socket.ticker.C:
			//			mqtt_socket.MutexRecvChan.Lock()
			//			db.GetDBSocket().TransactionChan <- mqtt_socket.RecvChan
			//			mqtt_socket.RecvChan = nil
			//			mqtt_socket.MutexRecvChan.Unlock()
		}
	}
}

func (mqtt_socket *Mqtt_srv_socket) ProccessSub() {
	for {
		select {
		case <-mqtt_socket.ticker.C:
			mqtt_socket.MutexRecvChan.Lock()
			db.GetDBSocket().TransactionChan <- mqtt_socket.RecvChan
			mqtt_socket.RecvChan = nil
			mqtt_socket.MutexRecvChan.Unlock()
		}
	}
}
func (mqtt_socket *Mqtt_srv_socket) Close() {
	mqtt_socket.Mqtt_socket.Disconnect(0)
	close(mqtt_socket.SendChan)
	close(mqtt_socket.RecvStringChan)
	mqtt_socket.ticker.Stop()

}

func GetMqttSocket() *Mqtt_srv_socket {
	return G_Mqtt_Socket
}