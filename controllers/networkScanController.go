package controllers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"netrunner/database"
	"netrunner/handlers"
	"netrunner/models"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ExecuteNetworkScan(task models.TaskStatus, params NetworkScanParams) error {
	report := fmt.Sprintf("report/networkscan/%s.xml", task.NumberTask)
	command := fmt.Sprintf("nmap  -sV -O --stats-every 5s -T%s -oX %s %s/16", params.Speed, report, params.NetworkBaseAddress)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", command)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("failed to start Network Scan command: %v", err)
	}

	progressRegex := regexp.MustCompile(`About (\d+(\.\d+)?)% done`)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[NetworkScan]: %s", line)

			if matches := progressRegex.FindStringSubmatch(line); matches != nil {
				percent := matches[1]
				percentValue, _ := strconv.ParseFloat(percent, 32)
				task.Percent = float32(percentValue)
				database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("percent", percentValue)
				BroadcastTask(task)
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("[NetworkScan STDERR]: %s", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("network scan execution failed: %v", err)
	}

	if err := handlers.ProcessNetworkScan(report); err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("failed to process Network Scan report: %v", err)
	}
	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  StatusCompleted,
		"percent": 100.0,
	})
	task.Status = StatusCompleted
	task.Percent = 100.0
	BroadcastTask(task)

	return nil
}

func GetNetworkJsonByNumberTask(c *gin.Context) {
	var task models.TaskStatus
	numberTask := c.Param("number_task")

	if err := database.DB.Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	report := fmt.Sprintf("report/networkscan/%s.xml.json", task.NumberTask)
	file, err := os.ReadFile(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Отчет не найден"})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(file))

}
