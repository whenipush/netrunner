package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"netrunner/models"
	"os"
	"strconv"
)

type NetworkScanReport struct {
	Info  GeneralInfo       `json:"general_info"`
	Hosts []NetworkScanHost `json:"hosts"`
}

type NetworkScanHost struct {
	Ip  string `json:"ip"`
	Mac string `json:"mac"`
	Os  string `json:"os"`
	Cpe string `json:"cpe"`
}

func ProcessNetworkScan(task models.TaskStatus) error {
	reportPath := fmt.Sprintf("report/networkscan/%s.xml", task.NumberTask)

	data, err := os.ReadFile(reportPath)
	if err != nil {
		return err
	}
	NmapData, err := Parse(data)
	if err != nil {
		return err
	}
	var report NetworkScanReport
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

	report = NetworkScanReport{
		Info: GeneralInfo{
			TaskName:    task.Name,
			StartTime:   NmapData.StartStr,
			EndTime:     NmapData.RunStats.Finished.TimeStr,
			Version:     NmapData.Version,
			TimeElapsed: fmt.Sprintf("%.2f", NmapData.RunStats.Finished.Elapsed),
			Summary:     NmapData.RunStats.Finished.Summary,
			UpHosts:     NmapData.RunStats.Hosts.Up,
			DownHosts:   NmapData.RunStats.Hosts.Down,
			TotalHosts:  NmapData.RunStats.Hosts.Total,
		},
		Hosts: scanHosts,
	}

	file, err := os.Create(reportPath + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
