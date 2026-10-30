package core

import (
    "errors"
    "net/url"
    "regexp"
    "strings"
)

var slugRE = regexp.MustCompile(`^\d+--[\w-]+$`)

func NormalizeMangaURL(raw string) (slug, canonical string, err error) {
    u, e := url.Parse(strings.TrimSpace(raw))
    if e != nil || u.Host == "" {
        return "", "", errors.New("Вставьте корректную ссылку на тайтл MangaLib")
    }
    host := strings.ToLower(u.Hostname())
    allowed := host == "mangalib.org" || strings.HasSuffix(host, ".mangalib.org") ||
        host == "mangalib.me" || strings.HasSuffix(host, ".mangalib.me")
    if u.Scheme != "https" || !allowed {
        return "", "", errors.New("Поддерживаются только HTTPS-ссылки на MangaLib")
    }

    segs := strings.Split(strings.Trim(u.Path, "/"), "/")
    for i, s := range segs {
        if s == "manga" && i+1 < len(segs) {
            slug = segs[i+1]
            break
        }
    }
    if slug == "" {
        for _, s := range segs {
            if slugRE.MatchString(s) {
                slug = s
                break
            }
        }
    }
    if slug == "" {
        return "", "", errors.New("Не удалось определить тайтл по этой ссылке")
    }
    return slug, "https://mangalib.org/ru/manga/" + slug, nil
}