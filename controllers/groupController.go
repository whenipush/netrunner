package controllers

import (
	"net/http"
	"netrunner/database"
	"netrunner/models"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&group).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(200, group)
}

func GetAllGroup(c *gin.Context) {
	var groups []models.Group = make([]models.Group, 0)
	if err := database.DB.Preload("Hosts").Find(&groups).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get groups"})
		return
	}

	c.JSON(200, groups)
}

func GetGroupByName(c *gin.Context) {
	name := c.Query("group")
	var group models.Group

	if err := database.DB.Preload("Hosts").Where("name = ?", name).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

func UpdateGroup(c *gin.Context) {
	var groupInput models.Group
	if err := c.ShouldBindJSON(&groupInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var group models.Group
	if err := database.DB.Where("id =?", id).First(&group).Error; err != nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}
	groupInput.ID = group.ID
	groupInput.CreatedAt = group.CreatedAt
	if err := database.DB.Save(&groupInput).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update group"})
		return
	}

	c.JSON(200, groupInput)
}

func DeleteGroup(c *gin.Context) {
	id := c.Params.ByName("id")
	var group models.Group

	if err := database.DB.Where("id =?", id).Delete(&group).Error; err != nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}
