//go:build linux

package cpu

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO: TEST LINUX CPU READING

func GetCpuPercent(interval time.Duration) (float64, error) {
	idle1, total1, err := readLinuxCPUStat()
	if err != nil {
		return 0, err
	}

	time.Sleep(interval)

	idle2, total2, err := readLinuxCPUStat()
	if err != nil {
		return 0, err
	}

	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)

	if totalDelta == 0 {
		return 0, nil
	}

	return (1.0 - idleDelta/totalDelta) * 100, nil
}

func readLinuxCPUStat() (idle, total uint64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			for i, field := range fields[1:] {
				val, _ := strconv.ParseUint(field, 10, 64)
				total += val
				if i == 3 {
					idle = val
				}
			}
			break
		}
	}
	return idle, total, scanner.Err()
}

func getDarwinCPUPercent(interval time.Duration) (float64, error) {
	return 0, nil
}

func getWindowsCPUPercent(interval time.Duration) (float64, error) {
	return 0, nil
}
