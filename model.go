package main

import "time"

type Port struct {
	Name        string
	Tx          uint
	Rx          uint
	LastTxBytes uint
	LastRxBytes uint
	LastUpdate  time.Time
}

func maxNameWidth(ports []*Port) int {
	max := 0
	for _, port := range ports {
		if len(port.Name) > max {
			max = len(port.Name)
		}
	}

	return max
}
