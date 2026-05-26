package solver

import (
	"fmt"
	"math"
)

// Вариант 8:
//
//	x1 =  0·x1  + 0.1·x2 - 0.1·x3 + 0.2·x4 - 1
//	x2 =  0.2·x1 + 0·x2  - 0.2·x3 + 0.1·x4 - 1
//	x3 =  0.13·x1 - 0.2·x2 + 0·x3 + 0.3·x4 + 2
//	x4 =  0.1·x1 - 0.1·x2 - 0.2·x3 + 0·x4  + 0.1

// C — матрица коэффициентов итерационной схемы x = Cx + d
var C = [4][4]float64{
	{0.0, 0.1, -0.1, 0.2},
	{0.2, 0.0, -0.2, 0.1},
	{0.13, -0.2, 0.0, 0.3},
	{0.1, -0.1, -0.2, 0.0},
}

// D — вектор свободных членов
var D = [4]float64{-1, -1, 2, 0.1}

// IterationStep выполняет одну итерацию схемы x = Cx + d
// и возвращает новый вектор x.
func IterationStep(x [4]float64) [4]float64 {
	var next [4]float64
	for i := 0; i < 4; i++ {
		s := D[i]
		for j := 0; j < 4; j++ {
			s += C[i][j] * x[j]
		}
		next[i] = s
	}
	return next
}

// MaxNorm возвращает чебышёвскую норму вектора (максимум абсолютных значений).
func MaxNorm(v [4]float64) float64 {
	m := 0.0
	for _, vi := range v {
		if a := math.Abs(vi); a > m {
			m = a
		}
	}
	return m
}

// MaxNormDiff возвращает норму разности двух векторов.
func MaxNormDiff(a, b [4]float64) float64 {
	var diff [4]float64
	for i := range diff {
		diff[i] = a[i] - b[i]
	}
	return MaxNorm(diff)
}

// SpectralRadius возвращает приближённый спектральный радиус матрицы C
// (максимум сумм абсолютных значений строк — достаточное условие сходимости).
func SpectralRadius() float64 {
	max := 0.0
	for i := 0; i < 4; i++ {
		s := 0.0
		for j := 0; j < 4; j++ {
			s += math.Abs(C[i][j])
		}
		if s > max {
			max = s
		}
	}
	return max
}

// IterationRecord хранит состояние одной итерации для истории сходимости.
type IterationRecord struct {
	Iter  int
	X     [4]float64
	Error float64 // ||x^(k+1) - x^(k)||
}

// Solve решает систему методом простой итерации x = Cx + d.
// Начальное приближение x0 = d (нулевой вектор можно передать явно).
// Возвращает: корень, историю итераций, число выполненных итераций, ошибку.
func Solve(x0 [4]float64, eps float64, maxIter int) ([4]float64, []IterationRecord, error) {
	history := []IterationRecord{}
	x := x0
	history = append(history, IterationRecord{Iter: 0, X: x, Error: math.Inf(1)})

	for k := 1; k <= maxIter; k++ {
		xNext := IterationStep(x)
		err := MaxNormDiff(xNext, x)

		history = append(history, IterationRecord{Iter: k, X: xNext, Error: err})

		if err < eps {
			return xNext, history, nil
		}
		x = xNext
	}
	return x, history, fmt.Errorf("превышено максимальное число итераций (%d)", maxIter)
}

// Residual вычисляет невязку r = ||Cx + d - x|| для проверки точности решения.
func Residual(x [4]float64) float64 {
	cx := IterationStep(x) // Cx + d
	return MaxNormDiff(cx, x)
}
