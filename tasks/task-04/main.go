package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-04/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-04/pkg/visualizer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Решение системы нелинейных уравнений методами МПИ и Ньютона ===")
	log.Info("Система (Вариант 8):")
	log.Info("  1) 2x - cos(y + 1) = 0")
	log.Info("  2) y + sin(x) = -0.4")
	log.Info("---------------------------------------------------------")

	// 1. Ввод параметров точности
	log.Prompt("Введите точность eps (например, 0.001): ")
	epsStr, _ := reader.ReadString('\n')
	epsStr = strings.TrimSpace(epsStr)
	eps, err := strconv.ParseFloat(epsStr, 64)
	if err != nil || eps <= 0 {
		log.Warn("Некорректная точность. Установлено по умолчанию eps = 0.001")
		eps = 0.001
	}

	// 2. Ввод начальной точки
	log.Prompt("Введите начальное приближение x0 (например, 0.5): ")
	x0Str, _ := reader.ReadString('\n')
	x0, _ := strconv.ParseFloat(strings.TrimSpace(x0Str), 64)

	log.Prompt("Введите начальное приближение y0 (например, -0.8): ")
	y0Str, _ := reader.ReadString('\n')
	y0, _ := strconv.ParseFloat(strings.TrimSpace(y0Str), 64)

	maxIter := 1000
	decimalPlaces := getDecimalPlaces(eps)

	// 3. Проверка условия сходимости
	isConvergent, norm := solver.CheckConvergence(x0, y0)
	log.Info(fmt.Sprintf("Анализ сходимости итерационной схемы в точке (%.2f, %.2f):", x0, y0))
	log.Info(fmt.Sprintf("  Норма Якобиана матрицы отображения ||J_phi||_inf = %.4f", norm))
	if isConvergent {
		log.Success("  Достаточное условие выполнено (норма < 1). МПИ сойдется.")
	} else {
		log.Warn("  Внимание: норма >= 1. Возможна расходимость МПИ!")
	}
	log.Info("---------------------------------------------------------")

	// 4. Расчет методом простой итерации
	log.Info("Запуск вычислений методом простой итерации...")
	mpiHistory, errMpi := solver.SolveSimpleIteration(x0, y0, eps, maxIter)
	if errMpi != nil {
		log.Error(fmt.Sprintf("Ошибка МПИ: %v", errMpi))
	} else {
		res := mpiHistory[len(mpiHistory)-1]
		log.Success(fmt.Sprintf("МПИ:    Корень найден! x = %.*f, y = %.*f (Итераций: %d)",
			decimalPlaces, res.X, decimalPlaces, res.Y, res.Iter))
	}

	// 5. Расчет методом Ньютона
	log.Info("Запуск вычислений методом Ньютона...")
	newtonHistory, errNewton := solver.SolveNewton(x0, y0, eps, maxIter)
	if errNewton != nil {
		log.Error(fmt.Sprintf("Ошибка метода Ньютона: %v", errNewton))
	} else {
		res := newtonHistory[len(newtonHistory)-1]
		log.Success(fmt.Sprintf("Ньютон: Корень найден! x = %.*f, y = %.*f (Итераций: %d)",
			decimalPlaces, res.X, decimalPlaces, res.Y, res.Iter))
	}

	// 6. Генерация красивых анимаций
	log.Info("---------------------------------------------------------")
	log.Info("Экспорт анимаций сходимости...")
	_ = os.Mkdir("docs", os.ModePerm)

	if errMpi == nil {
		path := "docs/convergence_mpi.gif"
		if err := visualizer.GenerateConvergenceGIF(path, "Simple Iteration Method", mpiHistory); err != nil {
			log.Warn(fmt.Sprintf("Ошибка GIF МПИ: %v", err))
		} else {
			log.Success(fmt.Sprintf("Сохранено: %s", path))
		}
	}

	if errNewton == nil {
		path := "docs/convergence_newton.gif"
		if err := visualizer.GenerateConvergenceGIF(path, "Newton Method (Systems)", newtonHistory); err != nil {
			log.Warn(fmt.Sprintf("Ошибка GIF Ньютона: %v", err))
		} else {
			log.Success(fmt.Sprintf("Сохранено: %s", path))
		}
	}
	log.Info("Вычисления успешно завершены.")
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
