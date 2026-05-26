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

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-04/pkg/solver"
)

const (
	width  = 800
	height = 550
	margin = 55
)

// Палитра Solarized Light из лабы 2/3
var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Фон
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: Оси координат
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: График 1-го уравнения (Синий)
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 3: График 2-го уравнения (Оранжевый)
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 4: Траектория итераций (Пурпурный)
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 5: Основной текст
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 6: Найденный корень (Зеленый)
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 7: Вспомогательная сетка
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

func GenerateConvergenceGIF(filePath string, title string, history []solver.IterationRecord) error {
	anim := gif.GIF{}

	// Границы отображения фазовой плоскости (окрестность корней)
	minX, maxX := -0.2, 1.2
	minY, maxY := -1.2, 0.2

	toX := func(x float64) int {
		return margin + int((x-minX)/(maxX-minX)*float64(width-2*margin))
	}
	toY := func(y float64) int {
		return height - margin - int((y-minY)/(maxY-minY)*float64(height-2*margin))
	}

	for frame := 1; frame <= len(history); frame++ {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)

		// Рендеринг заголовков
		dr := &font.Drawer{Dst: img, Src: image.NewUniform(palette[5]), Face: basicfont.Face7x13}
		dr.Dot = fixed.Point26_6{X: fixed.I(15), Y: fixed.I(25)}
		dr.DrawString(fmt.Sprintf("%s (Итерация %d/%d)", title, frame-1, len(history)-1))

		// Легенда графиков
		dr.Dot = fixed.Point26_6{X: fixed.I(width - 240), Y: fixed.I(25)}
		dr.DrawString("Синий:  2x - cos(y+1) = 0")
		dr.Dot = fixed.Point26_6{X: fixed.I(width - 240), Y: fixed.I(40)}
		dr.DrawString("Оранж:  y + sin(x) = -0.4")

		// Оси декартовой системы
		drawLine(img, margin, toY(0), width-margin, toY(0), 1)
		drawLine(img, toX(0), margin, toX(0), height-margin, 1)

		// Отрисовка нелинейных кривых системы
		for px := margin; px < width-margin; px++ {
			fx := minX + float64(px-margin)/float64(width-2*margin)*(maxX-minX)

			// 1) 2x - cos(y+1) = 0  => выразим x через y (поэтому прогоним по y плоскости)
			fy1 := minX + float64(px-margin)/float64(width-2*margin)*(maxX-minX)
			fx1 := 0.5 * math.Cos(fy1+1)
			pX1, pY1 := toX(fx1), toY(fy1)
			if pX1 >= margin && pX1 < width-margin && pY1 >= margin && pY1 < height-margin {
				img.SetColorIndex(pX1, pY1, 2)
			}

			// 2) y = -0.4 - sin(x)
			fy2 := -0.4 - math.Sin(fx)
			pY2 := toY(fy2)
			if pY2 >= margin && pY2 < height-margin {
				img.SetColorIndex(px, pY2, 3)
			}
		}

		// Отрисовка пути сходимости (траектория)
		for i := 1; i < frame; i++ {
			x0, y0 := toX(history[i-1].X), toY(history[i-1].Y)
			x1, y1 := toX(history[i].X), toY(history[i].Y)
			drawLine(img, x0, y0, x1, y1, 4) // Шаг
			drawPoint(img, x0, y0, 2, 4)
		}

		// Подсветка текущего положения алгоритма
		curr := history[frame-1]
		drawPoint(img, toX(curr.X), toY(curr.Y), 4, 6)

		// Строка статуса внизу
		dr.Dot = fixed.Point26_6{X: fixed.I(margin), Y: fixed.I(height - 15)}
		errVal := curr.Error
		if math.IsInf(errVal, 1) {
			errVal = 0.0
		}
		dr.DrawString(fmt.Sprintf("Текущая точка: x = %.5f,  y = %.5f,  ||Δx|| = %.2e", curr.X, curr.Y, errVal))

		delay := 25
		if frame == len(history) {
			delay = 250
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Фиксация последнего кадра
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
