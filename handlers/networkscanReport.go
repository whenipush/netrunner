package handlers

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type NetworkScanReport struct {
	Hosts []NetworkScanHost `json:"hosts"`
}

type NetworkScanHost struct {
	Ip  string `json:"ip"`
	Mac string `json:"mac"`
	Os  string `json:"os"`
	Cpe string `json:"cpe"`
}

func ProcessNetworkScan(reportPath string) error {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return err
	}
	NmapData, err := Parse(data)
	if err != nil {
		return err
	}
	var scanHosts []NetworkScanHost = make([]NetworkScanHost, 0)
	for _, host := range NmapData.Hosts {
		var scanHost NetworkScanHost
		for _, addr := range host.Addresses {
			if addr.AddrType == "ipv4" {
				scanHost.Ip = addr.Addr
			} else if addr.AddrType == "mac" {
				scanHost.Mac = addr.Addr
			}
		}
		for _, os := range host.Os.OsMatches {
			acc, err := strconv.Atoi(os.Accuracy)
			if err != nil {
				log.Printf("Error converting accuracy to int")
				continue
			}
			if acc > 90 {
				scanHost.Os = os.Name
				for _, osclass := range os.OsClasses {
					if len(osclass.CPEs) > 0 {
						scanHost.Cpe = string(osclass.CPEs[0])
					}
				}
			}
		}
		scanHosts = append(scanHosts, scanHost)
	}
	file, err := os.Create(reportPath + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(scanHosts)
}
