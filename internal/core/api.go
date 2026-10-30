package core

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "net/http"
    "net/url"
    "regexp"
    "sort"
    "strconv"
    "time"
    "strings"
)

const APIOrigin = "https://api.cdnlibs.org"

var httpClient = &http.Client{Timeout: 90 * time.Second}

func getJSON(ctx context.Context, rawurl string) (any, error) {
    var lastErr error
    for attempt := 1; attempt <= 3; attempt++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
        if err != nil {
            return nil, err
        }
        req.Header.Set("Referer", "https://mangalib.org/")
        req.Header.Set("User-Agent", "Mozilla/5.0 mangrab/1.0")

        resp, err := httpClient.Do(req)
        if err == nil {
            var decoded any
            decErr := func() error {
                defer resp.Body.Close()
                if resp.StatusCode != http.StatusOK {
                    return fmt.Errorf("MangaLib вернул код %d", resp.StatusCode)
                }
                return json.NewDecoder(resp.Body).Decode(&decoded)
            }()
            if decErr == nil {
                return decoded, nil
            }
            lastErr = decErr
        } else {
            lastErr = err
        }

        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(time.Duration(attempt) * 900 * time.Millisecond):
        }
    }
    if lastErr == nil {
        lastErr = errors.New("Не удалось получить данные MangaLib")
    }
    return nil, lastErr
}

func findChapterList(v any) []map[string]any {
    var walk func(any) []map[string]any
    walk = func(x any) []map[string]any {
        switch t := x.(type) {
        case []any:
            out := make([]map[string]any, 0, len(t))
            for _, it := range t {
                if m, ok := it.(map[string]any); ok {
                    if _, has := m["number"]; has {
                        out = append(out, m)
                    }
                }
            }
            if len(out) > 0 {
                return out
            }
            for _, it := range t {
                if r := walk(it); len(r) > 0 {
                    return r
                }
            }
        case map[string]any:
            for _, it := range t {
                if r := walk(it); len(r) > 0 {
                    return r
                }
            }
        }
        return nil
    }
    return walk(v)
}

func findPages(v any) []map[string]any {
    var walk func(any) []map[string]any
    walk = func(x any) []map[string]any {
        if m, ok := x.(map[string]any); ok {
            if rawPages, has := m["pages"]; has {
                if arr, ok := rawPages.([]any); ok {
                    out := make([]map[string]any, 0, len(arr))
                    okAll := true
                    for _, it := range arr {
                        if pm, ok := it.(map[string]any); ok {
                            out = append(out, pm)
                        } else {
                            okAll = false
                            break
                        }
                    }
                    if okAll && len(out) > 0 {
                        return out
                    }
                }
            }
            for _, it := range m {
                if r := walk(it); len(r) > 0 {
                    return r
                }
            }
        } else if arr, ok := x.([]any); ok {
            for _, it := range arr {
                if r := walk(it); len(r) > 0 {
                    return r
                }
            }
        }
        return nil
    }
    return walk(v)
}

var numRE = regexp.MustCompile(`-?\d+(?:\.\d+)?`)

func numberKey(s string) (float64, string) {
    s = strings.ReplaceAll(s, ",", ".")
    if f, err := strconv.ParseFloat(s, 64); err == nil {
        return f, s
    }
    if m := numRE.FindString(s); m != "" {
        if f, err := strconv.ParseFloat(m, 64); err == nil {
            return f, s
        }
    }
    return math.Inf(1), s
}

func compareNumberLike(a, b string) int {
    af, as := numberKey(a)
    bf, bs := numberKey(b)
    if af < bf {
        return -1
    }
    if af > bf {
        return 1
    }
    return strings.Compare(as, bs)
}

func sortChapters(chapters []Chapter) {
    sort.SliceStable(chapters, func(i, j int) bool {
        if c := compareNumberLike(chapters[i].Volume, chapters[j].Volume); c != 0 {
            return c < 0
        }
        return compareNumberLike(chapters[i].Number, chapters[j].Number) < 0
    })
}

func sortPages(pages []Page) {
    sort.SliceStable(pages, func(i, j int) bool {
        return compareNumberLike(pageKey(pages[i]), pageKey(pages[j])) < 0
    })
}

func pageKey(p Page) string {
    if p.Slug != "" {
        return p.Slug
    }
    if p.Page != nil {
        return fmt.Sprint(p.Page)
    }
    return "0"
}

func strOr(m map[string]any, key string) string {
    if v, ok := m[key]; ok && v != nil {
        return fmt.Sprint(v)
    }
    return ""
}

func FetchChapters(ctx context.Context, mangaURL string) (*Manga, error) {
    slug, canonical, err := NormalizeMangaURL(mangaURL)
    if err != nil {
        return nil, err
    }
    raw, err := getJSON(ctx, APIOrigin+"/api/manga/"+url.PathEscape(slug)+"/chapters")
    if err != nil {
        return nil, err
    }
    list := findChapterList(raw)
    if len(list) == 0 {
        return nil, errors.New("MangaLib не вернул список глав для этого тайтла")
    }
    chapters := make([]Chapter, 0, len(list))
    for _, m := range list {
        chapters = append(chapters, Chapter{
            Number: strOr(m, "number"),
            Volume: strOr(m, "volume"),
            Name:   strOr(m, "name"),
        })
    }
    sortChapters(chapters)
    return &Manga{Slug: slug, URL: canonical, Chapters: chapters}, nil
}

func FetchChapterPages(ctx context.Context, slug string, ch Chapter) ([]Page, error) {
    q := url.Values{}
    if ch.Number == "" {
        q.Set("number", "0")
    } else {
        q.Set("number", ch.Number)
    }
    if ch.Volume == "" {
        q.Set("volume", "1")
    } else {
        q.Set("volume", ch.Volume)
    }
    endpoint := APIOrigin + "/api/manga/" + url.PathEscape(slug) + "/chapter?" + q.Encode()
    raw, err := getJSON(ctx, endpoint)
    if err != nil {
        return nil, err
    }
    rawPages := findPages(raw)
    if len(rawPages) == 0 {
        return nil, fmt.Errorf("не найдены страницы главы %s", ch.Number)
    }
    pages := make([]Page, 0, len(rawPages))
    for _, m := range rawPages {
        pages = append(pages, Page{
            Slug:  strOr(m, "slug"),
            URL:   strOr(m, "url"),
            Image: strOr(m, "image"),
            Page:  m["page"],
        })
    }
    sortPages(pages)
    return pages, nil
}
