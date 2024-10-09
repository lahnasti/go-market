package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Чтение количества тестов
	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var s string
		fmt.Fscan(reader, &s)

		// Длина строки
		n := len(s)

		// Найдем первую цифру, которую имеет смысл удалить
		found := false
		for j := 0; j < n-1; j++ {
			if s[j] < s[j+1] {
				// Удаляем текущую цифру s[j]
				fmt.Fprintln(writer, s[:j]+s[j+1:])
				found = true
				break
			}
		}

		// Если не нашли, то удаляем последнюю цифру
		if !found {
			fmt.Fprintln(writer, s[:n-1])
		}
	}
}
