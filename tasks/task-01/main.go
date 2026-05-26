package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	platform "github.com/ZeroD1vision/Numerical-Methods/internal/logger"
	solver "github.com/ZeroD1vision/Numerical-Methods/tasks/task-01/pkg/dichotomy"
)

// main является легковесной точкой входа: он отвечает лишь за ввод/вывод и оркестрацию.
// Бизнес-логика (поиск интервала, вычисление корня) вынесена в пакет `solver`.
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

	// создаём платформенный логгер, который сам решит, включать ли цвета
	logger := platform.NewLogger(os.Stdout, useJSON)

	// интерактивный ввод через stdin
	reader := bufio.NewReader(os.Stdin)

	logger.Info("=== Поиск корня методом дихотомии ===")
	logger.Info("Уравнение: sin(x + pi/3) - 0.5x = 0")
	logger.Info("---------------------------------------------------------")

	// 1) ввод точности
	logger.Prompt("Введите точность eps (например, 0.001): ")
	epsStr, _ := reader.ReadString('\n')
	epsStr = strings.TrimSpace(epsStr)
	eps, err := strconv.ParseFloat(epsStr, 64)
	if err != nil || eps <= 0 {
		logger.Warn("Некорректное значение точности. Установлено eps = 0.001")
		eps = 0.001
	}

	// 2) ввод границ (опционально)
	logger.Prompt("Введите левую границу 'a' (или нажмите Enter для автопоиска): ")
	aStr, _ := reader.ReadString('\n')
	aStr = strings.TrimSpace(aStr)

	var a, b float64
	useAutoSearch := false

	if aStr == "" {
		useAutoSearch = true
	} else {
		logger.Prompt("Введите правую границу 'b': ")
		bStr, _ := reader.ReadString('\n')
		bStr = strings.TrimSpace(bStr)

		var errA, errB error
		a, errA = strconv.ParseFloat(aStr, 64)
		b, errB = strconv.ParseFloat(bStr, 64)
		if errA != nil || errB != nil {
			logger.Warn("Некорректный ввод границ. Включён автопоиск...")
			useAutoSearch = true
		} else if a >= b {
			logger.Warn("Левая граница должна быть строго меньше правой. Включён автопоиск...")
			useAutoSearch = true
		}
	}

	// 3) локализация интервала
	if useAutoSearch {
		logger.Info("Запуск автоматического поиска интервала с изменением знака...")
		var errInterval error
		a, b, errInterval = solver.FindInterval(0.0, 1.0, 1000)
		if errInterval != nil {
			logger.Error(fmt.Sprintf("Критическая ошибка: %v", errInterval))
			return
		}
		logger.Success(fmt.Sprintf("Интервал успешно найден автоматически: [%.4f, %.4f]", a, b))
	} else {
		if solver.F(a)*solver.F(b) >= 0 {
			logger.Warn(fmt.Sprintf("На отрезке [%.4f, %.4f] функция не меняет знак (f(a)=%.4f, f(b)=%.4f).", a, b, solver.F(a), solver.F(b)))
			logger.Info("Метод дихотомии не может быть применён без изменения знака.")
			return
		}
	}

	// 4) вычисления — вызываем экспортированную функцию Bisection
	logger.Info("Запуск вычислений методом дихотомии...")
	root, steps := solver.DichotomyMethod(a, b, eps)

	// 5) вывод результатов
	logger.Info("---------------------------------------------------------")
	logger.Success(fmt.Sprintf("Корень найден: x = %.*f", getDecimalPlaces(eps), root))
	logger.Info(fmt.Sprintf("Значение функции в корне: f(x) = %e", solver.F(root)))
	logger.Info(fmt.Sprintf("Количество итераций: %d", steps))
}

// вспомогательная функция: определяет количество знаков после запятой
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
