package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-02/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-02/pkg/visualizer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Поиск корня методами Ньютона и Хорд ===")
	log.Info("Уравнение: x^3 + 3x - 1 = 0")
	log.Info("---------------------------------------------------------")

	// 1. Ввод точности eps
	log.Prompt("Введите точность eps (например, 0.001): ")
	epsStr, _ := reader.ReadString('\n')
	epsStr = strings.TrimSpace(epsStr)
	eps, err := strconv.ParseFloat(epsStr, 64)
	if err != nil || eps <= 0 {
		log.Warn("Некорректное значение точности. Установлено eps = 0.001")
		eps = 0.001
	}

	// 2. Ввод границ отрезка
	log.Prompt("Введите левую границу 'a' (или Enter для автопоиска): ")
	aStr, _ := reader.ReadString('\n')
	aStr = strings.TrimSpace(aStr)

	var a, b float64
	useAutoSearch := false

	if aStr == "" {
		useAutoSearch = true
	} else {
		log.Prompt("Введите правую границу 'b': ")
		bStr, _ := reader.ReadString('\n')
		bStr = strings.TrimSpace(bStr)

		var errA, errB error
		a, errA = strconv.ParseFloat(aStr, 64)
		b, errB = strconv.ParseFloat(bStr, 64)

		if errA != nil || errB != nil {
			log.Warn("Некорректный ввод параметров. Включен автопоиск интервала...")
			useAutoSearch = true
		} else if a >= b {
			log.Warn("Левая граница должна быть строго меньше правой. Включен автопоиск...")
			useAutoSearch = true
		}
	}

	// 3. Определение интервала локализации
	if useAutoSearch {
		log.Info("Запуск автоматической локализации корня...")
		var errInterval error
		a, b, errInterval = solver.FindInterval(0.0, 0.5, 2000)
		if errInterval != nil {
			log.Error("Критическая ошибка локализации интервала: %v", errInterval)
			return
		}
		log.Success(fmt.Sprintf("Интервал успешно определен: [%.4f, %.4f]", a, b))
	} else {
		if solver.F(a)*solver.F(b) >= 0 {
			log.Error(fmt.Sprintf("Ошибка: критерий Больцано-Коши не выполнен. На [%.4f, %.4f] нет смены знаков функции.", a, b))
			return
		}
	}

	decimalPlaces := getDecimalPlaces(eps)

	// 4. Расчет методом Ньютона
	log.Info("Запуск вычислений методом Ньютона (касательных)...")
	newtonRoot, newtonHistory, errNewton := solver.NewtonMethod(a, b, eps)
	if errNewton != nil {
		log.Error(fmt.Sprintf("Ошибка метода Ньютона: %v", errNewton))
	} else {
		log.Success(fmt.Sprintf("Метод Ньютона: x = %.*f (Итераций: %d, f(x) = %e)", decimalPlaces, newtonRoot, len(newtonHistory)-1, solver.F(newtonRoot))) // <-- Изменено тут
	}

	// 5. Расчет методом Хорд
	log.Info("Запуск вычислений методом хорд (секущих)...")
	secantRoot, secantHistory, errSecant := solver.SecantMethod(a, b, eps)
	if errSecant != nil {
		log.Error(fmt.Sprintf("Ошибка метода хорд: %v", errSecant))
	} else {
		log.Success(fmt.Sprintf("Метод хорд:    x = %.*f (Итераций: %d, f(x) = %e)", decimalPlaces, secantRoot, len(secantHistory)-1, solver.F(secantRoot)))
	}

	// 6. Генерация графики
	if errNewton == nil && errSecant == nil {
		log.Info("Генерация анимаций сходимости...")

		// 6.1. Анимация метода Ньютона (касательные)
		if err := visualizer.GenerateNewtonGIF("docs/convergence_newton.gif", solver.F, solver.D1, a, b, newtonHistory); err != nil {
			log.Warn(fmt.Sprintf("Ошибка генерации GIF Ньютона: %v", err))
		} else {
			log.Success("Сохранено: docs/convergence_newton.gif")
		}

		// 6.2. Анимация метода Хорд (секущие)
		if err := visualizer.GenerateSecantGIF("docs/convergence_secant.gif", solver.F, a, b, secantHistory); err != nil {
			log.Warn(fmt.Sprintf("Ошибка генерации GIF Хорд: %v", err))
		} else {
			log.Success("Сохранено: docs/convergence_secant.gif")
		}

		// 6.3. Общая анимация сравнения методов
		if err := visualizer.GenerateCombinedGIF("docs/convergence_combined.gif", solver.F, a, b, newtonHistory, secantHistory); err != nil {
			log.Warn(fmt.Sprintf("Ошибка генерации общего GIF: %v", err))
		} else {
			log.Success("Сохранено: docs/convergence_combined.gif")
		}
	}

	log.Info("---------------------------------------------------------")
	log.Info("Вычисление успешно завершено.")
}

func getDecimalPlaces(eps float64) int {
	if eps <= 0 {
		return 3
	}
	lg := math.Log10(eps)
	if lg < 0 {
		return int(math.Ceil(math.Abs(lg)))
	}
	return 1
}
