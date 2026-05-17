package solver

import (
	"fmt"
	"math"
)

// F возвращает значение функции f(x) = sin(x + pi/3) - 0.5x
func F(x float64) float64 {
	return math.Sin(x+math.Pi/3) - 0.5*x
}

// FindInterval локализует интервал, на котором функция меняет знак
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
	return 0, 0, fmt.Errorf("no root found in the given range for %d steps", maxSteps)
}

// Bisection выполняет метод половинного деления (дихотомии)
func DichotomyMethod(a, b, eps float64) (float64, int) {
	iterations := 0
	var mid float64
	for (b - a) > eps {
		iterations++
		mid = (a + b) / 2
		if F(mid) == 0 {
			return mid, iterations
		} else if F(a)*F(mid) < 0 {
			b = mid
		} else {
			a = mid
		}
	}
	return mid, iterations
}
