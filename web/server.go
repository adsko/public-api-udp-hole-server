package web

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"server/hub"
	"time"
)

var r *rand.Rand // Rand for this package.
var clientId uint64 = 1

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Run (h hub.Hub) error {
	port := os.Getenv("PORT")

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	http.HandleFunc("/register", registerAPI(h))
	http.HandleFunc("/connect", connectAPI(h))
	http.HandleFunc("/health-check", healthCheck)
	return http.ListenAndServe(":" + port, nil)
}

func healthCheck(w http.ResponseWriter, r * http.Request) {
	data := []byte("Server works")
	w.WriteHeader(200)
	w.Write(data)
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

		h.OnInform(registerData.Name, func(addr []string, port string, clientId uint64, clientSecret, serverSecret string) {
			oc := openConnection{Addr: addr, Port: port, ClientID: clientId, ClientSecret: clientSecret, ServerSecret: serverSecret}
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
			ClientSecret: RandomString(16),
			ServerSecret: RandomString(16),
		}

		responseJson, err := json.Marshal(response)

		if err != nil {
			log.Panic(err)
		}

		h.Inform(m.Server, m.Addr, m.Port, clientId, response.ClientSecret, response.ServerSecret)
		clientId += 1
		w.WriteHeader(200)
		w.Write(responseJson)
	}
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}