package main

import (
	ninja "github.com/ninjasphere/go-ninja"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	hid "github.com/GeertJohan/go.hid"
	"log"
	"os"
	"os/signal"
	"strings"
	"bytes"
	"encoding/binary"
)

type Payload struct {
	Buttons      uint32
	LeftX        int8
	LeftY        int8
	RightX       int8
	RightY       int8
}

const NONE = ""
var BUTTONS = []string {
	NONE, NONE, NONE, NONE, NONE, NONE,
	"plus", "minus", "home",
	NONE, NONE,
	"dup", "ddown", "dleft", "dright",
	"a", "b", "x", "y",
	"l", "r",
	"zl", "zr",
	"leftstick", "rightstick",
}

func main() {

	conn, err := ninja.Connect("com.ninjablocks.wii")

	bus, err := conn.AnnounceDriver("com.ninjablocks.wii", "driver-wii", getCurDir())
	if err != nil {
		log.Fatalf("Could not get driver bus: %s", err)
	}

	devices, err := hid.Enumerate(0, 0)
	if err != nil {
		log.Fatalf("Could not list devices: %s", err)
	}

	for _, deviceInfo := range devices {
		if strings.Contains(deviceInfo.Product, "Wiimote") {
			spew.Dump("Found a wiimote", deviceInfo)
			device, err := deviceInfo.Device()
			check(err)

			controller, _ := createController(device, bus)
			controller.Read()
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}

type XboxController struct {
  device *hid.Device
}

func createController(device *hid.Device, bus *ninja.DriverBus) (*XboxController, error) {
	controller := &XboxController{
		device: device,
	}

	return controller, nil
}

func (this XboxController) Read() {

	var payloadBytes [10]byte

	_, err := this.device.Read(payloadBytes[:])

	check(err)

	var payload Payload

  err = binary.Read(bytes.NewReader(payloadBytes[:]), binary.LittleEndian, &payload)

	spew.Dump(payload)

	log.Print("----")

	for index,element := range BUTTONS {
		pressed := payload.Buttons>>uint(index)&1
		if pressed > 0 {
			log.Printf("Pressed Button %s", element)
		}
	}

	this.Read()
}

func getCurDir() string {
	pwd, err := os.Getwd()
	check(err)
	return pwd + "/"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
