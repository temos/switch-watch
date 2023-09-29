package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	RefreshDelay = time.Second * 3
)

const (
	EscMoveHome        = "\x1b[H"
	EscEraseRestOfLine = "\x1b[0K"
	EscEraseScreen     = "\x1b[2J"
)

type Port struct {
	Name       string
	Alias      string
	Tx         uint
	Rx         uint
	LastTx     uint
	LastRx     uint
	LastUpdate time.Time
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <ip> <community>\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]
	community := os.Args[2]

	snmp := gosnmp.GoSNMP{
		Target:    target,
		Port:      161,
		Timeout:   time.Second * 2,
		Retries:   2,
		MaxOids:   64,
		Community: community,
		Version:   gosnmp.Version2c,
	}
	err := snmp.Connect()
	if err != nil {
		fmt.Println("failed to connect to snmp: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("getting hostname")
	hostname, err := GetHostname(&snmp)
	if err != nil {
		fmt.Println("failed to get hostname: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("detecting ports")
	ports, err := DetectPorts(&snmp)
	if err != nil {
		fmt.Println("failed to detect ports: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Waiting for data")

	app, updateUI := createApp(target, hostname, ports)

	go func() {
		for {
			time.Sleep(RefreshDelay)
			if err := UpdateRxTx(&snmp, ports); err != nil {
				fmt.Print(EscEraseScreen)
				fmt.Println("failed to update RX/TX: " + err.Error())
				os.Exit(1)
			}

			updateUI()
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}