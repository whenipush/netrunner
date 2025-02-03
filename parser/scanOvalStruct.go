package parser

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type OvalDefinitions struct {
	XMLName     xml.Name `xml:"oval_definitions"`
	Definitions struct {
		Definition []OvalDefinition `xml:"definition"`
	} `xml:"definitions"`
}

type OvalDefinition struct {
	Class    string       `xml:"class,attr"`
	ID       string       `xml:"id,attr"`
	Version  string       `xml:"version,attr"`
	Metadata OvalMetadata `xml:"metadata"`
	Criteria OvalCriteria `xml:"criteria"`
}

type OvalMetadata struct {
	Title    string `xml:"title"`
	Affected struct {
		Family   string   `xml:"family,attr"`
		Platform []string `xml:"platform"`
		Product  []string `xml:"product"`
	} `xml:"affected"`
	Reference []struct {
		Source string `xml:"source,attr"`
		RefURL string `xml:"ref_url,attr"`
		RefID  string `xml:"ref_id,attr"`
	} `xml:"reference"`
	Description string `xml:"description"`
	Bdu         struct {
		Severity    string `xml:"severity"`
		Cwe         string `xml:"cwe"`
		Remediation string `xml:"remediation"`
		Cvssv20     string `xml:"cvssv20"`
	} `xml:"bdu"`
	Ovaldb struct {
		Cpeb string `xml:"cpeb,attr"`
	} `xml:"ovaldb"`
}

type OvalCriteria struct {
	Operator         string `xml:"operator,attr,omitempty"`
	Comment          string `xml:"comment,attr,omitempty"`
	ExtendDefinition []struct {
		Text          string `xml:",chardata"`
		Comment       string `xml:"comment,attr"`
		DefinitionRef string `xml:"definition_ref,attr"`
	} `xml:"extend_definition,omitempty"`
	Criteria  []OvalCriteria `xml:"criteria,omitempty"`
	Criterion []struct {
		Comment string `xml:"comment,attr"`
		TestRef string `xml:"test_ref,attr"`
	} `xml:"criterion"`
}

var ScanOvalDatabase OvalDefinitions

func ParseScanOval(filepath string) error {

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Error openig ScanOval file: %s", err)
	}
	xmlParser := xml.NewDecoder(file)
	token, err := xmlParser.Token()
	if err != nil {
		return fmt.Errorf("Failed to parse xml token: %s", err.Error())
	}
	log.Printf("%s", token)

	return xmlParser.Decode(&ScanOvalDatabase)
}
