package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
)

type VulnerabilityEntry struct {
	Source Vulnerability `json:"_source"`
}

type VulnDatabase map[string]Vulnerability

var Database VulnDatabase = make(VulnDatabase)

//var VulnDatabase map[string]Vulnerability = make(map[string]Vulnerability)

type Vulnerability struct {
	Id               string           `json:"id"`
	Type             string           `json:"type"`
	CWE              []string         `json:"cwe"`
	Descrption       string           `json:"description"`
	Cpe              []CPE            `json:"cpe"`
	Cpe23            []CPE            `json:"cpe23"`
	CpeConfiguration CPEConfiguration `json:"cpeConfiguration"`
	Cvss             CVSS             `json:"cvss"`
	Cvss3            map[string]CVSS3 `json:"cvss3"`
	Link             string           `json:"href"`
	Exploits         []Details        `json:"exploits"`
	Solutions        []Details        `json:"solutions"`
	Workarounds      []Details        `json:"workarounds"`
}

type CPEMatch struct {
	Criteria   CPE    `json:"criteria"`
	VersionEnd string `json:"versionEndIncluding"`
	Vulnerable bool   `json:"vulnerable"`
}

type CPEConfiguration struct {
	CPEMatches []CPEMatch
}

func (c *CPEConfiguration) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	nodes := raw["nodes"]
	if nodes == nil {
		*c = CPEConfiguration{
			CPEMatches: make([]CPEMatch, 0),
		}
		return nil
	}
	var cpes []CPEMatch = make([]CPEMatch, 0)
	for _, n := range nodes.([]interface{}) {
		node := n.(map[string]interface{})
		if node == nil {
			*c = CPEConfiguration{
				CPEMatches: make([]CPEMatch, 0),
			}
			return nil
		}
		cpeMatches := node["cpeMatch"].([]interface{})
		for _, cc := range cpeMatches {
			cpeMatch := cc.(map[string]interface{})
			var cpe CPEMatch
			text, _ := json.Marshal(cpeMatch)
			json.Unmarshal(text, &cpe)
			cpes = append(cpes, cpe)
		}
	}
	*c = CPEConfiguration{
		CPEMatches: cpes,
	}
	return nil
}

type CPE struct {
	CPEVersion string
	Vendor     string
	Product    string
	Version    string
}

func NewCpe(str string) (*CPE, error) {
	tokens := strings.Split(str, ":")
	if len(tokens) == 1 {
		return nil, fmt.Errorf("string is not valid cpe string")
	}
	if tokens[0] != "cpe" {
		return nil, fmt.Errorf("string is not valid cpe string")
	}
	if tokens[1] == "2.3" {
		c := &CPE{
			CPEVersion: "2.3",
			Vendor:     tokens[3],
			Product:    tokens[4],
			Version:    tokens[5],
		}
		return c, nil
	} else if strings.Contains(tokens[1], "/") {
		c := &CPE{
			CPEVersion: "2.2",
			Vendor:     tokens[2],
			Product:    tokens[3],
		}
		if len(tokens) == 5 {
			c.Version = tokens[4]
		}
		return c, nil
	}
	return nil, fmt.Errorf("string is not valid cpe string")
}

func (c CPE) MatchProduct(c1 CPE) bool {
	return c.Product == c1.Product && c.Vendor == c1.Vendor
}

func (c *CPE) UnmarshalJSON(data []byte) error {
	var text string
	err := json.Unmarshal(data, &text)
	if err != nil {
		return err
	}
	cpe, err := NewCpe(text)
	if err != nil {
		return err
	} else {
		*c = *cpe
		return nil
	}
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
		Database[source.Id] = source
	}
	log.Printf("Parsed all CVE's")
}

func (v Vulnerability) MatchCPE(cpe CPE) bool {
	if cpe.Version == "*" || cpe.Version == "" {
		//log.Printf("[cpeParser] Tried to match cpe without version")
		return false
	}
	cpeVersion, _ := version.NewVersion(cpe.Version)
	found := false
	for _, c := range v.Cpe {
		match := c.MatchProduct(cpe)
		if match && (c.Version != "*" && c.Version != "") {
			cVersion, err := version.NewVersion(c.Version)
			if err != nil {
				log.Printf("Error parsing version: %v %v", err, c.Version)
				log.Printf("%v", c)
			}
			return cVersion.GreaterThanOrEqual(cpeVersion)
		}
		found = match || found
	}
	if !found {
		return false
	}
	for _, cpes := range v.CpeConfiguration.CPEMatches {
		if !cpes.Criteria.MatchProduct(cpe) {
			continue
		}
		if cpes.VersionEnd == "" {
			if cpes.Criteria.Version == "*" {
				log.Printf("%v", v)
				return true
			}
			cVersion, err := version.NewVersion(cpes.Criteria.Version)
			if err != nil {
				log.Printf("Error parsing version: %v %v", err, cpes.Criteria.Version)
				log.Printf("%v", cpes.Criteria)
				return false
			}
			return cVersion.GreaterThanOrEqual(cpeVersion)
		}
		cpesVersion, err := version.NewVersion(cpes.VersionEnd)
		if err != nil {
			log.Printf("Error parsing version: %v %v", err, cpes.VersionEnd)
			log.Printf("%v", cpes)
		}
		log.Printf("%v, %v: %v", cpesVersion, cpeVersion, cpesVersion.GreaterThanOrEqual(cpeVersion))
		return cpesVersion.GreaterThanOrEqual(cpeVersion)
	}
	return false
}

func (vdb VulnDatabase) FindCve(cpe CPE) []Vulnerability {
	var vulners []Vulnerability = make([]Vulnerability, 0)
	for _, v := range vdb {
		if v.MatchCPE(cpe) {
			vulners = append(vulners, v)
		}
	}
	return vulners
}
