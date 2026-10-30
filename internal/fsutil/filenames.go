package fsutil

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "github.com/progachev/mangrab/internal/core"
)

var badChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

func SafeFilename(s string) string {
    s = badChars.ReplaceAllString(s, "_")
    s = strings.TrimRight(strings.TrimSpace(s), ". ")
    if s == "" {
        return "Chapter"
    }
    return s
}

func ChapterFilename(ch core.Chapter, used map[string]bool) string {
    base := SafeFilename("Chapter " + ch.Number)
    name := base + ".pdf"
    key := strings.ToLower(name)
    if used[key] {
        name = base + " (Volume " + SafeFilename(ch.Volume) + ").pdf"
        key = strings.ToLower(name)
    }
    used[key] = true
    return name
}

func UniquePath(dir, name string) string {
    path := filepath.Join(dir, name)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return path
    }
    ext := filepath.Ext(name)
    base := strings.TrimSuffix(name, ext)
    for i := 2; i < 10000; i++ {
        candidate := filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base, i, ext))
        if _, err := os.Stat(candidate); os.IsNotExist(err) {
            return candidate
        }
    }
    return path
}
