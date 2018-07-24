package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/go-pmset"
	"github.com/caseymrm/go-smc"
	"github.com/caseymrm/menuet"
)

func setMenu() {
	smc.OpenSMC()
	temp := smc.ReadTemperature()
	speeds := smc.ReadFanSpeeds()
	smc.CloseSMC()
	celsius := menuet.Defaults().Boolean("celsius")
	text := fmt.Sprintf("%.01f°C", temp)
	if !celsius {
		text = fmt.Sprintf("%.01f°F", temp*1.8+32)
	}
	average := 0
	for _, speed := range speeds {
		average += speed
	}
	averageText := ""
	if len(speeds) > 0 {
		average = average / len(speeds)
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
	smc.OpenSMC()
	temp := smc.ReadTemperature()
	speeds := smc.ReadFanSpeeds()
	smc.CloseSMC()
	text := fmt.Sprintf("%.01f°C", temp)
	if !celsius {
		text = fmt.Sprintf("%.01f°F", temp*1.8+32)
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
	for _, speed := range speeds {
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("%d RPM", speed),
		})
	}
	if len(speeds) == 0 {
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

func watchCPU() {
	ticker := time.NewTicker(3 * time.Second)
	for ; true; <-ticker.C {
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

var lastCPULimit int

func cpuSpeedLimit() int {
	thermal := pmset.GetThermalConditions()
	return thermal["CPU_Speed_Limit"]
}

func monitorThermalChanges(channel chan bool) {
	lastCPULimit = cpuSpeedLimit()
	for range channel {
		newLimit := cpuSpeedLimit()
		if newLimit == lastCPULimit {
			continue
		}
		if newLimit == 100 {
			menuet.App().Notification(menuet.Notification{
				Title:   "CPU no longer throttled",
				Message: fmt.Sprintf("Previously limited to %d%%", lastCPULimit),
			})
		} else {
			menuet.App().Notification(menuet.Notification{
				Title:   "CPU being throttled",
				Message: fmt.Sprintf("Limit to %d%%", newLimit),
			})
		}
		lastCPULimit = newLimit
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
