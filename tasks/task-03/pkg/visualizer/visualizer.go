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

	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-03/pkg/solver"
)

const (
	width  = 750
	height = 480
)

// Палитра в стиле Solarized Light — единая со вторым заданием
var palette = []color.Color{
	color.RGBA{0xfd, 0xf6, 0xe3, 0xff}, // 0: фон
	color.RGBA{0x58, 0x6e, 0x75, 0xff}, // 1: оси / сетка
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // 2: x1 (синий)
	color.RGBA{0x2a, 0xa1, 0x98, 0xff}, // 3: x2 (бирюзовый)
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // 4: x3 (оранжевый)
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // 5: x4 (пурпурный)
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // 6: норма ошибки (зелёный)
	color.RGBA{0x07, 0x36, 0x42, 0xff}, // 7: текст / заголовки
	color.RGBA{0x93, 0xa1, 0xa1, 0xff}, // 8: вспомогательные линии
	color.RGBA{0xee, 0x8d, 0x2d, 0xff}, // 9: опорная eps-линия
}

// ── Примитивы рисования ──────────────────────────────────────────────────────

func drawLine(img *image.Paletted, x0, y0, x1, y1 int, ci uint8) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
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

func drawText(img *image.Paletted, x, y int, text string, ci uint8) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(img.Palette[ci]),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

// ── Макет: оси + легенда ─────────────────────────────────────────────────────

func newFrame(title string) *image.Paletted {
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, palette)
	drawText(img, 12, 18, title, 7)
	// легенда
	drawPoint(img, 12, 38, 4, 2)
	drawText(img, 22, 42, "x1", 2)
	drawPoint(img, 55, 38, 4, 3)
	drawText(img, 65, 42, "x2", 3)
	drawPoint(img, 98, 38, 4, 4)
	drawText(img, 108, 42, "x3", 4)
	drawPoint(img, 141, 38, 4, 5)
	drawText(img, 151, 42, "x4", 5)
	drawPoint(img, 184, 38, 4, 6)
	drawText(img, 194, 42, "err", 6)
	return img
}

// ── Вспомогательный маппинг: итерация → пиксельX ────────────────────────────

const (
	marginLeft   = 60
	marginRight  = 20
	marginTop    = 60
	marginBottom = 40
)

func plotWidth() int  { return width - marginLeft - marginRight }
func plotHeight() int { return height - marginTop - marginBottom }

func iterToX(iter, maxIter int) int {
	if maxIter <= 1 {
		return marginLeft
	}
	return marginLeft + iter*plotWidth()/(maxIter-1)
}

func valToY(v, vMin, vMax float64) int {
	if math.Abs(vMax-vMin) < 1e-12 {
		return marginTop + plotHeight()/2
	}
	frac := (v - vMin) / (vMax - vMin)
	return marginTop + plotHeight() - int(frac*float64(plotHeight()))
}

func drawAxes(img *image.Paletted, yMin, yMax float64, maxIter int) {
	// горизонтальная ось (y=0 или нижняя граница)
	yZeroFrac := (0.0 - yMin) / (yMax - yMin)
	yZeroPx := marginTop + plotHeight() - int(yZeroFrac*float64(plotHeight()))
	if yZeroPx < marginTop {
		yZeroPx = marginTop
	}
	if yZeroPx > marginTop+plotHeight() {
		yZeroPx = marginTop + plotHeight()
	}
	drawLine(img, marginLeft, yZeroPx, marginLeft+plotWidth(), yZeroPx, 1)

	// вертикальная ось
	drawLine(img, marginLeft, marginTop, marginLeft, marginTop+plotHeight(), 1)

	// засечки по Y
	for _, v := range []float64{yMin, (yMin + yMax) / 2, yMax} {
		py := valToY(v, yMin, yMax)
		drawLine(img, marginLeft-4, py, marginLeft, py, 1)
		drawText(img, 2, py+4, fmt.Sprintf("%.2f", v), 7)
	}
	// засечки по X
	for step := 0; step <= maxIter-1; step += max1(1, (maxIter-1)/6) {
		px := iterToX(step, maxIter)
		drawLine(img, px, marginTop+plotHeight(), px, marginTop+plotHeight()+4, 1)
		drawText(img, px-4, marginTop+plotHeight()+16, fmt.Sprintf("%d", step), 7)
	}
	drawText(img, marginLeft+plotWidth()+4, marginTop+plotHeight()+4, "iter", 7)
}

func max1(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ── Генератор 1: траектории компонент x1..x4 ────────────────────────────────

// GenerateConvergenceGIF анимирует как компоненты решения сходятся к x*.
// Каждый кадр добавляет одну итерацию к уже нарисованным линиям.
func GenerateConvergenceGIF(filePath string, history []solver.IterationRecord) error {
	if len(history) < 2 {
		return fmt.Errorf("недостаточно итераций для анимации")
	}

	// Диапазоны по Y
	yMin, yMax := math.Inf(1), math.Inf(-1)
	for _, rec := range history {
		for _, v := range rec.X {
			if v < yMin {
				yMin = v
			}
			if v > yMax {
				yMax = v
			}
		}
	}
	pad := (yMax - yMin) * 0.15
	yMin -= pad
	yMax += pad

	maxIter := len(history)
	anim := gif.GIF{}

	for frameEnd := 1; frameEnd <= len(history); frameEnd++ {
		img := newFrame(fmt.Sprintf("Simple Iteration: x components  (iter %d / %d)", frameEnd-1, len(history)-1))
		drawAxes(img, yMin, yMax, maxIter)

		colors := []uint8{2, 3, 4, 5}
		labels := []string{"x1", "x2", "x3", "x4"}

		for comp := 0; comp < 4; comp++ {
			ci := colors[comp]
			for k := 1; k < frameEnd; k++ {
				x0 := iterToX(k-1, maxIter)
				y0 := valToY(history[k-1].X[comp], yMin, yMax)
				x1 := iterToX(k, maxIter)
				y1 := valToY(history[k].X[comp], yMin, yMax)
				drawLine(img, x0, y0, x1, y1, ci)
			}
			// точка текущей итерации
			if frameEnd > 0 {
				px := iterToX(frameEnd-1, maxIter)
				py := valToY(history[frameEnd-1].X[comp], yMin, yMax)
				drawPoint(img, px, py, 3, ci)
				drawText(img, px+4, py-2, labels[comp], ci)
			}
		}

		// подпись текущих значений
		if frameEnd > 0 {
			rec := history[frameEnd-1]
			drawText(img, marginLeft, height-22,
				fmt.Sprintf("x=[%.4f, %.4f, %.4f, %.4f]  err=%.2e",
					rec.X[0], rec.X[1], rec.X[2], rec.X[3], rec.Error), 7)
		}

		delay := 15
		if frameEnd == len(history) {
			delay = 120
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	return saveGIF(filePath, &anim)
}

// ── Генератор 2: сходимость нормы ошибки ────────────────────────────────────

// GenerateErrorGIF анимирует убывание ||x^(k+1) - x^(k)|| по итерациям.
func GenerateErrorGIF(filePath string, history []solver.IterationRecord, eps float64) error {
	if len(history) < 2 {
		return fmt.Errorf("недостаточно итераций для анимации")
	}

	// Собираем конечные значения ошибок (пропускаем первую запись Inf)
	errs := make([]float64, 0, len(history))
	for _, r := range history {
		if !math.IsInf(r.Error, 1) {
			errs = append(errs, r.Error)
		}
	}

	eMax := 0.0
	for _, e := range errs {
		if e > eMax {
			eMax = e
		}
	}
	eMin := 0.0
	pad := eMax * 0.1
	yMax := eMax + pad
	yMin := eMin

	maxIter := len(history)
	anim := gif.GIF{}

	for frameEnd := 1; frameEnd <= len(errs); frameEnd++ {
		img := newFrame(fmt.Sprintf("Simple Iteration: error norm  (iter %d / %d)", frameEnd, len(errs)))
		drawAxes(img, yMin, yMax, maxIter)

		// eps-линия
		epsPy := valToY(eps, yMin, yMax)
		drawLine(img, marginLeft, epsPy, marginLeft+plotWidth(), epsPy, 9)
		drawText(img, marginLeft+plotWidth()-55, epsPy-4, fmt.Sprintf("eps=%.3f", eps), 9)

		// кривая ошибки
		for k := 1; k < frameEnd; k++ {
			// +1 сдвиг потому что errs[0] соответствует iter=1 (первый шаг)
			x0 := iterToX(k, maxIter)
			y0 := valToY(errs[k-1], yMin, yMax)
			x1 := iterToX(k+1, maxIter)
			y1 := valToY(errs[k], yMin, yMax)
			drawLine(img, x0, y0, x1, y1, 6)
		}
		if frameEnd > 0 {
			px := iterToX(frameEnd, maxIter)
			py := valToY(errs[frameEnd-1], yMin, yMax)
			drawPoint(img, px, py, 4, 6)
			drawText(img, marginLeft, height-22,
				fmt.Sprintf("iter %d  ||Δx|| = %.6f  (eps = %.3f)", frameEnd, errs[frameEnd-1], eps), 7)
		}

		delay := 15
		if frameEnd == len(errs) {
			delay = 120
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	return saveGIF(filePath, &anim)
}

// ── Генератор 3: совмещённый — компоненты + ошибка бок о бок ────────────────

// GenerateCombinedGIF отрисовывает оба графика на одном широком холсте
// (левая половина — траектории x_i, правая — норма ошибки).
func GenerateCombinedGIF(filePath string, history []solver.IterationRecord, eps float64) error {
	if len(history) < 2 {
		return fmt.Errorf("недостаточно итераций для анимации")
	}

	const (
		wFull  = 1100
		hFull  = 480
		wHalf  = wFull / 2
		mLeft  = 55
		mRight = 15
		mTop   = 60
		mBot   = 40
	)

	pW := wHalf - mLeft - mRight
	pH := hFull - mTop - mBot

	// диапазон Y для компонент
	yMin, yMax := math.Inf(1), math.Inf(-1)
	for _, rec := range history {
		for _, v := range rec.X {
			if v < yMin {
				yMin = v
			}
			if v > yMax {
				yMax = v
			}
		}
	}
	pad := (yMax - yMin) * 0.15
	yMin -= pad
	yMax += pad

	// диапазон Y для ошибок
	errs := make([]float64, 0)
	for _, r := range history[1:] {
		errs = append(errs, r.Error)
	}
	eMax := 0.0
	for _, e := range errs {
		if e > eMax {
			eMax = e
		}
	}
	eyMin := 0.0
	eyMax := eMax * 1.1

	maxIter := len(history)

	toX := func(iter, panel int) int {
		base := panel * wHalf
		if maxIter <= 1 {
			return base + mLeft
		}
		return base + mLeft + iter*pW/(maxIter-1)
	}
	toY := func(v, vMin, vMax float64) int {
		if math.Abs(vMax-vMin) < 1e-12 {
			return mTop + pH/2
		}
		frac := (v - vMin) / (vMax - vMin)
		return mTop + pH - int(frac*float64(pH))
	}

	drawAxisPair := func(img *image.Paletted) {
		for panel := 0; panel < 2; panel++ {
			base := panel * wHalf
			// вертикальная
			drawLine(img, base+mLeft, mTop, base+mLeft, mTop+pH, 1)
			// горизонтальная (y=0 для левой, низ для правой)
			vMin, vMax := yMin, yMax
			if panel == 1 {
				vMin, vMax = eyMin, eyMax
			}
			yZeroFrac := (0.0 - vMin) / (vMax - vMin)
			yZeroPx := mTop + pH - int(yZeroFrac*float64(pH))
			if yZeroPx < mTop {
				yZeroPx = mTop
			}
			if yZeroPx > mTop+pH {
				yZeroPx = mTop + pH
			}
			drawLine(img, base+mLeft, yZeroPx, base+mLeft+pW, yZeroPx, 1)
			// засечки Y
			for _, v := range []float64{vMin, (vMin + vMax) / 2, vMax} {
				py := toY(v, vMin, vMax)
				drawText(img, base+2, py+4, fmt.Sprintf("%.2f", v), 7)
			}
		}
	}

	pal := make([]color.Color, len(palette))
	copy(pal, palette)

	anim := gif.GIF{}
	compColors := []uint8{2, 3, 4, 5}

	for frameEnd := 1; frameEnd <= len(history); frameEnd++ {
		rect := image.Rect(0, 0, wFull, hFull)
		img := image.NewPaletted(rect, pal)

		// заголовки
		drawText(img, mLeft, 18, fmt.Sprintf("Component trajectories  (iter %d / %d)", frameEnd-1, len(history)-1), 7)
		drawText(img, wHalf+mLeft, 18, "Convergence: ||Δx|| per iteration", 7)

		// легенда левой панели
		for i, ci := range compColors {
			ox := mLeft + i*55
			drawPoint(img, ox, 38, 4, ci)
			drawText(img, ox+8, 42, fmt.Sprintf("x%d", i+1), ci)
		}
		// eps-метка правой панели
		drawText(img, wHalf+mLeft, 38, fmt.Sprintf("— eps=%.3f", eps), 9)

		drawAxisPair(img)

		// левая панель: траектории
		for comp := 0; comp < 4; comp++ {
			ci := compColors[comp]
			for k := 1; k < frameEnd; k++ {
				x0 := toX(k-1, 0)
				y0 := toY(history[k-1].X[comp], yMin, yMax)
				x1 := toX(k, 0)
				y1 := toY(history[k].X[comp], yMin, yMax)
				drawLine(img, x0, y0, x1, y1, ci)
			}
			if frameEnd > 0 {
				px := toX(frameEnd-1, 0)
				py := toY(history[frameEnd-1].X[comp], yMin, yMax)
				drawPoint(img, px, py, 3, ci)
			}
		}

		// правая панель: ошибка
		epsPy := toY(eps, eyMin, eyMax)
		drawLine(img, wHalf+mLeft, epsPy, wHalf+mLeft+pW, epsPy, 9)

		errFrame := frameEnd - 1 // errs[0] = history[1].Error
		if errFrame > len(errs) {
			errFrame = len(errs)
		}
		for k := 1; k < errFrame; k++ {
			x0 := toX(k, 1)
			y0 := toY(errs[k-1], eyMin, eyMax)
			x1 := toX(k+1, 1)
			y1 := toY(errs[k], eyMin, eyMax)
			drawLine(img, x0, y0, x1, y1, 6)
		}
		if errFrame > 0 {
			px := toX(errFrame, 1)
			py := toY(errs[errFrame-1], eyMin, eyMax)
			drawPoint(img, px, py, 4, 6)
		}

		// нижняя строка
		if frameEnd > 0 {
			rec := history[frameEnd-1]
			errVal := rec.Error
			if math.IsInf(errVal, 1) {
				errVal = 0
			}
			drawText(img, 10, hFull-10,
				fmt.Sprintf("x=[%.4f, %.4f, %.4f, %.4f]  ||Δx||=%.2e",
					rec.X[0], rec.X[1], rec.X[2], rec.X[3], errVal), 7)
		}

		delay := 15
		if frameEnd == len(history) {
			delay = 150
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	return saveGIF(filePath, &anim)
}

// ── Сохранение ───────────────────────────────────────────────────────────────

func saveGIF(filePath string, anim *gif.GIF) error {
	// 4 дополнительных кадра-заморозки в конце
	for i := 0; i < 4; i++ {
		last := anim.Image[len(anim.Image)-1]
		anim.Image = append(anim.Image, last)
		anim.Delay = append(anim.Delay, 150)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, anim)
}
