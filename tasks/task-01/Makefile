# Конфигурация цветового вывода (проверка NO_COLOR и флага colors=no)
ifneq ($(colors),no)
ifneq ($(NO_COLOR),1)
	COLOR_RESET   := \033[0m
	COLOR_GREEN   := \033[32m
	COLOR_YELLOW  := \033[33m
	COLOR_CYAN    := \033[36m
	COLOR_MAGENTA := \033[35m
endif
endif

# Имя результирующего бинарного файла
BINARY_NAME=numerical-methods-task01
# Директория для сборки артефактов (корень монорепозитория)
BUILD_DIR=../../bin

.PHONY: all build run test fmt lint clean help

all: help

## build: Сборка оптимизированного бинарного файла в общий каталог bin
build: fmt test
	@printf "$(COLOR_CYAN)[make] Compiling static binary to $(BUILD_DIR)/$(BINARY_NAME)...$(COLOR_RESET)\n"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@printf "$(COLOR_GREEN)[make] Build static binary successful!$(COLOR_RESET)\n"

## run: Корректный запуск приложения из исходного кода (доступен флаг json=1)
run: fmt
	@printf "$(COLOR_CYAN)[make] Executing application...$(COLOR_RESET)\n"
	@if [ "$(json)" = "1" ]; then \
		go run main.go -json; \
	else \
		go run main.go; \
	fi

## fmt: Форматирование исходного кода текущей задачи по канонам Go
fmt:
	@printf "$(COLOR_MAGENTA)[make] Formatting source code...$(COLOR_RESET)\n"
	@go fmt ./...

## clean: Очистка скомпилированных артефактов текущей задачи из папки bin
clean:
	@printf "$(COLOR_YELLOW)[make] Purging build artifacts...$(COLOR_RESET)\n"
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	@go clean
	@printf "$(COLOR_GREEN)[make] Clean workflow complete.$(COLOR_RESET)\n"

## help: Вывод текущей интерактивной справки по командам проекта
help:
	@printf "$(COLOR_CYAN)=== Доступные инженерные команды ===$(COLOR_RESET)\n"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | sort | awk -F: '{printf "$(COLOR_GREEN)%-10s$(COLOR_RESET) %s\n", $$1, $$2}'