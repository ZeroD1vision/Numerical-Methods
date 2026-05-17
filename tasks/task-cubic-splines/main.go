package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	spline "github.com/ZeroD1vision/Numerical-Methods/tasks/task-cubic-splines/pkg/spline"
)

func main() {
	// определяем, хотим ли вывод в JSON (CI или переменная LOG_JSON)
	useJSON := false
	if strings.ToLower(os.Getenv("LOG_JSON")) == "1" || strings.ToLower(os.Getenv("LOG_JSON")) == "true" {
		useJSON = true
	}
	// если запускаемся в CI, по умолчанию используем JSON-логи
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		useJSON = true
	}
	logger := platform.NewLogger(os.Stdout, useJSON)
	reader := bufio.NewReader(os.Stdin)

	logger.Info("=== Интерполяция кубическими сплайнами ===")
	logger.Info("Демонстрация на опорных точках из лабораторного примера:")
	// Тестовые данные из лабораторного конспекта (ПР.1)
	xNodes := []float64{0.0, 1.0, 2.0, 3.0}
	yNodes := []float64{1.0, 2.0, 1.5, 3.0}

	for i := 0; i < len(xNodes); i++ {
		logger.Info("  Узел %d: X = %.1f, Y = %.1f", i, xNodes[i], yNodes[i])
	}
	logger.Info("---------------------------------------------------------")

	// Построение сплайна
	segments, err := spline.BuildSpline(xNodes, yNodes)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка построения сплайна: %v", err))
		return
	}

	logger.Success("Коэффициенты сегментов сплайна успешно вычислены:")
	for i, seg := range segments {
		logger.Info("  Интервал %d [%.1f, %.1f]:", i+1, seg.XMin, seg.XMax)
		logger.Info("    ã_%d = %6.3f (высота)", i, seg.A)
		logger.Info("    b̃_%d = %6.3f (наклон)", i, seg.B)
		logger.Info("    c̃_%d = %6.3f (кривизна)", i, seg.C)
		logger.Info("    d̃_%d = %6.3f (изменение кривизны)", i, seg.D)
	}
	logger.Info("---------------------------------------------------------")

	// Интерактивный режим интерполяции
	for {
		logger.Prompt("Введите точку X для интерполяции (или 'exit' для выхода): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" || input == "" {
			break
		}

		targetX, err := strconv.ParseFloat(input, 64)
		if err != nil {
			logger.Warn("Некорректный формат числа. Попробуйте снова.")
			continue
		}

		// Поиск подходящего сегмента для targetX
		found := false
		for _, seg := range segments {
			if targetX >= seg.XMin && targetX <= seg.XMax {
				resultY := seg.Interp(targetX)
				logger.Success(fmt.Sprintf("Значение сплайна S(%.4f) = %.4f", targetX, resultY))
				found = true
				break
			}
		}

		if !found {
			logger.Warn(fmt.Sprintf("Точка %.4f выходит за границы интервала сплайна [%.1f, %.1f]", targetX, xNodes[0], xNodes[len(xNodes)-1]))
		}
	}
	logger.Info("Выполнение программы завершено.")
}
