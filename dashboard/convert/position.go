package convert

import (
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func enhancePos(pos *gqlmodel.DashboardEntryPos, size string) *gqlmodel.DashboardEntryPos {
	switch size {
	case "mobile":
		pos.MinH = 2
		pos.MinW = 1
	case "desktop":
		pos.MinH = 3
		pos.MinW = 2
	default:
		panic("unknown size")
	}
	if pos.MinW > pos.W {
		pos.W = pos.MinW
	}
	if pos.MinH > pos.H {
		pos.H = pos.MinH
	}
	return pos
}

// EmptyPos returns a empty position.
func EmptyPos() string {
	internal := internalDashboardPositionV1{
		Version: "v1",
		H:       0,
		W:       0,
		X:       0,
		Y:       0,
	}

	if bytes, err := json.Marshal(&internal); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

func toInternalPosition(external gqlmodel.InputDashboardEntryPos) (string, error) {
	internal := internalDashboardPositionV1{Version: "v1"}
	if err := copier.Copy(&internal, &external); err != nil {
		return "", err
	}
	bytes, err := json.Marshal(&internal)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ApplyPos copies the position to the external position.
func ApplyPos(entry *model.DashboardEntry, pos *gqlmodel.InputResponsiveDashboardEntryPos) error {
	if pos == nil {
		return nil
	}

	if pos.Mobile != nil {
		lpos, err := toInternalPosition(*pos.Mobile)
		if err != nil {
			return err
		}
		entry.MobilePosition = lpos
	}
	if pos.Desktop != nil {
		mpos, err := toInternalPosition(*pos.Desktop)
		if err != nil {
			return err
		}
		entry.DesktopPosition = mpos
	}
	return nil
}

func toExternalPosition(pos string) (*gqlmodel.DashboardEntryPos, error) {
	if pos == "" {
		return nil, nil
	}
	internal := internalDashboardPositionV1{}
	if err := json.Unmarshal([]byte(pos), &internal); err != nil {
		return nil, err
	}

	if internal.Version != "v1" {
		log.Error().Interface("pos", internal).Msg("Invalid position version")
	}
	external := gqlmodel.DashboardEntryPos{}

	if err := copier.Copy(&external, &internal); err != nil {
		return nil, err
	}

	return &external, nil
}

type internalDashboardPositionV1 struct {
	Version string
	W       int `json:"w"`
	H       int `json:"h"`
	X       int `json:"x"`
	Y       int `json:"y"`
}
