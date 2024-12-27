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
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: ДОДЕЛАТЬ СОЗДАНИЕ ЗАДАЧИ // СДЕЛАНО
// TODO: Выбор скрипта
// TODO: Переписать получение отчета
type NmapRequest struct {
	Hosts   []string `form:"hosts" binding:"required"`
	Group   string   `form:"group" binding:"required"`
	Script  string   `form:"script" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Status  string   `json:"status"`
	Percent float32  `json:"percent"`
}

func startNmap(host string, script string, taskID uint) {
	log_name := fmt.Sprintf("TASK-%d", taskID)
	logger, err := handlers.LogNmapError(log_name)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	// Получение текущей даты
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	// Формирование имени файла отчёта
	report := fmt.Sprintf("report/report-task%d-%s.xml", taskID, currentDate)

	// Формирование команды NMap с --stats-every
	nmapCmd := exec.Command("nmap", "-sV", "--stats-every", "5s", "-oX", report, host)

	// Захват stdout и stderr
	stdout, err := nmapCmd.StdoutPipe()
	if err != nil {
		log.Println("Ошибка получения stdout:", err)
		logger.Println("Ошибка получения stdout:", err)
		return
	}

	stderr, err := nmapCmd.StderrPipe()
	if err != nil {
		log.Println("Ошибка получения stderr:", err)
		logger.Println("Ошибка получения stderr:", err)
		return
	}

	// Запуск команды
	if err := nmapCmd.Start(); err != nil {
		log.Println("Ошибка запуска NMap:", err)
		logger.Println("Ошибка запуска NMap:", err)
		return
	}

	// Регулярное выражение для извлечения процентов выполнения
	percentRegex := regexp.MustCompile(`About (\d+(\.\d+)?)%`)

	// Потоковая обработка вывода stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if matches := percentRegex.FindStringSubmatch(line); matches != nil {
				percent := matches[1] // Достаем процент выполнения
				fmt.Printf("Прогресс: %s%% завершено\n", percent)

				// Обновляем процент выполнения в базе данных
				if err := database.DB.Model(&models.TaskStatus{}).Where("id = ?", taskID).Update("percent", percent).Error; err != nil {
					log.Println("Ошибка обновления процента выполнения:", err)
					logger.Println("Ошибка обновления процента выполнения:", err)
				}
			}
		}
	}()

	// Потоковая обработка вывода stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("[NMap STDERR]:", scanner.Text())
		}
	}()

	// Ожидание завершения команды
	if err := nmapCmd.Wait(); err != nil {
		log.Println("Ошибка выполнения команды NMap:", err)
		logger.Println("Ошибка выполнения команды NMap:", err)
		database.DB.Model(&models.TaskStatus{}).Where("id = ?", taskID).Update("status", "error")
		return
	}

	// Обновляем статус и процент выполнения на 100% после завершения
	database.DB.Model(&models.TaskStatus{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"percent": 100,
		"status":  "completed",
	})

	fmt.Println("Результаты NMap сохранены в файл:", report)
	fmt.Println("Скрипт:", script)
}

func ProcessNmapRequest(c *gin.Context) {
	var input struct {
		Name   string   `json:"name" binding:"required"`
		Hosts  []string `json:"hosts" binding:"required"` // Список IP хостов
		Script string   `json:"script" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования всех хостов
	var existingHosts []models.Host
	database.DB.Where("ip IN ?", input.Hosts).Find(&existingHosts)

	// Собираем список отсутствующих хостов
	existingIPs := make(map[string]bool)
	for _, host := range existingHosts {
		existingIPs[host.IP] = true
	}

	var missingHosts []string
	for _, ip := range input.Hosts {
		if !existingIPs[ip] {
			missingHosts = append(missingHosts, ip)
		}
	}

	// Если есть отсутствующие хосты, возвращаем ошибку
	if len(missingHosts) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Some hosts do not exist",
			"missing_hosts": missingHosts,
		})
		return
	}

	// Создание задачи
	task := models.TaskStatus{
		Name:   input.Name,
		Status: "pending",
		Hosts:  existingHosts, // Привязываем только существующие хосты
		Script: input.Script,
	}

	// Сохранение задачи в базе данных
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go func() {
		for ip := range existingIPs {
			startNmap(ip, input.Script, task.ID)
		}
	}()
	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully", "task": task.ID})

}

func UploadScript(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, fmt.Sprintf("get form file err: %v", err))
		return
	}

	if file.Filename == "" {
		c.JSON(400, "File name is required")
		return
	}
	filesplit := strings.Split(file.Filename, ".")
	fileExt := filesplit[len(filesplit)-1]

	if fileExt != "lua" {
		c.JSON(400, gin.H{"Error": "Only lua files are allowed"})
		return
	}

	filepath := "./scripts/" + file.Filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, fmt.Sprintf("upload file err: %v", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}

func GetLastNmap(c *gin.Context) {
	// Путь к папке с файлами
	dir := "reports"

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	// Найти самый последний файл
	var latestFile os.DirEntry
	var latestModTime time.Time

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(latestModTime) {
			latestFile = file
			latestModTime = info.ModTime()
		}
	}

	if latestFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	filePath := filepath.Join(dir, latestFile.Name())
	c.File(filePath)
}

func GetAllNmap(c *gin.Context) {
	// Путь к папке с файлами
	dir := "reports"

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	// Сортируем файлы по дате модификации
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Создаем массив с путями к файлам
	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, filepath.Join(dir, file.Name()))
	}

	c.JSON(http.StatusOK, gin.H{"files": filePaths})
}

func GetTaskStatus(c *gin.Context) {
	// Получаем значение NumberTask из параметров запроса
	numberTask := c.Param("number_task")

	// Инициализируем объект для задачи
	var task models.TaskStatus

	// Ищем задачу в базе данных
	if err := database.DB.Preload("Hosts").Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Если задача найдена, возвращаем её в JSON
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	// Получаем параметр number_task из URL
	numberTask := c.Param("number_task")

	// Проверяем, передан ли параметр
	if numberTask == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "number_task is required"})
		return
	}

	// Проверяем, существует ли задача
	var task models.TaskStatus
	if err := database.DB.Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Удаляем задачу
	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "number_task": numberTask})
}

func GetTaskAll(c *gin.Context) {

	var task []models.TaskStatus

	// Ищем задачу в базе данных
	if err := database.DB.Preload("Hosts").Find(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tasks not found"})
		return
	}

	// Если задача найдена, возвращаем её в JSON
	c.JSON(http.StatusOK, task)
}
