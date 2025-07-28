package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Wayland XML schema

type Protocol struct {
	Copyright  string      `xml:"copyright"`
	Interfaces []Interface `xml:"interface"`
}

type Interface struct {
	Name     string   `xml:"name,attr"`
	Version  int      `xml:"version,attr"`
	Requests []Method `xml:"request"`
	Events   []Method `xml:"event"`
	Enums    []Enum   `xml:"enum"`
}

type Method struct {
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Description Description `xml:"description"`
	Args        []Argument  `xml:"arg"`
}

type Argument struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	Interface string `xml:"interface,attr"`
	Enum      string `xml:"enum,attr"`
	Summary   string `xml:"summary,attr"`
}

type Enum struct {
	Name        string      `xml:"name,attr"`
	Description Description `xml:"description"`
	Entries     []Entry     `xml:"entry"`
}

type Entry struct {
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
	Summary string `xml:"summary,attr"`
}

type Description struct {
	Summary string `xml:"summary,attr"`
	Full    string `xml:",chardata"`
}

func toPascalCase(str string) string {
	var builder strings.Builder

	parts := strings.Split(str, "_")
	for _, part := range parts {
		runes := []rune(part)

		if len(runes) > 0 {
			builder.WriteString(strings.ToUpper(string(runes[0:1])))
			builder.WriteString(string(runes[1:]))
		}
	}

	return builder.String()
}

func toCamelCase(str string) string {
	var builder strings.Builder
	parts := strings.Split(str, "_")
	for i, part := range parts {
		if i == 0 {
			builder.WriteString(part)
			continue
		}

		runes := []rune(part)

		if len(runes) > 0 {
			builder.WriteString(strings.ToUpper(string(runes[0:1])))
			builder.WriteString(string(runes[1:]))
		}
	}

	result := builder.String()
	if result == "range" {
		return "rnge"
	} else if result == "interface" {
		return "iface"
	}
	return result
}

func main() {
	var builder strings.Builder
	builder.WriteString(`package wlclient

import "git.whizanth.com/go/wayland"

type Object struct {
	client  *Client
	id      uint32
	iface   string
	version uint32
}

func (object Object) Id() uint32 {
	return object.id
}

func (object Object) Interface() string {
	return object.iface
}

func New() (*Client, error) {
	client, err := wayland.NewClient()
	if err != nil {
		return nil, err
	}

	result := &Client{Client: client}
	result.display.client = result
	result.display.id = result.NewObjectId()
	result.display.iface = "wl_display"
	return result, nil
}

type Client struct {
	*wayland.Client
	display WlDisplay
}

func (client *Client) GetDisplay() WlDisplay {
	return client.display
}

`)

	generateClient(&builder, filepath.Join("wayland", "protocol", "wayland.xml"))

	for _, protocolStage := range []string{"stable", "staging", "experimental"} {
		protocolDirs, _ := os.ReadDir(filepath.Join("wayland-protocols", protocolStage))
		for _, dir := range protocolDirs {
			files, _ := os.ReadDir(filepath.Join("wayland-protocols", protocolStage, dir.Name()))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".xml") {
					generateClient(&builder, filepath.Join("wayland-protocols", protocolStage, dir.Name(), file.Name()))
				}
			}
		}
	}

	os.WriteFile(filepath.Join("wlclient", "generated.go"), []byte(builder.String()), 0755)
}

func generateClient(builder *strings.Builder, xmlPath string) {
	content, err := os.ReadFile(xmlPath)
	if err != nil {
		fmt.Printf("unable to read %s: %v\n", xmlPath, err)
		return
	}

	var proto Protocol
	if err := xml.Unmarshal(content, &proto); err != nil {
		fmt.Printf("unable to unmarshal %s: %v\n", xmlPath, err)
		return
	}

	for _, iface := range proto.Interfaces {
		builder.WriteString("type " + toPascalCase(iface.Name) + " Object\n")
		builder.WriteString("\n")

		for opCode, request := range iface.Requests {
			var argsBuilder strings.Builder
			var returnsBuilder strings.Builder
			var newsBuilder strings.Builder
			var msgArgsBuilder strings.Builder
			var fdBuilder strings.Builder
			var returnBuilder strings.Builder

			fd := 0
			newId := 0
			args := 0
			returns := 0

			for _, arg := range request.Args {
				if arg.Type == "new_id" {
					if newId > 0 {
						returnsBuilder.WriteString(", ")
						returnBuilder.WriteString(", ")
					}

					if arg.Interface == "" {
						arg.Interface = "Object"

						if args > 0 {
							argsBuilder.WriteString(", ")
						}
						argsBuilder.WriteString("iface string, version uint32")
						args++
						msgArgsBuilder.WriteString(", iface, version")
					}

					returnsBuilder.WriteString("" + toPascalCase(arg.Interface))
					returnBuilder.WriteString("" + toCamelCase(arg.Name))
					returns++

					newsBuilder.WriteString("	" + toCamelCase(arg.Name) + " := ")

					if arg.Interface != "Object" {
						newsBuilder.WriteString(toPascalCase(arg.Interface) + "(")
					}

					newsBuilder.WriteString("Object{\n")
					newsBuilder.WriteString("		client: object.client,\n")
					newsBuilder.WriteString("		id: object.client.NewObjectId(),\n")
					newsBuilder.WriteString(`		iface: "` + arg.Interface + "\",\n")
					newsBuilder.WriteString("	}")

					if arg.Interface != "Object" {
						newsBuilder.WriteString(")")
					}

					newsBuilder.WriteString("\n\n")

					msgArgsBuilder.WriteString(", " + toCamelCase(arg.Name) + ".id")
					newId++

					continue

				}

				if args > 0 {
					argsBuilder.WriteString(", ")
				}

				argsBuilder.WriteString(toCamelCase(arg.Name) + " ")

				if arg.Type == "string" {
					argsBuilder.WriteString("string")
				} else if arg.Type == "uint" {
					argsBuilder.WriteString("uint32")
				} else if arg.Type == "int" {
					argsBuilder.WriteString("int32")
				} else if arg.Type == "enum" {
					argsBuilder.WriteString("uint32")
				} else if arg.Type == "object" {
					argsBuilder.WriteString(toPascalCase(arg.Interface))
				} else if arg.Type == "fixed" {
					argsBuilder.WriteString("wayland.Fixed")
				} else if arg.Type == "fd" {
					argsBuilder.WriteString("int")

					if fd > 0 {
						fdBuilder.WriteString(", ")
					}

					fdBuilder.WriteString(toCamelCase(arg.Name))
					fd++
					args++
					continue
				} else {
					fmt.Println("(!) unsupported type: " + arg.Type)
				}

				msgArgsBuilder.WriteString(", " + toCamelCase(arg.Name))
				if arg.Type == "object" {
					msgArgsBuilder.WriteString(".id")
				}

				args++
			}

			builder.WriteString("func (object " + toPascalCase(iface.Name) + ") " + toPascalCase(request.Name) + "(")
			builder.WriteString(argsBuilder.String())
			builder.WriteString(") ")

			if returns == 1 {
				builder.WriteString(returnsBuilder.String() + " ")
			} else if returns > 1 {
				builder.WriteString("(" + returnsBuilder.String() + ") ")
			}

			builder.WriteString("{\n")
			builder.WriteString(newsBuilder.String())

			if fd > 0 {
				builder.WriteString("	object.client.Write(wayland.NewMessage(object.id, " + strconv.Itoa(opCode) + msgArgsBuilder.String() + ").WithFds(" + fdBuilder.String() + "))\n")
			} else {
				builder.WriteString("	object.client.Write(wayland.NewMessage(object.id, " + strconv.Itoa(opCode) + msgArgsBuilder.String() + "))\n")
			}

			if returns > 0 {
				builder.WriteString("\n")
				builder.WriteString("	return " + returnBuilder.String() + "\n")
			}

			builder.WriteString("}\n")
			builder.WriteString("\n")
		}

		for opCode, event := range iface.Events {
			var args1Builder strings.Builder
			var args2Builder strings.Builder

			args := 0

			for _, arg := range event.Args {
				if args > 0 {
					args1Builder.WriteString(", ")
					args2Builder.WriteString(", ")
				}

				args1Builder.WriteString(toCamelCase(arg.Name) + " ")

				if arg.Type == "string" {
					args1Builder.WriteString("string")
					args2Builder.WriteString("message.ReadString()")
				} else if arg.Type == "uint" {
					args1Builder.WriteString("uint32")
					args2Builder.WriteString("message.ReadUint32()")
				} else if arg.Type == "int" || arg.Type == "enum" {
					args1Builder.WriteString("int32")
					args2Builder.WriteString("message.ReadInt32()")
				} else if arg.Type == "fixed" {
					args1Builder.WriteString("wayland.Fixed")
					args2Builder.WriteString("message.ReadFixed()")
				} else if arg.Type == "object" {
					if arg.Interface != "" {
						args1Builder.WriteString(toPascalCase(arg.Interface))
						args2Builder.WriteString(toPascalCase(arg.Interface) + "(Object{object.client, message.ReadUint32(), message.ReadString(), message.ReadUint32()})")
					} else {
						args1Builder.WriteString("Object")
						args2Builder.WriteString("Object{object.client, message.ReadUint32(), message.ReadString(), message.ReadUint32()}")
					}
				} else if arg.Type == "fd" {
					args1Builder.WriteString("int")
					args2Builder.WriteString("message.ReadFd()")
				} else if arg.Type == "array" {
					args1Builder.WriteString("[]uint32")
					args2Builder.WriteString("message.ReadArray()")
				} else if arg.Type == "new_id" {
					args1Builder.WriteString(toPascalCase(arg.Interface))
					args2Builder.WriteString(toPascalCase(arg.Interface) + `(Object{client: object.client, id: object.client.NewObjectId(), iface: "` + arg.Interface + `"})`)
				}

				args++
			}

			builder.WriteString("func (object " + toPascalCase(iface.Name) + ") On" + toPascalCase(event.Name) + "(listener func(")
			builder.WriteString(args1Builder.String())
			builder.WriteString(")) chan struct{} {\n")
			builder.WriteString("	return object.client.On(object.id, " + strconv.Itoa(opCode) + ", func(message *wayland.Message) {\n")
			builder.WriteString("		listener(" + args2Builder.String() + ")\n")
			builder.WriteString("	})\n")
			builder.WriteString("}\n")
			builder.WriteString("\n")
		}
	}
}
