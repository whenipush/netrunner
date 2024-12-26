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

type TaskStatus struct {
	gorm.Model
	Status  string `json:"status"`
	Host    string `json:"host,omitempty"`
	Message string `json:"message,omitempty"`
}
