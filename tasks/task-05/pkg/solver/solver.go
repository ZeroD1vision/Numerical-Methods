package solver

import "math"

// Вариант 8: α = -7
// Матрица A = | 1  2  α |   = | 1  2 -7 |
//             | 2  3  4 |     | 2  3  4 |
//             | α  4  5 |     |-7  4  5 |

const Alpha = -7.0

// Matrix3x3 — симметричная матрица 3×3
type Matrix3x3 [3][3]float64

// IterationRecord хранит состояние алгоритма на каждом шаге
type IterationRecord struct {
	Iter       int
	Eigenvalue float64 // доминантное собственное значение
	Error      float64 // ||A*v - λ*v||
	Vector     [3]float64
}

// BuildMatrix строит матрицу A для варианта 8
func BuildMatrix() Matrix3x3 {
	return Matrix3x3{
		{1, 2, Alpha},
		{2, 3, 4},
		{Alpha, 4, 5},
	}
}

// mulMV умножает матрицу на вектор
func mulMV(A Matrix3x3, v [3]float64) [3]float64 {
	var r [3]float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r[i] += A[i][j] * v[j]
		}
	}
	return r
}

// norm2 вычисляет евклидову норму вектора
func norm2(v [3]float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x * x
	}
	return math.Sqrt(s)
}

// normalize нормализует вектор
func normalize(v [3]float64) [3]float64 {
	n := norm2(v)
	if n < 1e-14 {
		return v
	}
	return [3]float64{v[0] / n, v[1] / n, v[2] / n}
}

// dot вычисляет скалярное произведение
func dot(a, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

// PowerIteration — метод степенных итераций для нахождения наибольшего по модулю СЗ
func PowerIteration(A Matrix3x3, eps float64, maxIter int) []IterationRecord {
	// Начальный вектор
	v := normalize([3]float64{1, 1, 1})
	history := []IterationRecord{{Iter: 0, Eigenvalue: 0, Error: math.Inf(1), Vector: v}}

	var lambda float64
	for k := 1; k <= maxIter; k++ {
		w := mulMV(A, v)
		lambda = dot(w, v) // Rayleigh quotient
		newV := normalize(w)

		// Невязка ||A*v - λ*v||
		residual := [3]float64{
			w[0] - lambda*v[0],
			w[1] - lambda*v[1],
			w[2] - lambda*v[2],
		}
		err := norm2(residual)
		history = append(history, IterationRecord{Iter: k, Eigenvalue: lambda, Error: err, Vector: newV})

		v = newV
		if err < eps {
			break
		}
	}
	return history
}

// CharPolyRoots вычисляет все 3 СЗ через характеристический полином
func AllEigenvalues(A Matrix3x3) [3]float64 {
	// Метод Якоби для симметричной матрицы
	mat := A
	const maxSweeps = 100
	const eps = 1e-10

	eigenvals := [3]float64{0, 0, 0}

	for sweep := 0; sweep < maxSweeps; sweep++ {
		// Найти наибольший внедиагональный элемент
		maxVal := 0.0
		p, q := 0, 1
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if math.Abs(mat[i][j]) > maxVal {
					maxVal = math.Abs(mat[i][j])
					p, q = i, j
				}
			}
		}
		if maxVal < eps {
			break
		}
		// Вычислить угол поворота
		theta := 0.5 * math.Atan2(2*mat[p][q], mat[q][q]-mat[p][p])
		c, s := math.Cos(theta), math.Sin(theta)

		// Построить матрицу вращения и применить
		newMat := mat
		newMat[p][p] = c*c*mat[p][p] - 2*s*c*mat[p][q] + s*s*mat[q][q]
		newMat[q][q] = s*s*mat[p][p] + 2*s*c*mat[p][q] + c*c*mat[q][q]
		newMat[p][q] = 0
		newMat[q][p] = 0
		for r := 0; r < 3; r++ {
			if r != p && r != q {
				newMat[p][r] = c*mat[p][r] - s*mat[q][r]
				newMat[r][p] = newMat[p][r]
				newMat[q][r] = s*mat[p][r] + c*mat[q][r]
				newMat[r][q] = newMat[q][r]
			}
		}
		mat = newMat
	}
	eigenvals[0] = mat[0][0]
	eigenvals[1] = mat[1][1]
	eigenvals[2] = mat[2][2]

	// Сортировка по убыванию
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if eigenvals[j] > eigenvals[i] {
				eigenvals[i], eigenvals[j] = eigenvals[j], eigenvals[i]
			}
		}
	}
	return eigenvals
}
