package handlers

//import (
//	"encoding/csv"
//	"encoding/xml"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//	"regexp"
//
//	"github.com/xuri/excelize/v2"
//)
//
//// -------------------------НЕ ТРОГАТЬ-------------------------------------------------
//// Структуры для работы с XML
//type Vulnerabilities struct {
//	XMLName xml.Name `xml:"vulnerabilities"`
//	Vul     []Vul    `xml:"vul"`
//}
//
//type Vul struct {
//	Identifier     string         `xml:"identifier"`
//	IdentifierCVE  string         `xml:"identifier-cve,omitempty"`
//	Name           string         `xml:"name"`
//	Description    string         `xml:"description"`
//	VulnerableSoft VulnerableSoft `xml:"vulnerable_software"`
//	Environment    Environment    `xml:"environment"`
//	CWE            CWE            `xml:"cwe"`
//	IdentifyDate   string         `xml:"identify_date"`
//	CVSS           CVSS           `xml:"cvss"`
//	CVSS3          CVSS3          `xml:"cvss3"`
//	Severity       string         `xml:"severity"`
//	Solution       string         `xml:"solution"`
//	VulStatus      string         `xml:"vul_status"`
//	ExploitStatus  string         `xml:"exploit_status"`
//	FixStatus      string         `xml:"fix_status"`
//	Sources        string         `xml:"sources"`
//	Other          string         `xml:"other"`
//	VulIncident    string         `xml:"vul_incident"`
//	VulClass       string         `xml:"vul_class"`
//}
//
//type VulnerableSoft struct {
//	Soft []Soft `xml:"soft"`
//}
//
//type Soft struct {
//	Vendor   string `xml:"vendor"`
//	Name     string `xml:"name"`
//	Version  string `xml:"version"`
//	Platform string `xml:"platform"`
//	Types    Types  `xml:"types"`
//}
//
//type Types struct {
//	Type string `xml:"type"`
//}
//
//type Environment struct {
//	OS OS `xml:"os"`
//}
//
//type OS struct {
//	Vendor   string `xml:"vendor"`
//	Name     string `xml:"name"`
//	Version  string `xml:"version"`
//	Platform string `xml:"platform"`
//}
//
//type CWE struct {
//	Identifier string `xml:"identifier"`
//}
//
//type CVSS struct {
//	Vector CVSSVector `xml:"vector"`
//}
//
//type CVSSVector struct {
//	Score string `xml:"score,attr"`
//	Value string `xml:",chardata"`
//}
//
//type CVSS3 struct {
//	Vector CVSS3Vector `xml:"vector"`
//}
//
//type CVSS3Vector struct {
//	Score string `xml:"score,attr"`
//}
//
//// UpdateDatabase - Основная функция модуля
//func UpdateDatabaseBDU(excelPath, xmlPath string) (string, error) {
//	csvPath := "cve_mapping.csv"           // Временный файл для CSV
//	outputXML := "vulners/bdu_updated.xml" // Итоговый обновленный XML
//
//	// Генерация CSV из Excel
//	if err := generateCSVFromExcel(excelPath, csvPath); err != nil {
//		return "", fmt.Errorf("ошибка генерации CSV: %v", err)
//	}
//
//	// Нормализация XML
//	if err := normalizeXML(xmlPath, outputXML, csvPath); err != nil {
//		return "", fmt.Errorf("ошибка нормализации XML: %v", err)
//	}
//
//	// Удаляем временный CSV файл
//	_ = os.Remove(csvPath)
//
//	return outputXML, nil
//}
//
//// Генерация CSV из Excel
//func generateCSVFromExcel(excelPath, csvPath string) error {
//	// Проверяем формат файла
//	if filepath.Ext(excelPath) != ".xlsx" {
//		return fmt.Errorf("некорректный формат файла: %s. Требуется .xlsx", filepath.Ext(excelPath))
//	}
//
//	// Открываем файл Excel
//	f, err := excelize.OpenFile(excelPath)
//	if err != nil {
//		return fmt.Errorf("ошибка открытия файла Excel: %v", err)
//	}
//
//	// Получаем первую страницу
//	sheetName := f.GetSheetName(0)
//	rows, err := f.GetRows(sheetName)
//	if err != nil {
//		return fmt.Errorf("ошибка получения строк из Excel: %v", err)
//	}
//
//	// Регулярное выражение для поиска CVE
//	cveRegex := regexp.MustCompile(`CVE-\d{4}-\d{4,7}`)
//
//	// Создаём CSV для записи
//	outputFile, err := os.Create(csvPath)
//	if err != nil {
//		return fmt.Errorf("ошибка создания CSV файла: %v", err)
//	}
//	defer outputFile.Close()
//
//	writer := csv.NewWriter(outputFile)
//	defer writer.Flush()
//
//	// Обрабатываем строки Excel
//	for i, row := range rows {
//		if i < 3 { // Пропускаем первые 3 строки
//			continue
//		}
//
//		if len(row) < 8 { // Проверяем длину строки
//			continue
//		}
//
//		bdu := row[0]
//		otherSystem := row[len(row)-7]
//		cve := cveRegex.FindString(otherSystem)
//		if cve == "" {
//			cve = "N/A"
//		}
//
//		// Записываем в CSV
//		err := writer.Write([]string{bdu, cve})
//		if err != nil {
//			return fmt.Errorf("ошибка записи в CSV: %v", err)
//		}
//	}
//
//	return nil
//}
//
//// Нормализация XML
//func normalizeXML(xmlPath, outputXML, csvPath string) error {
//	// Загружаем соответствия BDU и CVE
//	cveMapping, err := loadCSV(csvPath)
//	if err != nil {
//		return err
//	}
//
//	// Читаем исходный XML
//	xmlFile, err := os.Open(xmlPath)
//	if err != nil {
//		return fmt.Errorf("ошибка открытия XML файла: %v", err)
//	}
//	defer xmlFile.Close()
//
//	xmlData, err := ioutil.ReadAll(xmlFile)
//	if err != nil {
//		return fmt.Errorf("ошибка чтения XML файла: %v", err)
//	}
//
//	var vulnerabilities Vulnerabilities
//	err = xml.Unmarshal(xmlData, &vulnerabilities)
//	if err != nil {
//		return fmt.Errorf("ошибка парсинга XML файла: %v", err)
//	}
//
//	// Обновляем данные
//	for i, vul := range vulnerabilities.Vul {
//		if cve, exists := cveMapping[vul.Identifier]; exists {
//			vulnerabilities.Vul[i].IdentifierCVE = cve
//		} else {
//			vulnerabilities.Vul[i].IdentifierCVE = "N/A"
//		}
//	}
//
//	// Сохраняем обновленный XML
//	output, err := xml.MarshalIndent(vulnerabilities, "", "  ")
//	if err != nil {
//		return fmt.Errorf("ошибка создания XML: %v", err)
//	}
//
//	err = ioutil.WriteFile(outputXML, []byte(xml.Header+string(output)), 0644)
//	if err != nil {
//		return fmt.Errorf("ошибка записи XML файла: %v", err)
//	}
//	return nil
//}
//
//// Чтение CSV и преобразование в карту
//func loadCSV(filePath string) (map[string]string, error) {
//	csvFile, err := os.Open(filePath)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка открытия CSV файла: %v", err)
//	}
//	defer csvFile.Close()
//
//	reader := csv.NewReader(csvFile)
//	csvData, err := reader.ReadAll()
//	if err != nil {
//		return nil, fmt.Errorf("ошибка чтения CSV файла: %v", err)
//	}
//
//	mapping := make(map[string]string)
//	for _, row := range csvData {
//		if len(row) >= 2 {
//			mapping[row[0]] = row[1]
//		}
//	}
//	return mapping, nil
//}
//
//// // -------------------------КОНЕЦ НЕ ТРОГАТЬ-------------------------------------------------
//
//func GenerateReport() {
//
//}
