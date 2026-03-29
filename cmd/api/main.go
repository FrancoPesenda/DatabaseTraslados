package main

import (
	"log"
	"os"

	"database/cmd/api"
)

func main() {
	if err := api.NewFxApp().Run(); err != nil {
		log.New(os.Stdout, "", log.LstdFlags|log.LUTC).Fatalf("run app: %v", err)
	}
}

