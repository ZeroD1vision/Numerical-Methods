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

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-08/pkg/solver"
)

const (
	width  = 800
	height = 550
	margin = 60
)

var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Фон
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: Оси
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: Синий — Метод Эйлера
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 3: Оранжевый — Модифицированный Эйлера
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 4: Основной текст
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 5: Зеленый — "Точное" решение (RK4)
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 6: Вспомогательная сетка
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

// GenerateODEGIF создает анимацию решения ОДУ с добавлением точного решения
func GenerateODEGIF(filePath string, eulerPts, modEulerPts, exactPts []solver.Point) error {
	anim := gif.GIF{}

	// Ищем экстремумы среди всех трех графиков для масштабирования
	minY, maxY := eulerPts[0].Y, eulerPts[0].Y
	allPoints := [][]solver.Point{eulerPts, modEulerPts, exactPts}
	for _, pts := range allPoints {
		for _, p := range pts {
			if p.Y < minY {
				minY = p.Y
			}
			if p.Y > maxY {
				maxY = p.Y
			}
		}
	}

	minT, maxT := eulerPts[0].T, eulerPts[len(eulerPts)-1].T

	yPadding := (maxY - minY) * 0.1
	if yPadding == 0 {
		yPadding = 1.0
	}
	minY -= yPadding
	maxY += yPadding

	plotW := width - 2*margin
	plotH := height - 2*margin

	toX := func(t float64) int {
		return margin + int((t-minT)/(maxT-minT)*float64(plotW))
	}
	toY := func(y float64) int {
		return margin + plotH - int((y-minY)/(maxY-minY)*float64(plotH))
	}

	nIter := len(eulerPts)

	for frame := 1; frame <= nIter; frame++ {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)

		drawText(img, margin, margin-20, "Task 08 | ODE Cauchy Problem | Euler vs Mod. Euler vs Exact", palette[4])
		drawText(img, margin, height-margin+30, "t", palette[4])
		drawText(img, margin-20, margin-10, "y", palette[4])

		// Обновленная легенда
		drawText(img, width-margin-150, margin-35, "— Exact (RK4)", palette[5])
		drawText(img, width-margin-150, margin-20, "— Euler", palette[2])
		drawText(img, width-margin-150, margin-5, "— Mod. Euler", palette[3])

		// Оси
		drawLine(img, margin, margin, margin, margin+plotH, 1)
		drawLine(img, margin, margin+plotH, margin+plotW, margin+plotH, 1)

		// Сетка
		for i := 0; i <= 10; i++ {
			t := minT + float64(i)*(maxT-minT)/10.0
			xx := toX(t)
			for yy := margin; yy <= margin+plotH; yy += 5 {
				img.SetColorIndex(xx, yy, 6)
			}
			drawText(img, xx-10, margin+plotH+15, fmt.Sprintf("%.1f", t), palette[4])
		}

		for i := 0; i <= 5; i++ {
			y := minY + float64(i)*(maxY-minY)/5.0
			yy := toY(y)
			for xx := margin; xx <= margin+plotW; xx += 5 {
				img.SetColorIndex(xx, yy, 6)
			}
			drawText(img, margin-45, yy+5, fmt.Sprintf("%.2f", y), palette[4])
		}

		// 1. Рисуем НАСТОЯЩУЮ (точную) функцию на текущем промежутке времени (зеленый)
		// Она рисуется пиксель в пиксель для идеальной плавности
		currentMaxT := eulerPts[frame-1].T
		var prevX, prevY int
		first := true
		for xx := margin; xx <= toX(currentMaxT); xx++ {
			// Переводим экранный X назад в значение t
			t := minT + float64(xx-margin)/float64(plotW)*(maxT-minT)
			yExact := solver.GetExactY(exactPts, t)
			yy := toY(yExact)

			if !first {
				drawLine(img, prevX, prevY, xx, yy, 5)
			}
			first = false
			prevX, prevY = xx, yy
		}

		// 2. Рисуем метод Эйлера (синий)
		for i := 1; i < frame; i++ {
			x0, y0 := toX(eulerPts[i-1].T), toY(eulerPts[i-1].Y)
			x1, y1 := toX(eulerPts[i].T), toY(eulerPts[i].Y)
			drawLine(img, x0, y0, x1, y1, 2)
			drawPoint(img, x1, y1, 3, 2)
		}
		drawPoint(img, toX(eulerPts[0].T), toY(eulerPts[0].Y), 3, 2)

		// 3. Рисуем модифицированный Эйлер (оранжевый)
		for i := 1; i < frame; i++ {
			x0, y0 := toX(modEulerPts[i-1].T), toY(modEulerPts[i-1].Y)
			x1, y1 := toX(modEulerPts[i].T), toY(modEulerPts[i].Y)
			drawLine(img, x0, y0, x1, y1, 3)
			drawPoint(img, x1, y1, 3, 3)
		}
		drawPoint(img, toX(modEulerPts[0].T), toY(modEulerPts[0].Y), 3, 3)

		delay := 15
		if frame == nIter {
			delay = 300
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Стоп-кадры в конце
	for i := 0; i < 3; i++ {
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
