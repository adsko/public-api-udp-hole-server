package main

import (
	"log"
	"os"
	"server/hub"
	"server/web"
)

func main() {
	h := hub.GetHub()
	switch os.Args[1] {
	case "run-server":
		err := web.Run(h)

		if err != nil {
			log.Panic("Could not run server: ", err)
		}

	}
}
