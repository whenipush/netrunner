package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netrunner/database"
	"netrunner/models"
)

func GetTaskStatus(c *gin.Context) {
	// Получаем значение NumberTask из параметров запроса
	numberTask := c.Param("number_task")

	// Инициализируем объект для задачи
	var task models.TaskStatus

	// Ищем задачу в базе данных
	if err := database.DB.Preload("Hosts").Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	// Если задача найдена, возвращаем её в JSON
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	// Получаем параметр number_task из URL
	numberTask := c.Param("number_task")

	// Проверяем, передан ли параметр
	if numberTask == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number_task is required"})
		return
	}

	// Проверяем, существует ли задача
	var task models.TaskStatus
	if err := database.DB.Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Удаляем задачу
	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "number_task": numberTask})
}

func GetTaskAll(c *gin.Context) {

	var task []models.TaskStatus

	// Ищем задачу в базе данных
	if err := database.DB.Preload("Hosts").Find(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tasks not found"})
		return
	}

	// Если задача найдена, возвращаем её в JSON
	c.JSON(http.StatusOK, task)
}

// / Константы для статусов задач
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusError     = "error"
)

// Типы параметров для задач
type NmapParams struct {
	Ports  string `json:"ports"`
	Script string `json:"script"`
}

type SQLMapParams struct {
	TargetURL string `json:"target_url"`
	Cookies   string `json:"cookies,omitempty"`
	Level     int    `json:"level"`
	Risk      int    `json:"risk"`
}

type DDosParams struct {
	Ports       string `json:"port"`
	PacketType  string `json:"packet_type"`
	Speed       string `json:"speed"`
	PacketCount string `json:"packet_count"`
}

type PentestParams struct {
	Ports string `json:"ports"`
	Speed string `json:"speed"`
}

// decodeParams декодирует JSON-параметры задачи
func decodeParams(task models.TaskStatus, dest interface{}) error {
	params, err := json.Marshal(task.Params)
	if err != nil {
		return fmt.Errorf("Ошибка декодирования Params: %v", err)
	}
	if err := json.Unmarshal([]byte(params), dest); err != nil {
		return fmt.Errorf("Ошибка декодирования Params: %v", err)
	}
	return nil
}

// ProcessNmapRequest обрабатывает запрос на создание задачи
func CreateTask(c *gin.Context) {
	var input struct {
		Name   string      `json:"name" binding:"required"`
		Hosts  []string    `json:"hosts" binding:"required"`
		Type   string      `json:"type" binding:"required"`
		Params models.Json `json:"params" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)

	// Проверяем существование хостов в базе данных
	var existingHosts []models.Host
	if err := database.DB.Where("ip IN ?", input.Hosts).Find(&existingHosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hosts", "details": err.Error()})
		return
	}

	// Собираем список существующих IP
	existingIPs := make(map[string]bool)
	for _, host := range existingHosts {
		existingIPs[host.IP] = true
	}

	// Определяем отсутствующие IP
	var missingHosts []string
	for _, ip := range input.Hosts {
		if !existingIPs[ip] {
			missingHosts = append(missingHosts, ip)
		}
	}

	// Если есть отсутствующие IP, возвращаем ошибку
	if len(missingHosts) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Some hosts do not exist",
			"missing_hosts": missingHosts,
		})
		return
	}

	// Создаём задачу с привязанными хостами
	task := models.TaskStatus{
		Name:   input.Name,
		Type:   input.Type,
		Status: StatusPending,
		Params: input.Params,
		Hosts:  existingHosts, // Привязываем только существующие хосты
	}

	// Сохраняем задачу в базе данных
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	BroadcastTask(task) // Уведомляем WebSocket клиентов

	// Асинхронная обработка задачи
	go func() {
		if err := executeTask(task); err != nil {
			log.Printf("Ошибка выполнения задачи: %v", err)
			database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		}
		//BroadcastTask(task) // Уведомляем о завершении
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully", "task_id": task.ID})
}

// executeTask выполняет задачу в зависимости от её типа
func executeTask(task models.TaskStatus) error {
	switch task.Type {
	case "nmap":
		var params NmapParams
		if err := decodeParams(task, &params); err != nil {
			database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
			return err
		}
		return executeNmap(task, params)
	case "sqlmap":
		var params SQLMapParams
		if err := decodeParams(task, &params); err != nil {
			database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
			return err
		}
		return RunSQL(task, params)
	case "ddos":
		var params DDosParams
		if err := decodeParams(task, &params); err != nil {
			database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
			return err
		}
		return ExecuteDDos(task, params)
	case "pentest":
		var params PentestParams
		if err := decodeParams(task, &params); err != nil {
			database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
			return err
		}
		log.Printf("Выполняется задача Pentest с параметрами: %+v", params)
		return ExecutePentest(task, params)

	default:
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("unsupported task type: %s", task.Type)
	}
}
