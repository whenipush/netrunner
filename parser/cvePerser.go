package parser

import (
	"encoding/json"
	"log"
	"os"
)

type VulnerabilityEntry struct {
	Source Vulnerability `json:"_source"`
}

type Vulnerability struct {
	Id          string           `json:"id"`
	Type        string           `json:"type"`
	CWE         []string         `json:"cwe"`
	Descrption  string           `json:"description"`
	Cpe         []string         `json:"cpe"`
	Cpe23       []string         `json:"cpe23"`
	Cvss        CVSS             `json:"cvss"`
	Cvss3       map[string]CVSS3 `json:"cvss3"`
	Link        string           `json:"href"`
	Exploits    []Details        `json:"exploits"`
	Solutions   []Details        `json:"solutions"`
	Workarounds []Details        `json:"workarounds"`
}

type Details struct {
	Language string `json:"lang"`
	Text     string `json:"value"`
}

type CVSS struct {
	Score    float32 `json:"score"`
	Severity string  `json:"severity"`
	Vector   string  `json:"vector"`
	Version  string  `json:"version"`
}

type CVSS3 struct {
	AttackComplexity      string  `json:"attackComplexity"`
	AttackVector          string  `json:"attackVector"`
	AvailabilityImpact    string  `json:"availabilityImpact"`
	BaseScore             float32 `json:"baseScore"`
	BaseSeverity          string  `json:"baseSeverity"`
	ConfidentialityImpact string  `json:"confidentialityImpact"`
	IntegrityImpact       string  `json:"integrityImpact"`
	PrivelegesRequired    string  `json:"privelegesRequired"`
	Scope                 string  `json:"scope"`
	UserInteraction       string  `json:"userInteraction"`
	VectorString          string  `json:"vectorString"`
	Version               string  `json:"version"`
}

func ParseCVE(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Failed to open file")
		return
	}
	jsonParse := json.NewDecoder(file)
	token, err := jsonParse.Token()
	if err != nil {
		log.Printf("Failed to parse json token: %s", err.Error())
		return
	}
	log.Printf("%s", token)
	for jsonParse.More() {
		var entry VulnerabilityEntry
		err := jsonParse.Decode(&entry)
		if err != nil {
			log.Printf("failed to parse vulnerability entry: %s", err.Error())
			return
		}

		source := entry.Source
		if len(source.Solutions) > 0 || len(source.Workarounds) > 0 || len(source.Exploits) > 0 {

			log.Printf("%v\n%v\n%v\n", source.Solutions, source.Workarounds, source.Exploits)

		}
	}
}
