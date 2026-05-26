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

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-07/pkg/solver"
)

const (
	width  = 800
	height = 550
	margin = 60
)

var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Фон (Solarized Light)
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: Оси
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: Синий — Исходная функция f(x)
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 3: Оранжевый — Квадратичные параболы Симпсона
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 4: Текст
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 5: Зелёный — Точки разметки / узлы
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 6: Сетка
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 7: Пурпурный — Границы интегрирования
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

// GenerateSimpsonGIF создаёт визуализацию процесса интегрирования методом Симпсона
func GenerateSimpsonGIF(filePath string, a, b float64, history []solver.IntegrationRecord) error {
	anim := gif.GIF{}

	// Границы графической области для красивого отображения f(x) на [1, 2]
	xMin, xMax := 0.7, 2.3
	yMin, yMax := 0.0, 1.2

	plotW := width - 2*margin
	plotH := height - margin - 80

	toX := func(x float64) int {
		return margin + int((x-xMin)/(xMax-xMin)*float64(plotW))
	}
	toY := func(y float64) int {
		return margin + 20 + int((1-(y-yMin)/(yMax-yMin))*float64(plotH))
	}

	for frame, rec := range history {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)

		// Заголовок кадра
		drawText(img, 15, 18, fmt.Sprintf("Task 07 | Simpson's Integration Rule | N = %d", rec.N), palette[4])

		// Сетка координат
		for gx := 0.8; gx <= 2.2; gx += 0.2 {
			xx := toX(gx)
			for yy := margin + 20; yy < margin+20+plotH; yy += 4 {
				if xx >= 0 && xx < width && yy >= 0 && yy < height {
					img.SetColorIndex(xx, yy, 6)
				}
			}
		}
		for gy := 0.2; gy <= 1.0; gy += 0.2 {
			yy := toY(gy)
			for xx := margin; xx < margin+plotW; xx += 4 {
				if xx >= 0 && xx < width && yy >= 0 && yy < height {
					img.SetColorIndex(xx, yy, 6)
				}
			}
		}

		// Рисуем оси OX и OY
		drawLine(img, margin, toY(0), margin+plotW, toY(0), 1)
		drawLine(img, toX(1.0), margin+20, toX(1.0), margin+20+plotH, 1)

		// Метки осей
		for xi := 1.0; xi <= 2.0; xi += 0.5 {
			xx := toX(xi)
			drawLine(img, xx, toY(0)-3, xx, toY(0)+3, 1)
			drawText(img, xx-8, toY(0)+14, fmt.Sprintf("%.1f", xi), palette[4])
		}
		for yi := 0.0; yi <= 1.0; yi += 0.2 {
			yy := toY(yi)
			drawLine(img, toX(1.0)-3, yy, toX(1.0)+3, yy, 1)
			drawText(img, margin-32, yy+4, fmt.Sprintf("%.1f", yi), palette[4])
		}

		// 1. Отрисовка точной исходной функции f(x) синим цветом
		first := true
		var prevX, prevY int
		for px := 0; px <= plotW; px++ {
			x := xMin + float64(px)/float64(plotW)*(xMax-xMin)
			y := solver.Variant8Func(x)
			sx, sy := margin+px, toY(y)
			if !first && sy >= margin && sy < margin+plotH+20 {
				drawLine(img, prevX, prevY, sx, sy, 2)
			}
			prevX, prevY = sx, sy
			first = false
		}

		// 2. Построение составных парабол Симпсона для текущего N
		h := rec.H
		for i := 0; i < rec.N; i += 2 {
			x0 := a + float64(i)*h
			x1 := x0 + h
			x2 := x0 + 2*h

			y0 := solver.Variant8Func(x0)
			y1 := solver.Variant8Func(x1)
			y2 := solver.Variant8Func(x2)

			// Линии разбиения вниз до оси OX
			drawLine(img, toX(x0), toY(0), toX(x0), toY(y0), 6)
			if i == rec.N-2 {
				drawLine(img, toX(x2), toY(0), toX(x2), toY(y2), 6)
			}

			// Интерполяционная парабола Лагранжа по 3 точкам Симпсона
			pFirst := true
			var pPrevX, pPrevY int
			startX, endX := toX(x0), toX(x2)

			for px := startX; px <= endX; px++ {
				x := xMin + float64(px-margin)/float64(plotW)*(xMax-xMin)
				// Базисные полиномы
				l0 := ((x - x1) * (x - x2)) / ((x0 - x1) * (x0 - x2))
				l1 := ((x - x0) * (x - x2)) / ((x1 - x0) * (x1 - x2))
				l2 := ((x - x0) * (x - x1)) / ((x2 - x0) * (x2 - x1))
				yParabola := y0*l0 + y1*l1 + y2*l2

				sy := toY(yParabola)
				if !pFirst {
					drawLine(img, pPrevX, pPrevY, px, sy, 3)
				}
				pPrevX, pPrevY = px, sy
				pFirst = false
			}

			// Выделяем узловые точки интеграции зеленым цветом
			drawPoint(img, toX(x0), toY(y0), 3, 5)
			drawPoint(img, toX(x1), toY(y1), 2, 5)
			drawPoint(img, toX(x2), toY(y2), 3, 5)
		}

		// Вертикальные границы интегрирования [a, b]
		drawLine(img, toX(a), margin+20, toX(a), margin+20+plotH, 7)
		drawLine(img, toX(b), margin+20, toX(b), margin+20+plotH, 7)
		drawText(img, toX(a)+4, margin+35, "a=1.0", palette[7])
		drawText(img, toX(b)-42, margin+35, "b=2.0", palette[7])

		// Легенда и текущие численные результаты на кадре
		drawText(img, width-220, margin+30, fmt.Sprintf("Шаг h: %.6f", rec.H), palette[4])
		drawText(img, width-220, margin+50, fmt.Sprintf("Значение I: %.6f", rec.Value), palette[5])
		if frame > 0 {
			drawText(img, width-220, margin+70, fmt.Sprintf("Погр. Рунге: %.6f", rec.Error), palette[3])
		} else {
			drawText(img, width-220, margin+70, "Погр. Рунге: —", palette[4])
		}

		// Нижняя статус-строка
		drawText(img, margin, height-25,
			fmt.Sprintf("Итерация %d | Текущий интеграл = %.6f", rec.Step, rec.Value),
			palette[5])

		delay := 120
		if frame == len(history)-1 {
			delay = 400 // Задерживаем финальный точный кадр
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Зацикливание стоп-кадров в конце анимации
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
