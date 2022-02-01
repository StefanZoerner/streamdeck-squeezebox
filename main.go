package main

import (
	"context"
	"github.com/StefanZoerner/streamdeck-squeezebox/plugin"
	"log"
)

func main() {
	ctx := context.Background()
	if err := plugin.Run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}
