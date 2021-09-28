package hub

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type IpEntries []string

type hub struct {
	mu sync.Mutex
	v  map[string][]string
	p map[string]string
	c map[string] informCallback
}

func (h *hub) Register(name string, ips []string, port string) error {
	log.Printf("Registering new API: %s with IP: %s on port: %s", name, ips, port)

	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.v[name]
	if ok {
		return errors.New(fmt.Sprintf("api: %s is registered", name))
	}


	h.v[name] = ips
	h.p[name] = port

	return nil
}

func (h *hub) Unregister(name string) {
	log.Printf("Unregiser API: %s", name)
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.v, name)
	delete(h.p, name)
	delete(h.c, name)
}

func (h *hub) GetConnection(name string) ([]string, string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	v, ok := h.v[name]

	if !ok {
		return nil, ""
	}

	return v, h.p[name]
}

func (h *hub) Inform(name string, addr []string, port string, clientId uint64, clientSecret, serverSecret string) {
	c, ok := h.c[name]

	if ok {
		c(addr, port, clientId, clientSecret, serverSecret)
	}
}

func (h *hub) OnInform(name string, callback informCallback) {
	h.c[name] = callback
}

type informCallback func([]string, string, uint64, string, string)


type Hub interface {
	Register(name string, ip []string, port string) error
	Unregister(name string)
	GetConnection(name string) ([]string, string)
	Inform(name string, addr []string, port string, clientId uint64, clientSecret, serverSecret string)
	OnInform(name string, callback informCallback)
}

func GetHub () Hub {
	return &hub{
		v: map[string][]string{},
		c: map[string]informCallback{},
		p: map[string]string{},
	}
}
