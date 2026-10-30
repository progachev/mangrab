# mangrab

TUI-приложение на Go для скачивания глав манги с MangaLib и сборки их в PDF — только Go, HTTP и PDF.

```
 ███╗   ███╗ █████╗ ███╗   ██╗ ██████╗ ██████╗  █████╗ ██████╗
 ████╗ ████║██╔══██╗████╗  ██║██╔════╝ ██╔══██╗██╔══██╗██╔══██╗
 ██╔████╔██║███████║██╔██╗ ██║██║  ███╗██████╔╝███████║██████╔╝
 ██║╚██╔╝██║██╔══██║██║╚██╗██║██║   ██║██╔══██╗██╔══██║██╔══██╗
 ██║ ╚═╝ ██║██║  ██║██║ ╚████║╚██████╔╝██║  ██║██║  ██║██████╔╝
 ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝
```

## Возможности

- 🔍 **Загрузка списка глав** по ссылке на тайтл MangaLib (`https://mangalib.org/ru/manga/...`).
- ✅ **Выбор отдельных глав или всех сразу** — `space` по одной, `a` — всё, `n` — снять.
- 📄 **Генерация PDF** — 150 DPI, ограничение размера страницы 14000 pt, авто-ориентация по EXIF, «сплющивание» альфа-канала на белый фон.
- 🖼 **Обход нескольких зеркал изображений** — `cdnlibs.org`, `mixlib.me`, `imglib.info`, `crops.mangalib.me`, автоматический перебор при ошибках.
- 🔄 **Ретраи и «залипание» на рабочем зеркале** — приоритет сервера, который уже сработал в текущей главе.
- 🗂 **Безопасные имена файлов** — недопустимые символы заменяются, дубликаты не перезаписываются, существующие файлы получают суффикс `(2)`, `(3)`…
- 🖥 **Прогресс-бар и спиннер** прямо в терминале.
- ⏹ **Отмена скачивания** по `Esc` / `Ctrl+C` через `context.Context`.
- 💾 **Локальный `config.json`** рядом с бинарником: папка сохранения, язык, список источников, последняя ссылка.
- 🌐 **Локализация интерфейса** — `ru` / `en`.
- ⚙️ **Экран настроек прямо в TUI** (`s`) — редактирование папки, языка и источников без правки файла.
- 🎨 **Цветные хоткеи** в подсказках и экран результатов со скроллом по созданным файлам и ошибкам.

## Требования

- Go 1.22+

## Установка

```bash
git clone https://github.com/progachev/mangrab.git
cd mangrab
go mod tidy
go build -o mangrab .
```

```bash
go run .
```

## Использование

### Основные клавиши

| Клавиша | Действие |
| ------------- | ------------- |
| Enter | Загрузить главы по введённой ссылке |
| ↑ / ↓ или k / j | Перемещение по списку глав |
| Space | Отметить / снять главу |
| a | Отметить все |
| n | Снять выделение |
| o | Выбрать папку (см. ниже) |
| d | Начать скачивание выбранных глав |
| s | Открыть настройки (папка, язык, источники) |
| Esc | Отменить скачивание / вернуться на главный экран |
| q / Ctrl+C | Выход |

> В режиме ввода ссылки `s` не перехватывается — попадает в текст. Открыть настройки в этот момент нельзя.

### Экран настроек (`s`)

Позволяет менять всё прямо из TUI:

| Поле | Описание |
| ------------- | ------------- |
| Output folder | Папка сохранения PDF |
| Language | `ru` или `en` — интерфейс перерисовывается сразу после сохранения |
| Image sources | Список зеркал; порядок = порядок попыток. Можно включать/отключать, добавлять, редактировать, удалять |

Навигация: `↑↓` / `Tab`, `space` — вкл/выкл источник, `a` — добавить, `e` — изменить, `d` — удалить, `Enter` — сохранить и применить, `Esc` — отмена.

### Выбор папки

TUI не использует системный диалог выбора папки. Источники значения:

1. `outputDir` из `config.json`.
2. Дефолт — `~/Downloads` (или `$HOME`, если `Downloads` нет).

Нажми `o` — применится значение из `MLLIB_DIR` (или текущее из конфига / дефолт). Нажми `s` → поле «Output folder» — чтобы задать произвольный путь вручную.

### Конфигурация

Файл `config.json` лежит **рядом с бинарником** (при `go run` — в текущей директории). Путь можно переопределить переменной `MANGRAB_CONFIG`:

```bash
MANGRAB_CONFIG=/etc/mangrab.json mangrab
```

Порядок загрузки:

1. Пользовательский `~/.config/mangrab/config.json`.
2. Локальный `config.json` рядом с бинарником — **переопределяет** пользовательский.

Пример `config.json`:

```json
{
  "lastMangaUrl": "https://mangalib.org/ru/manga/7559--darren-shan",
  "outputDir": "/home/user/Books/Manga",
  "lang": "ru",
  "sources": [
    "https://img3.cdnlibs.org",
    "https://img4.cdnlibs.org",
    "https://img5.cdnlibs.org",
    "https://imglib.info",
    "https://crops.mangalib.me",
    "https://img3.mixlib.me",
    "https://img4.mixlib.me",
    "https://img5.mixlib.me"
  ],
  "sourcesOff": [
    "https://imglib.info"
  ]
}
```

| Поле | Назначение |
| ------------- | ------------- |
| `lastMangaUrl` | Последняя открытая ссылка (подставляется при старте) |
| `outputDir` | Папка сохранения PDF |
| `lang` | `ru` / `en` |
| `sources` | Включённые зеркала изображений; порядок важен — пробуются по порядку |
| `sourcesOff` | Отключённые зеркала (сохраняются, но не используются) |

Если `sources` пуст, используется встроенный набор зеркал по умолчанию. Файл создаётся автоматически при сохранении из экрана настроек или при старте скачивания (сохраняется `lastMangaUrl`).

## Структура проекта

```text
mangrab/
├── go.mod
├── go.sum
├── main.go
├── config.json                 // пример локального конфига
└── internal/
    ├── core/                     // ядро: API, изображения, PDF
    │   ├── api.go                // FetchChapters, FetchChapterPages
    │   ├── images.go             // зеркала, FetchImage, SetImageServers
    │   ├── pdf.go                // CreateChapterPDF, normalizeToJPEG
    │   ├── types.go              // Chapter, Page, Manga
    │   └── url.go                // NormalizeMangaURL
    ├── fsutil/                   // имена файлов
    │   └── filenames.go          // SafeFilename, ChapterFilename, UniquePath
    └── ui/                       // TUI на Bubble Tea
        ├── config.go             // локальный config.json
        ├── i18n.go               // i18nTable, hintsTable, t(), Tr()
        ├── logo.go               // логотип
        ├── model.go              // Model, Update, handleKey, команды
        ├── settings.go           // состояние экрана настроек
        ├── styles.go             // lipgloss-стили
        └── views.go              // рендер экранов
```

## Архитектура

Приложение построено по схеме Model–Update–View (Bubble Tea):

```text
┌──────────────────────┐
│  main.go             │  LoadConfig → NewModel → tea.Program
└──────────┬───────────┘
           ▼
┌──────────────────────┐
│  internal/ui         │  Model → Update → View
│  ─────────────────   │
│  handleKey           │  клавиатура
│  handleSettings      │  экран настроек (s)
│  findChaptersCmd     │  tea.Cmd → API
│  downloadCmd         │  tea.Cmd → PDF + файлы
└──────────┬───────────┘
           ▼
┌──────────────────────┐
│  internal/core       │  HTTP, зеркала, PDF-сборка
│  ─────────────────   │
│  FetchChapters       │
│  FetchChapterPages   │
│  SetImageServers     │  применяет список источников из конфига
│  FetchImage          │
│  CreateChapterPDF    │
└──────────┬───────────┘
           ▼
┌──────────────────────┐
│  internal/fsutil     │  Chapter N.pdf, уникализация имён
└──────────────────────┘
```

## Поток данных

```text
config.json (рядом с бинарником)
   │
   ▼
ui.LoadConfig ──► core.SetImageServers(cfg.Sources)
   │
URL MangaLib
   │
   ▼
NormalizeMangaURL ──► slug
   │
   ▼
FetchChapters ──► [Chapter, ...]
   │
   ▼ (для каждой выбранной главы)
FetchChapterPages ──► [Page, ...]
   │
   ▼
CreateChapterPDF
   ├─ FetchImage  (перебор зеркал из cfg.Sources, preferredServer)
   ├─ normalizeToJPEG
   └─ gopdf.Write  ──► []byte
   │
   ▼
os.WriteFile(dir/Chapter N.pdf)
```

## Отмена

Все долгие операции принимают `context.Context`. Нажатие `Esc` / `Ctrl+C` вызывает `m.cancel()`, что отменяет `FetchChapterPages` и `CreateChapterPDF` на ближайшей проверке `ctx.Done()`. Уже созданные PDF при этом сохраняются, а экран результатов показывает, что успело скачаться и что — нет.

## Зависимости

| Модуль | Назначение |
| ------------- | ------------- |
| github.com/charmbracelet/bubbletea | TUI-фреймворк (MVU) |
| github.com/charmbracelet/bubbles | Готовые компоненты: textinput, spinner, progress |
| github.com/charmbracelet/lipgloss | Стили для терминала |
| github.com/signintech/gopdf | Генерация PDF |
| github.com/disintegration/imaging | Декодирование изображений и конвертация в JPEG |

## Разработка

### Сборка под разные ОС

```powershell
$env:GOOS="linux";   $env:GOARCH="amd64"; go build -o dist/mangrab-linux-amd64 .
$env:GOOS="darwin";  $env:GOARCH="arm64"; go build -o dist/mangrab-darwin-arm64 .
$env:GOOS="darwin";  $env:GOARCH="amd64"; go build -o dist/mangrab-darwin-amd64 .
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o dist/mangrab-windows-amd64.exe .
```

## Ограничения

- Поддерживаются только **HTTPS-ссылки** на домены `mangalib.org` / `mangalib.me` (и поддомены).
- Приложение не обходит авторизацию и ограничения доступа MangaLib — работает с публично доступными данными.
- Для некоторых тайтлов/глав API может не вернуть страницы; такие главы попадают в список ошибок на экране результатов.

## TO-DO

- [ ] Переделать интерфейс
- [ ] Оптимизации
- [ ] Работа над ошибками
- [ ] Полноценный prompt выбора папки внутри TUI (с автодополнением по Tab).
- [ ] Параллельная загрузка глав через errgroup с ограничением concurrency.
- [ ] Поддержка диапазона `1-10` при выборе глав.
- [ ] Тесты для `core` (URL, сортировки, кандидаты картинок).
- [ ] CI-сборка артефактов под Linux/macOS/Windows.
- [ ] Горячая клавиша `u` — редактировать ссылку на экране глав без возврата на Home.

## Лицензия

MIT — см. [LICENSE](https://github.com/progachev/mangrab/blob/main/LICENSE).

## Благодарности

- [Charm](https://charm.sh/) за Bubble Tea, Bubbles и Lipgloss.
- [signintech](https://github.com/signintech) за gopdf.
- [disintegration](https://github.com/disintegration) за imaging.