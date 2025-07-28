package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"git.whizanth.com/go/wayland/wlclient"
	"golang.org/x/sys/unix"
)

func main() {
	// Connect to compositor
	client, err := wlclient.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go client.Listen()
	defer client.Close()

	// Get wl_display
	display := client.GetDisplay()

	// Get required global objects from registry
	registry := display.GetRegistry()
	registryDoneCallback := display.Sync()

	globals := make(map[string]wlclient.Object)

	registry.OnGlobal(func(name uint32, iface string, version uint32) {
		if iface == "wl_compositor" || iface == "xdg_wm_base" || iface == "wl_shm" {
			globals[iface] = registry.Bind(name, iface, version)
		}
	})

	<-registryDoneCallback.OnDone(func(callbackData uint32) {})

	for _, ext := range []string{"wl_compositor", "xdg_wm_base", "wl_shm"} {
		if _, ok := globals[ext]; !ok {
			fmt.Println("required compositor extensions missing")
			os.Exit(1)
		}
	}

	compositor := wlclient.WlCompositor(globals["wl_compositor"])
	xdgWmBase := wlclient.XdgWmBase(globals["xdg_wm_base"])
	shm := wlclient.WlShm(globals["wl_shm"])

	// Create toplevel surface
	surface := compositor.CreateSurface()
	xdgSurface := xdgWmBase.GetXdgSurface(surface)
	xdgToplevel := xdgSurface.GetToplevel()
	surface.Commit()

	display.OnError(func(objectId wlclient.Object, code uint32, message string) {
		fmt.Println("error:", objectId.Interface(), objectId.Id(), code, message)
	})

	<-xdgSurface.OnConfigure(func(serial uint32) {
		xdgSurface.AckConfigure(serial)
	})

	// Framebuffer size
	width := 800
	height := 600

	// Create shared memory
	fd, err := unix.MemfdCreate("wayland-framebuffer", 0)
	if err != nil {
		fmt.Printf("can't create shared memory: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)
	unix.Ftruncate(fd, int64(width*height*4))

	// Memory map framebuffer
	framebuffer, err := syscall.Mmap(fd, 0, width*height*4, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Munmap(framebuffer)

	// Fill framebuffer
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pix := (y*width + x) * 4
			framebuffer[pix+0] = 0xff // B
			framebuffer[pix+1] = 0x00 // G
			framebuffer[pix+2] = 0x00 // R
			framebuffer[pix+3] = 0xff // A
		}
	}

	// Send wl_shm::create_pool
	shmPool := shm.CreatePool(fd, int32(width*height*4))
	buffer := shmPool.CreateBuffer(0, int32(width), int32(height), int32(width*4), 0)
	surface.Attach(buffer, 0, 0)
	surface.Damage(0, 0, int32(width), int32(height))

	surface.Commit()

	xdgWmBase.OnPing(func(serial uint32) {
		xdgWmBase.Pong(serial)
	})

	<-xdgToplevel.OnClose(func() {})
}
