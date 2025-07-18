package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату на основе текущей даты now, начальной даты dstart и правила повторения repeat
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата начала %s: %w", dstart, err)
	}

	// парсинг правил повторения
	repeatParts := strings.Split(repeat, " ")
	if len(repeatParts) == 0 {
		return "", fmt.Errorf("некорректный формат правила повторения: %s", repeat)
	}

	// выбор функции обработки в зависимости от типа повторения
	switch repeatParts[0] {
	case "d":
		return nextDateDayRepeat(now, start, repeatParts)
	case "y":
		return nextDateYearRepeat(now, start, repeatParts)
	case "w":
		return nextDateWeekRepeat(now, start, repeatParts)
	case "m":
		return nextDateMonthRepeat(now, start, repeatParts)
	default:
		return "", fmt.Errorf("некорректный формат правила повторения: %s", repeatParts[0])
	}
}

func nextDateDayRepeat(now, start time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) != 2 {
		return "", fmt.Errorf("некорректный формат правила повторения: %s", strings.Join(repeatParts, " "))
	}
	daysInterval, err := strconv.Atoi(repeatParts[1])
	if err != nil || daysInterval < 1 || daysInterval > 400 {
		return "", fmt.Errorf("неккоректный интервал повторения %q: должен быть от 1 до 400", repeatParts[1])
	}
	nextDate := start.AddDate(0, 0, daysInterval)
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(0, 0, daysInterval)
	}
	return nextDate.Format(DateFormat), nil
}

func nextDateYearRepeat(now, start time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) != 1 {
		return "", fmt.Errorf("некорректный формат правила повторения: %s", strings.Join(repeatParts, " "))
	}
	nextDate := start.AddDate(1, 0, 0)
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format(DateFormat), nil
}

func nextDateWeekRepeat(now, start time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) != 2 {
		return "", fmt.Errorf("некорректный формат правила повторения: %s", strings.Join(repeatParts, " "))
	}
	days := strings.Split(repeatParts[1], ",")
	var isValidWeekday [8]bool
	for _, dayString := range days {
		dayNum, err := strconv.Atoi(dayString)
		if err != nil || dayNum < 1 || dayNum > 7 {
			return "", fmt.Errorf("некорректный день недели: %s", dayString)
		}
		isValidWeekday[dayNum] = true
	}

	// если start > now, то возвращается start, иначе ищется ближайшая дата now + 1
	var nextDate time.Time
	if afterNow(start, now) {
		nextDate = start
	} else {
		nextDate = now.AddDate(0, 0, 1)
	}

	for {
		weekday := int(nextDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if isValidWeekday[weekday] {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

func nextDateMonthRepeat(now, start time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) < 2 || len(repeatParts) > 3 {
		return "", fmt.Errorf("некорректный формат правила повторения: %s", strings.Join(repeatParts, " "))
	}

	// парсинг дней месяца
	days := strings.Split(repeatParts[1], ",")
	var isValidDay [32]bool
	var allowLast, allowSecondLast bool
	for _, dayString := range days {
		dayNum, err := strconv.Atoi(dayString)
		if err != nil {
			return "", fmt.Errorf("некорректный день месяца: %s", dayString)
		}
		switch {
		case dayNum >= 1 && dayNum <= 31:
			isValidDay[dayNum] = true
		case dayNum == -1:
			allowLast = true
		case dayNum == -2:
			allowSecondLast = true
		default:
			return "", fmt.Errorf("некорректный день месяца: %s", dayString)
		}
	}

	// парсинг месяцев
	var applyMonthRepeat bool
	var isValidMonth [13]bool
	if len(repeatParts) == 3 {
		applyMonthRepeat = true
		months := strings.Split(repeatParts[2], ",")
		for _, monthString := range months {
			monthNum, err := strconv.Atoi(monthString)
			if err != nil || monthNum < 1 || monthNum > 12 {
				return "", fmt.Errorf("некорректный месяц: %s", monthString)
			}
			isValidMonth[monthNum] = true
		}
	}

	// поиск первой подходящей даты
	var nextDate time.Time
	if afterNow(start, now) {
		nextDate = start
	} else {
		nextDate = now.AddDate(0, 0, 1)
	}

	for {
		y, month, day := nextDate.Date()
		if !applyMonthRepeat || isValidMonth[int(month)] {
			lastDay := time.Date(y, month+1, 0, 0, 0, 0, 0, nextDate.Location()).Day()
			secondLast := lastDay - 1

			if isValidDay[day] || (allowLast && day == lastDay) || (allowSecondLast && day == secondLast) {
				return nextDate.Format(DateFormat), nil
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

// afterNow возвращает true, если date > now (игнорируя время)
func afterNow(date, now time.Time) bool {
	return now.Format(DateFormat) < date.Format(DateFormat)
}
