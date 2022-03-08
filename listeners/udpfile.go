package listeners

import (
	"log"
	"net"

	"github.com/cyops-se/dd-proxy/db"
	"github.com/nats-io/nats.go"
)

type UDPFileListener struct {
	Port int `json:"port"`
	nc   *nats.Conn
}

func (listener *UDPFileListener) InitListener() {
	var err error
	listener.nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Printf("Unable to connect NATS")
	}

	listeners = append(listeners, listener)
	go listener.run()
}

func (listener *UDPFileListener) run() {
	port := 4358
	p := make([]byte, 2048*1024*8)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		db.Error("UDPFileListener error", "Failed to listen at address: %s, error: %s", addr.String(), err.Error())
		return
	}

	log.Println("UDP listening for FILE messages...")
	for {
		n, _, err := ser.ReadFromUDP(p)
		if err != nil {
			db.Error("UDPFileListener error", "Failed to read UDP file (ReadFromUDP), error: %s", err.Error())
			continue
		}

		// Send it over NATS
		listener.nc.Publish("file", p[:n])
		// engine.NotifySubscribers("file", string(p[:n]))
	}
}
