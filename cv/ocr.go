package cv

import (
	"bytes"
	"image"
	"image/png"

	"github.com/pkg/errors"
)

func ocr(img image.Image, whitelist string) (string, error) {
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img)); err != nil {
		return "", errors.Wrap(err, "failed to encode image to png")
	}

	if err := client.SetWhitelist(whitelist); err != nil {
		return "", errors.Wrap(err, "failed to set ocr whitelist")
	}

	if err := client.SetImageFromBytes(buff.Bytes()); err != nil {
		return "", errors.Wrap(err, "failed to update ocr buffer")
	}

	text, err := client.Text()
	if err != nil {
		return "", errors.Wrap(err, "failed to ocr the image")
	}

	return text, nil
}
