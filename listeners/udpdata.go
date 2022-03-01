package listeners

import (
	"fmt"
	"log"
	"net"

	"github.com/cyops-se/dd-proxy/db"
	"github.com/cyops-se/dd-proxy/engine"
	"github.com/cyops-se/dd-proxy/types"
	"github.com/nats-io/nats.go"
)

var prevMsg *types.DataMessage

type UDPDataListener struct {
	Port int `json:"port"`
	nc   *nats.Conn
}

func (listener *UDPDataListener) InitListener() {
	var err error
	listener.nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Printf("Unable to connect NATS")
	}

	listeners = append(listeners, listener)
	go listener.run()
}

func (listener *UDPDataListener) run() {
	port := 4357
	p := make([]byte, 2048*1024*8)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		db.Error("UDPDataListener error", "Failed to listen at address: %s, error: %s", addr.String(), err.Error())
		return
	}

	log.Println("UDP listening for DATA messages...")
	for {
		n, _, err := ser.ReadFromUDP(p)
		if err != nil {
			db.Error("UDPDataListener error", "Failed to read UDP data (ReadFromUDP), error: %s", err.Error())
			continue
		}

		// Send it over NATS
		fmt.Println(n)
		listener.nc.Publish("data", p[:n])
		engine.NotifySubscribers("data", string(p[:n]))
	}
}
