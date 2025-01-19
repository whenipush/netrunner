package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name  string `json:"name"`
	Hosts []Host `gorm:"many2many:group_hosts" json:"hosts"`
}

type Host struct {
	gorm.Model
	IP       string       `gorm:"type:varchar(255);unique_index" json:"ip"`
	Groups   []Group      `gorm:"many2many:group_hosts"`
	TaskList []TaskStatus `gorm:"many2many:task_hosts" json:"task_list"`
}

type TaskStatus struct {
	gorm.Model
	NumberTask string  `gorm:"type:varchar(255);unique_index" json:"number_task"`
	Name       string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Type       string  `gorm:"type:varchar(50)" json:"type"` // Тип задачи
	Status     string  `gorm:"default:pending" json:"status"`
	Percent    float32 `gorm:"default:0" json:"percent"`
	Hosts      []Host  `gorm:"many2many:task_hosts" json:"hosts"` // Список хостов
	Params     string  `gorm:"type:text" json:"params"`           // JSON в виде строки
}

func (task *TaskStatus) BeforeCreate(tx *gorm.DB) (err error) {
	// Установка значения по умолчанию для статуса
	if task.Status == "" {
		task.Status = "pending"
	}
	// Установка значения по умолчанию для процента
	if task.Percent == 0 {
		task.Percent = 0.0
	}
	// Генерация уникального номера задачи
	var maxID int64
	tx.Model(&TaskStatus{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	task.NumberTask = fmt.Sprintf("TASK-%05d", maxID+1) // Пример: TASK-00001
	return
}
