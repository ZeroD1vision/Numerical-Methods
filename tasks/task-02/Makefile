ifneq ($(colors),no)
ifneq ($(NO_COLOR),1)
    COLOR_RESET   := \033[0m
    COLOR_GREEN   := \033[32m
    COLOR_YELLOW  := \033[33m
    COLOR_CYAN    := \033[36m
    COLOR_MAGENTA := \033[35m
endif
endif

BINARY_NAME=numerical-methods-task02
BUILD_DIR=../../bin

.PHONY: all build run fmt clean help

all: help

## build: Сборка приложения с оптимизацией
build: fmt
	@printf "$(COLOR_CYAN)[make] Building binary to $(BUILD_DIR)/$(BINARY_NAME)...$(COLOR_RESET)\n"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@printf "$(COLOR_GREEN)[make] Build successful!$(COLOR_RESET)\n"

## run: Автоматическое форматирование и запуск проекта
run: fmt
	@printf "$(COLOR_CYAN)[make] Running application...$(COLOR_RESET)\n"
	@go run main.go

## fmt: Строгое форматирование исходного кода
fmt:
	@printf "$(COLOR_MAGENTA)[make] Formatting source code...$(COLOR_RESET)\n"
	@go fmt ./...

## clean: Очистка артефактов компиляции
clean:
	@printf "$(COLOR_YELLOW)[make] Cleaning build artifacts...$(COLOR_RESET)\n"
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	@go clean
	@printf "$(COLOR_GREEN)[make] Clean complete.$(COLOR_RESET)\n"

## help: Отображение доступных команд
help:
	@printf "$(COLOR_CYAN)=== Доступные команды сборки ===$(COLOR_RESET)\n"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | sort | awk -F: '{printf "$(COLOR_GREEN)%-10s$(COLOR_RESET) %s\n", $$1, $$2}'