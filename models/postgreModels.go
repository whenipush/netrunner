package models

import (
	"netrunner/parser"

	"gorm.io/gorm"
)

type Vulnerability struct {
	ID          uint          `gorm:"primarykey"`
	Name        string        `gorm:"type:varchar(255);unique_index"`
	Description []Description `gorm:"foreignKey:VulnId"`
	CPE         []parser.CPE  `gorm:"many2many:vuln_cpe"`
	CWE         []CWE         `gorm:"many2many:vuln_cwe"`
	CVSS        parser.CVSS   `gorm:"foreignKey:VulnId"`
	CVSS3       parser.CVSS3  `gorm:"foreignKey:VulnId"`
	Link        string        `gorm:"type:varchar(255)"`
	Solutions   []Solutions   `gorm:"foreignKey:VulnId"`
	Workarounds []Workarounds `gorm:"foreignKey:VulnId"`
	Exploits    []Exploits    `gorm:"foreignKey:VulnId"`
}
type Description parser.Details
type Solutions parser.Details
type Workarounds parser.Details
type Exploits parser.Details

type CWE struct {
	Id  uint   `gorm:"primarykey"`
	CWE string `gorm:"unique;not null;type:varchar(16)"`
}

func FindVuln(c parser.CPE) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Model(&Vulnerability{}).Where("id IN (?)", d.Table("vuln_cpe").
			Select("vulnerability_id").
			Where("cpe_id IN (?)", d.Model(&parser.CPE{}).
				Order("version desc").
				Select("id").
				Where("vendor = ? AND product = ? AND version = ?", c.Vendor, c.Product, c.Version)))
	}
}
