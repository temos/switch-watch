package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	RefreshDelay = time.Second * 2
)

const (
	EscMoveHome                  = "\x1b[H"
	EscEraseRestOfLine           = "\x1b[0K"
	EscEraseScreen               = "\x1b[2J"
	EscSwitchToAlternateScreen   = "\x1b[?1049h"
	EscSwitchFromAlternateScreen = "\x1b[?1049l"
)

type Port struct {
	Name        string
	Alias       string
	TxBytes     uint
	RxBytes     uint
	LastTxBytes uint
	LastRxBytes uint
	LastUpdate  time.Time
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <ip> <community>\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]
	community := os.Args[2]

	fmt.Print(EscSwitchToAlternateScreen, EscEraseScreen, EscMoveHome)
	defer fmt.Print(EscSwitchFromAlternateScreen)

	snmp, hostname, ports, err := detect(target, community)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer snmp.Conn.Close()

	app, updateUI := createApp(target, hostname, ports)

	go func() {
		for {
			time.Sleep(RefreshDelay)
			if err := UpdateRxTx(snmp, ports); err != nil {
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

func detect(target string, community string) (*gosnmp.GoSNMP, string, []*Port, error) {
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
		return nil, "", nil, fmt.Errorf("failed to connect to snmp: %w", err)
	}

	fmt.Println("getting hostname")
	hostname, err := GetHostname(&snmp)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	fmt.Println("detecting ports")
	ports, err := DetectPorts(&snmp)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to detect ports: %w", err)
	}

	return &snmp, hostname, ports, nil
}
