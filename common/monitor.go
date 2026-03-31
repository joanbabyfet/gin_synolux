package common

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// 获取cpu使用率
func GetCpuPercent() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return percent[0], nil
}

// 获取内存使用率
func GetRamPercent() (float64, error) {
	ram_info, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return ram_info.UsedPercent, nil
}

// 获取cpu温度
func GetCpuTemp() (int, error) {
	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	tempStr := strings.Replace(out.String(), "\n", "", -1)
	temp, err := strconv.Atoi(tempStr)
	if err != nil {
		return 0, err
	}
	temp = temp / 1000
	return temp, nil
}