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

	// Get wl_display
	display := client.GetDisplay()

	// Get required global objects from registry
	registry := display.GetRegistry()
	registryDoneCallback := display.Sync()

	globals := make(map[string]wlclient.Object)
	for {
		if msg := client.Read(); msg.ObjectId == wlclient.Object(registry).Id() && msg.OpCode == 0 {
			// Received wl_registry::global
			name := msg.ReadUint32()
			iface := msg.ReadString()
			version := msg.ReadUint32()

			if iface == "wl_compositor" || iface == "xdg_wm_base" || iface == "wl_shm" {
				globals[iface] = registry.Bind(name, iface, version)
			}
		} else if msg.ObjectId == wlclient.Object(registryDoneCallback).Id() && msg.OpCode == 0 {
			// Received wl_callback::done
			break
		}
	}

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

	for {
		if msg := client.Read(); msg.ObjectId == 1 && msg.OpCode == 0 {
			// Received wl_display::error
			fmt.Println("error:", msg.ReadUint32(), msg.ReadUint32(), msg.ReadString())
			break
		} else if msg.ObjectId == wlclient.Object(xdgSurface).Id() && msg.OpCode == 0 {
			// Received xdg_surface::configure
			xdgSurface.AckConfigure(msg.ReadUint32())
			break
		}
	}

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

	for {
		if msg := client.Read(); msg.ObjectId == 1 && msg.OpCode == 0 {
			// Received wl_display::error
			fmt.Println("error:", msg.ReadUint32(), msg.ReadUint32(), msg.ReadString())
		} else if msg.ObjectId == wlclient.Object(xdgWmBase).Id() && msg.OpCode == 0 {
			// Received xdg_wm_base::ping
			xdgWmBase.Pong(msg.ReadUint32())
		} else if msg.ObjectId == wlclient.Object(xdgToplevel).Id() && msg.OpCode == 1 {
			// Received xdg_toplevel::close
			break
		}
	}
}
