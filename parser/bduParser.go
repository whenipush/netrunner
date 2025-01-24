package parser

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

var BDUDatabase Vulnerabilities

func ParseBDU(filepath string) error {

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Error openig BDU file: %s", err)
	}
	xmlParser := xml.NewDecoder(file)
	token, err := xmlParser.Token()
	if err != nil {
		return fmt.Errorf("Failed to parse json token: %s", err.Error())
	}
	log.Printf("%s", token)

	return xmlParser.Decode(&BDUDatabase)
}
