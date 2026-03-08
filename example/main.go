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
	display.OnError(func(objectId wlclient.Object, code uint32, message string) {
		fmt.Println("error:", objectId.Interface(), objectId.Id(), code, message)
	})

	// Bind global objects from registry
	globals := make(map[string]wlclient.Object)

	registry := display.GetRegistry()
	registry.OnGlobal(func(name uint32, iface string, version uint32) {
		globals[iface] = registry.Bind(name, iface, version)
	})

	<-display.Sync().OnDone(func(callbackData uint32) {})

	// Required global objects
	for _, ext := range []string{"wl_compositor", "xdg_wm_base", "wl_shm"} {
		if _, ok := globals[ext]; !ok {
			fmt.Println("required compositor extensions missing")
			os.Exit(1)
		}
	}

	compositor := wlclient.WlCompositor(globals["wl_compositor"])
	xdgWmBase := wlclient.XdgWmBase(globals["xdg_wm_base"])
	shm := wlclient.WlShm(globals["wl_shm"])

	// Prevent the application from being marked as "Not responding"
	xdgWmBase.OnPing(func(serial uint32) {
		xdgWmBase.Pong(serial)
	})

	// Create toplevel surface
	surface := compositor.CreateSurface()

	xdgSurface := xdgWmBase.GetXdgSurface(surface)
	xdgSurface.OnConfigure(func(serial uint32) {
		xdgSurface.AckConfigure(serial)
	})

	xdgToplevel := xdgSurface.GetToplevel()
	xdgToplevel.SetTitle("My First Wayland Application")
	xdgToplevel.SetAppId("example")

	surface.Commit()

	// Try adding server-side decorations
	if decorationManager, ok := globals["zxdg_decoration_manager_v1"]; ok {
		wlclient.ZxdgDecorationManagerV1(decorationManager).GetToplevelDecoration(xdgToplevel).SetMode(2)
	}

	// Create framebuffer
	width := 800
	height := 600

	// Allocate shared memory
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

	// Attach framebuffer to surface
	surface.Attach(shm.CreatePool(fd, int32(width*height*4)).CreateBuffer(0, int32(width), int32(height), int32(width*4), 0), 0, 0)
	surface.Commit()

	// Wait until window is closed
	<-xdgToplevel.OnClose(func() {})
}
