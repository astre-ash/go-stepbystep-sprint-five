// Package actioninfo предоставляет общий интерфейс для обработки данных
// о физической активности и функцию для массовой обработки этих данных.
package actioninfo

import (
	"fmt"
	"log"
)

// DataParser — интерфейс, который должны реализовывать типы данных о тренировках.
// Parse отвечает за разбор входящей строки с данными.
// ActionInfo возвращает текстовый отчет о результатах активности.
type DataParser interface {
	Parse(string) error
	ActionInfo() (string, error)
}

// Info последовательно обрабатывает срез строк с данными (dataset) с помощью парсера dp.
// Для каждой корректно распарсенной строки выводит результат ActionInfo в консоль.
// В случае ошибок парсинга или расчета — логирует их и переходит к следующему элементу.
func Info(dataset []string, dp DataParser) {
	for i, data := range dataset {
		// Попытка разобрать текущую строку данных
		if err := dp.Parse(data); err != nil {
			log.Printf("parsing error at index %d for data %q: %v", i, data, err)
			continue
		}
		// Попытка получить отформатированную информацию об активности
		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("failed to get action info for data %q: %v", data, err)
			continue
		}
		// Вывод результата в стандартный поток
		fmt.Println(info)

	}
}
