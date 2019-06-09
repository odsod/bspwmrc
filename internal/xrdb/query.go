package xrdb

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type Resources struct {
	Bspwm Bspwm
	Dunst Dunst
}

type Bspwm struct {
	BorderWidth        int
	WindowGap          int
	NormalBorderColor  string
	ActiveBorderColor  string
	FocusedBorderColor string
}

type Dunst struct {
	Geometry string
}

func (rs *Resources) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if err := rs.UnmarshalLine(sc.Text()); err != nil {
			return xerrors.Errorf("xrdb query: %w", err)
		}
	}
	if sc.Err() != nil {
		return xerrors.Errorf("xrdb query: %w", sc.Err())
	}
	return nil
}

func (rs *Resources) UnmarshalLine(l string) error {
	parts := strings.SplitN(l, ":", 2)
	if len(parts) != 2 {
		return xerrors.Errorf("unmarshal: malformed line: %+v", l)
	}
	key := parts[0]
	value := strings.TrimSpace(parts[1])
	switch key {
	case "bspwm.borderWidth":
		i, err := strconv.Atoi(value)
		if err != nil {
			return xerrors.Errorf("unmarshal: %w", err)
		}
		rs.Bspwm.BorderWidth = i
	case "bspwm.windowGap":
		i, err := strconv.Atoi(value)
		if err != nil {
			return xerrors.Errorf("unmarshal: %w", err)
		}
		rs.Bspwm.WindowGap = i
	case "bspwm.normalBorderColor":
		rs.Bspwm.NormalBorderColor = value
	case "bspwm.activeBorderColor":
		rs.Bspwm.ActiveBorderColor = value
	case "bspwm.focusedBorderColor":
		rs.Bspwm.FocusedBorderColor = value
	case "dunst.geometry":
		rs.Dunst.Geometry = value
	}
	return nil
}

func Query() (*Resources, error) {
	cmd := exec.Command("xrdb", "-query", "-all")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, xerrors.Errorf("xrdb query: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, xerrors.Errorf("xrdb query: %w", err)
	}
	var resources Resources
	if err := resources.Read(stdout); err != nil {
		return nil, xerrors.Errorf("xrdb query: %w", err)
	}
	return &resources, nil
}
