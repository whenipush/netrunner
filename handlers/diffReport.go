package handlers

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"netrunner/database"
	"netrunner/models"
	"os"
)

func Hash(str string) uint64 {
	alg := fnv.New64a()
	alg.Write([]byte(str))
	return alg.Sum64()
}

func diffReports(oldReport PentestReport, newReport PentestReport, hostIP string) PentestDiff {
	var diff PentestDiff

	newReportVulns := newReport.Hosts[hostIP].Vulns
	oldReportVulns := oldReport.Hosts[hostIP].Vulns
	add, rem := differentiate(oldReportVulns, newReportVulns)
	diff.Added = add
	diff.Removed = rem

	return diff
}

func (p *PentestReportController) diffReport() map[string]PentestDiff {
	var diff map[string]PentestDiff = make(map[string]PentestDiff)

	for _, host := range p.report.Hosts {
		var hostDB models.Host

		if err := database.DB.Preload("TaskList").Where("ip = ?", host.Ip).First(&hostDB).Error; err != nil {
			log.Printf("Failed to find host entry %s in database: %v", host.Ip, err)
			continue
		}
		// Если задачи не было до этого то все текущие уязвимости новые
		if len(hostDB.TaskList) < 2 {
			diff[host.Ip] = PentestDiff{
				Added: host.Vulns,
			}
		} else {
			last_task := hostDB.TaskList[len(hostDB.TaskList)-2]
			file, err := os.ReadFile(fmt.Sprintf("report/pentest/%s.xml.json", last_task.NumberTask))
			if err != nil {
				log.Printf("[diffReport] Failed to read report file: %v", err)
				continue
			}
			var file_data PentestReport
			err = json.Unmarshal(file, &file_data)
			if err != nil {
				log.Printf("[diffReport] Failed to parse report file: %v", err)
				continue
			}
			diff[host.Ip] = diffReports(file_data, p.report, host.Ip)
		}
	}
	return diff

}

func differentiate[K interface{}](data1 map[string]K, data2 map[string]K) (map[string]K, map[string]K) {
	var add_arr map[string]K = make(map[string]K)
	var rem_arr map[string]K = make(map[string]K)
	for key, vul := range data1 {
		_, ok := data2[key]
		if !ok {
			rem_arr[key] = vul
		}
	}
	for key, vul := range data2 {
		_, ok := data1[key]
		if !ok {
			add_arr[key] = vul
		}
	}
	return add_arr, rem_arr
}
