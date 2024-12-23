package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name  string `json:"name"`
	Hosts []Host `gorm:"many2many:group_hosts" json:"hosts"`
}

type Host struct {
	gorm.Model
	IP     string  `gorm:"type:varchar(255);unique_index" json:"ip"`
	Groups []Group `gorm:"many2many:group_hosts"`
}
