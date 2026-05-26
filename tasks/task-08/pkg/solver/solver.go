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
