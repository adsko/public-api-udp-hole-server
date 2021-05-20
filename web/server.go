package web

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"server/hub"
)

var clientId uint64 = 1

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Run (h hub.Hub) error {
	log.Println("Running HUB server on port: 10000")
	http.HandleFunc("/register", registerAPI(h))
	http.HandleFunc("/connect", connectAPI(h))
	return http.ListenAndServe(":10000", nil)
}


func registerAPI(h hub.Hub) func(w http.ResponseWriter, r * http.Request) {
	return func (w http.ResponseWriter, r * http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		registerData := &register{}
		err = conn.ReadJSON(registerData)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.Register(registerData.Name, registerData.Addr, registerData.Port)

		if err != nil {
			log.Println(err)
			return
		}

		h.OnInform(registerData.Name, func(addr []string, port string, clientId uint64) {
			oc := openConnection{Addr: addr, Port: port, ClientID: clientId, Secret: "0123456789123456"}
			log.Printf("Opening connection from %s to %s", registerData.Name, addr)
			conn.WriteJSON(oc)
		})

		defer h.Unregister(registerData.Name)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}
}

func connectAPI(h hub.Hub) func(w http.ResponseWriter, r * http.Request) {
	return func (w http.ResponseWriter, r * http.Request) {
		var m connectBody
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(reqBody, &m)

		if err != nil {
			log.Panic(err)
		}

		if len(m.Server) == 0 {
			log.Panic("No server provided")
		}

		ips, port := h.GetConnection(m.Server)

		if ips == nil {
			w.WriteHeader(404)
			w.Write([]byte("{}"))
			return
		}

		response := connectResponse{
			Addr: ips,
			Port: port,
			ClientID: clientId,
			Secret: "0123456789123456",
		}

		responseJson, err := json.Marshal(response)

		if err != nil {
			log.Panic(err)
		}

		h.Inform(m.Server, m.Addr, m.Port, clientId)
		clientId += 1
		w.WriteHeader(200)
		w.Write(responseJson)
	}
}