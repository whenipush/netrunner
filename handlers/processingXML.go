package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

// Структуры, соответствующие структуре XML

type NmapRun struct {
	XMLName  xml.Name `xml:"nmaprun"`
	Args     string   `xml:"args,attr"`
	StartStr string   `xml:"startstr,attr"`
	Version  string   `xml:"version,attr"`
	Hosts    []Host   `xml:"host"`
	RunStats RunStats `xml:"runstats"`
}

type Host struct {
	StartTime string     `xml:"starttime,attr"`
	EndTime   string     `xml:"endtime,attr"`
	Status    HostStatus `xml:"status"`
	Address   Address    `xml:"address"`
	HostNames HostNames  `xml:"hostnames"`
	Ports     Ports      `xml:"ports"`
}

type HostStatus struct {
	State  string `xml:"state,attr"`
	Reason string `xml:"reason,attr"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

type HostNames struct {
	HostName []string `xml:"hostname"`
}

type Ports struct {
	ExtraPorts ExtraPorts `xml:"extraports"`
}

type ExtraPorts struct {
	State       string      `xml:"state,attr"`
	Count       int         `xml:"count,attr"`
	ExtraReason ExtraReason `xml:"extrareasons"`
}

type ExtraReason struct {
	Reason string `xml:"reason,attr"`
}

type RunStats struct {
	Finished Finished `xml:"finished"`
}

type Finished struct {
	TimeStr string `xml:"timestr,attr"`
	Summary string `xml:"summary,attr"`
	Elapsed string `xml:"elapsed,attr"`
	Exit    string `xml:"exit,attr"`
}

func ProcessingXML() {
	// Открытие XML файла
	xmlFile, err := os.Open("scan_results.xml")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer xmlFile.Close()

	// Декодирование XML в структуру Go
	var nmapRun NmapRun
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&nmapRun)
	if err != nil {
		fmt.Println("Ошибка при парсинге XML:", err)
		return
	}

	// Преобразование структуры Go в JSON
	jsonData, err := json.MarshalIndent(nmapRun, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при конвертации в JSON:", err)
		return
	}

	// Вывод JSON
	fmt.Println(string(jsonData))
}
