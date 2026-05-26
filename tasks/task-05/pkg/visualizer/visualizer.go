package visualizer

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-05/pkg/solver"
)

const (
	width  = 800
	height = 550
	margin = 60
)

var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Фон
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: Оси
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: Синий — λ(k)
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 3: Оранжевый — ошибка
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 4: Основной текст
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 5: Зелёный — сошлись
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 6: Вспомогательная сетка
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 7: Пурпурный — вектор
	color.RGBA{0x2a, 0xa1, 0x98, 0xff}, // 8: Циановый
	color.RGBA{0x6c, 0x71, 0xc4, 0xff}, // 9: Фиолетовый
}

func drawLine(img *image.Paletted, x0, y0, x1, y1 int, ci uint8) {
	dx, dy := int(math.Abs(float64(x1-x0))), int(math.Abs(float64(y1-y0)))
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	e := dx - dy
	for {
		if x0 >= 0 && x0 < width && y0 >= 0 && y0 < height {
			img.SetColorIndex(x0, y0, ci)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * e
		if e2 > -dy {
			e -= dy
			x0 += sx
		}
		if e2 < dx {
			e += dx
			y0 += sy
		}
	}
}

func drawPoint(img *image.Paletted, cx, cy, r int, ci uint8) {
	for x := cx - r; x <= cx+r; x++ {
		for y := cy - r; y <= cy+r; y++ {
			if x >= 0 && x < width && y >= 0 && y < height {
				img.SetColorIndex(x, y, ci)
			}
		}
	}
}

func drawText(img *image.Paletted, x, y int, text string, col color.Color) {
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	dr.DrawString(text)
}

// GenerateEigenGIF создаёт анимацию сходимости метода степенных итераций
func GenerateEigenGIF(filePath string, history []solver.IterationRecord, allEvals [3]float64) error {
	anim := gif.GIF{}

	// Диапазоны для графика λ(k)
	lambdaMin, lambdaMax := allEvals[2]-1, allEvals[0]+2
	if lambdaMin > -10 {
		lambdaMin = -10
	}
	nIter := len(history) - 1

	plotW := (width - 2*margin) / 2
	plotH := height - 2*margin - 40

	// left panel: λ convergence; right panel: log error
	for frame := 1; frame <= len(history); frame++ {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)

		// Заголовок
		drawText(img, 15, 20, fmt.Sprintf("Task 05 | Power Iteration | Variant 8  (α = -7)  — Step %d/%d", frame-1, nIter), palette[4])

		// --- Левая панель: график λ(k) ---
		lx0, ly0 := margin, margin+20
		drawText(img, lx0, ly0-5, "λ(k) — dominant eigenvalue convergence", palette[4])

		// Оси левой панели
		axX := lx0
		axY := ly0 + plotH
		drawLine(img, axX, ly0, axX, axY, 1)
		drawLine(img, axX, axY, axX+plotW, axY, 1)

		toX := func(k int) int {
			if nIter == 0 {
				return axX
			}
			return axX + int(float64(k)/float64(nIter)*float64(plotW))
		}
		toY := func(val float64) int {
			t := (val - lambdaMin) / (lambdaMax - lambdaMin)
			return axY - int(t*float64(plotH))
		}

		// Горизонтальные линии целевых СЗ
		for ei, ev := range allEvals {
			ci := uint8([]int{5, 2, 9}[ei])
			yy := toY(ev)
			if yy >= ly0 && yy <= axY {
				for xx := axX; xx <= axX+plotW; xx += 4 {
					if xx >= 0 && xx < width && yy >= 0 && yy < height {
						img.SetColorIndex(xx, yy, 6)
					}
				}
				drawText(img, axX+plotW+3, yy+4, fmt.Sprintf("λ%d=%.2f", ei+1, ev), palette[ci])
			}
		}

		// Кривая λ(k)
		for i := 1; i < frame; i++ {
			x0, y0 := toX(history[i-1].Iter), toY(history[i-1].Eigenvalue)
			x1, y1 := toX(history[i].Iter), toY(history[i].Eigenvalue)
			drawLine(img, x0, y0, x1, y1, 2)
		}
		if frame > 0 {
			cur := history[frame-1]
			drawPoint(img, toX(cur.Iter), toY(cur.Eigenvalue), 4, 5)
			drawText(img, lx0, axY+18, fmt.Sprintf("λ = %.6f", cur.Eigenvalue), palette[5])
		}

		// Тики по оси X
		for k := 0; k <= nIter; k += max(1, nIter/5) {
			xx := toX(k)
			drawLine(img, xx, axY, xx, axY+3, 1)
			drawText(img, xx-4, axY+14, fmt.Sprintf("%d", k), palette[4])
		}

		// --- Правая панель: log₁₀(||r||) ---
		rx0 := lx0 + plotW + margin
		drawText(img, rx0, ly0-5, "log₁₀(residual) — convergence rate", palette[4])

		eaxX := rx0
		eaxY := axY
		drawLine(img, eaxX, ly0, eaxX, eaxY, 1)
		drawLine(img, eaxX, eaxY, eaxX+plotW, eaxY, 1)

		// Собрать все конечные логарифмы
		logMin, logMax := -14.0, 2.0

		toEX := func(k int) int {
			if nIter == 0 {
				return eaxX
			}
			return eaxX + int(float64(k)/float64(nIter)*float64(plotW))
		}
		toEY := func(val float64) int {
			t := (val - logMin) / (logMax - logMin)
			return eaxY - int(t*float64(plotH))
		}

		// Горизонтальные уровни eps
		epsList := []float64{0, -2, -4, -6, -8, -10, -12}
		for _, e := range epsList {
			yy := toEY(e)
			if yy >= ly0 && yy <= eaxY {
				for xx := eaxX; xx <= eaxX+plotW; xx += 4 {
					if xx >= 0 && xx < width && yy >= 0 && yy < height {
						img.SetColorIndex(xx, yy, 6)
					}
				}
				drawText(img, eaxX+plotW+3, yy+4, fmt.Sprintf("1e%d", int(e)), palette[6])
			}
		}

		// Кривая log ошибки
		for i := 1; i < frame; i++ {
			if history[i].Error <= 0 || history[i-1].Error <= 0 || math.IsInf(history[i-1].Error, 0) {
				continue
			}
			logE0 := math.Log10(history[i-1].Error)
			logE1 := math.Log10(history[i].Error)
			x0, y0 := toEX(history[i-1].Iter), toEY(logE0)
			x1, y1 := toEX(history[i].Iter), toEY(logE1)
			drawLine(img, x0, y0, x1, y1, 3)
		}
		if frame > 1 {
			cur := history[frame-1]
			if cur.Error > 0 {
				drawPoint(img, toEX(cur.Iter), toEY(math.Log10(cur.Error)), 4, 3)
			}
			drawText(img, rx0, axY+18, fmt.Sprintf("||r|| = %.2e", cur.Error), palette[3])
		}

		// Тики по оси X (правая)
		for k := 0; k <= nIter; k += max(1, nIter/5) {
			xx := toEX(k)
			drawLine(img, xx, eaxY, xx, eaxY+3, 1)
			drawText(img, xx-4, eaxY+14, fmt.Sprintf("%d", k), palette[4])
		}

		// --- Нижняя строка: вектор ---
		cur := history[frame-1]
		drawText(img, margin, height-15,
			fmt.Sprintf("Eigen-vector: [%.4f, %.4f, %.4f]", cur.Vector[0], cur.Vector[1], cur.Vector[2]),
			palette[7])

		delay := 8
		if frame == len(history) {
			delay = 300
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Дополнительные стоп-кадры
	for i := 0; i < 4; i++ {
		anim.Image = append(anim.Image, anim.Image[len(anim.Image)-1])
		anim.Delay = append(anim.Delay, 100)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, &anim)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
