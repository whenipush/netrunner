package controllers

import (
	"bufio"
	"fmt"
	"log"
	"netrunner/database"
	"netrunner/models"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Доделать DDoS

// executeDDos выполняет задачу DDoS
func ExecuteDDos(task models.TaskStatus, params DDosParams) error {
	// TODO: Добавить обработчики

	ipList := []string{}
	for _, host := range task.Hosts {
		ipList = append(ipList, host.IP)
	}
	ip := strings.Join(ipList, " ")
	if ip == "" {
		return fmt.Errorf("no valid hosts specified")
	}

	if params.Ports == "" {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("missing required parameter 'ports'")
	}

	report := fmt.Sprintf("report/ddos/%s.xml", task.NumberTask)

	command := fmt.Sprintf(
		"nping --%s -p %s --rate %s --count %s %s",
		params.PacketType, params.Ports, params.Speed, params.PacketCount, ip,
	)
	cmd := exec.Command("powershell", "-Command", command)

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
		return fmt.Errorf("failed to start DDoS command: %v", err)
	}

	progressRegex := regexp.MustCompile(`About (\d+(\.\d+)?)% done`)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[DDoS]: %s", line)

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
			log.Printf("[DDoS STDERR]: %s", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("DDoS execution failed: %v", err)
	}

	// TODO: ОТЧЕТ
	log.Printf(report)
	// Обработка отчета
	//_, err = handlers.ProcessPentest(report, report+".json")
	//if err != nil {
	//	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
	//	return fmt.Errorf("failed to process DDoS report: %v", err)
	//}

	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  StatusCompleted,
		"percent": 100,
	})
	task.Status = StatusCompleted
	task.Percent = 100
	BroadcastTask(task)

	return nil
}
