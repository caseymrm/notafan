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
		items := []menuet.MenuItem{
			{
				Text:     "CPU Temp",
				FontSize: 10,
			},
			{
				Text: fmt.Sprintf("%.01f°C", temp),
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
		if len(speeds) > 0 {
			average = average / len(speeds)
		}
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: fmt.Sprintf("%.01f° %d", temp, average),
			Items: items,
		})
		time.Sleep(time.Second)
	}
}

func main() {
	go watchCPU()
	app := menuet.App()
	app.Name = "Why Awake?"
	app.Label = "com.github.caseymrm.whyawake"
	app.RunApplication()
}
