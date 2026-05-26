package solver

import (
	"fmt"
	"math"
)

// Исходная система нелинейных уравнений F(x, y) = 0
func F1(x, y float64) float64 { return 2*x - math.Cos(y+1) }
func F2(x, y float64) float64 { return y + math.Sin(x) + 0.4 }

// Эквивалентная форма x = Phi1(x, y) и y = Phi2(x, y) для МПИ
// Из 1-го уравнения: x = 0.5 * cos(y + 1)
// Из 2-го уравнения: y = -0.4 - sin(x)
func Phi1(x, y float64) float64 { return 0.5 * math.Cos(y+1) }
func Phi2(x, y float64) float64 { return -0.4 - math.Sin(x) }

// Проверка достаточного условия сходимости МПИ в заданной точке.
// Вычисляются частные производные отображения Phi:
// d(Phi1)/dx = 0,        d(Phi1)/dy = -0.5 * sin(y + 1)
// d(Phi2)/dx = -cos(x),   d(Phi2)/dy = 0
// Сходимость гарантирована, если чебышёвская норма матрицы Якоби меньше 1.
func CheckConvergence(x, y float64) (bool, float64) {
	dPhi1dy := math.Abs(-0.5 * math.Sin(y+1))
	dPhi2dx := math.Abs(-math.Cos(x))

	norm := math.Max(dPhi1dy, dPhi2dx)
	return norm < 1, norm
}

// Элементы матрицы Якоби системы F(x, y) для метода Ньютона:
// J = | dF1/dx  dF1/dy | = |    2      sin(y+1) |
//
//	| dF2/dx  dF2/dy |   |  cos(x)      1     |
func Jacobian(x, y float64) (float64, float64, float64, float64) {
	df1dx := 2.0
	df1dy := math.Sin(y + 1)
	df2dx := math.Cos(x)
	df2dy := 1.0
	return df1dx, df1dy, df2dx, df2dy
}

type IterationRecord struct {
	Iter  int
	X     float64
	Y     float64
	Error float64 // Максимальное отклонение текущего шага ||X_k - X_{k-1}||
}

// SolveSimpleIteration — Метод простой итерации
func SolveSimpleIteration(x0, y0 float64, eps float64, maxIter int) ([]IterationRecord, error) {
	history := []IterationRecord{{Iter: 0, X: x0, Y: y0, Error: math.Inf(1)}}
	cx, cy := x0, y0

	for k := 1; k <= maxIter; k++ {
		nextX := Phi1(cx, cy)
		nextY := Phi2(cx, cy)

		err := math.Max(math.Abs(nextX-cx), math.Abs(nextY-cy))
		history = append(history, IterationRecord{Iter: k, X: nextX, Y: nextY, Error: err})

		if err < eps {
			return history, nil
		}
		cx, cy = nextX, nextY
	}
	return history, fmt.Errorf("превышено ограничение в %d итераций", maxIter)
}

// SolveNewton — Метод Ньютона (касательных) для СНУ
func SolveNewton(x0, y0 float64, eps float64, maxIter int) ([]IterationRecord, error) {
	history := []IterationRecord{{Iter: 0, X: x0, Y: y0, Error: math.Inf(1)}}
	cx, cy := x0, y0

	for k := 1; k <= maxIter; k++ {
		f1 := F1(cx, cy)
		f2 := F2(cx, cy)

		a, b, c, d := Jacobian(cx, cy)
		det := a*d - b*c // Вычисляем определитель матрицы Якоби

		if math.Abs(det) < 1e-12 {
			return history, fmt.Errorf("матрица Якоби вырождена (det близко к 0)")
		}

		// Приращения по правилу Крамера (J * delta = -F)
		deltaX := (-f1*d - (-f2 * b)) / det
		deltaY := (a*(-f2) - c*(-f1)) / det

		nextX := cx + deltaX
		nextY := cy + deltaY

		err := math.Max(math.Abs(deltaX), math.Abs(deltaY))
		history = append(history, IterationRecord{Iter: k, X: nextX, Y: nextY, Error: err})

		if err < eps {
			return history, nil
		}
		cx, cy = nextX, nextY
	}
	return history, fmt.Errorf("превышено ограничение в %d итераций", maxIter)
}
