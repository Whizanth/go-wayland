# Go Wayland

This repository contains a pure Go implementation of the Wayland windowing protocol. It is explicitly free of any third‑party dependencies.

## Repository Layout

The `wayland` package implements the Wayland wire protocol and core types. `wayland.Client` can be used to send raw `wayland.Message` objects over the wire. This can be useful if you're trying to understand how the protocol works or are looking to build your own implementation (check out commit `b8d8973` to see the bare‑minimum amount of code needed to create a window) but in practice you'll want to use `wlclient.Client` instead.

The `wlclient` package provides idiomatic Go‑style bindings for the Wayland protocol. It is generated from the [Wayland specification XML files released by FreeDesktop](https://gitlab.freedesktop.org/wayland). The version contained in this repo might not always be up to date with the latest Wayland specifications, so you might want to generate it yourself.

The `scanner` package can be used to generate the `wlclient` package. It expects that the [wayland](https://gitlab.freedesktop.org/wayland/wayland) and [wayland‑protocols](https://gitlab.freedesktop.org/wayland/wayland‑protocols) directories in the current working directory contain the linked repositories.

The repository currently only contains an implementation of the client‑side of the Wayland protocol. A server‑side implementation might be developed later, but it is currently unclear to me whether a pure Go Wayland compositor could be practically viable due to missing graphics acceleration. While the Wayland protocol requires all compositors to support "dumb" memory‑based framebuffers (wl_shm), it doesn't require all clients to do so, so there might be some clients (perhaps games?) that only support EGLStreams or GBM.

## State of Development

This project is still work in progress; however, most features are implemented, and it should be ready for most clients.

What's left to implement:
* Server‑side protocol bindings for writing a compositor.
* More descriptive errors. This requires some refactoring to keep track of the names of opcodes, events, etc. during runtime.
* Support for arrays. A small but important change for feature completion.
* Automatic generation of up‑to‑date bindings using Actions.

No breaking changes are planned for how the API can be interacted with (mapping to objects, events, etc.), except *potentially* simplifying the way the connection to the server is initially established (which would be a simple and one‑time change).

## How To Use The Bindings

The Wayland wire protocol is fairly simple.

Actions are performed on objects by calling their methods. Objects have no properties/fields. Other implementations (like libwayland) might be different, but the Go bindings map Wayland objects into Go objects with actual methods.

Methods in the Wayland protocol can't have return values. Instead, they accept arguments of the generic `new_id` type. This is essentially a pointer to an already declared variable that the result will be written to, similar to how `json.Unmarshal` asks for a pointer to where the object should be unmarshaled to rather than returning the unmarshaled object. The Go bindings map `new_id` arguments into return values for the sake of cleaner and more idiomatic code.

Objects can also have events that can be listened to. The Go bindings make this really simple to do, by exposing `On` methods that accept a callback function. `On` methods also return a channel that you can use to await an event.

The `example` directory contains a minimal implementation of what's needed for a client to create a window.

## Contributions

This project is not looking for contributions, as it's an explicit goal not to use third‑party code.

## Licensing

All original work found in this repository is released under the Zero‑Clause BSD (0BSD) license, with the exception of the .git directory. The latter is not part of the project and must not be redistributed.

The Wayland protocol is licensed under the MIT license. It is unclear to me whether bindings generated from the XML files count as derivative work. However, as the author of the go‑wayland project, I don't impose any *additional* licensing limitations on them.

Note that the official implementation of the Go standard libraries is licensed under the 3‑Clause BSD License, meaning any binaries including it must retain attribution.
