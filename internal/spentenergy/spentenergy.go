// Package spentenergy содержит логику для расчета физических показателей:
// дистанции, средней скорости и сожженных калорий для различных видов активности.
package spentenergy

import (
	"fmt"
	"log"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// validateTrainingInputs проверяет входные параметры на корректность (положительные значения).
func validateTrainingInputs(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return fmt.Errorf("invalid steps: %d (must be > 0 zero)", steps)
	}
	if weight <= 0 {
		return fmt.Errorf("invalid weight: %.2f (must be > 0 zero)", weight)
	}
	if height <= 0 {
		return fmt.Errorf("invalid height: %.2f (must be > 0 zero)", height)
	}
	if duration <= 0 {
		return fmt.Errorf("invalid duration: %v ((must be > 0 zero)", duration)
	}
	return nil
}

// WalkingSpentCalories рассчитывает количество сожженных калорий при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateTrainingInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	v := MeanSpeed(steps, height, duration)
	m := duration.Minutes()
	baseCalories := (weight * v * m) / minInH
	return baseCalories * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories рассчитывает количество сожженных калорий при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateTrainingInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	m := duration.Minutes()
	v := MeanSpeed(steps, height, duration)
	return weight * m * v / minInH, nil
}

// MeanSpeed рассчитывает среднюю скорость движения в км/ч.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		log.Println("warning: duration must be positive to calculate mean speed")
		return 0
	}

	d := Distance(steps, height)
	hours := duration.Hours()
	return d / hours
}

// Distance рассчитывает пройденную дистанцию в километрах на основе количества шагов и роста.
func Distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / mInKm
}
