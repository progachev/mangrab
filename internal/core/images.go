package core

import (
    "context"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "sort"
    "strings"
    "time"
)

var defaultImageServers = []string{
	"https://img3.cdnlibs.org",
	"https://img4.cdnlibs.org",
	"https://img5.cdnlibs.org",
	"https://imglib.info",
	"https://crops.mangalib.me",
	"https://img3.mixlib.me",
	"https://img4.mixlib.me",
	"https://img5.mixlib.me",
}

var imageServers = append([]string(nil), defaultImageServers...)

// SetImageServers задаёт источники изображений (порядок важен — пробуются по порядку).
// Пустой список сбрасывает на дефолтный набор.
func SetImageServers(list []string) {
	clean := make([]string, 0, len(list))
	seen := map[string]struct{}{}
	for _, s := range list {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
			s = "https://" + s
		}
		s = strings.TrimRight(s, "/")
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		clean = append(clean, s)
	}
	if len(clean) == 0 {
		imageServers = append([]string(nil), defaultImageServers...)
		return
	}
	imageServers = clean
}

// ImageServers возвращает копию текущего списка источников.
func ImageServers() []string {
	out := make([]string, len(imageServers))
	copy(out, imageServers)
	return out
}

// DefaultImageServers возвращает копию дефолтного списка.
func DefaultImageServers() []string {
	out := make([]string, len(defaultImageServers))
	copy(out, defaultImageServers)
	return out
}

func hostOf(rawurl string) string {
    if u, err := url.Parse(rawurl); err == nil {
        return u.Scheme + "://" + u.Host
    }
    return ""
}

func imageCandidates(p Page) []string {
    raw := strings.TrimSpace(p.URL)
    if raw == "" {
        raw = strings.TrimSpace(p.Image)
    }
    if raw == "" {
        return nil
    }
    if u, err := url.Parse(raw); err == nil && u.IsAbs() {
        set := map[string]struct{}{raw: {}}
        for _, s := range imageServers {
            set[strings.TrimRight(s, "/")+u.Path] = struct{}{}
        }
        out := make([]string, 0, len(set))
        for k := range set {
            out = append(out, k)
        }
        return out
    }
    suffix := raw
    if !strings.HasPrefix(suffix, "/") {
        suffix = "/" + suffix
    }
    out := make([]string, 0, len(imageServers))
    for _, s := range imageServers {
        out = append(out, strings.TrimRight(s, "/")+suffix)
    }
    return out
}

type ImageResult struct {
    Bytes       []byte
    ContentType string
    URL         string
    Server      string
}

func FetchImage(ctx context.Context, p Page, preferred string) (*ImageResult, error) {
    candidates := imageCandidates(p)
    if len(candidates) == 0 {
        return nil, errors.New("нет кандидатов URL для изображения")
    }
    sort.SliceStable(candidates, func(i, j int) bool {
        hi := hostOf(candidates[i]) == preferred
        hj := hostOf(candidates[j]) == preferred
        return hi && !hj
    })

    var lastErr error
    for pass := 0; pass < 2; pass++ {
        for _, imageURL := range candidates {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            default:
            }
            req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
            if err != nil {
                lastErr = err
                continue
            }
            req.Header.Set("Referer", "https://mangalib.org/")
            req.Header.Set("User-Agent", "Mozilla/5.0 mangrab/1.0")

            resp, err := httpClient.Do(req)
            if err != nil {
                lastErr = err
                continue
            }
            body, readErr := io.ReadAll(resp.Body)
            resp.Body.Close()
            if readErr != nil {
                lastErr = readErr
                continue
            }
            if resp.StatusCode != http.StatusOK {
                lastErr = fmt.Errorf("Код %d", resp.StatusCode)
                continue
            }
            if len(body) < 500 {
                lastErr = errors.New("Пустое изображение")
                continue
            }
            return &ImageResult{
                Bytes:       body,
                ContentType: resp.Header.Get("Content-Type"),
                URL:         imageURL,
                Server:      hostOf(imageURL),
            }, nil
        }
        if pass < 1 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(time.Duration(pass+1) * 800 * time.Millisecond):
            }
        }
    }
    label := p.Slug
    if label == "" && p.Page != nil {
        label = fmt.Sprint(p.Page)
    }
    if label == "" {
        label = "?"
    }
    return nil, fmt.Errorf("Не удалось скачать страницу %s: %v", label, lastErr)
}
