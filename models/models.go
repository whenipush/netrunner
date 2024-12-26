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
	NumberTask string  `gorm:"type:varchar(255);unique_index" json:"number_task"` // Автоматическая генерация
	Name       string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Status     string  `gorm:"default:pending" json:"status"`
	Percent    float32 `gorm:"default:0" json:"percent"`
	Hosts      []Host  `gorm:"many2many:task_hosts" json:"hosts"`
	Script     string  `json:"script"`
}

func (task *TaskStatus) BeforeCreate(tx *gorm.DB) (err error) {
	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Percent == 0 {
		task.Percent = 0.0
	}
	return
}
