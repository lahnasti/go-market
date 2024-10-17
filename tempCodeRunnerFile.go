package main

import (
	"encoding/binary"
	"fmt"
	"bytes"
)

func main() {
	// Число для записи
	var num uint32 = 0x12345678

	// Буфер для записи
	buf := new(bytes.Buffer)

	// Записываем число в формате Big Endian
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		fmt.Println("Ошибка записи:", err)
	}

	fmt.Println("Big Endian:", buf.Bytes()) // Выведет байты в Big Endian

	// Очистим буфер для записи в другой порядок
	buf.Reset()

	// Записываем число в формате Little Endian
	err = binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("Ошибка записи:", err)
	}

	fmt.Println("Little Endian:", buf.Bytes()) // Выведет байты в Little Endian
}
