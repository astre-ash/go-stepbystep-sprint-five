// Package personaldata содержит структуры и методы для работы с личными
// антропометрическими данными пользователей.
package personaldata

import "fmt"

// Personal хранит основные физические параметры пользователя,
// необходимые для расчета спортивных показателей.
type Personal struct {
	Name   string
	Weight float64
	Height float64
}

// Print выводит антропометрические данные пользователя
// (имя, вес и рост) в стандартный поток вывода в читаемом формате.
func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f кг.\nРост: %.2f м.\n", p.Name, p.Weight, p.Height)
}
