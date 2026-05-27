package main

import (
	"fmt"
	"math"
	"os"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-08/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-08/pkg/visualizer"
)

func main() {
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Решение задачи Коши (Метод Эйлера и модифицированный) ===")
	log.Info("Вариант 8:")
	log.Info("Дифференциальное уравнение: y'(t) = cos(sqrt(t * y^2))")
	log.Info("Отрезок: t ∈ [1, 2]")
	log.Info("Начальное условие: y(1) = 2 (адаптировано с учетом опечатки в методичке y(0)=2)")
	log.Info("Количество шагов N = 10")
	log.Info("---------------------------------------------------------")

	t0 := 1.0
	y0 := 2.0
	T := 2.0
	N := 10

	// Вычисления
	log.Info("Запуск метода Эйлера...")
	eulerResult := solver.Euler(t0, y0, T, N)

	log.Info("Запуск модифицированного метода Эйлера...")
	modEulerResult := solver.ModifiedEuler(t0, y0, T, N)

	log.Info("Генерация эталонного точного решения (RK4, N=1000)...")
	exactResult := solver.RK4(t0, y0, T, 1000)

	// Вывод результатов в консоль
	log.Info("---------------------------------------------------------")
	log.Info(fmt.Sprintf("%-5s | %-8s | %-15s | %-15s | %-15s", "Шаг", "t", "y (Эйлер)", "y (Мод. Эйлер)", "y (Точное)"))
	log.Info("---------------------------------------------------------")
	for i := 0; i <= N; i++ {
		t := eulerResult[i].T
		ye := eulerResult[i].Y
		yme := modEulerResult[i].Y
		yExact := solver.GetExactY(exactResult, t)
		log.Info(fmt.Sprintf("%-5d | %-8.2f | %-15.6f | %-15.6f | %-15.6f", i, t, ye, yme, yExact))
	}
	log.Info("---------------------------------------------------------")

	exactFinalY := exactResult[len(exactResult)-1].Y
	eulerError := math.Abs(exactFinalY - eulerResult[N].Y)
	modEulerError := math.Abs(exactFinalY - modEulerResult[N].Y)

	log.Success(fmt.Sprintf("Конечные значения при t=%.1f:", T))
	log.Success(fmt.Sprintf("Эйлер:      y(%.1f) = %.6f (Погрешность: %.6f)", T, eulerResult[N].Y, eulerError))
	log.Success(fmt.Sprintf("Мод. Эйлер: y(%.1f) = %.6f (Погрешность: %.6f)", T, modEulerResult[N].Y, modEulerError))

	// Визуализация
	log.Info("---------------------------------------------------------")
	log.Info("Экспорт анимации...")
	_ = os.Mkdir("docs", os.ModePerm)

	path := "docs/ode_cauchy_methods.gif"
	if err := visualizer.GenerateODEGIF(path, eulerResult, modEulerResult, exactResult); err != nil {
		log.Error(fmt.Sprintf("Ошибка при создании GIF: %v", err))
	} else {
		log.Success(fmt.Sprintf("Анимация успешно сохранена: %s", path))
	}

	log.Info("Вычисления успешно завершены.")
}
