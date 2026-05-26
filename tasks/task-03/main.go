package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	// Внимание: замените пути ниже на актуальные для вашего go.mod, если они отличаются
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-03/pkg/solver"
	"github.com/ZeroD1vision/Numerical-Methods/tasks/task-03/pkg/visualizer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Решение СЛАУ методом простой итерации (Вариант 8) ===")

	// 1. Вывод достаточного условия сходимости
	rho := solver.SpectralRadius()
	fmt.Printf("Норма матрицы ||C||_inf (приближенный спектральный радиус): %.4f\n", rho)
	if rho < 1 {
		fmt.Println("Условие сходимости выполнено (||C|| < 1). Метод гарантированно сходится.")
	} else {
		fmt.Println("Внимание: Условие сходимости ||C|| < 1 не выполнено. Метод может расходиться.")
	}
	fmt.Println("---------------------------------------------------------")

	// 2. Ввод точности eps
	fmt.Print("Введите точность eps (например, 0.001): ")
	epsStr, _ := reader.ReadString('\n')
	epsStr = strings.TrimSpace(epsStr)
	eps, err := strconv.ParseFloat(epsStr, 64)
	if err != nil || eps <= 0 {
		fmt.Println("Некорректное значение. Установлено eps = 0.001")
		eps = 0.001
	}

	// 3. Начальное приближение x0 = d (согласно теории метода простой итерации)
	x0 := solver.D
	maxIter := 1000

	// 4. Запуск решения
	fmt.Println("Запуск итерационного процесса...")
	root, history, errSolve := solver.Solve(x0, eps, maxIter)
	if errSolve != nil {
		fmt.Printf("Ошибка: %v\n", errSolve)
		return
	}

	// 5. Вывод результатов
	fmt.Println("\nРешение успешно найдено!")
	fmt.Printf("Количество итераций: %d\n", len(history)-1)
	fmt.Printf("Полученный вектор x: [%.6f, %.6f, %.6f, %.6f]\n", root[0], root[1], root[2], root[3])

	residual := solver.Residual(root)
	fmt.Printf("Чебышёвская норма невязки ||Ax-b||: %e\n", residual)
	fmt.Println("---------------------------------------------------------")

	// 6. Генерация графиков
	fmt.Println("Генерация анимаций сходимости в папку docs/...")

	// Создаем папку docs, если её нет
	_ = os.Mkdir("docs", os.ModePerm)

	if err := visualizer.GenerateConvergenceGIF("docs/convergence_components.gif", history); err != nil {
		fmt.Printf("Ошибка генерации GIF компонентов: %v\n", err)
	} else {
		fmt.Println("Сохранено: docs/convergence_components.gif")
	}

	if err := visualizer.GenerateErrorGIF("docs/convergence_error.gif", history, eps); err != nil {
		fmt.Printf("Ошибка генерации GIF ошибки: %v\n", err)
	} else {
		fmt.Println("Сохранено: docs/convergence_error.gif")
	}

	if err := visualizer.GenerateCombinedGIF("docs/convergence_combined.gif", history, eps); err != nil {
		fmt.Printf("Ошибка генерации совмещенного GIF: %v\n", err)
	} else {
		fmt.Println("Сохранено: docs/convergence_combined.gif")
	}

	fmt.Println("Вычисление и визуализация успешно завершены.")
}
