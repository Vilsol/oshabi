package cv

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/data"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

	fullCacheDir := filepath.Join(cacheDir, "oshabi")
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

	wordsPath := filepath.Join(fullCacheDir, string(language)+".words")
	if err := client.SetVariable("user_words_file", wordsPath); err != nil {
		return errors.Wrap(err, "failed to set tessdata prefix")
	}

	return nil
}

var cleanRegex = regexp.MustCompile(`[",.%\d+]`)

func verifyLanguage(language config.Language) error {
	if err := client.SetLanguage(string(language)); err != nil {
		return errors.Wrap(err, "failed setting OCR language")
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.Wrap(err, "failed to find cache directory")
	}

	fullCacheDir := filepath.Join(cacheDir, "oshabi")
	if err := os.MkdirAll(fullCacheDir, 0777); err != nil {
		if !os.IsExist(err) {
			return errors.Wrap(err, "failed to create cache directory "+fullCacheDir)
		}
	}

	langPath := filepath.Join(fullCacheDir, string(language)+".traineddata")
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

	wordsPath := filepath.Join(fullCacheDir, string(language)+".words")
	_, err = os.Stat(wordsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to stat word file "+wordsPath)
		}

		wordMap := make(map[string]bool)
		for _, craft := range data.AllCrafts() {
			text := craft.Translations[language]
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

	patternsPath := filepath.Join(fullCacheDir, "digits.patterns")
	_, err = os.Stat(patternsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to stat patterns file "+patternsPath)
		}

		//		if err := os.WriteFile(patternsPath, []byte("\\d\\d\n"), 0777); err != nil {
		if err := os.WriteFile(patternsPath, []byte(`6\d
7\d
8\d
`), 0777); err != nil {
			return errors.Wrap(err, "failed to write digits pattern file")
		}
	}

	return nil
}

func ocr(img image.Image, whitelist string, mode gosseract.PageSegMode, crop bool, languageOverride bool) (string, error) {
	if languageOverride {
		if err := verifyLanguage(config.LanguageEnglish); err != nil {
			return "", err
		}
	} else {
		if err := verifyLanguage(config.Get().Language); err != nil {
			return "", err
		}
	}

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, PrepareForOCR(img, crop)); err != nil {
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

	if (config.Get().Language == config.LanguageChinese ||
		config.Get().Language == config.LanguageTaiwanese ||
		config.Get().Language == config.LanguageJapanese) && crop {
		if err := client.SetVariable("preserve_interword_spaces", "1"); err != nil {
			return "", errors.Wrap(err, "failed to set preserve_interword_spaces")
		}
	} else {
		if err := client.SetVariable("preserve_interword_spaces", "0"); err != nil {
			return "", errors.Wrap(err, "failed to set preserve_interword_spaces")
		}
	}

	if languageOverride {
		if err := client.SetVariable("user_patterns_file", filepath.Join(client.TessdataPrefix, "digits.patterns")); err != nil {
			return "", errors.Wrap(err, "failed to set user_patterns_file")
		}
	} else {
		if err := client.SetVariable("user_patterns_file", ""); err != nil {
			return "", errors.Wrap(err, "failed to set user_patterns_file")
		}
	}

	text, err := client.Text()
	if err != nil {
		return "", errors.Wrap(err, "failed to ocr the image")
	}

	log.Debug().Str("text", text).Str("lang", string(config.Get().Language)).Str("w", whitelist).Msg("ocr")

	return text, nil
}

const targetPixelCount = float64(150000)

func PrepareForOCR(img image.Image, crop bool) image.Image {
	out := imaging.Clone(img)

	size := float64(out.Bounds().Dx() * out.Bounds().Dy())
	if size < targetPixelCount {
		scale := math.Min(10, math.Max(2, math.Sqrt(targetPixelCount/size)))
		width := int(float64(out.Bounds().Dx()) * scale)
		height := int(float64(out.Bounds().Dy()) * scale)
		out = imaging.Resize(out, width, height, imaging.Linear)
	}

	// Looks like on any of these languages, grayscale actually hurts
	if !(config.Get().Language == config.LanguageKorean &&
		config.Get().Language == config.LanguageChinese &&
		config.Get().Language == config.LanguageJapanese &&
		config.Get().Language == config.LanguageTaiwanese) || !crop {
		out = imaging.Grayscale(out)
	}

	out = imaging.Invert(out)
	out = imaging.AdjustContrast(out, 30)
	out = imaging.Sharpen(out, 1)

	rightPixel := out.Bounds().Dx()
	bottomPixel := 0

	found := false
	for x := out.Bounds().Dx() - 1; x > 0; x-- {
		for y := out.Bounds().Dy() - 1; y > 0; y-- {
			px := out.NRGBAAt(x, y)
			if px.R < 150 || px.G < 150 || px.B < 150 {
				if !found {
					rightPixel = x
					found = true
				}

				if y > bottomPixel {
					bottomPixel = y
				}

				// This needs to be improved
				out.SetNRGBA(x, y, color.NRGBA{
					R: px.R / 2,
					G: px.G / 2,
					B: px.B / 2,
					A: px.A,
				})
			}
		}
	}

	rightPixel += ScaleN(15)
	bottomPixel += ScaleN(10)

	if crop {
		out = imaging.Crop(out, image.Rect(0, 0, rightPixel, bottomPixel))
	}

	return out
}

var whitelistCache = make(map[config.Language]string)

func GetWhitelist() string {
	if whitelist, ok := whitelistCache[config.Get().Language]; ok {
		return whitelist
	}

	characterMap := make(map[rune]bool)
	for _, craft := range data.AllCrafts() {
		text := craft.Translations[config.Get().Language]
		for _, s := range text {
			characterMap[s] = true
		}
	}

	var result strings.Builder
	for r := range characterMap {
		result.WriteRune(r)
	}

	s := []rune(result.String())
	sort.Slice(s, func(i int, j int) bool {
		return s[i] < s[j]
	})

	whitelistCache[config.Get().Language] = string(s)

	return string(s)
}
