package main

import (
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	//https://oidref.com/1.3.6.1.2.1.2.2.1.2
	OIDifDescr = ".1.3.6.1.2.1.2.2.1.2"

	//https://oidref.com/1.3.6.1.2.1.2.2.1.10
	OIDifInOctets = ".1.3.6.1.2.1.2.2.1.10"

	//https://oidref.com/1.3.6.1.2.1.2.2.1.16
	OIDifOutOctets = ".1.3.6.1.2.1.2.2.1.16"

	//https://oidref.com/1.3.6.1.2.1.1.5
	//the zero at the end is required even though it's not part of the OID
	OIDsysName = "1.3.6.1.2.1.1.5.0"

	//https://oidref.com/1.3.6.1.2.1.31.1.1.1.18.0
	OIDifAlias = "1.3.6.1.2.1.31.1.1.1.18"
)

var zeroTime time.Time

func DetectPorts(snmp *gosnmp.GoSNMP) ([]*Port, error) {
	ifacesResult, err := snmp.BulkWalkAll(OIDifDescr)
	if err != nil {
		return nil, err
	}

	aliasesResult, err := snmp.BulkWalkAll(OIDifAlias)
	if err != nil {
		return nil, err
	}

	rx, tx, updateTime, err := GetRxTx(snmp)
	if err != nil {
		return nil, err
	}

	ports := make([]*Port, len(ifacesResult))

	for i := range ports {
		ports[i] = &Port{
			Name:       string(ifacesResult[i].Value.([]byte)),
			Alias:      string(aliasesResult[i].Value.([]byte)),
			LastTx:     tx[i],
			LastRx:     rx[i],
			LastUpdate: updateTime,
		}
	}

	return ports, nil
}

func GetRxTx(snmp *gosnmp.GoSNMP) ([]uint, []uint, time.Time, error) {
	rxResult, err := snmp.BulkWalkAll(OIDifInOctets)
	if err != nil {
		return nil, nil, zeroTime, err
	}

	txResult, err := snmp.BulkWalkAll(OIDifOutOctets)
	if err != nil {
		return nil, nil, zeroTime, err
	}

	if len(rxResult) != len(txResult) {
		panic("expected RX and TX results to have the same length")
	}

	rx := make([]uint, len(rxResult))
	tx := make([]uint, len(txResult))
	for i := range rx {
		//'* 8' to convert from octets (bytes) to bits
		rx[i] = rxResult[i].Value.(uint) * 8
		tx[i] = txResult[i].Value.(uint) * 8
	}

	return rx, tx, time.Now(), nil
}

func UpdateRxTx(snmp *gosnmp.GoSNMP, ports []*Port) error {
	rx, tx, updateTime, err := GetRxTx(snmp)
	if err != nil {
		return err
	}

	for i, port := range ports {
		secondsDelta := updateTime.Sub(port.LastUpdate).Seconds()
		port.Rx = uint(float64(diffWithWrap(port.LastRx, rx[i])) / secondsDelta)
		port.Tx = uint(float64(diffWithWrap(port.LastTx, tx[i])) / secondsDelta)

		port.LastRx = rx[i]
		port.LastTx = tx[i]

		port.LastUpdate = updateTime
	}

	return nil
}

func GetHostname(snmp *gosnmp.GoSNMP) (string, error) {
	result, err := snmp.Get([]string{OIDsysName})
	if err != nil {
		return "", err
	}

	return string(result.Variables[0].Value.([]byte)), err
}

// calculates the difference between a base value and a new value accounting for counter wrapping
func diffWithWrap(base, new uint) uint {
	if new >= base {
		//no wrap
		return new - base
	}

	//wrap
	const MaxUint = ^uint(0)
	return (MaxUint - base) + new
}
