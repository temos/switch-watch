package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <ip> <community>\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]
	community := os.Args[2]
	fmt.Print(EscEraseScreen, EscMoveHome)

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

	err = snmp.BulkWalk("1.3.6.1.2.1.17", func(dataUnit gosnmp.SnmpPDU) error {
		if dataUnit.Type == gosnmp.OctetString {
			fmt.Println(dataUnit.Name, dataUnit.Type, string(dataUnit.Value.([]byte)))
		} else {
			fmt.Println(dataUnit.Name, dataUnit.Type, dataUnit.Value)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return

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

	for {
		time.Sleep(time.Second * 3)
		if err := UpdateRxTx(&snmp, ports); err != nil {
			fmt.Print(EscEraseScreen)
			fmt.Println("failed to update RX/TX: " + err.Error())
			os.Exit(1)
		}

		printUI(target, hostname, ports)
	}
}

func printUI(target string, hostname string, ports []*Port) {
	fmt.Print(EscMoveHome)

	fmt.Printf("Monitoring %s (%s)%s\n", target, hostname, EscEraseRestOfLine)

	width := maxNameWidth(ports)
	for _, port := range ports {
		namePadding := strings.Repeat(" ", width-len(port.Name))
		fmt.Printf("%s%s:\t%s%s\t%s%s\n", port.Name, namePadding, toReadable(port.Rx), EscEraseRestOfLine, toReadable(port.Tx), EscEraseRestOfLine)
	}
}
