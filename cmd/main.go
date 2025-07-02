package main

import (
	"fmt"
	"os"
	"time"

	"git.whizanth.com/go/wayland"
)

func main() {
	client, err := wayland.NewClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go func() {
		time.Sleep(time.Second)

		displayId := client.NewObjectId()
		registryId := client.NewObjectId()
		registryDoneCallbackId := client.NewObjectId()

		// get_registry
		client.Request(displayId, 1, registryId)

		// wl_display_sync
		client.Request(displayId, 0, registryDoneCallbackId)
	}()

	for {
		client.Read()
	}
}
