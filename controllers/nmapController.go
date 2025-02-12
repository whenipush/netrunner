package controllers

import (
	"bufio"
	"fmt"
	"log"
	"netrunner/database"
	"netrunner/models"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// TODO: ЭТО НАДО БУДЕТ СДЕЛАТЬ

// executeNmap выполняет задачу Nmap
func executeNmap(task models.TaskStatus, params NmapParams) error {
	if params.Ports == "" {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("missing required parameter 'ports'")
	}

	report := fmt.Sprintf("report/nmap/%s.xml", task.NumberTask)

	ipList := []string{}
	for _, host := range task.Hosts {
		ipList = append(ipList, host.IP)
	}
	ip := strings.Join(ipList, " ")
	if ip == "" {
		return fmt.Errorf("no valid hosts specified")
	}

	command := fmt.Sprintf("nmap -sV --stats-every 5s -p %s %s -oX %s", params.Ports, ip, report)
	if params.Script != "" {
		command = fmt.Sprintf("nmap -sV --stats-every 5s -p %s --script=%s %s -oX %s", params.Ports, params.Script, ip, report)
	}

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
		return fmt.Errorf("failed to start Nmap command: %v", err)
	}

	progressRegex := regexp.MustCompile(`About (\d+(\.\d+)?)% done`)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[Nmap]: %s", line)

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
			log.Printf("[Nmap STDERR]: %s", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Update("status", StatusError)
		return fmt.Errorf("nmap execution failed: %v", err)
	}

	database.DB.Model(&models.TaskStatus{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  StatusCompleted,
		"percent": 100,
	})
	task.Status = StatusCompleted
	task.Percent = 100
	BroadcastTask(task)

	return nil
}
