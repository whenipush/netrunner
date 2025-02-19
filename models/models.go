package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Json map[string]interface{}

func (j *Json) Scan(value interface{}) error {
	if value == nil {
		*j = make(Json)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Cant convert to []byte")
	}
	return json.Unmarshal(bytes, j)
}

func (j Json) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Group struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Hosts       []Host `gorm:"many2many:group_hosts" json:"hosts"`
}

type Host struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);unique" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	IP          string       `gorm:"type:varchar(255);unique" json:"ip"`
	Groups      []Group      `gorm:"many2many:group_hosts"`
	TaskList    []TaskStatus `gorm:"many2many:task_hosts" json:"task_list"`
}

type TaskStatus struct {
	gorm.Model
	NumberTask string  `gorm:"type:varchar(255);unique_index" json:"number_task"`
	Name       string  `gorm:"type:varchar(255);unique_index" json:"name"`
	Type       string  `gorm:"type:varchar(50)" json:"type"` // Тип задачи
	Status     string  `gorm:"default:pending" json:"status"`
	Percent    float32 `gorm:"default:0" json:"percent"`
	Hosts      []Host  `gorm:"many2many:task_hosts" json:"hosts"` // Список хостов
	Params     Json    `gorm:"type:json" json:"params"`           // JSON в виде строки
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
