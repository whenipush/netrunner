package controllers

import (
	"fmt"
	"log"
	"netrunner/database"
	"netrunner/models"
	"os/exec"
	"runtime"
)

// TODO: Доделать SQLMAP

type SQLmapRequest struct {
	Target  string            `json:"target" binding:"required"` // URL для сканирования
	Options map[string]string `json:"options"`                   // Дополнительные параметры для SQLmap
}

func RunSQL(task models.TaskStatus, params SQLMapParams) error {
	if params.TargetURL == "" {
		return fmt.Errorf("missing required parameter 'target_url'")
	}

	command := fmt.Sprintf("sqlmap -u %s --level=%d --risk=%d", params.TargetURL, params.Level, params.Risk)
	if params.Cookies != "" {
		command += fmt.Sprintf(" --cookie='%s'", params.Cookies)
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", command)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Ошибка выполнения SQLMap: %s", output)
		return fmt.Errorf("sqlmap execution failed: %v", err)
	}

	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  StatusCompleted,
		"percent": 100,
	})
	return nil
}
