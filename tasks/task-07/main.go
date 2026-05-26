package main

import (
	"fmt"
	"os"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-07/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-07/pkg/visualizer"
)

func main() {
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Численное интегрирование методом Симпсона ===")
	log.Info("Вариант 8:")
	log.Info("  Функция: f(x) = (x + 1) / (2 + ln(1 + x^2))")
	log.Info("  Пределы интегрирования: [a, b] = [1.0, 2.0]")
	log.Info("  Заданная точность: eps = 0.01")
	log.Info("---------------------------------------------------------")

	a := 1.0
	b := 2.0
	eps := 0.01

	// Вычисление истории сходимости по правилу Рунге
	history := solver.BuildIntegrationHistory(a, b, eps)

	log.Info("Ход выполнения итерационного процесса:")
	fmt.Printf(" %-5s | %-6s | %-10s | %-12s | %-12s\n", "Итер.", "N", "Шаг h", "Значение I", "Погр. Рунге")
	fmt.Println("---------------------------------------------------------")
	for _, rec := range history {
		errStr := "—"
		if rec.Step > 1 {
			errStr = fmt.Sprintf("%.6f", rec.Error)
		}
		fmt.Printf(" %-5d | %-6d | %-10.6f | %-12.6f | %-12s\n",
			rec.Step, rec.N, rec.H, rec.Value, errStr)
	}
	fmt.Println("---------------------------------------------------------")

	finalRec := history[len(history)-1]
	log.Success(fmt.Sprintf("Интеграл успешно вычислен!"))
	log.Success(fmt.Sprintf("Итоговое значение I = %.6f (при N = %d)", finalRec.Value, finalRec.N))
	log.Success(fmt.Sprintf("Оценка погрешности по Рунге: %.6f (<= eps)", finalRec.Error))

	// Экспорт анимации в docs
	log.Info("---------------------------------------------------------")
	log.Info("Экспорт анимации процесса интегрирования...")
	_ = os.MkdirAll("docs", os.ModePerm)

	path := "docs/simpson_integration.gif"
	if err := visualizer.GenerateSimpsonGIF(path, a, b, history); err != nil {
		log.Error(fmt.Sprintf("Ошибка генерации GIF: %v", err))
	} else {
		log.Success(fmt.Sprintf("Анимация сохранена: %s", path))
	}

	log.Info("Вычисления успешно завершены.")
}
