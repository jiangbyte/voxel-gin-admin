// Package auth 验证码 PNG 渲染。
//
// Author: Charlie
package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
)

const (
	captchaWidth      = 140
	captchaHeight     = 44
	noiseLines        = 6
	captchaAlphabet   = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	captchaKeyPrefix  = "captcha:"
	passwordKeyPrefix = "password:crypto:"
)

var captchaGlyphs = map[rune][]string{
	'2': {"11110", "00001", "00001", "11110", "10000", "10000", "11111"},
	'3': {"11110", "00001", "00001", "01110", "00001", "00001", "11110"},
	'4': {"10001", "10001", "10001", "11111", "00001", "00001", "00001"},
	'5': {"11111", "10000", "10000", "11110", "00001", "00001", "11110"},
	'6': {"01111", "10000", "10000", "11110", "10001", "10001", "01110"},
	'7': {"11111", "00001", "00010", "00100", "01000", "01000", "01000"},
	'8': {"01110", "10001", "10001", "01110", "10001", "10001", "01110"},
	'9': {"01110", "10001", "10001", "01111", "00001", "00001", "11110"},
	'A': {"01110", "10001", "10001", "11111", "10001", "10001", "10001"},
	'B': {"11110", "10001", "10001", "11110", "10001", "10001", "11110"},
	'C': {"01111", "10000", "10000", "10000", "10000", "10000", "01111"},
	'D': {"11110", "10001", "10001", "10001", "10001", "10001", "11110"},
	'E': {"11111", "10000", "10000", "11110", "10000", "10000", "11111"},
	'F': {"11111", "10000", "10000", "11110", "10000", "10000", "10000"},
	'G': {"01111", "10000", "10000", "10011", "10001", "10001", "01111"},
	'H': {"10001", "10001", "10001", "11111", "10001", "10001", "10001"},
	'J': {"00111", "00010", "00010", "00010", "10010", "10010", "01100"},
	'K': {"10001", "10010", "10100", "11000", "10100", "10010", "10001"},
	'L': {"10000", "10000", "10000", "10000", "10000", "10000", "11111"},
	'M': {"10001", "11011", "10101", "10101", "10001", "10001", "10001"},
	'N': {"10001", "11001", "10101", "10011", "10001", "10001", "10001"},
	'P': {"11110", "10001", "10001", "11110", "10000", "10000", "10000"},
	'Q': {"01110", "10001", "10001", "10001", "10101", "10010", "01101"},
	'R': {"11110", "10001", "10001", "11110", "10100", "10010", "10001"},
	'S': {"01111", "10000", "10000", "01110", "00001", "00001", "11110"},
	'T': {"11111", "00100", "00100", "00100", "00100", "00100", "00100"},
	'U': {"10001", "10001", "10001", "10001", "10001", "10001", "01110"},
	'V': {"10001", "10001", "10001", "10001", "01010", "01010", "00100"},
	'W': {"10001", "10001", "10001", "10101", "10101", "11011", "10001"},
	'X': {"10001", "01010", "00100", "00100", "00100", "01010", "10001"},
	'Y': {"10001", "01010", "00100", "00100", "00100", "00100", "00100"},
	'Z': {"11111", "00001", "00010", "00100", "01000", "10000", "11111"},
}

var (
	captchaBgColor   = color.RGBA{R: 248, G: 250, B: 252, A: 255}
	captchaLineColor = color.RGBA{R: 148, G: 163, B: 184, A: 255}
	captchaTextColor = color.RGBA{R: 15, G: 23, B: 42, A: 255}
)

// captchaPNGBase64 渲染点阵 PNG 验证码并返回 Base64。
func captchaPNGBase64(value string) string {
	img := image.NewRGBA(image.Rect(0, 0, captchaWidth, captchaHeight))
	for y := 0; y < captchaHeight; y++ {
		for x := 0; x < captchaWidth; x++ {
			img.Set(x, y, captchaBgColor)
		}
	}
	for i := 0; i < noiseLines; i++ {
		drawCaptchaLine(img, randBelow(captchaWidth), randBelow(captchaHeight),
			randBelow(captchaWidth), randBelow(captchaHeight), captchaLineColor)
	}
	for i, ch := range value {
		drawCaptchaGlyph(img, ch,
			18+i*28+randBelow(3),
			8+randBelow(4),
			4,
			captchaTextColor)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func drawCaptchaGlyph(img *image.RGBA, ch rune, x, y, scale int, c color.RGBA) {
	glyph, ok := captchaGlyphs[ch]
	if !ok {
		return
	}
	for rowIndex, row := range glyph {
		for columnIndex, cell := range row {
			if cell != '1' {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					setCaptchaPixel(img, x+columnIndex*scale+dx, y+rowIndex*scale+dy, c)
				}
			}
		}
	}
}

func drawCaptchaLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := absInt(x2 - x1)
	dy := -absInt(y2 - y1)
	sx := 1
	if x1 >= x2 {
		sx = -1
	}
	sy := 1
	if y1 >= y2 {
		sy = -1
	}
	err := dx + dy
	for {
		setCaptchaPixel(img, x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

func setCaptchaPixel(img *image.RGBA, x, y int, c color.RGBA) {
	if x < 0 || x >= captchaWidth || y < 0 || y >= captchaHeight {
		return
	}
	img.Set(x, y, c)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func randBelow(n int) int {
	if n <= 0 {
		return 0
	}
	b := make([]byte, 1)
	_, _ = rand.Read(b)
	return int(b[0]) % n
}
