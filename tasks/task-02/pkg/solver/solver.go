package solver

import (
	"fmt"
	"math"
)

// F исходная функция: x^3 + 3x - 1 = 0
func F(x float64) float64 {
	return x*x*x + 3*x - 1
}

// D1 первая производная: 3x^2 + 3
func D1(x float64) float64 {
	return 3*x*x + 3
}

// D2 вторая производная: 6x
func D2(x float64) float64 {
	return 6 * x
}

// FindInterval автоматически ищет отрезок локализации корня
func FindInterval(start, step float64, maxSteps int) (float64, float64, error) {
	a := start
	for i := 0; i < maxSteps; i++ {
		b := a + step
		if F(a)*F(b) < 0 {
			if a < b {
				return a, b, nil
			}
			return b, a, nil
		}
		a = b
	}

	// Поиск в отрицательную сторону
	a = start
	for i := 0; i < maxSteps; i++ {
		b := a - step
		if F(a)*F(b) < 0 {
			if b < a {
				return b, a, nil
			}
			return a, b, nil
		}
		a = b
	}
	return 0, 0, fmt.Errorf("интервал изменения знака не найден за %d шагов", maxSteps)
}

// NewtonMethod реализует метод касательных
func NewtonMethod(a, b, eps float64) (float64, []float64, error) {
	// Выбор начального приближения x0 по условию Фурье: f(x0) * f''(x0) > 0
	var x0 float64
	if F(a)*D2(a) > 0 {
		x0 = a
	} else if F(b)*D2(b) > 0 {
		x0 = b
	} else {
		// Если на краях не сработало, берем середину как дефолт
		x0 = (a + b) / 2
	}

	x := x0
	var history []float64
	history = append(history, x)
	iterations := 0
	maxIterations := 1000

	for i := 0; i < maxIterations; i++ {
		iterations++
		d1 := D1(x)
		if math.Abs(d1) < 1e-12 {
			return 0, nil, fmt.Errorf("производная близка к нулю, метод Ньютона расходится")
		}

		nextX := x - F(x)/d1
		history = append(history, nextX)

		if math.Abs(nextX-x) < eps {
			return nextX, history, nil
		}
		x = nextX
	}

	return x, history, fmt.Errorf("превышено максимальное число итераций метода Ньютона")
}

// SecantMethod реализует метод хорд
func SecantMethod(a, b, eps float64) (float64, []float64, error) {
	// Метод хорд требует двух начальных приближений.
	// Выбираем неподвижную точку (где знак f(x) совпадает со знаком f''(x))
	// Вторая точка будет итерироваться.
	var xFixed, x float64

	if F(a)*D2(a) > 0 {
		xFixed = a
		x = b
	} else {
		xFixed = b
		x = a
	}

	var history []float64
	history = append(history, x)

	iterations := 0
	maxIterations := 1000

	for i := 0; i < maxIterations; i++ {
		iterations++
		fx := F(x)
		fxFixed := F(xFixed)

		denom := fx - fxFixed
		if math.Abs(denom) < 1e-12 {
			return 0, nil, fmt.Errorf("знаменатель близок к нулю в методе хорд")
		}

		// Формула метода хорд с одной неподвижной границей
		nextX := x - fx*(x-xFixed)/denom
		history = append(history, nextX)

		if math.Abs(nextX-x) < eps {
			return nextX, history, nil
		}
		x = nextX
	}

	return x, history, fmt.Errorf("превышено максимальное число итераций метода хорд")
}
