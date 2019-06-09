package childprocess

import (
	"log"
	"os/exec"
	"syscall"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/xerrors"
)

type Processes struct {
	Dunst  *process.Process
	Sxhkd  *process.Process
	Urxvtd *process.Process
	Xcape  *process.Process
}

func LoadAll() (*Processes, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, xerrors.Errorf("load processes: %w", err)
	}
	var result Processes
	for _, p := range ps {
		name, err := p.Name()
		if err != nil {
			return nil, xerrors.Errorf("load processes: %w", err)
		}
		switch {
		case name == "dunst":
			result.Dunst = p
		case name == "sxhkd":
			result.Sxhkd = p
		case name == "urxvtd":
			result.Urxvtd = p
		case name == "xcape":
			result.Xcape = p
		}
	}
	return &result, nil
}

func (ps *Processes) reloadSxhkd() error {
	if ps.Sxhkd != nil {
		log.Println("reloading sxhkd")
		if err := ps.Sxhkd.SendSignal(syscall.SIGUSR1); err != nil {
			return xerrors.Errorf("reload sxhkd: %w", err)
		}
	} else {
		log.Println("starting sxhkd")
		sxhkdCmd := exec.Command("sxhkd", "-t", "1")
		if err := sxhkdCmd.Start(); err != nil {
			return xerrors.Errorf("reload sxhkd: %w", err)
		}
		if err := sxhkdCmd.Process.Release(); err != nil {
			return xerrors.Errorf("reload sxhkd: %w", err)
		}
	}
	return nil
}

func (ps *Processes) reloadDunst() error {
	if ps.Dunst != nil {
		log.Println("killing dunst")
		if err := ps.Dunst.Kill(); err != nil {
			return xerrors.Errorf("reload processes: %w", err)
		}
	}
	log.Println("starting dunst")
	dunstCmd := exec.Command("dunst", "-geometry", "200x5-30+30")
	if err := dunstCmd.Start(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	if err := dunstCmd.Process.Release(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	return nil
}

func (ps *Processes) reloadXcape() error {
	if ps.Xcape != nil {
		log.Println("killing xcape")
		if err := ps.Xcape.Kill(); err != nil {
			return xerrors.Errorf("reload processes: %w", err)
		}
	}
	log.Println("starting xcape")
	xcapeCmd := exec.Command("xcape", "-e", "Control_L=Escape;Hyper_L=Tab", "-t", "250")
	if err := xcapeCmd.Start(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	if err := xcapeCmd.Process.Release(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	return nil
}

func (ps *Processes) reloadUrxvtd() error {
	if ps.Urxvtd == nil {
		log.Println("starting urxvtd")
		urxvtdCmd := exec.Command("urxvtd", "--quiet", "--fork", "--opendisplay")
		if err := urxvtdCmd.Start(); err != nil {
			return xerrors.Errorf("reload processes: %w", err)
		}
		if err := urxvtdCmd.Process.Release(); err != nil {
			return xerrors.Errorf("reload processes: %w", err)
		}
	}
	return nil
}

func (ps *Processes) Reload() error {
	if err := ps.reloadSxhkd(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	if err := ps.reloadDunst(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	if err := ps.reloadXcape(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	if err := ps.reloadUrxvtd(); err != nil {
		return xerrors.Errorf("reload processes: %w", err)
	}
	return nil
}
