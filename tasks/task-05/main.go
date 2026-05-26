package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-05/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-05/pkg/visualizer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Нахождение собственных значений матрицы (Метод степенных итераций) ===")
	log.Info("Вариант 8: α = -7")
	log.Info("Матрица A:")
	log.Info("  |  1   2  -7 |")
	log.Info("  |  2   3   4 |")
	log.Info("  | -7   4   5 |")
	log.Info("---------------------------------------------------------")

	log.Prompt("Введите точность eps (например, 0.0001): ")
	epsStr, _ := reader.ReadString('\n')
	epsStr = strings.TrimSpace(epsStr)
	eps, err := strconv.ParseFloat(epsStr, 64)
	if err != nil || eps <= 0 {
		log.Warn("Некорректная точность. Установлено по умолчанию eps = 0.0001")
		eps = 0.0001
	}

	A := solver.BuildMatrix()

	// Метод степенных итераций
	log.Info("Запуск метода степенных итераций...")
	history := solver.PowerIteration(A, eps, 1000)
	last := history[len(history)-1]

	log.Success(fmt.Sprintf("Наибольшее по модулю СЗ: λ₁ = %.6f (итераций: %d)", last.Eigenvalue, last.Iter))
	log.Info(fmt.Sprintf("  Собственный вектор: [%.5f, %.5f, %.5f]", last.Vector[0], last.Vector[1], last.Vector[2]))
	log.Info(fmt.Sprintf("  Невязка ||A·v - λ·v|| = %.2e", last.Error))

	// Все собственные значения (метод Якоби)
	log.Info("---------------------------------------------------------")
	log.Info("Все собственные значения (метод Якоби):")
	allEvals := solver.AllEigenvalues(A)
	for i, ev := range allEvals {
		log.Success(fmt.Sprintf("  λ%d = %.6f", i+1, ev))
	}

	// Проверка: след матрицы = сумма СЗ, детерминант = произведение
	traceA := 1.0 + 3.0 + 5.0
	traceEvals := allEvals[0] + allEvals[1] + allEvals[2]
	log.Info("---------------------------------------------------------")
	log.Info(fmt.Sprintf("Проверка: tr(A) = %.1f,  Σλᵢ = %.4f  (погрешность: %.2e)",
		traceA, traceEvals, math.Abs(traceA-traceEvals)))

	// Визуализация
	log.Info("---------------------------------------------------------")
	log.Info("Экспорт анимации сходимости...")
	_ = os.Mkdir("docs", os.ModePerm)

	path := "docs/convergence_power_iteration.gif"
	if err := visualizer.GenerateEigenGIF(path, history, allEvals); err != nil {
		log.Error(fmt.Sprintf("Ошибка GIF: %v", err))
	} else {
		log.Success(fmt.Sprintf("Сохранено: %s", path))
	}

	log.Info("Вычисления успешно завершены.")
}
