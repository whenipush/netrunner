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
	var groups []models.Group
	if err := database.DB.Find(&groups).Error; err != nil {
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
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")
	if err := database.DB.Where("id =?", id).First(&group).Error; err != nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	if err := database.DB.Save(&group).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update group"})
		return
	}

	c.JSON(200, group)
}

func DeleteGroup(c *gin.Context) {
	id := c.Params.ByName("id")
	var group models.Group

	if err := database.DB.Where("id =?", id).Delete(&group).Error; err != nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
