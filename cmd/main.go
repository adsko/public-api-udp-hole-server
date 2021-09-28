package main

import (
	"log"
	"server/hub"
	"server/web"
)

func main() {
	h := hub.GetHub()
	err := web.Run(h)

	if err != nil {
		log.Panic("Could not run server: ", err)
	}

}
