package data

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"

	"github.com/vilsol/oshabi/types"

	"github.com/pkg/errors"
)

//go:embed crafts.json
var craftsJSON []byte
var crafts map[types.HarvestType]HarvestCraft

//go:embed menu_button.png
var menuButtonPNG []byte
var MenuButton image.Image

//go:embed info_button.png
var infoButtonPNG []byte
var InfoButton image.Image

//go:embed count_corner.png
var countCornerPNG []byte
var CountCorner image.Image

//go:embed horticrafting.png
var horticraftingPNG []byte
var Horticrafting image.Image

//go:embed inventory.png
var inventoryPNG []byte
var Inventory image.Image

//go:embed alert.mp3
var AlertMP3 []byte

func InitData() error {
	var err error
	if err = json.Unmarshal(craftsJSON, &crafts); err != nil {
		return errors.Wrap(err, "failed parsing crafts.json")
	}

	MenuButton, _, err = image.Decode(bytes.NewReader(menuButtonPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing menu_button.png")
	}

	InfoButton, _, err = image.Decode(bytes.NewReader(infoButtonPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing info_button.png")
	}

	CountCorner, _, err = image.Decode(bytes.NewReader(countCornerPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing count_corner.png")
	}

	Horticrafting, _, err = image.Decode(bytes.NewReader(horticraftingPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing horticrafting.png")
	}

	Inventory, _, err = image.Decode(bytes.NewReader(inventoryPNG))
	if err != nil {
		return errors.Wrap(err, "failed parsing inventory.png")
	}

	return nil
}
