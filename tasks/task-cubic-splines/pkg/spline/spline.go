package spline

import (
	"fmt"
)

// Segment содержит вычисленные коэффициенты для одного интервала [x_i, x_{i+1}]
type Segment struct {
	XMin float64
	XMax float64
	A    float64
	B    float64
	C    float64
	D    float64
}

// Interp вычисляет значение сплайна в заданной точке x
func (s Segment) Interp(x float64) float64 {
	dx := x - s.XMin
	return s.A + s.B*dx + s.C*dx*dx + s.D*dx*dx*dx
}

// BuildSpline строит кубический сплайн по точкам X и Y (условия на краях: m_0 = m_n = 0)
func BuildSpline(x, y []float64) ([]Segment, error) {
	n := len(x) - 1
	if n < 2 {
		return nil, fmt.Errorf("недостаточно точек для построения кубического сплайна")
	}

	// 1. Вычисление шагов h_i и запись в массив h
	h := make([]float64, n)
	for i := 0; i < n; i++ {
		h[i] = x[i+1] - x[i]
		if h[i] <= 0 {
			return nil, fmt.Errorf("узлы сетки X должны быть строго упорядочены по возрастанию")
		}
	}

	// 2. Формирование трехдиагональной СЛАУ для вторых производных m
	// Строки матрицы от 1 до n-1 (внутренние стыки)
	aSub := make([]float64, n+1)   // нижняя диагональ
	bDiag := make([]float64, n+1)  // главная диагональ
	cUp := make([]float64, n+1)    // верхняя диагональ
	dRight := make([]float64, n+1) // правая часть (изломы)

	// Условия на краях: m_0 = 0, m_n = 0
	bDiag[0] = 1.0
	bDiag[n] = 1.0

	for i := 1; i < n; i++ {
		aSub[i] = h[i-1]
		cUp[i] = h[i]
		bDiag[i] = 2.0 * (h[i-1] + h[i])
		dRight[i] = 6.0 * ((y[i+1]-y[i])/h[i] - (y[i]-y[i-1])/h[i-1])
	}

	// 3. Решение СЛАУ методом прогонки
	m := make([]float64, n+1)
	alpha := make([]float64, n+1)
	beta := make([]float64, n+1)

	alpha[1] = -cUp[0] / bDiag[0]
	beta[1] = dRight[0] / bDiag[0]

	for i := 1; i < n; i++ {
		denom := bDiag[i] + aSub[i]*alpha[i]
		alpha[i+1] = -cUp[i] / denom
		beta[i+1] = (dRight[i] - aSub[i]*beta[i]) / denom
	}

	m[n] = (dRight[n] - aSub[n]*beta[n]) / (bDiag[n] + aSub[n]*alpha[n])
	for i := n - 1; i >= 0; i-- {
		m[i] = alpha[i+1]*m[i+1] + beta[i+1]
	}

	// 4. Расчет конечных коэффициентов для каждого сегмента сплайна
	segments := make([]Segment, n)
	for i := 0; i < n; i++ {
		segments[i] = Segment{
			XMin: x[i],
			XMax: x[i+1],
			A:    y[i],
			C:    m[i] / 2.0,
			D:    (m[i+1] - m[i]) / (6.0 * h[i]),
			B:    (y[i+1]-y[i])/h[i] - h[i]*(2.0*m[i]+m[i+1])/6.0,
		}
	}

	return segments, nil
}
