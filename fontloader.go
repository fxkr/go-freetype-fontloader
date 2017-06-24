package fontloader

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
)

const (
	FONTCONFIG_MATCH_PROGRAM = "fc-match"
)

type FontCache interface {
	Load(fontName string) (*truetype.Font, error)
}

type fontCache struct {

	// Maps font names to paths.
	fontPathCache map[string]string

	// Maps font paths to loaded fonts.
	fontCache map[string]*truetype.Font
}

func NewFontCache() FontCache {
	return &fontCache{
		make(map[string]string),
		make(map[string]*truetype.Font),
	}
}

func (fc *fontCache) Load(fontName string) (*truetype.Font, error) {

	// Get font path from name (possibly from cache)
	fontPath, wasCached := fc.fontPathCache[fontName]
	if !wasCached {
		var err error
		fontPath, err = findFont(fontName)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not find non-cached font: %s", fontName))
		}
		fc.fontPathCache[fontName] = fontPath
	}

	// Get font from path (possibly from cache)
	font, wasCached := fc.fontCache[fontPath]
	if !wasCached {
		var err error
		font, err = loadFont(fontPath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not load non-cached font %s from %s", fontName, fontPath))
		}
		fc.fontCache[fontPath] = font
	}

	return font, nil
}

func Load(fontName string) (*truetype.Font, error) {

	// Get font path from name
	fontPath, err := findFont(fontName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not find font: %s", fontName))
	}

	// Get font from path
	font, err := loadFont(fontPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not load font %s from %s", fontName, fontPath))
	}

	return font, nil
}

func findFont(fontName string) (string, error) {

	// Makes this package usable without fc-match, so not just an optimization.
	if isAbsolutePath := strings.HasPrefix(fontName, "/"); isAbsolutePath {
		fontPath := fontName
		return fontPath, nil
	}

	cmd := exec.Command(FONTCONFIG_MATCH_PROGRAM, "--format=%{file}", fontName)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create pipe to fontconfig")
	}

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Failed to invoke fontconfig's fc-match command")
	}

	fontPathBytes, err := ioutil.ReadAll(pipe)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read output of fontconfig")
	}
	fontPath := string(fontPathBytes)

	if fontPath == "" {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to find font via fontconfig for: %s", fontName))
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.Wrap(err, "fontconfig's fc-match return with non-zero exit code")
	}

	return fontPath, nil
}

func loadFont(fontPath string) (*truetype.Font, error) {

	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to open font file: %s", fontPath))
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse truetype font file: %s", fontPath))
	}

	return font, nil
}
