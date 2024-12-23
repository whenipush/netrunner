package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Эндпоинт для запуска SQLmap
func RunSQLmap(c *gin.Context) {
	var req SQLmapRequest

	// Привязка JSON-запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка целевого URL
	if req.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target URL is required"})
		return
	}

	// Формируем имя файла для отчёта
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02_15-04-05")
	reportFile := fmt.Sprintf("reports/sqlmap_%s.xml", currentDate)

	// Создаём команду для запуска SQLmap
	args := []string{"-u", req.Target, "-o", reportFile, "--output-dir=./reports", "--batch"}
	for key, value := range req.Options {
		args = append(args, fmt.Sprintf("--%s=%s", key, value))
	}

	cmd := exec.Command("python3", "sqlmap.py")

	// Запуск команды
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running SQLmap: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to run SQLmap",
			"details": err.Error(),
			"output":  string(output),
		})
		return
	}

	// Возвращаем результат клиенту
	c.JSON(http.StatusOK, gin.H{
		"message":   "SQLmap completed successfully",
		"report":    reportFile,
		"sqlmapLog": string(output),
	})
}

// Эндпоинт для получения всех отчётов SQLmap
func GetSQLmapReports(c *gin.Context) {
	dir := "reports"

	// Сканируем папку
	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read reports directory"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No reports found"})
		return
	}

	// Составляем список файлов
	var reportFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".xml" {
			reportFiles = append(reportFiles, filepath.Join(dir, file.Name()))
		}
	}

	c.JSON(http.StatusOK, gin.H{"reports": reportFiles})
}
