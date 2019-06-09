package bspc

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

func resolveSocketPath() (string, error) {
	if socketPathFromEnv, ok := os.LookupEnv("BSPWM_SOCKET"); ok {
		return socketPathFromEnv, nil
	}
	displayStr := os.Getenv("DISPLAY")
	hostAndRest := strings.SplitN(displayStr, ":", 2)
	if len(hostAndRest) < 2 {
		return "", xerrors.Errorf("malformed DISPLAY: %s", displayStr)
	}
	host := hostAndRest[0]
	displayAndScreen := strings.SplitN(hostAndRest[1], ".", 2)
	display := displayAndScreen[0]
	screen := "0"
	if len(displayAndScreen) == 2 {
		screen = displayAndScreen[1]
	}
	return fmt.Sprintf("/tmp/bspwm%s_%s_%s-socket", host, display, screen), nil
}

func Run(args ...string) (response []byte, err error) {
	socketPath, err := resolveSocketPath()
	if err != nil {
		return nil, xerrors.Errorf("bspc run: %w", err)
	}
	socket, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, xerrors.Errorf("bspc run: %w", err)
	}
	defer func() {
		if errClose := socket.Close(); err != nil {
			response, err = nil, xerrors.Errorf("bspc run: %w", errClose)
		}
	}()
	_, err = socket.Write([]byte(strings.Join(args, "\x00") + "\x00"))
	if err != nil {
		return nil, xerrors.Errorf("bspc run: %w", err)
	}
	data, err := ioutil.ReadAll(socket)
	if err != nil {
		return nil, xerrors.Errorf("bspc run: %w", err)
	}
	return data, nil
}
