package api

import (
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/service"
)

// nextDateHandler - HTTP обработчик для вычисления следующей даты (/api/nextdate)
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// nowString - строка с текущей датой в формате YYYYMMDD
	nowString := r.URL.Query().Get("now")
	// dateString - строка начальной даты в формате YYYYMMDD
	dateString := r.URL.Query().Get("date")
	// repeatString - строка с правилом повторения
	repeatString := r.URL.Query().Get("repeat")

	// определение текущей даты
	// если nowString пустая, то текущая дата берется из системы
	var now time.Time
	if nowString == "" {
		t := time.Now()
		now = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowString)
		if err != nil {
			http.Error(w, fmt.Sprintf("некорректная текущая дата %q: %v", nowString, err), http.StatusBadRequest)
			return
		}
	}

	// проверка обязательных параметров
	if dateString == "" || repeatString == "" {
		http.Error(w, "пропущена дата или правило повторения", http.StatusBadRequest)
		return
	}

	// вычисление следующей даты
	next, err := service.NextDate(now, dateString, repeatString)
	if err != nil {
		fmt.Printf("ошибка вычисления следующей даты: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// запись результата в ResponseWriter
	w.Write([]byte(next))
}
