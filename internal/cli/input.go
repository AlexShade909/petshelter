package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var reader *bufio.Reader
var writer io.Writer

func Init(r io.Reader, w io.Writer) {
	reader = bufio.NewReader(r)
	writer = w
}

func Print(s ...any) {
	fmt.Fprint(writer, any(s))
}

func Println(s ...any) {
	fmt.Fprintln(writer, any(s))
}
func Printf(format string, args ...interface{}) {
	fmt.Fprintf(writer, format, args...)
}

func ValidationReadMenuChoice(min, max int) int {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Println("Ошибка чтения ввода, попробуйте снова")
			continue
		}
		line = strings.TrimSpace(line)
		choice, err := strconv.Atoi(line)
		if err != nil {
			Println("Введите целое число")
			continue
		}
		if choice < min || choice > max {
			Printf("Число должно быть от %d до %d\n", min, max)
			continue
		}
		return choice
	}
}

func InputValidString() string {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Println("Ошибка чтения ввода, попробуйте снова")
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			Println("Строка не может быть пустой")
			continue
		}
		return line
	}
}

func InputValidFloat() float64 {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Println("Ошибка чтения ввода, попробуйте снова")
			continue
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, ",", ".") // на случай ввода "7,1" вместо "7.1"
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			Println("Введите число, например 7.1")
			continue
		}
		return value
	}
}

func InputValidInt() int {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Println("Ошибка чтения ввода, попробуйте снова")
			continue
		}
		line = strings.TrimSpace(line)
		value, err := strconv.Atoi(line)
		if err != nil {
			Println("Введите целое число")
			continue
		}
		return value
	}
}

func ValidYesNo() bool {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Println("Ошибка чтения ввода, попробуйте снова")
			continue
		}
		line = strings.ToLower(strings.TrimSpace(line))
		switch line {
		case "да", "д", "y", "yes":
			return true
		case "нет", "н", "n", "no":
			return false
		default:
			Println("Введите 'да' или 'нет'")
		}
	}
}
