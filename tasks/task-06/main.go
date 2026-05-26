package main

import (
	"fmt"
	"os"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-06/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-06/pkg/visualizer"
)

func main() {
	log := platform.NewLogger(os.Stdout, false)

	log.Info("=== Интерполяционный многочлен Лагранжа ===")
	log.Info("Вариант 8:")
	log.Info("  Узлы: x = [0, 1, 2, 3],  y = [11, 12, 13, 14]")
	log.Info(fmt.Sprintf("  Точка интерполяции: x* = %.1f", solver.Variant8Target))
	log.Info("---------------------------------------------------------")

	points := solver.Variant8Data()
	xTarget := solver.Variant8Target

	// Вывод таблицы узлов
	log.Info("Узловые точки:")
	for _, p := range points {
		log.Info(fmt.Sprintf("  x = %.1f,  y = %.1f", p.X, p.Y))
	}
	log.Info("---------------------------------------------------------")

	// Результат интерполяции
	result := solver.LagrangeInterpolate(points, xTarget)
	log.Success(fmt.Sprintf("L(x* = %.1f) = %.6f", xTarget, result))

	// Вывод каждого базисного полинома
	log.Info("---------------------------------------------------------")
	log.Info("Значения базисных полиномов Лагранжа в точке x*:")
	for i, p := range points {
		li := solver.LagrangeBasis(points, i, xTarget)
		log.Info(fmt.Sprintf("  L%d(%.1f) = %.6f   (узел: x_%d = %.1f, y_%d = %.1f)",
			i, xTarget, li, i, p.X, i, p.Y))
	}

	// История для анимации
	history := solver.BuildConvergenceHistory(points, xTarget)
	log.Info("---------------------------------------------------------")
	log.Info("Последовательное добавление узлов:")
	for _, rec := range history {
		log.Info(fmt.Sprintf("  k=%d узл(а): L(x*) = %.6f", rec.Used, rec.Value))
	}

	// Визуализация
	log.Info("---------------------------------------------------------")
	log.Info("Экспорт анимации интерполяции...")
	_ = os.Mkdir("docs", os.ModePerm)

	path := "docs/lagrange_interpolation.gif"
	if err := visualizer.GenerateLagrangeGIF(path, points, xTarget, history); err != nil {
		log.Error(fmt.Sprintf("Ошибка GIF: %v", err))
	} else {
		log.Success(fmt.Sprintf("Сохранено: %s", path))
	}

	log.Info("Вычисления успешно завершены.")
}
