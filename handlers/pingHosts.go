package handlers

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func pingHost(host string) bool {
	var cmd *exec.Cmd
	// Определяем ОС и параметры для ping
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", host) // для Windows, -w - тайм-аут
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host) // для Linux/Mac
	}

	err := cmd.Run()
	return err == nil
}

func ScanNetwork(ipRange string) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var activeHosts []string

	// Преобразуем ipRange в список адресов
	ipParts := strings.Split(ipRange, ".")
	baseIP := fmt.Sprintf("%s.%s.%s.", ipParts[0], ipParts[1], ipParts[2])

	// Сканируем все IP в пределах локальной сети
	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", baseIP, i)
		wg.Add(1)

		go func(ip string) {
			defer wg.Done()

			if pingHost(ip) {
				mu.Lock()
				activeHosts = append(activeHosts, ip)
				mu.Unlock()
				fmt.Printf("Host found: %s\n", ip)
			}
		}(ip)

		// Добавляем небольшую задержку, чтобы избежать перегрузки
		// Например, 10 горутин одновременно
		if i%10 == 0 {
			time.Sleep(200 * time.Millisecond)
		}
	}

	// Ждем завершения всех горутин
	wg.Wait()
	return activeHosts
}
