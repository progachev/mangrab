package ui

import "fmt"

type lang string

const (
	langRU lang = "ru"
	langEN lang = "en"
)

var i18nTable = map[lang]map[string]string{
	langRU: {
		"home.sub":           "качай мангу · собирай pdf · оставайся в терминале",
		"home.hint":          "вставьте ссылку MangaLib · Enter — загрузить · s настройки · Ctrl+C — выход",
		"home.loading":       "загрузка…",
		"home.empty_url":     "вставьте ссылку на тайтл MangaLib",
		"home.version":       "v0.1.0 · MangaLib → PDF",
		"chapters.label":     "ГЛАВЫ",
		"chapters.found":     "найдено глав: %d",
		"chapters.none":      "нет глав",
		"chapters.tools":     "space выбрать · a все · n снять · o папка · d скачать · s настройки",
		"chapters.meta":      "%d глав · %d выбрано",
		"chapters.pick_one":  "отметьте хотя бы одну главу",
		"chapters.no_dir":    "сначала выберите папку (клавиша o)",
		"chapters.mkdir":     "не удалось создать папку: %v",
		"chapters.range":     "  … %d–%d из %d",
		"output.label":       "ВЫВОД",
		"output.none":        "не выбрана (нажмите o)",
		"source.label":       "ИСТОЧНИК",
		"status.search":      "ищу главы…",
		"status.cancel":      "отменяю…",
		"status.home":        "вставьте ссылку и нажмите Enter",
		"status.folder":      "Папка: %s",
		"status.done":        "готово: %d PDF",
		"status.done_err":    "готово с ошибками: %d PDF, ошибок %d",
		"status.aborted":     "скачивание прервано: %v",
		"status.loaded":      "найдено глав: %d",
		"status.error_list":  "ошибка при загрузке списка глав",
		"download.getpages":  "получаю страницы: ",
		"download.done":      "готово: ",
		"download.page":      "Страница %d из %d",
		"download.chapter":   "глава %s",
		"header.working":     "работа",
		"header.ready":       "готово",
		"header.error":       "ошибка",
		"header.title":       "mangrab",
		"header.subtitle":    "тайтл: %s",
		"done.title":         "РЕЗУЛЬТАТ",
		"done.created":       "создано %d PDF",
		"done.failed":        "ошибок %d",
		"done.nothing":       "ничего не скачано",
		"done.folder":        "папка: ",
		"done.hint":          "↑↓ прокрутка · Enter — назад к главам · s — настройки · q — выход",
		"footer.workspace":   "↑↓ навигация · space выбрать · a все · n снять · o папка · d скачать · s настройки · Esc домой · q выход",
		"settings.title":     "НАСТРОЙКИ",
		"settings.output":    "Папка сохранения",
		"settings.lang":      "Язык (ru/en)",
		"settings.sources":   "Источники изображений",
		"settings.sources.empty": "список пуст · a — добавить",
		"settings.sources.hint":  "space вкл/выкл · a добавить · e изменить · d удалить",
		"settings.hint":      "↑↓ поле/строка · Tab сменить поле · space вкл/выкл · a добавить · e изменить · d удалить · Enter сохранить · Esc отмена",
		"settings.saved":     "настройки сохранены",
		"settings.bad_lang":  "язык: только ru или en",
		"app.error":          "ошибка: %v",
	},
	langEN: {
		"home.sub":           "grab manga · build pdf · stay in terminal",
		"home.hint":          "paste a MangaLib link · Enter to load · s settings · Ctrl+C to quit",
		"home.loading":       "loading…",
		"home.empty_url":     "insert a link to the title MangaLib",
		"home.version":       "v0.1.0 · MangaLib → PDF",
		"chapters.label":     "CHAPTERS",
		"chapters.found":     "found %d chapters",
		"chapters.none":      "no chapters",
		"chapters.tools":     "space toggle · a all · n none · o folder · d download · s settings",
		"chapters.meta":      "%d chapters · %d selected",
		"chapters.pick_one":  "select at least one chapter",
		"chapters.no_dir":    "choose a folder first (key o)",
		"chapters.mkdir":     "cannot create folder: %v",
		"chapters.range":     "  … %d–%d of %d",
		"output.label":       "OUTPUT",
		"output.none":        "not selected (press o)",
		"source.label":       "SOURCE",
		"status.search":      "loading chapters…",
		"status.cancel":      "cancelling…",
		"status.home":        "paste a link and press Enter",
		"status.folder":      "Folder: %s",
		"status.done":        "done: %d PDF",
		"status.done_err":    "done with errors: %d PDF, %d failed",
		"status.aborted":     "download aborted: %v",
		"status.loaded":      "found %d chapters",
		"status.error_list":  "error while loading chapters",
		"download.getpages":  "get pages: ",
		"download.done":      "done: ",
		"download.page":      "Page %d of %d",
		"download.chapter":   "chapter %s",
		"header.working":     "working",
		"header.ready":       "ready",
		"header.error":       "error",
		"header.title":       "mangrab",
		"header.subtitle":    "title: %s",
		"done.title":         "RESULT",
		"done.created":       "created %d PDF",
		"done.failed":        "failed %d",
		"done.nothing":       "nothing was downloaded",
		"done.folder":        "folder: ",
		"done.hint":          "↑↓ scroll · Enter — back to chapters · s settings · q — quit",
		"footer.workspace":   "↑↓ navigate · space toggle · a all · n none · o folder · d download · s settings · Esc home · q quit",
		"settings.title":     "SETTINGS",
		"settings.output":    "Output folder",
		"settings.lang":      "Language (ru/en)",
		"settings.sources":   "Image sources",
		"settings.sources.empty": "list is empty · a to add",
		"settings.sources.hint":  "space toggle · a add · e edit · d delete",
		"settings.hint":      "↑↓ field/row · Tab switch field · space toggle · a add · e edit · d delete · Enter save · Esc cancel",
		"settings.saved":     "settings saved",
		"settings.bad_lang":  "language: only ru or en",
		"app.error":          "error: %v",
	},
}

// hint — один сегмент подсказки: клавиша (или комбинация) и её описание.
type hint struct {
	key  string
	desc string
}

// hintsTable — структурированные подсказки с горячими клавишами.
var hintsTable = map[lang]map[string][]hint{
	langRU: {
		"home.hint": {
			{key: "Enter", desc: "загрузить"},
			{key: "s", desc: "настройки"},
			{key: "Ctrl+C", desc: "выход"},
		},
		"chapters.tools": {
			{key: "space", desc: "выбрать"},
			{key: "a", desc: "все"},
			{key: "n", desc: "снять"},
			{key: "o", desc: "папка"},
			{key: "d", desc: "скачать"},
			{key: "s", desc: "настройки"},
		},
		"footer.workspace": {
			{key: "↑↓", desc: "навигация"},
			{key: "space", desc: "выбрать"},
			{key: "a", desc: "все"},
			{key: "n", desc: "снять"},
			{key: "o", desc: "папка"},
			{key: "d", desc: "скачать"},
			{key: "s", desc: "настройки"},
			{key: "Esc", desc: "домой"},
			{key: "q", desc: "выход"},
		},
		"done.hint": {
			{key: "↑↓", desc: "прокрутка"},
			{key: "Enter", desc: "назад к главам"},
			{key: "s", desc: "настройки"},
			{key: "q", desc: "выход"},
		},
		"settings.hint": {
			{key: "↑↓", desc: "поле/строка"},
			{key: "Tab", desc: "сменить поле"},
			{key: "space", desc: "вкл/выкл"},
			{key: "a", desc: "добавить"},
			{key: "e", desc: "изменить"},
			{key: "d", desc: "удалить"},
			{key: "Enter", desc: "сохранить"},
			{key: "Esc", desc: "отмена"},
		},
		"settings.sources.hint": {
			{key: "space", desc: "вкл/выкл"},
			{key: "a", desc: "добавить"},
			{key: "e", desc: "изменить"},
			{key: "d", desc: "удалить"},
		},
	},
	langEN: {
		"home.hint": {
			{key: "Enter", desc: "load"},
			{key: "s", desc: "settings"},
			{key: "Ctrl+C", desc: "quit"},
		},
		"chapters.tools": {
			{key: "space", desc: "toggle"},
			{key: "a", desc: "all"},
			{key: "n", desc: "none"},
			{key: "o", desc: "folder"},
			{key: "d", desc: "download"},
			{key: "s", desc: "settings"},
		},
		"footer.workspace": {
			{key: "↑↓", desc: "navigate"},
			{key: "space", desc: "toggle"},
			{key: "a", desc: "all"},
			{key: "n", desc: "none"},
			{key: "o", desc: "folder"},
			{key: "d", desc: "download"},
			{key: "s", desc: "settings"},
			{key: "Esc", desc: "home"},
			{key: "q", desc: "quit"},
		},
		"done.hint": {
			{key: "↑↓", desc: "scroll"},
			{key: "Enter", desc: "back to chapters"},
			{key: "s", desc: "settings"},
			{key: "q", desc: "quit"},
		},
		"settings.hint": {
			{key: "↑↓", desc: "field/row"},
			{key: "Tab", desc: "switch field"},
			{key: "space", desc: "toggle"},
			{key: "a", desc: "add"},
			{key: "e", desc: "edit"},
			{key: "d", desc: "delete"},
			{key: "Enter", desc: "save"},
			{key: "Esc", desc: "cancel"},
		},
		"settings.sources.hint": {
			{key: "space", desc: "toggle"},
			{key: "a", desc: "add"},
			{key: "e", desc: "edit"},
			{key: "d", desc: "delete"},
		},
	},
}

func (m *Model) t(key string, args ...any) string {
	l := lang(m.lang)
	table, ok := i18nTable[l]
	if !ok {
		table = i18nTable[langRU]
	}
	s, ok := table[key]
	if !ok {
		s = i18nTable[langRU][key]
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

// hints возвращает структурированные сегменты подсказки для текущего языка.
func (m *Model) hints(name string) []hint {
	tbl, ok := hintsTable[lang(m.lang)]
	if !ok {
		tbl = hintsTable[langRU]
	}
	h, ok := tbl[name]
	if !ok {
		h = hintsTable[langRU][name]
	}
	return h
}

// Tr — перевод для произвольного места вне Model (например, в main).
func Tr(l, key string, args ...any) string {
	tbl, ok := i18nTable[lang(l)]
	if !ok {
		tbl = i18nTable[langRU]
	}
	s, ok := tbl[key]
	if !ok {
		s = i18nTable[langRU][key]
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

func trRaw(l, key string) string {
	tbl, ok := i18nTable[lang(l)]
	if !ok {
		tbl = i18nTable[langRU]
	}
	s, ok := tbl[key]
	if !ok {
		s = i18nTable[langRU][key]
	}
	return s
}