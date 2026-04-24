// Package trainings предоставляет инструменты для работы с данными о тренировках,
// включая парсинг входных данных и расчет показателей активности.
package trainings

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

// Проверка реализации интерфейса actioninfo.DataParser типом *Training.
var _ actioninfo.DataParser = (*Training)(nil)

// Training представляет структуру данных о тренировке.
// Содержит количество шагов, тип активности, длительность и антропометрические данные пользователя.
type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

// Parse анализирует входную строку формата "шаги,тип,длительность" (например, "5000,Бег,1h30m").
// Возвращает ошибку, если формат данных некорректен или значения невалидны.
func (t *Training) Parse(datastring string) (err error) {
	if datastring == "" {
		return errors.New("input string cannot be empty")
	}

	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("invalid format: expected 3 parts (steps, duration, trainingType), but got %d in %q", len(parts), datastring)

	}

	// Парсинг количества шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	if steps <= 0 {
		return fmt.Errorf("invalid steps: %d (must be > 0)", steps)
	}

	// Парсинг длительности (например, "1h", "30m")
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return err
	}
	if duration <= 0 {
		return fmt.Errorf("invalid duration: %v (must be > 0)", duration)
	}

	// Тип физической активности
	t.TrainingType = parts[1]
	if t.TrainingType == "" {
		return errors.New("training type is required")
	}
	t.Steps = steps
	t.Duration = duration

	return nil
}

// ActionInfo рассчитывает метрики тренировки (дистанцию, скорость, калории)
// и возвращает отформатированную строку с результатами.
// Поддерживаемые типы тренировок: "Бег", "Ходьба".
func (t Training) ActionInfo() (string, error) {
	d := spentenergy.Distance(t.Steps, t.Height)
	v := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)
	h := t.Duration.Hours()

	var calories float64
	var err error

	// Выбор алгоритма расчета калорий в зависимости от типа активности
	switch t.TrainingType {

	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)

	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)

	default:
		return "", fmt.Errorf("unknown activity type: %v", t.TrainingType)
	}

	if err != nil {
		return "", fmt.Errorf("error calculating calories for %v: %w", t.TrainingType, err)
	}

	// Формирование итогового отчета
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, h, d, v, calories), nil
}
