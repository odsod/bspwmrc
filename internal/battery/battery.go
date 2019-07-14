package battery

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"

	"golang.org/x/xerrors"
)

const powerSupplyDirname = "/sys/class/power_supply"

type B struct {
	Name       string
	Status     string
	EnergyNow  int
	EnergyFull int
}

func (b *B) Charge() float64 {
	if b.EnergyFull == 0 {
		return 0
	}
	return float64(b.EnergyNow) / float64(b.EnergyFull)
}

func listBatteries() ([]string, error) {
	filenames, err := ioutil.ReadDir(powerSupplyDirname)
	if err != nil {
		return nil, xerrors.Errorf("list batteries: %w", err)
	}
	result := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		powerSupply := filepath.Join(powerSupplyDirname, filename.Name())
		t, err := ioutil.ReadFile(filepath.Join(powerSupply, "type"))
		if err != nil {
			return nil, xerrors.Errorf("list batteries: %w", err)
		}
		if bytes.Equal(t, []byte("Battery\n")) {
			result = append(result, powerSupply)
		}
	}
	return result, nil
}

func readStr(f string) (string, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return "", xerrors.Errorf("read str: %w", err)
	}
	return string(bytes.TrimSpace(data)), nil
}

func readInt(f string) (int, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, xerrors.Errorf("read int: %w", err)
	}
	i, err := strconv.Atoi(string(bytes.TrimSpace(data)))
	if err != nil {
		return 0, xerrors.Errorf("read int: %w", err)
	}
	return i, nil
}

func Load(f string) (*B, error) {
	status, err := readStr(filepath.Join(f, "status"))
	if err != nil {
		return nil, xerrors.Errorf("load battery: %w", err)
	}
	energyFull, err := readInt(filepath.Join(f, "energy_full"))
	if err != nil {
		return nil, xerrors.Errorf("load battery: %w", err)
	}
	energyNow, err := readInt(filepath.Join(f, "energy_now"))
	if err != nil {
		return nil, xerrors.Errorf("load battery: %w", err)
	}
	return &B{
		Name:       path.Base(f),
		Status:     status,
		EnergyFull: energyFull,
		EnergyNow:  energyNow,
	}, nil
}

func LoadAll() ([]*B, error) {
	bs, err := listBatteries()
	if err != nil {
		return nil, xerrors.Errorf("load all batteries: %w", err)
	}
	result := make([]*B, 0, len(bs))
	for _, bb := range bs {
		b, err := Load(bb)
		if err != nil {
			return nil, xerrors.Errorf("load all batteries: %w", err)
		}
		result = append(result, b)
	}
	return result, nil
}
