package cv

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/vilsol/oshabi/data"

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

	language := config.Get().Language
	if err := verifyLanguage(language); err != nil {
		return err
	}

	wordsPath := path.Join(fullCacheDir, string(language)+".words")
	if err := client.SetVariable("user_words_file", wordsPath); err != nil {
		return errors.Wrap(err, "failed to set tessdata prefix")
	}

	return nil
}

var cleanRegex = regexp.MustCompile(`[",.%\d+]`)

func verifyLanguage(language config.Language) error {
	if err := client.SetLanguage(string(config.Get().Language)); err != nil {
		return errors.Wrap(err, "failed setting OCR language")
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

	langPath := path.Join(fullCacheDir, string(language)+".traineddata")
	_, err = os.Stat(langPath)
	if err != nil {
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
	}

	wordsPath := path.Join(fullCacheDir, string(language)+".words")
	_, err = os.Stat(wordsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to stat word file "+wordsPath)
		}

		wordMap := make(map[string]bool)
		for _, craft := range data.AllCrafts() {
			text := craft.Translations[config.Get().Language]
			clean := cleanRegex.ReplaceAllString(text, " ")
			for _, s := range strings.Split(clean, " ") {
				if s != "" {
					wordMap[s] = true
				}
			}
		}

		serialized := ""
		for s := range wordMap {
			serialized += s + "\n"
		}

		if err := os.WriteFile(wordsPath, []byte(serialized), 0777); err != nil {
			return errors.Wrap(err, "failed to write word file")
		}
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
	src := img

	size := src.Bounds().Dx() * src.Bounds().Dy()
	if size < 10000 {
		scale := math.Max(2, (10000/float64(size))/4)
		width := int(float64(src.Bounds().Dx()) * scale)
		height := int(float64(src.Bounds().Dy()) * scale)
		src = imaging.Resize(src, width, height, imaging.Linear)
	}

	src = imaging.Grayscale(src)
	src = imaging.Invert(src)
	src = imaging.AdjustContrast(src, 30)
	src = imaging.Sharpen(src, 1)

	return src
}
