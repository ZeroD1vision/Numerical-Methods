package solver

import "math"

// Функция f(t, y) для варианта 8: y' = cos(sqrt(t * y^2))
func F(t, y float64) float64 {
	// y^2 всегда положительно, t на отрезке [1, 2] тоже положительно
	return math.Cos(math.Sqrt(t * y * y))
}

// Point представляет точку на графике (t, y)
type Point struct {
	T float64
	Y float64
}

// Euler решает задачу Коши методом ломаных Эйлера
func Euler(t0, y0, T float64, N int) []Point {
	h := (T - t0) / float64(N)
	pts := make([]Point, N+1)
	pts[0] = Point{T: t0, Y: y0}

	for i := 0; i < N; i++ {
		t := pts[i].T
		y := pts[i].Y
		pts[i+1] = Point{
			T: t + h,
			Y: y + h*F(t, y),
		}
	}
	return pts
}

// ModifiedEuler решает задачу Коши модифицированным методом Эйлера (метод средней точки)
func ModifiedEuler(t0, y0, T float64, N int) []Point {
	h := (T - t0) / float64(N)
	pts := make([]Point, N+1)
	pts[0] = Point{T: t0, Y: y0}

	for i := 0; i < N; i++ {
		t := pts[i].T
		y := pts[i].Y

		// y_{k+1} = y_k + h * f(t_k + h/2, y_k + h/2 * f(t_k, y_k))
		halfH := h / 2.0
		halfK := halfH * F(t, y)

		pts[i+1] = Point{
			T: t + h,
			Y: y + h*F(t+halfH, y+halfK),
		}
	}
	return pts
}

// RK4 решает ОДУ классическим методом Рунге-Кутты 4-го порядка.
// Используется для генерации эталонного ("точного") решения с мелким шагом.
func RK4(t0, y0, T float64, N int) []Point {
	h := (T - t0) / float64(N)
	pts := make([]Point, N+1)
	pts[0] = Point{T: t0, Y: y0}

	for i := 0; i < N; i++ {
		t := pts[i].T
		y := pts[i].Y

		k1 := F(t, y)
		k2 := F(t+h/2.0, y+h/2.0*k1)
		k3 := F(t+h/2.0, y+h/2.0*k2)
		k4 := F(t+h, y+h*k3)

		pts[i+1] = Point{
			T: t + h,
			Y: y + (h/6.0)*(k1+2*k2+2*k3+k4),
		}
	}
	return pts
}

// GetExactY принимает массив точек высокой точности и возвращает интерполированное
// значение Y для конкретного t (нужно для отрисовки плавного точного графика).
func GetExactY(exactPts []Point, t float64) float64 {
	if t <= exactPts[0].T {
		return exactPts[0].Y
	}
	if t >= exactPts[len(exactPts)-1].T {
		return exactPts[len(exactPts)-1].Y
	}

	// Бинарный поиск нужного отрезка
	low, high := 0, len(exactPts)-1
	for high-low > 1 {
		mid := (low + high) / 2
		if exactPts[mid].T <= t {
			low = mid
		} else {
			high = mid
		}
	}

	// Линейная интерполяция между соседними точками эталона
	p0, p1 := exactPts[low], exactPts[high]
	if p1.T == p0.T {
		return p0.Y
	}
	return p0.Y + (t-p0.T)*(p1.Y-p0.Y)/(p1.T-p0.T)
}
