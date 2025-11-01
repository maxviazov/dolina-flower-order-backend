package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Анализируем файл заказа
	analyzeFile("testfiles/VVA KENYA АМС 02.11 по 03.11 Строго.xlsx")

	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// Анализируем мастер-таблицу
	analyzeFile("testfiles/02.10 -05.11 BLANK 2025 AMS.xlsx")
}

func analyzeFile(filename string) {
	fmt.Printf("Анализ файла: %s\n", filename)
	fmt.Println(strings.Repeat("-", 50))

	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v", filename, err)
		return
	}
	defer f.Close()

	// Получаем список листов
	sheets := f.GetSheetList()
	fmt.Printf("Листы в файле: %v\n", sheets)

	// Анализируем каждый лист
	for _, sheetName := range sheets {
		fmt.Printf("\n--- Лист: %s ---\n", sheetName)

		// Получаем все строки
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Ошибка чтения листа %s: %v", sheetName, err)
			continue
		}

		if len(rows) == 0 {
			fmt.Println("Лист пустой")
			continue
		}

		// Показываем первые 10 строк для анализа структуры
		maxRows := len(rows)
		if maxRows > 10 {
			maxRows = 10
		}

		for i, row := range rows[:maxRows] {
			fmt.Printf("Строка %d: ", i+1)
			for j, cell := range row {
				if j > 10 { // Ограничиваем количество колонок
					fmt.Print("...")
					break
				}
				if cell != "" {
					fmt.Printf("[%d]='%s' ", j+1, cell)
				}
			}
			fmt.Println()
		}

		fmt.Printf("Всего строк в листе: %d\n", len(rows))
		if len(rows) > 0 {
			fmt.Printf("Максимальное количество колонок: %d\n", len(rows[0]))
		}
	}
}
