package wm

import (
	"encoding/json"

	"github.com/odsod/bspwmrc/internal/bspc"
	"golang.org/x/xerrors"
)

type State struct {
	FocusedMonitorID int        `json:"focusedMonitorId"`
	PrimaryMonitorID int        `json:"primaryMonitorId"`
	ClientsCount     int        `json:"clientsCount"`
	Monitors         []*Monitor `json:"monitors"`
}

type Monitor struct {
	Name             string     `json:"name"`
	ID               int        `json:"id"`
	RandrID          int        `json:"randrId"`
	Wired            bool       `json:"wired"`
	StickyCount      int        `json:"stickyCount"`
	WindowGap        int        `json:"windowGap"`
	BorderWidth      int        `json:"borderWidth"`
	FocusedDesktopID int        `json:"focusedDesktopID"`
	Padding          Padding    `json:"padding"`
	Rectangle        Rectangle  `json:"rectangle"`
	Desktops         []*Desktop `json:"desktops"`
}

type Desktop struct {
	Name          string  `json:"name"`
	ID            int     `json:"id"`
	Layout        string  `json:"layout"`
	WindowGap     int     `json:"windowGap"`
	BorderWidth   int     `json:"borderWidth"`
	FocusedNodeID int     `json:"focusedNodeID"`
	Padding       Padding `json:"padding"`
	Root          *Node   `json:"root"`
}

type Node struct {
	ID            int         `json:"id"`
	SplitType     string      `json:"splitType"`
	SplitRatio    float64     `json:"splitRatio"`
	BirthRotation float64     `json:"birthRotation"`
	Vacant        bool        `json:"vacant"`
	Hidden        bool        `json:"hidden"`
	Sticky        bool        `json:"sticky"`
	Private       bool        `json:"private"`
	Locked        bool        `json:"locked"`
	Presel        interface{} `json:"-"`
	Rectangle     Rectangle   `json:"rectangle"`
	Constraints   Constraints `json:"constraints"`
	Client        *Client     `json:"client"`
	FirstChild    *Node       `json:"firstChild"`
	SecondChild   *Node       `json:"secondChild"`
}

type Client struct {
	ClassName         string    `json:"className"`
	InstanceName      string    `json:"instanceName"`
	BorderWidth       int       `json:"borderWidth"`
	State             string    `json:"state"`
	LastState         string    `json:"lastState"`
	Layer             string    `json:"layer"`
	LastLayer         string    `json:"lastLayer"`
	Urgent            bool      `json:"urgent"`
	Shown             bool      `json:"shown"`
	TiledRectangle    Rectangle `json:"tiledRectangle"`
	FloatingRectangle Rectangle `json:"floatingRectangle"`
}

type Rectangle struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Padding struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

type Constraints struct {
	MinWidth  int `json:"min_width"`
	MinHeight int `json:"min_height"`
}

func LoadState() (*State, error) {
	response, err := bspc.Run("wm", "-d")
	if err != nil {
		return nil, xerrors.Errorf("load wm state: %w", err)
	}
	var state *State
	if err := json.Unmarshal(response, &state); err != nil {
		return nil, xerrors.Errorf("load wm state: %w", err)
	}
	return state, nil
}

func (s *State) FocusedDesktop() (*Desktop, error) {
	for _, m := range s.Monitors {
		if m.ID == s.FocusedMonitorID {
			for _, d := range m.Desktops {
				if d.ID == m.FocusedDesktopID {
					return d, nil
				}
			}
			return nil, xerrors.Errorf("no desktop for id: %v", m.FocusedDesktopID)
		}
	}
	return nil, xerrors.Errorf("no monitor for id: %v", s.FocusedMonitorID)
}
