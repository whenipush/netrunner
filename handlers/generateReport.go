package handlers

//
//import (
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"strings"
//)
//
//// TODO:ВОТ ЭТИ 2 ХЕРНИ В КОММЕНТЕ ИДУТ ИЗ proccesingBDU.go, надо все структуры в отдельным файл пихнуть для понимания
//
////type Vulnerabilities struct {
////	XMLName xml.Name `xml:"vulnerabilities"`
////	Vul     []Vul    `xml:"vul"`
////}
////
////type Vul struct {
////	Identifier    string `xml:"identifier"`
////	IdentifierCVE string `xml:"identifier-cve,omitempty"`
////	Name          string `xml:"name"`
////	Description   string `xml:"description"`
////	Severity      string `xml:"severity"`
////	Solution      string `xml:"solution"`
////}
//
//// Структуры для работы с XML (BDU и Nmap)
//type NmapRun struct {
//	XMLName  xml.Name `xml:"nmaprun"`
//	StartStr string   `xml:"startstr,attr"`
//	Version  string   `xml:"version,attr"`
//	RunStats RunStats `xml:"runstats"`
//	Hosts    []Host   `xml:"host"`
//}
//
//type RunStats struct {
//	Finished Finished `xml:"finished"`
//	Hosts    Hosts    `xml:"hosts"`
//}
//
//type Finished struct {
//	TimeStr string `xml:"timestr,attr"`
//	Elapsed string `xml:"elapsed,attr"`
//	Summary string `xml:"summary,attr"`
//}
//
//type Hosts struct {
//	Up    string `xml:"up,attr"`
//	Down  string `xml:"down,attr"`
//	Total string `xml:"total,attr"`
//}
//
//type Host struct {
//	Address []Address `xml:"address"`
//	Ports   Ports     `xml:"ports"`
//}
//
//type Address struct {
//	Addr string `xml:"addr,attr"`
//}
//
//type Ports struct {
//	Port []Port `xml:"port"`
//}
//
//type Port struct {
//	PortID  string   `xml:"portid,attr"`
//	Scripts []Script `xml:"script"`
//}
//
//type Script struct {
//	Tables []Table `xml:"table"`
//}
//
//type Table struct {
//	Key string `xml:"key,attr"`
//}
//
//type Report struct {
//	GeneralInfo GeneralInfo      `json:"general_info"`
//	Details     []HostVulDetails `json:"details"`
//}
//
//type GeneralInfo struct {
//	Start   string `json:"start"`
//	End     string `json:"end"`
//	Version string `json:"version"`
//	Elapsed string `json:"elapsed"`
//	Summary string `json:"summary"`
//	Up      string `json:"up"`
//	Down    string `json:"down"`
//	Total   string `json:"total"`
//}
//
//type HostVulDetails struct {
//	IP          string  `json:"ip"`
//	Port        string  `json:"port"`
//	CVE         string  `json:"cve"`
//	Identifier  *string `json:"identifier,omitempty"` // Новый элемент
//	Name        *string `json:"name,omitempty"`
//	Description *string `json:"description,omitempty"`
//	Severity    *string `json:"severity,omitempty"`
//	Solution    *string `json:"solution,omitempty"`
//}
//
//// Чтение BDU XML и загрузка в карту
//func loadBDUXML(xmlPath string) (map[string]Vul, error) {
//	file, err := os.Open(xmlPath)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка открытия BDU XML файла: %v", err)
//	}
//	defer file.Close()
//
//	data, err := ioutil.ReadAll(file)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка чтения BDU XML файла: %v", err)
//	}
//
//	var vulnerabilities Vulnerabilities
//	err = xml.Unmarshal(data, &vulnerabilities)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка парсинга BDU XML файла: %v", err)
//	}
//
//	bduMap := make(map[string]Vul)
//	for _, vul := range vulnerabilities.Vul {
//		if vul.IdentifierCVE != "" {
//			cveKey := strings.ToUpper(strings.TrimSpace(vul.IdentifierCVE))
//			bduMap[cveKey] = vul
//		}
//	}
//
//	// Отладка загруженных данных
//	fmt.Println("Загруженные CVE из BDU XML:")
//	for cve, vul := range bduMap {
//		fmt.Printf("CVE=%s, Name=%s\n", cve, vul.Name)
//	}
//
//	return bduMap, nil
//}
//
//// Парсинг Nmap XML
//func parseNmapXML(xmlFile string) ([][3]string, []string, error) {
//	file, err := os.Open(xmlFile)
//	if err != nil {
//		return nil, nil, fmt.Errorf("ошибка открытия Nmap XML файла: %v", err)
//	}
//	defer file.Close()
//
//	decoder := xml.NewDecoder(file)
//	var nmapRun NmapRun
//	err = decoder.Decode(&nmapRun)
//	if err != nil {
//		return nil, nil, fmt.Errorf("ошибка парсинга Nmap XML: %v", err)
//	}
//
//	// Извлекаем общую информацию
//	startScan := nmapRun.StartStr
//	endScan := nmapRun.RunStats.Finished.TimeStr
//	versionNmap := nmapRun.Version
//	elapsed := nmapRun.RunStats.Finished.Elapsed
//	summary := nmapRun.RunStats.Finished.Summary
//
//	hostsInfo := nmapRun.RunStats.Hosts
//	up := hostsInfo.Up
//	down := hostsInfo.Down
//	total := hostsInfo.Total
//
//	info := []string{startScan, endScan, versionNmap, elapsed, summary, up, down, total}
//
//	// Отладка общей информации
//	fmt.Printf("Общая информация: Start=%s, End=%s, Version=%s, Elapsed=%s, Summary=%s, Up=%s, Down=%s, Total=%s\n",
//		startScan, endScan, versionNmap, elapsed, summary, up, down, total)
//
//	var data [][3]string
//	for _, host := range nmapRun.Hosts {
//		ip := ""
//		if len(host.Address) > 0 {
//			ip = host.Address[0].Addr
//		}
//
//		for _, port := range host.Ports.Port {
//			cveFound := false
//			for _, script := range port.Scripts {
//				for _, table := range script.Tables {
//					if len(table.Key) > 0 && strings.HasPrefix(table.Key, "CVE") {
//						data = append(data, [3]string{ip, port.PortID, table.Key})
//						cveFound = true
//					}
//				}
//			}
//
//			if !cveFound {
//				data = append(data, [3]string{ip, port.PortID, "N/A"})
//			}
//		}
//	}
//
//	// Отладка извлеченных данных
//	for _, entry := range data {
//		fmt.Printf("Nmap Data: IP=%s, Port=%s, CVE=%s\n", entry[0], entry[1], entry[2])
//	}
//
//	return data, info, nil
//}
//
//// Генерация JSON отчета
//func generateJSONReport(nmapData [][3]string, bduMap map[string]Vul, info []string, outputFile string) error {
//	report := Report{
//		GeneralInfo: GeneralInfo{
//			Start:   info[0],
//			End:     info[1],
//			Version: info[2],
//			Elapsed: info[3],
//			Summary: info[4],
//			Up:      info[5],
//			Down:    info[6],
//			Total:   info[7],
//		},
//		Details: []HostVulDetails{},
//	}
//
//	for _, entry := range nmapData {
//		ip, port, cve := entry[0], entry[1], strings.ToUpper(strings.TrimSpace(entry[2]))
//
//		if cve != "N/A" {
//			if cveData, exists := bduMap[cve]; exists {
//				// Если CVE найдено в BDU
//				identifier := cveData.Identifier
//				report.Details = append(report.Details, HostVulDetails{
//					IP:          ip,
//					Port:        port,
//					CVE:         cve,
//					Identifier:  &identifier, // Добавляем Identifier
//					Name:        &cveData.Name,
//					Description: &cveData.Description,
//					Severity:    &cveData.Severity,
//					Solution:    &cveData.Solution,
//				})
//			} else {
//				// Если CVE отсутствует в BDU
//				msg := fmt.Sprintf("Информация для CVE %s отсутствует в BDU", cve)
//				report.Details = append(report.Details, HostVulDetails{
//					IP:          ip,
//					Port:        port,
//					CVE:         cve,
//					Description: &msg,
//				})
//			}
//		} else {
//			// Если CVE не найдено (N/A)
//			report.Details = append(report.Details, HostVulDetails{
//				IP:   ip,
//				Port: port,
//				CVE:  "N/A",
//			})
//		}
//	}
//
//	file, err := os.Create(outputFile)
//	if err != nil {
//		return fmt.Errorf("ошибка создания JSON файла: %v", err)
//	}
//	defer file.Close()
//
//	encoder := json.NewEncoder(file)
//	encoder.SetIndent("", "  ")
//	if err := encoder.Encode(report); err != nil {
//		return fmt.Errorf("ошибка записи JSON файла: %v", err)
//	}
//
//	fmt.Printf("JSON отчет сохранен в %s\n", outputFile)
//	return nil
//}
//
//// TODO: ПОПОЗЖЕ
//func generateHTMLReport(jsonFile, htmlFile string) error {
//	// Чтение JSON-файла
//	file, err := os.Open(jsonFile)
//	if err != nil {
//		return fmt.Errorf("ошибка открытия JSON файла: %v", err)
//	}
//	defer file.Close()
//
//	data, err := ioutil.ReadAll(file)
//	if err != nil {
//		return fmt.Errorf("ошибка чтения JSON файла: %v", err)
//	}
//
//	var report Report
//	if err := json.Unmarshal(data, &report); err != nil {
//		return fmt.Errorf("ошибка парсинга JSON: %v", err)
//	}
//
//	// Открываем HTML-файл для записи
//	html, err := os.Create(htmlFile)
//	if err != nil {
//		return fmt.Errorf("ошибка создания HTML файла: %v", err)
//	}
//	defer html.Close()
//	//<script src="https://cdn.tailwindcss.com"></script>
//	// Генерация HTML-структуры
//	html.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
//	html.WriteString("<meta charset=\"UTF-8\">\n<title>NETRUNNER Report</title>\n")
//	html.WriteString("<script src=\"https://cdn.tailwindcss.com\"></script>\n")
//	html.WriteString("</head>\n<body class=\"bg-gray-100 text-gray-800\">\n")
//	html.WriteString("<div class=\"bg-blue-900 text-white py-6 px-6 flex items-center justify-center shadow-md\">\n")
//	html.WriteString("<h1 class=\"text-3xl font-extrabold uppercase tracking-widest\">NETRUNNER</h1>\n</div>\n")
//
//	// Общая информация
//	html.WriteString("<div class=\"p-6\">\n<h2 class=\"text-lg font-bold\">Информация о сканировании</h2>\n")
//	html.WriteString(fmt.Sprintf("<p>Начало сканирования: %s</p>\n", report.GeneralInfo.Start))
//	html.WriteString(fmt.Sprintf("<p>Завершение сканирования: %s</p>\n", report.GeneralInfo.End))
//	html.WriteString(fmt.Sprintf("<p>Версия NMap: %s</p>\n", report.GeneralInfo.Version))
//	html.WriteString(fmt.Sprintf("<p>Время сканирования (секунды): %s</p>\n", report.GeneralInfo.Elapsed))
//	html.WriteString(fmt.Sprintf("<p>Активные хосты: %s</p>\n", report.GeneralInfo.Up))
//	html.WriteString(fmt.Sprintf("<p>Неактивные хосты: %s</p>\n", report.GeneralInfo.Down))
//	html.WriteString(fmt.Sprintf("<p>Всего хостов: %s</p>\n", report.GeneralInfo.Total))
//	html.WriteString("</div>\n")
//
//	// Детали хостов
//	html.WriteString("<div class=\"p-6\">\n<h2 class=\"text-lg font-bold\">Детальная информация</h2>\n")
//	for _, detail := range report.Details {
//		html.WriteString("<div class=\"bg-white shadow-md rounded-lg p-4 mb-6\">\n")
//		html.WriteString(fmt.Sprintf("<p><strong>IP:</strong> %s</p>\n", detail.IP))
//		html.WriteString(fmt.Sprintf("<p><strong>Port:</strong> %s</p>\n", detail.Port))
//		html.WriteString(fmt.Sprintf("<p><strong>CVE:</strong> %s</p>\n", detail.CVE))
//		if detail.Identifier != nil {
//			html.WriteString(fmt.Sprintf("<p><strong>Identifier:</strong> %s</p>\n", *detail.Identifier))
//		}
//		if detail.Name != nil {
//			html.WriteString(fmt.Sprintf("<p><strong>Name:</strong> %s</p>\n", *detail.Name))
//		}
//		if detail.Description != nil {
//			html.WriteString(fmt.Sprintf("<p><strong>Description:</strong> %s</p>\n", *detail.Description))
//		}
//		if detail.Severity != nil {
//			html.WriteString(fmt.Sprintf("<p><strong>Severity:</strong> %s</p>\n", *detail.Severity))
//		}
//		if detail.Solution != nil {
//			html.WriteString(fmt.Sprintf("<p><strong>Solution:</strong> %s</p>\n", *detail.Solution))
//		}
//		html.WriteString("</div>\n")
//	}
//	html.WriteString("</div>\n")
//
//	// Завершение HTML
//	html.WriteString("<div class=\"bg-blue-900 text-white text-center py-4 mt-auto w-full\">&copy; 2025 NETRUNNER. All rights reserved.</div>\n")
//	html.WriteString("</body>\n</html>")
//
//	fmt.Printf("HTML отчет успешно создан: %s\n", htmlFile)
//	return nil
//}
//
//func GenerateReportBDU() {
//	bduXMLPath := "bdu_updated.xml"
//	nmapXMLPath := "nmap.xml"
//	outputJSON := "report.json"
//
//	// Загружаем данные из BDU XML
//	bduMap, err := loadBDUXML(bduXMLPath)
//	if err != nil {
//		fmt.Println("Ошибка загрузки BDU XML:", err)
//		return
//	}
//
//	// Парсим Nmap XML
//	nmapData, info, err := parseNmapXML(nmapXMLPath)
//	if err != nil {
//		fmt.Println("Ошибка парсинга Nmap XML:", err)
//		return
//	}
//
//	// Генерируем JSON отчет
//	err = generateJSONReport(nmapData, bduMap, info, outputJSON)
//	if err != nil {
//		fmt.Println("Ошибка создания JSON отчета:", err)
//		return
//	}
//
//	jsonFile := "report.json"
//	htmlFile := "report.html"
//
//	if err := generateHTMLReport(jsonFile, htmlFile); err != nil {
//		fmt.Println("Ошибка генерации HTML отчета:", err)
//	}
//}
