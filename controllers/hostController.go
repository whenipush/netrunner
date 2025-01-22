package controllers

import (
	"net"
	"net/http"
	"netrunner/database"
	"netrunner/models"

	"github.com/gin-gonic/gin"
)

// CreateHost - создает новый хост.

func CreateHost(c *gin.Context) {
	// Create new host
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if net.ParseIP(host.IP) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
		return
	}

	if err := database.DB.Where("ip =?", host.IP).First(&host).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Host with this IP already exists"})
		return
	}

	if err := database.DB.Create(&host).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, host)
}

// GetAllHost - получает все хосты.

func GetAllHost(c *gin.Context) {
	var hosts []models.Host
	if err := database.DB.Preload("Groups").Find(&hosts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

func GetHostByID(c *gin.Context) {
	ip := c.Query("ip")
	var host models.Host

	if err := database.DB.Preload("Groups").Where("ip = ?", ip).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}
	c.JSON(http.StatusOK, host)
}

// UpdateHost - изменяет хост по ID.

func UpdateHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")
	if err := database.DB.Where("id = ?", id).First(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	if err := database.DB.Save(&host).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, host)
}

// DeleteHost - удаляет хост по ID.

func DeleteHost(c *gin.Context) {
	id := c.Params.ByName("id")
	var host models.Host

	if err := database.DB.Where("id =?", id).Delete(&host).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id #" + id: "deleted"})
}

// AddHostToGroupHandler - добавляет хосты к группам.

func AddHostToGroup(c *gin.Context) {

	var input struct {
		HostIDs  []uint `json:"host_ids" binding:"required"`  // Массив ID хостов
		GroupIDs []uint `json:"group_ids" binding:"required"` // Массив ID групп
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()}) // Ошибка валидации
		return
	}

	var hosts []models.Host
	if err := database.DB.Find(&hosts, input.HostIDs).Error; err != nil || len(hosts) == 0 {
		c.JSON(404, gin.H{"error": "Some hosts not found"})
		return
	}

	var groups []models.Group
	if err := database.DB.Find(&groups, input.GroupIDs).Error; err != nil || len(groups) == 0 {
		c.JSON(404, gin.H{"error": "Some groups not found"})
		return
	}

	for _, host := range hosts {
		if err := database.DB.Model(&host).Association("Groups").Append(groups); err != nil {
			c.JSON(500, gin.H{"error": "Failed to add some hosts to groups"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Hosts added to groups successfully"})
}
