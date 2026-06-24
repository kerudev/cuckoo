package utils

import (
	"strconv"
	"unicode"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// SortAlphabetically sorts alphabetically, including numbers.
// It is meant to be used inside functions like `sort.Sort`.
//
// Regular sort: "1", "10", "2" (see https://stackoverflow.com/a/35087122).
// This sort: 	 "1", "2", "10"
func SortAlphabetically(a, b string) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		iChar, bChar := a[i], b[j]

		// If both characters are a digit, sort by number
		if unicode.IsDigit(rune(iChar)) && unicode.IsDigit(rune(bChar)) {
			iStart, iEnd := extractNumber(a, i)
			jStart, jEnd := extractNumber(b, j)

			// Compare numbers as integers
			if iStart != jStart {
				return iStart < jStart
			}
			i, j = iEnd, jEnd
			continue
		}

		// Regular character comparison
		if iChar != bChar {
			return iChar < bChar
		}
		i++
		j++
	}

	// Shorter strings come first
	return len(a) < len(b)
}

func extractNumber(s string, start int) (int, int) {
	end := start
	for end < len(s) && unicode.IsDigit(rune(s[end])) {
		end++
	}
	n, _ := strconv.Atoi(s[start:end])
	return n, end
}

func CalcBucket(value int, segment int) int {
	if value == 0 {
		return 0
	}
	return (value/segment)*segment + (segment - 1)
}

func CountDuplicates[T comparable](arr []T) map[T]int {
	res := map[T]int{}

	for _, item := range arr {
		_, ok := res[item]
		if !ok {
			res[item] = 1
		} else {
			res[item]++
		}
	}

	return res
}

func All[T comparable](arr []T, pred func(T) bool) bool {
	for _, el := range arr {
		if pred(el) {
			continue
		}
		return false
	}
	return true
}

// MinF32 returns the smaller of x or y.
func MinF32(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
}

// f <  0.5  -> darken color
// f == 0.5 -> same color
// f >  0.5  -> brighten color
func lerpRGB(r uint8, g uint8, b uint8, a uint8, f float32) (uint8, uint8, uint8, uint8) {
	f = max(0.0, min(1.0, f))

	rf := float32(r)
	gf := float32(g)
	bf := float32(b)
	af := float32(a)

	r2 := uint8(0)
	g2 := uint8(0)
	b2 := uint8(0)
	a2 := uint8(0)

	if f < 0.5 {
		// Darken color
		factor := f / 0.5
		r2 = uint8(rf * factor)
		g2 = uint8(gf * factor)
		b2 = uint8(bf * factor)
		a2 = uint8(af * factor)
	} else {
		// Brighten color
		factor := (f - 0.5) / 0.5
		r2 = uint8(rf + (255-rf)*factor)
		g2 = uint8(gf + (255-gf)*factor)
		b2 = uint8(bf + (255-bf)*factor)
		a2 = uint8(af + (255-af)*factor)
	}

	return r2, g2, b2, a2
}

func LerpColor(color rl.Color, f float32) rl.Color {
	return rl.NewColor(lerpRGB(color.R, color.G, color.B, color.A, f))
}

func LerpHex(hex rg.PropertyValue, f float32) rg.PropertyValue {
	return rg.NewColorPropertyValue(LerpColor(hex.AsColor(), f))
}

func LerpColorToHex(color rl.Color, f float32) rg.PropertyValue {
	return rg.NewColorPropertyValue(LerpColor(color, f))
}
