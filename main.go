package main

import (
	"flag"
	"fmt"
	"github.com/Banrai/PiScan/scanner" 
	"github.com/Unknwon/macaron"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var Printcode string
 
 

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}
func worker(start chan bool) {
	heartbeat := time.Tick(3 * time.Second)
	for {
		select {
		// … do some stuff
		case <-heartbeat:
			broadcastMessage(websocket.TextMessage, []byte("heartbeat"))
		}
	}
}

func main() {
	m := macaron.Classic()
	m.Get("/", func() string {
		return `<html><body><script src='//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js'></script>
    <ul id=messages></ul><form><input id=message><input type="submit" id=send value=Send></form>
    <script>
    var c=new WebSocket('ws://localhost:4000/wsbarcode');
    c.onopen = function(){
      c.onmessage = function(response){
        console.log(response.data);
        var newMessage = $('<li>').text(response.data);
        $('#messages').append(newMessage);
        $('#message').val('');
      };
      $('form').submit(function(){
        c.send($('#message').val());
        return false;
      });
    }
    </script></body></html>`
	})

	m.Get("/wsbarcode", func(w http.ResponseWriter, r *http.Request) {
		log.Println(ActiveClients)
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		client := ws.RemoteAddr()
		sockCli := ClientConn{ws, client}
		addClient(sockCli)

		for {

			log.Println(len(ActiveClients), ActiveClients)
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				deleteClient(sockCli)
				log.Println("bye")
				log.Println(err)
				return
			}

			broadcastMessage(messageType, p)

		}
	})

	startheartbeat := make(chan bool)
	go worker(startheartbeat)
 

 
	var (
		device string
	)

	flag.StringVar(&device, "device", scanner.SCANNER_DEVICE, fmt.Sprintf("The '/dev/input/event' device associated with your scanner (defaults to '%s')", scanner.SCANNER_DEVICE))

	processScanFn := func(barcode string) {
		fmt.Println("newcode:" + barcode)
		Printcode = barcode
		broadcastMessage(websocket.TextMessage, []byte(barcode))
	}

	errorFn := func(e error) {
		fmt.Println(e)
	}
	fmt.Println("capturing barcode scanner")
	go scanner.ScanForever(device, processScanFn, errorFn)

	fmt.Println("web server running")
	Printcode = "test"
 
	m.Get("/httpbarcode", func() (barcode string) {
		barcode = Printcode
		Printcode = ""
		return
	})

	m.Run()

}
