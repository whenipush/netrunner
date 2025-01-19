package controllers

import (
	"fmt"
	"log"
	"netrunner/database"
	"netrunner/models"
	"os/exec"
)

// executeDDos выполняет задачу DDoS
func ExecuteDDos(task models.TaskStatus, params DDosParams) error {
	if params.Target == "" {
		return fmt.Errorf("missing required parameter 'target'")
	}

	command := fmt.Sprintf("nping --%s -p %d --rate %d --count %d %s", params.PacketType, params.Port, params.Speed, params.PacketCount, params.Target)
	cmd := exec.Command("sh", "-c", command)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Ошибка выполнения DDoS: %s", output)
		return fmt.Errorf("ddos simulation failed: %v", err)
	}

	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  StatusCompleted,
		"percent": 100,
	})
	return nil
}
