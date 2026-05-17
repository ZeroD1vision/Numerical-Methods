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
)

const (
	width  = 700
	height = 450
)

// Расширенная палитра для детальной графики
var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: Светлый фон (Solarized Light)
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: График функции (темно-серый)
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 2: Оси координат и сетка
	color.RGBA{0x2a, 0xa1, 0x98, 0xff}, // 3: Метод Ньютона / Касательные (Бирюзовый)
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 4: Метод Хорд / Секущие (Пурпурный)
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 5: Текст заголовков и подписей (Оранжевый)
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 6: Дополнительный цвет для точек
}

// Вспомогательная структура для масштабирования осей
type Plot struct {
	xMin, xMax float64
	yMin, yMax float64
}

func (p *Plot) toPixelX(x float64) int {
	return int((x - p.xMin) / (p.xMax - p.xMin) * float64(width))
}

func (p *Plot) toPixelY(y float64) int {
	return height - int((y-p.yMin)/(p.yMax-p.yMin)*float64(height))
}

// Алгоритм Брезенхема для рисования линий
func drawLine(img *image.Paletted, x0, y0, x1, y1 int, colorIndex uint8) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < width && y0 >= 0 && y0 < height {
			img.SetColorIndex(x0, y0, colorIndex)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// Отрисовка текста стандартным шрифтом
func drawText(img *image.Paletted, x, y int, text string, colorIndex uint8) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(img.Palette[colorIndex]),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

// Отрисовка точки (заполненного квадрата)
func drawPoint(img *image.Paletted, cx, cy, size int, colorIndex uint8) {
	for x := cx - size; x <= cx+size; x++ {
		for y := cy - size; y <= cy+size; y++ {
			if x >= 0 && x < width && y >= 0 && y < height {
				img.SetColorIndex(x, y, colorIndex)
			}
		}
	}
}

// Базовая разметка осей и сетки координат
func drawBaseLayout(img *image.Paletted, p *Plot, title string) {
	yZero := p.toPixelY(0.0)
	xZero := p.toPixelX(0.0)

	// Ось X и Ось Y
	drawLine(img, 0, yZero, width, yZero, 2)
	drawLine(img, xZero, 0, xZero, height, 2)

	// Стрелки осей и подписи координат
	drawText(img, width-15, yZero-5, "X", 5)
	drawText(img, xZero+8, 15, "F(x)", 5)
	drawText(img, 15, 25, title, 5)

	// Засечки границ интервала локализации
	drawText(img, p.toPixelX(0.0), yZero+15, "0.0", 2)
	drawText(img, p.toPixelX(0.5), yZero+15, "0.5", 2)
	drawLine(img, p.toPixelX(0.5), yZero-5, p.toPixelX(0.5), yZero+5, 2)
	drawLine(img, p.toPixelX(0.0), yZero-5, p.toPixelX(0.0), yZero+5, 2)
}

// 1. Метод Ньютона: Отрисовка графика функции и касательных прямых
func GenerateNewtonGIF(filePath string, f, d1 func(float64) float64, a, b float64, history []float64) error {
	anim := gif.GIF{}
	p := Plot{xMin: a - 0.1, xMax: b + 0.1, yMin: -1.2, yMax: 1.2}

	for idx, xCurr := range history {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)
		drawBaseLayout(img, &p, fmt.Sprintf("Newton Method: Iteration %d (x = %.4f)", idx, xCurr))

		// График функции
		for px := 1; px < width; px++ {
			x1 := p.xMin + (float64(px-1)/float64(width))*(p.xMax-p.xMin)
			x2 := p.xMin + (float64(px)/float64(width))*(p.xMax-p.xMin)
			drawLine(img, px-1, p.toPixelY(f(x1)), px, p.toPixelY(f(x2)), 1)
		}

		// Построение касательной: y - f(xCurr) = f'(xCurr) * (x - xCurr)
		// Пересечение с осью X произойдет в точке следующего шага (или конечного корня)
		yCurr := f(xCurr)
		pX, pY := p.toPixelX(xCurr), p.toPixelY(yCurr)
		drawPoint(img, pX, pY, 4, 3) // Текущая точка на кривой

		// Отрисовка касательной линии через две точки по оси X
		tangentYStart := yCurr + d1(xCurr)*(p.xMin-xCurr)
		tangentYEnd := yCurr + d1(xCurr)*(p.xMax-xCurr)
		drawLine(img, 0, p.toPixelY(tangentYStart), width, p.toPixelY(tangentYEnd), 3)

		// Вертикальный пунктир от оси до кривой
		drawLine(img, pX, p.toPixelY(0.0), pX, pY, 2)

		anim.Delay = append(anim.Delay, 120)
		anim.Image = append(anim.Image, img)
	}
	return saveGIF(filePath, &anim)
}

// 2. Метод Хорд: Отрисовка графика функции и секущих прямых
func GenerateSecantGIF(filePath string, f func(float64) float64, a, b float64, history []float64) error {
	anim := gif.GIF{}
	p := Plot{xMin: a - 0.1, xMax: b + 0.1, yMin: -1.2, yMax: 1.2}

	// По условию Фурье для нашего варианта, неподвижная точка — это b (0.5)
	xFixed := b
	yFixed := f(xFixed)

	for idx, xCurr := range history {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)
		drawBaseLayout(img, &p, fmt.Sprintf("Secant Method: Iteration %d (x = %.4f)", idx, xCurr))

		// График функции
		for px := 1; px < width; px++ {
			x1 := p.xMin + (float64(px-1)/float64(width))*(p.xMax-p.xMin)
			x2 := p.xMin + (float64(px)/float64(width))*(p.xMax-p.xMin)
			drawLine(img, px-1, p.toPixelY(f(x1)), px, p.toPixelY(f(x2)), 1)
		}

		// Точки хорды
		yCurr := f(xCurr)
		pX, pY := p.toPixelX(xCurr), p.toPixelY(yCurr)
		pFX, pFY := p.toPixelX(xFixed), p.toPixelY(yFixed)

		drawPoint(img, pX, pY, 4, 4)   // Подвижная точка
		drawPoint(img, pFX, pFY, 4, 6) // Неподвижная точка опорная

		// Хорда (прямая линия через подвижную и неподвижную точки)
		drawLine(img, pX, pY, pFX, pFY, 4)
		// Вертикальная линия проекции
		drawLine(img, pX, p.toPixelY(0.0), pX, pY, 2)

		anim.Delay = append(anim.Delay, 120)
		anim.Image = append(anim.Image, img)
	}
	return saveGIF(filePath, &anim)
}

// 3. Совмещенный график: Одновременное сравнение траекторий двух маркеров
func GenerateCombinedGIF(filePath string, f func(float64) float64, a, b float64, newtonHistory, secantHistory []float64) error {
	anim := gif.GIF{}
	p := Plot{xMin: a - 0.1, xMax: b + 0.1, yMin: -1.2, yMax: 1.2}

	maxFrames := len(newtonHistory)
	if len(secantHistory) > maxFrames {
		maxFrames = len(secantHistory)
	}

	for frameIdx := 0; frameIdx < maxFrames; frameIdx++ {
		rect := image.Rect(0, 0, width, height)
		img := image.NewPaletted(rect, palette)
		drawBaseLayout(img, &p, fmt.Sprintf("Comparison: Iteration %d (Green: Newton, Purple: Secant)", frameIdx))

		// График функции
		for px := 1; px < width; px++ {
			x1 := p.xMin + (float64(px-1)/float64(width))*(p.xMax-p.xMin)
			x2 := p.xMin + (float64(px)/float64(width))*(p.xMax-p.xMin)
			drawLine(img, px-1, p.toPixelY(f(x1)), px, p.toPixelY(f(x2)), 1)
		}

		// Маркер Ньютона
		nIdx := frameIdx
		if nIdx >= len(newtonHistory) {
			nIdx = len(newtonHistory) - 1
		}
		nX := newtonHistory[nIdx]
		pNX, pNY := p.toPixelX(nX), p.toPixelY(f(nX))
		drawPoint(img, pNX, pNY, 5, 3)
		drawText(img, pNX-10, pNY-12, "N", 3)

		// Маркер Хорд
		sIdx := frameIdx
		if sIdx >= len(secantHistory) {
			sIdx = len(secantHistory) - 1
		}
		sX := secantHistory[sIdx]
		pSX, pSY := p.toPixelX(sX), p.toPixelY(f(sX))
		drawPoint(img, pSX, pSY, 5, 4)
		drawText(img, pSX-10, pSY-12, "S", 4)

		anim.Delay = append(anim.Delay, 120)
		anim.Image = append(anim.Image, img)
	}
	return saveGIF(filePath, &anim)
}

// Запись файла и зацикливание кадра в конце анимации
func saveGIF(filePath string, anim *gif.GIF) error {
	for i := 0; i < 4; i++ {
		anim.Delay = append(anim.Delay, 150)
		anim.Image = append(anim.Image, anim.Image[len(anim.Image)-1])
	}
	fOut, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fOut.Close()
	return gif.EncodeAll(fOut, anim)
}
