package data

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"

	"github.com/pkg/errors"
)

//go:embed crafts.json
var craftsJSON []byte
var crafts map[string]HarvestCraft

//go:embed menu_button.png
var menuButtonPNG []byte
var MenuButton image.Image

//go:embed top_right_x.png
var topRightXPNG []byte
var TopRightX image.Image

//go:embed scroll_up.png
var scrollUpPNG []byte
var ScrollUp image.Image

//go:embed scroll_down.png
var scrollDownPNG []byte
var ScrollDown image.Image

//go:embed info_button.png
var infoButtonPNG []byte
var InfoButton image.Image

//go:embed count_corner.png
var countCornerPNG []byte
var CountCorner image.Image

func InitData() error {
	var err error
	if err = json.Unmarshal(craftsJSON, &crafts); err != nil {
		return errors.Wrap(err, "failed parsing crafts.json")
	}

	MenuButton, _, err = image.Decode(bytes.NewReader(menuButtonPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing menu_button.png")
	}

	TopRightX, _, err = image.Decode(bytes.NewReader(topRightXPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing top_right_x.png")
	}

	ScrollUp, _, err = image.Decode(bytes.NewReader(scrollUpPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing scroll_up.png")
	}

	ScrollDown, _, err = image.Decode(bytes.NewReader(scrollDownPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing scroll_down.png")
	}

	InfoButton, _, err = image.Decode(bytes.NewReader(infoButtonPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing info_button.png")
	}

	CountCorner, _, err = image.Decode(bytes.NewReader(countCornerPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing count_corner.png")
	}

	return nil
}
