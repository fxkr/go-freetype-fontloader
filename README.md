# go-freetype-fontloader

Finds and loads TrueType font files for use with [github.com/golang/freetype](https://github.com/golang/freetype/).

Includes a simple font cache.

* License: revised BSD, see LICENSE
* Docs: [godoc](https://godoc.org/github.com/fxkr/go-freetype-fontloader)
* Dependencies: `fc-match` from the [`fontconfig`](https://www.freedesktop.org/wiki/Software/fontconfig/) package.


# Usage

Directly:

```go
font, err := Load("sans")
```

With cache:

```go
cache := NewFontCache()
font, err := cache.Load("sans")
```

Absolute paths are also supported. This provides a fallback for users who don't have fontconfig installed:

```go
font, err := Load("/usr/share/fonts/dejavu/DejaVuSans.ttf")
```
