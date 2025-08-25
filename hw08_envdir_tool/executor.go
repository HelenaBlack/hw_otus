package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd запускает команду с аргументами (cmd) с переменными окружения из env.
// Возвращает код завершения процесса.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Проверяем, что команда не пуста
	if len(cmd) == 0 {
		return 111
	}
	cmdName := cmd[0]
	// Проверяем, что команда существует в PATH
	if _, err := exec.LookPath(cmdName); err != nil {
		fmt.Fprintf(os.Stderr, "Command not found: %s\n", cmdName)
		return 111
	}
	// Фильтруем пустые аргументы
	var filteredArgs []string
	for _, arg := range cmd[1:] {
		if strings.TrimSpace(arg) != "" {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	command := exec.Command(cmdName, filteredArgs...)

	// Формируем карту переменных окружения на основе текущего окружения процесса
	envMap := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.IndexByte(e, '='); i >= 0 {
			envMap[e[:i]] = e[i+1:]
		}
	}

	// Модифицируем окружение в соответствии с env:
	// - если NeedRemove, удаляем переменную
	// - иначе устанавливаем новое значение
	for k, v := range env {
		if v.NeedRemove {
			delete(envMap, k)
		} else {
			envMap[k] = v.Value
		}
	}

	// Преобразуем карту обратно в срез строк вида "ключ=значение"
	envList := make([]string, 0, len(envMap))
	for k, v := range envMap {
		envList = append(envList, k+"="+v)
	}
	command.Env = envList

	// Перенаправляем стандартные потоки ввода/вывода/ошибок на текущий процесс
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Запускаем команду и обрабатываем возможные ошибки
	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		// Если это ошибка завершения процесса, возвращаем соответствующий код выхода
		if errors.As(err, &exitErr) {
			// Для UNIX-систем: получаем код выхода через ExitStatus()
			if status, ok := exitErr.Sys().(interface{ ExitStatus() int }); ok {
				return status.ExitStatus()
			}
			// Для других случаев: используем ExitCode()
			return exitErr.ExitCode()
		}
		// Для других ошибок выводим сообщение и возвращаем 111
		fmt.Fprintf(os.Stderr, "RunCmd error: %v\n", err)
		return 111
	}
	// Если команда завершилась успешно, возвращаем 0
	return 0
}
