package cv

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/disintegration/imaging"

	"github.com/otiai10/gosseract/v2"
	"github.com/vilsol/oshabi/config"

	"github.com/pkg/errors"
)

var client *gosseract.Client

func InitOCR() error {
	client = gosseract.NewClient()

	if err := client.DisableOutput(); err != nil {
		return errors.Wrap(err, "failed to disable tesseract output")
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.Wrap(err, "failed to find cache directory")
	}

	fullCacheDir := path.Join(cacheDir, "oshabi")
	if err := os.MkdirAll(fullCacheDir, 0777); err != nil {
		if !os.IsExist(err) {
			return errors.Wrap(err, "failed to create cache directory "+fullCacheDir)
		}
	}

	if err := client.SetTessdataPrefix(fullCacheDir); err != nil {
		return errors.Wrap(err, "failed to set tessdata prefix")
	}

	if err := client.SetLanguage("eng"); err != nil {
		return errors.Wrap(err, "failed setting OCR language")
	}

	if err := verifyLanguage(config.Get().Language); err != nil {
		return err
	}

	return nil
}

func verifyLanguage(language config.Language) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.Wrap(err, "failed to find cache directory")
	}

	fullCacheDir := path.Join(cacheDir, "oshabi")
	if err := os.MkdirAll(fullCacheDir, 0777); err != nil {
		if !os.IsExist(err) {
			return errors.Wrap(err, "failed to create cache directory "+fullCacheDir)
		}
	}

	langPath := path.Join(fullCacheDir, string(language)+".traineddata")
	_, err = os.Stat(langPath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to stat language file "+langPath)
	}

	out, err := os.Create(langPath)
	if err != nil {
		return errors.Wrap(err, "failed creating language file "+langPath)
	}
	defer out.Close()

	resp, err := http.Get(fmt.Sprintf("https://github.com/tesseract-ocr/tessdata/raw/main/%s.traineddata", string(language)))
	if err != nil {
		return errors.Wrap(err, "failed making request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read request body")
	}

	return nil
}

func ocr(img image.Image, whitelist string, mode gosseract.PageSegMode) (string, error) {
	if err := verifyLanguage(config.Get().Language); err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img)); err != nil {
		return "", errors.Wrap(err, "failed to encode image to png")
	}

	if err := client.SetWhitelist(whitelist); err != nil {
		return "", errors.Wrap(err, "failed to set ocr whitelist")
	}

	if err := client.SetPageSegMode(mode); err != nil {
		return "", errors.Wrap(err, "failed to set page segmentation mode")
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

func PrepareForOCR(img image.Image) image.Image {
	src := imaging.Grayscale(img)
	src = imaging.Invert(src)
	src = imaging.AdjustContrast(src, 30)
	return imaging.Sharpen(src, 1)
}
