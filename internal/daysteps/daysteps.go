// Package daysteps предоставляет функционал для обработки данных о ежедневной
// активности пользователя (количество шагов за определенный период).
package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/actioninfo"
	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Проверка реализации интерфейса actioninfo.DataParser типом *DaySteps.
var _ actioninfo.DataParser = (*DaySteps)(nil)

// DaySteps представляет структуру, которая хранит информацию о пройденных шагах за определенное время
// и личные данные пользователя для расчетов.
type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

// Parse анализирует строку формата "шаги,длительность" (например, "10000,2h30m").
// Возвращает ошибку, если строка пуста, имеет неверный формат или неположительные значения.
func (ds *DaySteps) Parse(datastring string) (err error) {
	if datastring == "" {
		return errors.New("input string cannot be empty")
	}

	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: expected 2 parts (steps,duration), but got %d in %q", len(parts), datastring)
	}

	// Парсинг количества шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("faile to parse steps: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("invalid steps: %d (must be > 0)", steps)
	}

	// Парсинг длительности (например, "1h", "30m")
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("invalid duration: %v (must be > 0)", duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// ActionInfo рассчитывает пройденную дистанцию и сожженные калории на основе шагов.
// Возвращает отформатированный отчет в виде строки.
func (ds DaySteps) ActionInfo() (string, error) {
	if ds.Steps <= 0 {
		return "", fmt.Errorf("invalid steps: %d (must be > 0)", ds.Steps)
	}

	// Расчет дистанции и калорий через пакет spentenergy
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("failed to calculate calories: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories), nil
}
