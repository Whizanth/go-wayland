package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"git.whizanth.com/go/wayland"
	"golang.org/x/sys/unix"
)

func main() {
	// Connect to compositor
	client, err := wayland.NewClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// First object is always the wl_display
	displayId := client.NewObjectId()

	// Send wl_display::get_registry
	registryId := client.NewObjectId()
	client.Request(displayId, 1, registryId)

	// Send wl_display::sync
	registryDoneCallbackId := client.NewObjectId()
	client.Request(displayId, 0, registryDoneCallbackId)

	// Query registry
	registry := make(map[string]uint32)
	for {
		if msg := client.Read(); msg.ObjectId == registryId && msg.OpCode == 0 {
			// Received wl_registry::global
			name := msg.ReadUint32()
			iface := msg.ReadString()
			version := msg.ReadUint32()

			// Send wl_registry::bind
			registry[iface] = client.NewObjectId()
			client.Request(registryId, 0, name, iface, version, registry[iface])
		} else if msg.ObjectId == registryDoneCallbackId && msg.OpCode == 0 {
			// Received wl_callback::done

			break
		}
	}

	// Check for required extensions
	for _, ext := range []string{"wl_compositor", "xdg_wm_base", "wl_shm"} {
		if _, ok := registry[ext]; !ok {
			fmt.Println("required compositor extensions missing")
			os.Exit(1)
		}
	}

	// Retrieve extensions
	compositorId := registry["wl_compositor"]
	xdgWmBaseId := registry["xdg_wm_base"]
	shmId := registry["wl_shm"]

	// Send wl_compositor::create_surface
	surfaceId := client.NewObjectId()
	client.Request(compositorId, 0, surfaceId)

	// Send xdg_wm_base::get_xdg_surface
	xdgSurfaceId := client.NewObjectId()
	client.Request(xdgWmBaseId, 2, xdgSurfaceId, surfaceId)

	// Send xdg_surface::get_toplevel
	xdgToplevelId := client.NewObjectId()
	client.Request(xdgSurfaceId, 1, xdgToplevelId)

	// Send wl_surface::commit
	client.Request(surfaceId, 6)

	for {
		if msg := client.Read(); msg.ObjectId == 1 && msg.OpCode == 0 {
			// Received wl_display::error

			fmt.Println("error:", msg.ReadUint32(), msg.ReadUint32(), msg.ReadString())
			break
		} else if msg.ObjectId == xdgSurfaceId && msg.OpCode == 0 {
			// Received xdg_surface::configure

			// Send xdg_surface::ack_configure
			client.Request(xdgSurfaceId, 4, msg.ReadUint32())

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
	shmPoolId := client.NewObjectId()
	client.WriteMsgUnix(wayland.NewMessage(shmId, 0, shmPoolId, int32(width*height*4)), fd)
	if err != nil {
		fmt.Printf("error creating pool: %v\n", err)
		os.Exit(1)
	}

	// Send wl_shm_pool::create_buffer
	bufferId := client.NewObjectId()
	client.Request(shmPoolId, 0, bufferId, 0, width, height, width*4, 0)

	// Send wl_surface::attach
	client.Request(surfaceId, 1, bufferId, 0, 0)

	// Send wl_surface::damage
	client.Request(surfaceId, 2, 0, 0, width, height)

	/*
		if decorationManagerId, ok := registry["zxdg_decoration_manager_v1"]; ok {
			// Send zxdg_decoration_manager_v1::get_toplevel_decoration
			toplevelDecorationId := client.NewObjectId()
			client.Request(decorationManagerId, 1, toplevelDecorationId, xdgToplevelId)

			// Send zxdg_toplevel_decoration_v1::set_mode
			client.Request(toplevelDecorationId, 1, uint32(2))

			for {
				if msg := client.Read(); msg.ObjectId == toplevelDecorationId && msg.OpCode == 0 {
					// Received xdg_toplevel_decoration::configure

					// Send xdg_surface::ack_configure
					client.Request(xdgSurfaceId, 4, msg.ReadUint32())
					break
				}
			}
		}
	*/

	// Send wl_surface::commit
	client.Request(surfaceId, 6)

	for {
		if msg := client.Read(); msg.ObjectId == 1 && msg.OpCode == 0 {
			// Received wl_display::error

			fmt.Println("error:", msg.ReadUint32(), msg.ReadUint32(), msg.ReadString())
		} else if msg.ObjectId == xdgWmBaseId && msg.OpCode == 0 {
			// Received xdg_wm_base::ping

			// Send xdg_wm_base::pong
			client.Request(xdgWmBaseId, 3, msg.ReadUint32())
		} else if msg.ObjectId == xdgToplevelId && msg.OpCode == 1 {
			// Received xdg_toplevel::close

			break
		}
	}
}
