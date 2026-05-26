package solver

import "math"

// IntegrationRecord хранит информацию о шаге интегрирования для анимации и логов
type IntegrationRecord struct {
	Step  int
	N     int
	H     float64
	Value float64
	Error float64
}

// Variant8Func возвращает значение интегрируемой функции для варианта 8:
// f(x) = (x + 1) / (2 + ln(1 + x^2))
func Variant8Func(x float64) float64 {
	return (x + 1.0) / (2.0 + math.Log(1.0+x*x))
}

// SimpsonIntegrate выполняет составное интегрирование методом Симпсона для заданного N
func SimpsonIntegrate(f func(float64) float64, a, b float64, n int) float64 {
	if n%2 != 0 {
		n++ // Метод Симпсона требует четного числа разбиений
	}
	h := (b - a) / float64(n)
	sum := f(a) + f(b)

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2.0 * f(x)
		} else {
			sum += 4.0 * f(x)
		}
	}
	return sum * h / 3.0
}

// BuildIntegrationHistory итеративно удваивает число разбиений N,
// оценивая погрешность по правилу Рунге, пока не будет достигнута точность eps
func BuildIntegrationHistory(a, b, eps float64) []IntegrationRecord {
	var history []IntegrationRecord
	n := 2
	step := 1

	for {
		valCurrent := SimpsonIntegrate(Variant8Func, a, b, n)
		errEst := 0.0

		if step > 1 {
			valPrev := history[step-2].Value
			// Для метода Симпсона порядок точности p=4, коэффициент Рунге = 1 / (2^4 - 1) = 1/15
			errEst = math.Abs(valCurrent-valPrev) / 15.0
		}

		h := (b - a) / float64(n)
		history = append(history, IntegrationRecord{
			Step:  step,
			N:     n,
			H:     h,
			Value: valCurrent,
			Error: errEst,
		})

		if step > 1 && errEst <= eps {
			break
		}

		// Предохранитель от бесконечного зацикливания
		if n > 8192 {
			break
		}

		n *= 2
		step++
	}
	return history
}
