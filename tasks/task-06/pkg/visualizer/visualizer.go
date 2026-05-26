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

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-06/pkg/solver"
)

const (
	width  = 800
	height = 550
	margin = 60
)

var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Фон
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: Оси
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: Синий — полином
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 3: Оранжевый — узлы
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 4: Текст
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 5: Зелёный — результат
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 6: Сетка
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 7: Пурпурный
	color.RGBA{0x2a, 0xa1, 0x98, 0xff}, // 8: Циановый
	color.RGBA{0x6c, 0x71, 0xc4, 0xff}, // 9: Базисный полином 1
}

// Цвета для базисных полиномов
var basisColors = []uint8{2, 3, 7, 8}

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

// GenerateLagrangeGIF создаёт анимацию интерполяционного многочлена Лагранжа
func GenerateLagrangeGIF(filePath string, points []solver.DataPoint, xTarget float64, history []solver.IterationRecord) error {
	anim := gif.GIF{}

	// Диапазон отображения
	xMin, xMax := -0.5, 3.5
	yMin, yMax := 9.0, 16.0

	plotW := width - 2*margin
	plotH := height - margin - 80

	toX := func(x float64) int {
		return margin + int((x-xMin)/(xMax-xMin)*float64(plotW))
	}
	toY := func(y float64) int {
		return margin + 20 + int((1-(y-yMin)/(yMax-yMin))*float64(plotH))
	}

	nFrames := len(history)

	// Дополнительный финальный кадр с подписями
	for frame := 0; frame < nFrames; frame++ {
		rec := history[frame]
		subset := points[:rec.Used]

		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)

		// Заголовок
		drawText(img, 15, 18, fmt.Sprintf("Task 06 | Lagrange Interpolation | Variant 8  — %d/%d узл(ов)", rec.Used, len(points)), palette[4])

		// Сетка
		for gx := xMin; gx <= xMax; gx += 0.5 {
			xx := toX(gx)
			for yy := margin + 20; yy < margin+20+plotH; yy += 2 {
				if xx >= 0 && xx < width && yy >= 0 && yy < height {
					img.SetColorIndex(xx, yy, 6)
				}
			}
		}
		for gy := yMin; gy <= yMax; gy += 1 {
			yy := toY(gy)
			for xx := margin; xx < margin+plotW; xx += 2 {
				if xx >= 0 && xx < width && yy >= 0 && yy < height {
					img.SetColorIndex(xx, yy, 6)
				}
			}
		}

		// Оси
		drawLine(img, margin, toY(0), margin+plotW, toY(0), 1)
		drawLine(img, toX(0), margin+20, toX(0), margin+20+plotH, 1)

		// Подписи осей
		for xi := 0; xi <= 3; xi++ {
			xx := toX(float64(xi))
			drawLine(img, xx, toY(0)-3, xx, toY(0)+3, 1)
			drawText(img, xx-3, toY(0)+14, fmt.Sprintf("%d", xi), palette[4])
		}
		for yi := int(yMin); yi <= int(yMax); yi += 2 {
			yy := toY(float64(yi))
			drawLine(img, toX(0)-3, yy, toX(0)+3, yy, 1)
			drawText(img, margin-28, yy+4, fmt.Sprintf("%d", yi), palette[4])
		}

		// Базисные полиномы L_i(x) для текущего подмножества
		for i := range subset {
			ci := basisColors[i%len(basisColors)]
			prevX, prevY := 0, 0
			first := true
			for px := 0; px <= plotW; px++ {
				x := xMin + float64(px)/float64(plotW)*(xMax-xMin)
				lx := solver.LagrangeBasis(subset, i, x)
				sy := toY(lx)
				if !first && sy >= margin && sy < margin+plotH+20 {
					drawLine(img, margin+px-1, prevY, margin+px, sy, ci)
				}
				prevX, prevY = margin+px, sy
				_ = prevX
				first = false
			}
			// Легенда
			drawText(img, width-180, margin+30+i*15,
				fmt.Sprintf("L%d(x), L%d(x*)=%.3f", i, i, rec.Basis[i]),
				palette[ci])
		}

		// Полином Лагранжа L(x)
		prevX2, prevY2 := 0, 0
		first2 := true
		for px := 0; px <= plotW; px++ {
			x := xMin + float64(px)/float64(plotW)*(xMax-xMin)
			lx := solver.LagrangeInterpolate(subset, x)
			sy := toY(lx)
			if !first2 && sy >= margin && sy < height-50 {
				drawLine(img, margin+px-1, prevY2, margin+px, sy, 2)
			}
			prevX2, prevY2 = margin+px, sy
			_ = prevX2
			first2 = false
		}

		// Исходные узлы (все)
		for _, p := range points {
			drawPoint(img, toX(p.X), toY(p.Y), 5, 3)
			drawText(img, toX(p.X)+7, toY(p.Y)-5, fmt.Sprintf("(%.0f,%.0f)", p.X, p.Y), palette[4])
		}

		// Активные узлы подмножества — подсветка
		for _, p := range subset {
			drawPoint(img, toX(p.X), toY(p.Y), 3, 5)
		}

		// Вертикаль x*
		xt := toX(xTarget)
		drawLine(img, xt, margin+20, xt, margin+20+plotH, 7)
		drawText(img, xt+4, margin+25, fmt.Sprintf("x*=%.1f", xTarget), palette[7])

		// Точка результата
		drawPoint(img, xt, toY(rec.Value), 6, 5)
		drawLine(img, xt-6, toY(rec.Value), xt+6, toY(rec.Value), 5)

		// Строка результата
		drawText(img, margin, height-20,
			fmt.Sprintf("L(x*=%.1f) = %.6f   (узлов: %d)", xTarget, rec.Value, rec.Used),
			palette[5])

		delay := 80
		if frame == nFrames-1 {
			delay = 400
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Стоп-кадры
	for i := 0; i < 5; i++ {
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
