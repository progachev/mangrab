package core

import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "image"
    "image/draw"
    "image/jpeg"
    _ "image/png"
    "math"

    "github.com/disintegration/imaging"
    "github.com/signintech/gopdf"
)

const (
    maxPDFDimension = 14000.0
    pdfDPI          = 150.0
)

func CreateChapterPDF(
    ctx context.Context,
    ch Chapter,
    pages []Page,
    progress func(pageIndex, totalPages int, message string),
) ([]byte, error) {
    if len(pages) == 0 {
        return nil, errors.New("Нет страниц для PDF")
    }

    pdf := gopdf.GoPdf{}
    pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}})
    pdf.SetInfo(gopdf.PdfInfo{
        Title:   "Chapter " + ch.Number,
        Creator: "mangrab",
    })

    preferred := ""
    for i, p := range pages {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        if progress != nil {
            progress(i, len(pages), fmt.Sprintf("Страница %d из %d", i+1, len(pages)))
        }
        img, err := FetchImage(ctx, p, preferred)
        if err != nil {
            return nil, err
        }
        preferred = img.Server

        data, err := normalizeToJPEG(img.Bytes)
        if err != nil {
            return nil, fmt.Errorf("Страница %v: %w", p.Slug, err)
        }
        cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
        if err != nil {
            return nil, err
        }
        wPt := float64(cfg.Width) * 72.0 / pdfDPI
        hPt := float64(cfg.Height) * 72.0 / pdfDPI
        scale := math.Min(1, maxPDFDimension/math.Max(wPt, hPt))
        wPt = math.Max(1, wPt*scale)
        hPt = math.Max(1, hPt*scale)

        pdf.AddPageWithOption(gopdf.PageOption{
            PageSize: &gopdf.Rect{W: wPt, H: hPt},
        })
        holder, err := gopdf.ImageHolderByBytes(data)
        if err != nil {
            return nil, err
        }
        if err := pdf.ImageByHolder(holder, 0, 0, &gopdf.Rect{W: wPt, H: hPt}); err != nil {
            return nil, err
        }
    }

    var out bytes.Buffer
    if err := pdf.Write(&out); err != nil {
        return nil, err
    }
    return out.Bytes(), nil
}

func normalizeToJPEG(raw []byte) ([]byte, error) {
    if len(raw) >= 3 && raw[0] == 0xFF && raw[1] == 0xD8 && raw[2] == 0xFF {
        return raw, nil
    }
    img, err := imaging.Decode(bytes.NewReader(raw), imaging.AutoOrientation(true))
    if err != nil {
        return nil, fmt.Errorf("Не удалось распознать изображение: %w", err)
    }
    bounds := img.Bounds()
    flat := image.NewRGBA(bounds)
    for i := range flat.Pix {
        flat.Pix[i] = 0xFF
    }
    draw.Draw(flat, flat.Bounds(), img, bounds.Min, draw.Over)

    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, flat, &jpeg.Options{Quality: 92}); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}