# Not a Fan

Not a Fan is a macOS menu bar app that monitors your CPU temperature and fan state, and notifies you when your CPU is being throttled because of excessive heat.

Works on **Intel and Apple Silicon** (M1/M2/M3/M4…) Macs as of v1.0.0 — sensor reads go through [go-smc](https://github.com/caseymrm/go-smc), which now supports both architectures. Fanless Macs (MacBook Air, Mac mini base, Mac Studio) show "No fans!" — they have no fans to read.

## Install

* [Download the app](https://github.com/caseymrm/notafan/releases/latest)
* Unzip it
* Put it in Applications
* Run it
* To run every time you log in, select "Start at login" from the menu bar

## Screenshots

![CPU and Fans](https://github.com/caseymrm/notafan/raw/master/notafan.png)

![Throttled](https://github.com/caseymrm/notafan/raw/master/throttled.png)

![Not Throttled](https://github.com/caseymrm/notafan/raw/master/notthrottled.png)


## Built with

* [Menuet](https://github.com/caseymrm/menuet) — build menu bar apps in Go
* [go-smc](https://github.com/caseymrm/go-smc) — read your Mac's temperature and fan speed (Intel + Apple Silicon)
* [go-pmset](https://github.com/caseymrm/go-pmset) — subscribe to changes in your Mac's temperature and throttling state

## License

MIT
