package scratchpad

import (
	"os/exec"
	"os/user"
	"strconv"

	"github.com/odsod/bspwmrc/internal/bspc"
	"github.com/odsod/bspwmrc/internal/wm"
	"golang.org/x/xerrors"
)

func All() map[string]*S {
	return map[string]*S{
		"u": {
			Name:         "terminal",
			Cmd:          []string{"urxvtmux", "scratchpad"},
			ClassName:    "URxvt",
			InstanceName: "scratchpad",
		},
		"d": {
			Name:         "arandr",
			Cmd:          []string{"arandr"},
			ClassName:    "Arandr",
			InstanceName: "arandr",
		},
		"h": {
			Name: "browser",
			Cmd: []string{
				"google-chrome",
				"--user-data-dir=.local/share/browser-scratchpad",
			},
			ClassName:    "Google-chrome",
			InstanceName: "google-chrome (.local/share/browser-scratchpad)",
		},
		"t": {
			Name:         "keepassxc",
			Cmd:          []string{"keepassxc"},
			ClassName:    "keepassxc",
			InstanceName: "keepassxc",
		},
		"n": {
			Name:         "spotify",
			Cmd:          []string{"spotify"},
			ClassName:    "",
			InstanceName: "",
		},
		"s": {
			Name:         "slack",
			Cmd:          []string{"slack"},
			ClassName:    "Slack",
			InstanceName: "slack",
		},
		"g": {
			Name: "mail",
			Cmd: []string{
				"google-chrome",
				"--user-data-dir=.local/share/browser-scratchpad",
				"--app=https://mail.google.com",
			},
			ClassName:    "Google-chrome",
			InstanceName: "mail.google.com",
		},
		"c": {
			Name: "calendar",
			Cmd: []string{
				"google-chrome",
				"--user-data-dir=.local/share/browser-scratchpad",
				"--app=https://calendar.google.com",
			},
			ClassName:    "Google-chrome",
			InstanceName: "calendar.google.com",
		},
		"r": {
			Name: "drive",
			Cmd: []string{
				"google-chrome",
				"--user-data-dir=.local/share/browser-scratchpad",
				"--app=https://drive.google.com",
			},
			ClassName:    "Google-chrome",
			InstanceName: "drive.google.com",
		},
		"l": {
			Name: "meet",
			Cmd: []string{
				"google-chrome",
				"--user-data-dir=.local/share/browser-scratchpad",
				"--app=https://meet.google.com",
			},
			ClassName:    "Google-chrome",
			InstanceName: "meet.google.com",
		},
	}
}

type S struct {
	Name         string
	Cmd          []string
	ClassName    string
	InstanceName string
}

type SearchResult struct {
	Node    *wm.Node
	Desktop *wm.Desktop
	Monitor *wm.Monitor
}

func (s *SearchResult) Toggle(state *wm.State) error {
	if s.IsFocused(state) {
		if _, err := bspc.Run("node", strconv.Itoa(s.Node.ID), "--flag", "hidden=on"); err != nil {
			return xerrors.Errorf("toggle scratchpad: %w", err)
		}
		return nil
	}
	focusedDesktop, err := state.FocusedDesktop()
	if err != nil {
		return xerrors.Errorf("toggle scratchpad: %w", err)
	}
	if s.Desktop.ID != focusedDesktop.ID {
		if _, err := bspc.Run("node", strconv.Itoa(s.Node.ID), "--to-desktop", strconv.Itoa(focusedDesktop.ID)); err != nil {
			return xerrors.Errorf("toggle scratchpad: %w", err)
		}
	}
	if _, err := bspc.Run("node", strconv.Itoa(s.Node.ID), "--flag", "hidden=off", "--focus"); err != nil {
		return xerrors.Errorf("toggle scratchpad: %w", err)
	}
	return nil
}

func (s *SearchResult) IsFocused(state *wm.State) bool {
	if state.FocusedMonitorID != s.Monitor.ID {
		return false
	}
	if s.Monitor.FocusedDesktopID != s.Desktop.ID {
		return false
	}
	if s.Desktop.FocusedNodeID != s.Node.ID {
		return false
	}
	return true
}

func (s *S) Start() error {
	if len(s.Cmd) == 0 {
		return xerrors.New("empty command")
	}
	cmd := exec.Command(s.Cmd[0], s.Cmd[1:]...)
	currentUser, err := user.Current()
	if err != nil {
		return xerrors.Errorf("start scratchpad: %w", err)
	}
	cmd.Dir = currentUser.HomeDir
	if err := cmd.Start(); err != nil {
		return xerrors.Errorf("start scratchpad: %w", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return xerrors.Errorf("start scratchpad: %w", err)
	}
	return nil
}

func (s *S) SearchState(state *wm.State) (*SearchResult, bool) {
	for _, m := range state.Monitors {
		for _, d := range m.Desktops {
			if n, ok := s.SearchNode(d.Root); ok {
				return &SearchResult{
					Node:    n,
					Desktop: d,
					Monitor: m,
				}, true
			}
		}
	}
	return nil, false
}

func (s *S) SearchNode(root *wm.Node) (*wm.Node, bool) {
	if root == nil {
		return nil, false
	}
	if root.Client != nil {
		if root.Client.ClassName == s.ClassName && root.Client.InstanceName == s.InstanceName {
			return root, true
		}
	}
	if child, ok := s.SearchNode(root.FirstChild); ok {
		return child, true
	}
	if child, ok := s.SearchNode(root.SecondChild); ok {
		return child, true
	}
	return nil, false
}
