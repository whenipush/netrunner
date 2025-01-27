package parser

import "encoding/xml"

type VulnSoftware struct {
	Soft []Soft `xml:"soft"`
}

type Soft struct {
	Text     string `xml:",chardata"`
	Vendor   string `xml:"vendor"`
	Name     string `xml:"name"`
	Version  string `xml:"version"`
	Platform string `xml:"platform"`
	Types    Types  `xml:"types"`
}
type Types struct {
	Type []string `xml:"type"`
}

type Environment struct {
	Os []Os `xml:"os"`
}

type Os struct {
	Vendor   string `xml:"vendor"`
	Name     string `xml:"name"`
	Version  string `xml:"version"`
	Platform string `xml:"platform"`
}

type CWE struct {
	Identifier string `xml:"identifier"`
}

type BDUCVSS struct {
	Vector CVSSVector `xml:"vector"`
}

type CVSSVector struct {
	Vector string `xml:",chardata"`
	Score  string `xml:"score,attr"`
}

type BDUCVSS3 struct {
	Vector CVSSVector `xml:"vector"`
}

type Identifiers struct {
	Identifier []Identifier `xml:"identifier"`
}

type Identifier struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
	Link string `xml:"link,attr"`
}

type Vul struct {
	XMLName            xml.Name     `xml:"vul"`
	Identifier         string       `xml:"identifier"`
	Name               string       `xml:"name"`
	Description        string       `xml:"description"`
	VulnerableSoftware VulnSoftware `xml:"vulnerable_software"`
	Environment        Environment  `xml:"environment"`
	Cwe                CWE          `xml:"cwe"`
	IdentifyDate       string       `xml:"identify_date"`
	Cvss               BDUCVSS      `xml:"cvss"`
	Cvss3              BDUCVSS3     `xml:"cvss3"`
	Severity           string       `xml:"severity"`
	Solution           string       `xml:"solution"`
	VulStatus          string       `xml:"vul_status"`
	ExploitStatus      string       `xml:"exploit_status"`
	FixStatus          string       `xml:"fix_status"`
	Sources            string       `xml:"sources"`
	Identifiers        Identifiers  `xml:"identifiers,omitempty"`
	Other              string       `xml:"other"`
	VulIncident        string       `xml:"vul_incident"`
	VulClass           string       `xml:"vul_class"`
}

type Vulnerabilities struct {
	XMLName xml.Name `xml:"vulnerabilities"`
	Vulns   []Vul    `xml:"vul"`
}
