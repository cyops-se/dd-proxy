package listeners

import (
	"log"
	"reflect"
	"time"

	"github.com/cyops-se/dd-proxy/db"
	"github.com/cyops-se/dd-proxy/engine"
	"github.com/cyops-se/dd-proxy/types"
)

var typeRegistry = make(map[string]reflect.Type)
var listeners []types.IListener

func RegisterType(i types.IListener) {
	typename := reflect.TypeOf(i).String()
	log.Println("Registering listener type:", typename, i)
	typeRegistry[typename] = reflect.TypeOf(i)
}

func RunDispatch() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C
		engine.NotifySubscribers("heartbeat", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func Init() {
	var listeners []types.Listener

	db.DB.Find(&listeners)
	udpdata := &UDPDataListener{}
	udpdata.InitListener()

	udpmeta := &UDPMetaListener{}
	udpmeta.InitListener()

	udpfile := &UDPFileListener{}
	udpfile.InitListener()

	go RunDispatch()
}
