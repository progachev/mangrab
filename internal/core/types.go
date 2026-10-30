package core

type Chapter struct {
    Number   string
    Volume   string
    Name     string
    Selected bool
}

type Page struct {
    Slug  string
    Page  any
    URL   string
    Image string
}

type Manga struct {
    Slug     string
    URL      string
    Chapters []Chapter
}
