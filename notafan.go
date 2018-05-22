package main

import (
	"fmt"
	"time"

	"github.com/caseymrm/go-smc"
	"github.com/caseymrm/menuet"
)

func watchCPU() {
	ticker := time.NewTicker(time.Second)
	for ; true; <-ticker.C {
		smc.OpenSMC()
		temp := smc.ReadTemperature()
		speeds := smc.ReadFanSpeeds()
		smc.CloseSMC()
		celsius := menuet.Defaults().Boolean("celsius")
		text := fmt.Sprintf("%.01f°C", temp)
		if !celsius {
			text = fmt.Sprintf("%.01f°F", temp*1.8+32)
		}

		items := []menuet.MenuItem{
			{
				Text:     "CPU Temp",
				FontSize: 10,
			},
			{
				Text: text,
			},
			{
				Type: menuet.Separator,
			},
			{
				Text:     "Fan speeds",
				FontSize: 10,
			},
		}
		average := 0
		for _, speed := range speeds {
			average += speed
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("%d RPM", speed),
			})
		}
		averageText := ""
		if len(speeds) > 0 {
			average = average / len(speeds)
			averageText = fmt.Sprintf(" %d", average)
		} else {
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("No fans!"),
			})
		}
		items = append(items, menuet.MenuItem{
			Type: menuet.Separator,
		})
		items = append(items, menuet.MenuItem{
			Text: "Preferences",
			Children: []menuet.MenuItem{
				{
					Text:     "Fahrenheit",
					Callback: "fahrenheit",
					State:    !celsius,
				},
				{
					Text:     "Celsius",
					Callback: "celsius",
					State:    celsius,
				},
				{
					Type: menuet.Separator,
				},
			},
		})
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: text + averageText,
			Items: items,
		})
		time.Sleep(time.Second)
	}
}

func handleClicks(callback chan string) {
	for clicked := range callback {
		switch clicked {
		case "celsius":
			menuet.Defaults().SetBoolean("celsius", true)
		case "fahrenheit":
			menuet.Defaults().SetBoolean("celsius", false)
		}
	}
}

func main() {
	go watchCPU()
	app := menuet.App()
	app.Name = "Not a Fan"
	app.Label = "com.github.caseymrm.notafan"
	clickChannel := make(chan string)
	go handleClicks(clickChannel)
	app.Clicked = clickChannel
	app.RunApplication()
}
