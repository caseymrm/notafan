package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/go-pmset"
	smc "github.com/caseymrm/go-smc"
	"github.com/caseymrm/menuet"
)

func setMenu() {
	celsius := menuet.Defaults().Boolean("celsius")
	text := fmt.Sprintf("%.01f°C", lastTemp)
	if !celsius {
		text = fmt.Sprintf("%.01f°F", lastTemp*1.8+32)
	}
	average := 0
	for _, speed := range lastSpeeds {
		average += speed
	}
	averageText := ""
	if len(lastSpeeds) > 0 {
		average = average / len(lastSpeeds)
		averageText = fmt.Sprintf(" %d", average)
	}
	if lastCPULimit != 100 {
		text += fmt.Sprintf(" %d%%", lastCPULimit)
	}
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: text + averageText,
	})
	menuet.App().MenuChanged()
}

func menuItems(key string) []menuet.MenuItem {
	celsius := menuet.Defaults().Boolean("celsius")
	if key == "units" {
		return []menuet.MenuItem{
			{
				Text:  "Fahrenheit",
				Key:   "fahrenheit",
				State: !celsius,
			},
			{
				Text:  "Celsius",
				Key:   "celsius",
				State: celsius,
			},
		}
	}
	text := fmt.Sprintf("%.01f°C", lastTemp)
	if !celsius {
		text = fmt.Sprintf("%.01f°F", lastTemp*1.8+32)
	}
	throttleText := "Not throttled"
	if lastCPULimit != 100 {
		throttleText = fmt.Sprintf("Throttled to %d%%", lastCPULimit)

	}
	items := []menuet.MenuItem{
		{
			Text:     "CPU",
			FontSize: 9,
		},
		{
			Text: text,
		},
		{
			Text: throttleText,
		},
		{
			Type: menuet.Separator,
		},
		{
			Text:     "Fan speeds",
			FontSize: 9,
		},
	}
	for _, speed := range lastSpeeds {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("%d RPM", speed),
		})
	}
	if len(lastSpeeds) == 0 {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("No fans!"),
		})
	}
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})
	items = append(items, menuet.MenuItem{
		Key:      "units",
		Text:     "Units",
		Children: true,
	})
	if lastCPULimit != 100 {
		text += fmt.Sprintf(" %d%%", lastCPULimit)
	}
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

func handleClick(clicked string) {
	switch clicked {
	case "celsius":
		menuet.Defaults().SetBoolean("celsius", true)
		setMenu()
	case "fahrenheit":
		menuet.Defaults().SetBoolean("celsius", false)
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
	app.Clicked = handleClick
	app.MenuOpened = menuItems
	app.AutoUpdate.Version = "v0.1"
	app.AutoUpdate.Repo = "caseymrm/notafan"
	app.RunApplication()
}
