package main

import (
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/odsod/bspwmrc/internal/battery"
	"github.com/odsod/bspwmrc/internal/bspc"
	"github.com/odsod/bspwmrc/internal/childprocess"
	"github.com/odsod/bspwmrc/internal/notify"
	"github.com/odsod/bspwmrc/internal/scratchpad"
	"github.com/odsod/bspwmrc/internal/wm"
	"github.com/odsod/bspwmrc/internal/xrdb"
)

func main() {
	s, err := syslog.New(syslog.LOG_DEBUG, "bspwmrc")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			panic(err)
		}
	}()
	logger := log.New(s, "", log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("panic: %+v", r)
		}
	}()
	args := os.Args[1:]
	switch {
	case len(args) == 0:
		config(logger)
	case args[0] == "toggle-scratchpad" && len(args) == 2:
		toggleScratchpad(logger, args[1])
	case args[0] == "cron":
		cron(logger)
	case args[0] == "run":
		run(logger)
	case args[0] == "clock":
		clock(logger)
	case args[0] == "battery-charge":
		batteryCharge(logger)
	case args[0] == "prev":
		prev(logger)
	default:
		logger.Printf("unhandled: %+v", args)
	}
}

func config(logger *log.Logger) {
	logger.Printf("config")
	// Load xresources
	xresources, err := xrdb.Query()
	if err != nil {
		panic(err)
	}
	// Configure bspwm
	for _, cmd := range [][]string{
		{"config", "focus_follows_pointer", "true"},
		{"config", "pointer_follows_focus", "true"},
		{"config", "pointer_follows_monitor", "true"},
		{"config", "borderless_monocle", "true"},
		{"config", "paddingless_monocle", "true"},
		{"config", "gapless_monocle", "true"},
		{"config", "single_monocle", "true"},
		{"config", "pointer_modifier", "mod3"},
		{"config", "remove_unplugged_monitors", "true"},
		{"config", "window_gap", strconv.Itoa(xresources.Bspwm.WindowGap)},
		{"config", "border_width", strconv.Itoa(xresources.Bspwm.BorderWidth)},
		{"config", "normal_border_color", xresources.Bspwm.NormalBorderColor},
		{"config", "active_border_color", xresources.Bspwm.ActiveBorderColor},
		{"config", "focused_border_color", xresources.Bspwm.FocusedBorderColor},
	} {
		logger.Printf("bspc %v", cmd)
		if _, err := bspc.Run(cmd...); err != nil {
			panic(err)
		}
	}
	// Load processes
	cps, err := childprocess.LoadAll()
	if err != nil {
		panic(err)
	}
	if err := cps.Reload(); err != nil {
		panic(err)
	}
	for _, cmd := range [][]string{
		{"setxkbmap", "custom"},
		{"xsetroot", "-cursor_name", "left_ptr"},
		{"feh", "--bg-scale", "/usr/share/backgrounds/ubuntu-default-greyscale-wallpaper.png"},
	} {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			panic(err)
		}
	}
	if err := notify.Send("bspwmrc", "desktop reloaded", time.Second); err != nil {
		panic(err)
	}
}

func toggleScratchpad(logger *log.Logger, name string) {
	logger.Printf("toggle-scratchpad name=%s", name)
	sp, ok := scratchpad.All()[name]
	if !ok {
		logger.Printf("no such scratchpad: %v", name)
		return
	}
	state, err := wm.LoadState()
	if err != nil {
		panic(err)
	}
	searchResult, ok := sp.SearchState(state)
	if !ok {
		if err := sp.Start(); err != nil {
			panic(err)
		}
		return
	}
	if err := searchResult.Toggle(state); err != nil {
		panic(err)
	}
}

func cron(logger *log.Logger) {
	logger.Printf("cron")
	bs, err := battery.LoadAll()
	if err != nil {
		panic(err)
	}
	for _, b := range bs {
		if b.Charge() < 0.1 {
			if err := notify.Send(
				fmt.Sprintf("%s charge", b.Name),
				fmt.Sprintf("%.2f%%", b.Charge()*100),
				2*time.Second,
			); err != nil {
				panic(err)
			}
		}
	}
}

func run(logger *log.Logger) {
	logger.Printf("run")
	cmd := exec.Command("rofi", "-show", "run", "-display-run", "", "-theme-str", "#window { border: 5; }")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func prev(logger *log.Logger) {
	logger.Printf("prev")
	if _, err := bspc.Run("node", "--focus", "prev"); err != nil {
		panic(err)
	}
}

func clock(logger *log.Logger) {
	logger.Printf("clock")
	now := time.Now()
	nowDate := now.Format("Mon Jan _2")
	nowTime := now.Format("15:04:05 MST")
	if err := notify.Send(nowDate, nowTime, 2*time.Second); err != nil {
		panic(err)
	}
}

func batteryCharge(logger *log.Logger) {
	logger.Printf("battery-charge")
	batteries, err := battery.LoadAll()
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	for _, b := range batteries {
		if _, err := fmt.Fprintf(&buf, "%s %s:\n%.2f%%\n", b.Name, b.Status, b.Charge()*100); err != nil {
			panic(err)
		}
	}
	if err := notify.Send("Battery", buf.String(), 2*time.Second); err != nil {
		panic(err)
	}
}
