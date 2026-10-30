package ui

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	LastMangaURL string   `json:"lastMangaUrl,omitempty"`
	OutputDir    string   `json:"outputDir,omitempty"`
	Lang         string   `json:"lang,omitempty"`
	Sources      []string `json:"sources,omitempty"`
    SourcesOff   []string `json:"sourcesOff,omitempty"`
}

// localConfigPath — конфиг рядом с бинарником (или cwd при `go run`).
func localConfigPath() string {
	if env := strings.TrimSpace(os.Getenv("MANGRAB_CONFIG")); env != "" {
		return env
	}
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		// при `go run` бинарник лежит во временной папке — используем cwd
		if strings.Contains(dir, os.TempDir()) {
			if wd, err := os.Getwd(); err == nil {
				dir = wd
			}
		}
		return filepath.Join(dir, "config.json")
	}
	if wd, err := os.Getwd(); err == nil {
		return filepath.Join(wd, "config.json")
	}
	return "config.json"
}

func userConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	full := filepath.Join(dir, "mangrab")
	if err := os.MkdirAll(full, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(full, "config.json"), nil
}

// LoadConfig: сначала пользовательский, затем локальный (переопределяет).
func LoadConfig() config {
	var c config

	if p, err := userConfigPath(); err == nil {
		if raw, err := os.ReadFile(p); err == nil {
			_ = json.Unmarshal(raw, &c)
		}
	}

	if raw, err := os.ReadFile(localConfigPath()); err == nil {
		var local config
		if json.Unmarshal(raw, &local) == nil {
			if local.LastMangaURL != "" {
				c.LastMangaURL = local.LastMangaURL
			}
			if local.OutputDir != "" {
				c.OutputDir = local.OutputDir
			}
			if local.Lang != "" {
				c.Lang = local.Lang
			}
			if len(local.Sources) > 0 {
				c.Sources = local.Sources
			}
            if len(local.SourcesOff) > 0 {
                c.SourcesOff = local.SourcesOff
            }
		}
	}

	if c.Lang != "en" && c.Lang != "ru" {
	c.Lang = "ru"
    }
	return c
}

// SaveConfig пишет в локальный файл рядом с бинарником (атомарно).
func SaveConfig(c config) error {
	raw, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	path := localConfigPath()
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func DefaultOutputDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Не удалось определить домашнюю папку")
	}
	downloads := filepath.Join(home, "Downloads")
	if _, err := os.Stat(downloads); err == nil {
		return downloads, nil
	}
	return home, nil
}