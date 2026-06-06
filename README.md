> ## ⚠️ Archived — Intel Macs only
>
> **This app does not work on Apple Silicon (M1/M2/M3/…).** It reads CPU temperature
> and fan speed via [go-smc](https://github.com/caseymrm/go-smc), which uses Intel
> SMC keys that don't exist on ARM Macs. The bundled `NotAFan.app` is an x86_64-only
> binary from 2018.
>
> **For a modern Mac, use [Stats](https://github.com/exelban/stats)** — an actively
> maintained, open-source menu bar monitor that supports Apple Silicon out of the box.
> [iStat Menus](https://bjango.com/mac/istatmenus/) is a great paid alternative.
>
> This repo is kept for historical reference.

---

# Not a Fan
Not a Fan is an OSX menu bar app that monitors your CPU temperature and fan state, and notifies you when your CPU is being throttled because of excessive heat

## Install

* [Download the app](https://github.com/caseymrm/notafan/releases/download/v0.1/NotAFan.app.zip)
* Unzip it
* Put it in Applications
* Run it
* To run every time you login, select "start at login" from the menu bar

## Screenshots

![CPU and Fans](https://github.com/caseymrm/notafan/raw/master/notafan.png)

![Throttled](https://github.com/caseymrm/notafan/raw/master/throttled.png)

![Not Throttled](https://github.com/caseymrm/notafan/raw/master/notthrottled.png)


## Built with

* [Menuet](https://github.com/caseymrm/menuet) - build menu bar apps in Go
* [go-smc](https://github.com/caseymrm/go-smc) - get your Mac's temperature and fan speed
* [go-pmset](https://github.com/caseymrm/go-pmset) - subscribe to changes in your Mac's temperature and throttling state

## License

MIT
