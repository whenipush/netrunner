package controllers

import (
	"github.com/gin-gonic/gin"
	"netrunner/handlers"
)

func PingHosts(c *gin.Context) {
	// Структура данных для запроса
	var request struct {
		IPRange string `json:"ip"` // Поле для диапазона IP
	}

	if request.IPRange == "" {
		c.JSON(200, gin.H{"error": "Not found IPs"})
	}

	// Прочитаем тело запроса и заполним структуру
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Сканируем сеть по указанному диапазону
	activeHosts := handlers.ScanNetwork(request.IPRange)

	// Отправляем список активных хостов в ответ
	c.JSON(200, gin.H{
		"activeHosts": activeHosts,
	})
}
