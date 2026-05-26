package solver

import "math"

// Вариант 8:
// x: 0  1  3  5,  x* = 1
// y: 11 12 13 11
// (из таблицы: строка 8, первая половина — x:[0,1,3,5], y:[11,12,13,11])

// DataPoint — узел интерполяции
type DataPoint struct {
	X, Y float64
}

// Variant8Data возвращает узловые точки для варианта 8
func Variant8Data() []DataPoint {
	return []DataPoint{
		{0, 11},
		{1, 12},
		{3, 13},
		{5, 11},
	}
}

// Variant8Target — точка, в которой вычисляется f(x*)
const Variant8Target = 1.0

// LagrangeBasis вычисляет i-й базисный полином L_i(x)
func LagrangeBasis(points []DataPoint, i int, x float64) float64 {
	n := len(points)
	result := 1.0
	for j := 0; j < n; j++ {
		if j != i {
			result *= (x - points[j].X) / (points[i].X - points[j].X)
		}
	}
	return result
}

// LagrangeInterpolate вычисляет значение интерполяционного многочлена Лагранжа в точке x
func LagrangeInterpolate(points []DataPoint, x float64) float64 {
	result := 0.0
	for i, p := range points {
		result += p.Y * LagrangeBasis(points, i, x)
	}
	return result
}

// IterationRecord — запись для анимации
type IterationRecord struct {
	Step  int
	X     float64   // текущая точка вычисления
	Value float64   // L(x) на текущем шаге
	Basis []float64 // значения базисных полиномов
	Used  int       // сколько узлов задействовано (1..n)
}

// BuildConvergenceHistory строит историю добавления узлов (от 1 до n)
// для анимации: сначала 1 узел, потом 2, ..., потом n
func BuildConvergenceHistory(points []DataPoint, xTarget float64) []IterationRecord {
	var history []IterationRecord
	n := len(points)

	for k := 1; k <= n; k++ {
		subset := points[:k]
		val := LagrangeInterpolate(subset, xTarget)
		basis := make([]float64, k)
		for i := range subset {
			basis[i] = LagrangeBasis(subset, i, xTarget)
		}
		history = append(history, IterationRecord{
			Step:  k,
			X:     xTarget,
			Value: val,
			Basis: basis,
			Used:  k,
		})
	}
	return history
}

// EvalPolynomial вычисляет значения многочлена на отрезке для построения графика
func EvalPolynomial(points []DataPoint, xMin, xMax float64, nPts int) ([]float64, []float64) {
	xs := make([]float64, nPts)
	ys := make([]float64, nPts)
	for i := 0; i < nPts; i++ {
		x := xMin + float64(i)/float64(nPts-1)*(xMax-xMin)
		xs[i] = x
		ys[i] = LagrangeInterpolate(points, x)
	}
	return xs, ys
}

// AbsFloat64 — абсолютное значение
func AbsFloat64(x float64) float64 { return math.Abs(x) }
