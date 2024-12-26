package controllers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"netrunner/database"
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

// TODO: ДОДЕЛАТЬ СОЗДАНИЕ ЗАДАЧИ
type NmapRequest struct {
	Hosts   []string `form:"hosts" binding:"required"`
	Script  string   `form:"script" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Status  string   `json:"status"`
	Percent float32  `json:"percent"`
}

func startNmap(host string, script string) {
	// Получение текущей даты
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	// Формирование имени файла отчёта
	report := fmt.Sprintf("report/report_%s.xml", currentDate)

	// Формирование команды NMap с --stats-every
	nmapCmd := exec.Command("nmap", "-sV", "--stats-every", "5s", "-oX", report, host)

	// Захват stdout и stderr
	stdout, err := nmapCmd.StdoutPipe()
	if err != nil {
		log.Println("Ошибка получения stdout:", err)
		return
	}

	stderr, err := nmapCmd.StderrPipe()
	if err != nil {
		log.Println("Ошибка получения stderr:", err)
		return
	}

	// Запуск команды
	if err := nmapCmd.Start(); err != nil {
		log.Println("Ошибка запуска NMap:", err)
		return
	}
	percentRegex := regexp.MustCompile(`About (\d+(\.\d+)?)%`)

	// Потоковая обработка вывода stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			//fmt.Println("[NMap STDOUT]:", scanner.Text())
			line := scanner.Text()
			if matches := percentRegex.FindStringSubmatch(line); matches != nil {
				percent := matches[1] // Достаем процент выполнения
				fmt.Printf("Прогресс: %s%% завершено\n", percent)
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
		return
	}

	// Сообщение об успешном завершении
	fmt.Println("100%")
	fmt.Println("Результаты NMap сохранены в файл:", report)
	fmt.Println("Скрипт:", script)
}

func ProcessNmapRequest(c *gin.Context) {
	var input NmapRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "Задача отправлена"})

	go func() {
		for _, host := range input.Hosts {
			startNmap(host, input.Script)

		}
		log.Println("Сделано!")
	}()

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
	numberTask := c.Param("number_task")
	var task models.TaskStatus

	if err := database.DB.Where("number_task = ?", numberTask).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
