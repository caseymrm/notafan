package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/go-pmset"
	smc "github.com/caseymrm/go-smc"
	"github.com/caseymrm/menuet"
)

func formatTemperature(tempC float64, celsius bool) string {
	if celsius {
		return fmt.Sprintf("%.01f°C", tempC)
	}
	return fmt.Sprintf("%.01f°F", tempC*1.8+32)
}

func averageFanSpeed(speeds []int) (avg int, ok bool) {
	if len(speeds) == 0 {
		return 0, false
	}
	sum := 0
	for _, s := range speeds {
		sum += s
	}
	return sum / len(speeds), true
}

func formatTitle(tempC float64, speeds []int, cpuLimit int, celsius bool) string {
	title := formatTemperature(tempC, celsius)
	if cpuLimit != 100 {
		title += fmt.Sprintf(" %d%%", cpuLimit)
	}
	if avg, ok := averageFanSpeed(speeds); ok {
		title += fmt.Sprintf(" %d", avg)
	}
	return title
}

func formatThrottleStatus(cpuLimit int) string {
	if cpuLimit == 100 {
		return "Not throttled"
	}
	return fmt.Sprintf("Throttled to %d%%", cpuLimit)
}

func setMenu() {
	celsius := menuet.Defaults().Boolean("celsius")
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: formatTitle(lastTemp, lastSpeeds, lastCPULimit, celsius),
	})
	menuet.App().MenuChanged()
}

func menuItems() []menuet.MenuItem {
	celsius := menuet.Defaults().Boolean("celsius")
	items := []menuet.MenuItem{
		{Text: "CPU", FontSize: 9},
		{Text: formatTemperature(lastTemp, celsius)},
		{Text: formatThrottleStatus(lastCPULimit)},
		{Type: menuet.Separator},
		{Text: "Fan speeds", FontSize: 9},
	}
	for _, speed := range lastSpeeds {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("%d RPM", speed),
		})
	}
	if len(lastSpeeds) == 0 {
		items = append(items, menuet.MenuItem{Text: "No fans!"})
	}
	items = append(items, menuet.MenuItem{Type: menuet.Separator})
	items = append(items, menuet.MenuItem{
		Text: "Units",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				{
					Text: "Fahrenheit",
					Clicked: func() {
						menuet.Defaults().SetBoolean("celsius", false)
						setMenu()
					},
					State: !celsius,
				},
				{
					Text: "Celsius",
					Clicked: func() {
						menuet.Defaults().SetBoolean("celsius", true)
						setMenu()
					},
					State: celsius,
				},
			}
		},
	})
	return items
}

var lastTemp float64
var lastSpeeds []int

func readTempAndFanSpeeds() (float64, []int) {
	smc.OpenSMC()
	temp := smc.ReadTemperature()
	speeds := smc.ReadFanSpeeds()
	smc.CloseSMC()
	return temp, speeds
}

func watchCPU() {
	lastTemp, lastSpeeds = readTempAndFanSpeeds()
	ticker := time.NewTicker(3 * time.Second)
	for ; true; <-ticker.C {
		lastTemp, lastSpeeds = readTempAndFanSpeeds()
		setMenu()
	}
}

var lastCPULimit int

func cpuSpeedLimit() int {
	thermal := pmset.GetThermalConditions()
	return thermal["CPU_Speed_Limit"]
}

func monitorThermalChanges(channel chan bool) {
	lastCPULimit = cpuSpeedLimit()
	lastNotification := time.Now()
	for range channel {
		newLimit := cpuSpeedLimit()
		if newLimit == lastCPULimit {
			continue
		}
		if lastNotification.Add(time.Second).After(time.Now()) {
			continue
		}
		if newLimit == 100 {
			menuet.App().Notification(menuet.Notification{
				Title: "CPU no longer throttled",
			})
		} else {
			menuet.App().Notification(menuet.Notification{
				Title: fmt.Sprintf("CPU being throttled to %d%%", newLimit),
			})
		}
		lastCPULimit = newLimit
		lastNotification = time.Now()
		setMenu()
	}
}

func main() {
	thermalChannel := make(chan bool)
	pmset.SubscribeThermalChanges(thermalChannel)
	go monitorThermalChanges(thermalChannel)
	go watchCPU()
	app := menuet.App()
	app.Name = "Not a Fan"
	app.Label = "com.github.caseymrm.notafan"
	app.Children = menuItems
	app.AutoUpdate.Version = "v1.0.0"
	app.AutoUpdate.Repo = "caseymrm/notafan"
	app.RunApplication()
}
