package main

import "testing"

func TestFormatTemperature(t *testing.T) {
	cases := []struct {
		name    string
		tempC   float64
		celsius bool
		want    string
	}{
		{"celsius round", 50.0, true, "50.0°C"},
		{"celsius fractional", 42.37, true, "42.4°C"},
		{"fahrenheit round", 100.0, false, "212.0°F"},
		{"fahrenheit zero c", 0.0, false, "32.0°F"},
		{"fahrenheit body temp", 37.0, false, "98.6°F"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := formatTemperature(c.tempC, c.celsius)
			if got != c.want {
				t.Errorf("formatTemperature(%v, %v) = %q, want %q", c.tempC, c.celsius, got, c.want)
			}
		})
	}
}

func TestAverageFanSpeed(t *testing.T) {
	cases := []struct {
		name    string
		speeds  []int
		wantAvg int
		wantOK  bool
	}{
		{"empty (fanless Mac)", []int{}, 0, false},
		{"nil", nil, 0, false},
		{"single fan", []int{2000}, 2000, true},
		{"two fans", []int{2000, 3000}, 2500, true},
		{"three fans integer truncation", []int{1000, 1000, 1001}, 1000, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			avg, ok := averageFanSpeed(c.speeds)
			if avg != c.wantAvg || ok != c.wantOK {
				t.Errorf("averageFanSpeed(%v) = (%d, %v), want (%d, %v)", c.speeds, avg, ok, c.wantAvg, c.wantOK)
			}
		})
	}
}

func TestFormatTitle(t *testing.T) {
	cases := []struct {
		name     string
		tempC    float64
		speeds   []int
		cpuLimit int
		celsius  bool
		want     string
	}{
		{"normal C with fans", 50.0, []int{2000, 4000}, 100, true, "50.0°C 3000"},
		{"normal F with fans", 50.0, []int{2000, 4000}, 100, false, "122.0°F 3000"},
		{"fanless Mac, no throttle", 45.0, nil, 100, true, "45.0°C"},
		{"fanless Mac, throttled", 90.0, nil, 80, true, "90.0°C 80%"},
		{"throttled with fans (order: temp, throttle, fans)", 90.0, []int{5000}, 50, true, "90.0°C 50% 5000"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := formatTitle(c.tempC, c.speeds, c.cpuLimit, c.celsius)
			if got != c.want {
				t.Errorf("formatTitle(%v, %v, %d, %v) = %q, want %q",
					c.tempC, c.speeds, c.cpuLimit, c.celsius, got, c.want)
			}
		})
	}
}

func TestFormatThrottleStatus(t *testing.T) {
	cases := []struct {
		name     string
		cpuLimit int
		want     string
	}{
		{"not throttled", 100, "Not throttled"},
		{"half", 50, "Throttled to 50%"},
		{"zero (worst case)", 0, "Throttled to 0%"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := formatThrottleStatus(c.cpuLimit)
			if got != c.want {
				t.Errorf("formatThrottleStatus(%d) = %q, want %q", c.cpuLimit, got, c.want)
			}
		})
	}
}
