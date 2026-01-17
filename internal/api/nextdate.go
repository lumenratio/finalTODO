package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lumenratio/finalTODO/internal/logger"
)

type repeatData struct {
	date    time.Time
	curTime time.Time
	raw     string
	repType string
	repNum  int
	year    int
	//month   [13]bool
}

const dateFormat string = "20060102"

func (r *repeatData) repeatCheck() bool {
	repeatSlice := strings.Split(r.raw, " ")
	if len(repeatSlice) == 0 {
		return false
	}

	first := repeatSlice[0]
	if len(repeatSlice) == 1 && first == "y" {
		r.repType = "y"
		return true
	}

	if len(repeatSlice) == 2 && first == "d" {
		// check days number
		num, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return false
		}
		if num > 400 || num < 1 {
			return false
		}
		r.repType = "d"
		r.repNum = num
		return true
	}
	return false
}

func (r *repeatData) dateCalc() {
	for {
		r.date = r.date.AddDate(r.year, 0, r.repNum)
		if r.date.After(r.curTime) {
			break
		}
	}
}

func NextDate(nowReq string, dstart string, repeat string) (string, error) {
	now, err := time.Parse(dateFormat, nowReq)
	if err != nil {
		return "", err
	}
	r := repeatData{raw: repeat, curTime: now}
	if !r.repeatCheck() {
		return "", fmt.Errorf("got malformed \"repeat\" parameter in request /api/nextdate: %s", repeat)
	}

	// date parameter convert and check
	// Parse nowReq date to time.Time type and check
	t, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("got malformed \"date\" parameter in request /api/nextdate: %s", dstart)
	}
	r.date = t

	// это надо переделать, выглядит ужасно
	if r.repType == "d" {
		r.dateCalc()
		return r.date.Format(dateFormat), err
	}

	// знаем, что остался только тип "Y"
	// такое завершение - не нравится
	r.year = 1
	r.dateCalc()
	return r.date.Format(dateFormat), err
}

func nextDayHandler(w http.ResponseWriter, req *http.Request) {
	// check parameters
	nowReq, dateReq, repeatReq := req.FormValue("now"), req.FormValue("date"), req.FormValue("repeat")
	if nowReq == "" || dateReq == "" || repeatReq == "" {
		logger.Info.Println(nowReq, dateReq, repeatReq)
		logger.Err.Println("parameters can't be empty or not use. Malformed GET request /api/nextdate:", req.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// вычисляем дату
	newDate, err := NextDate(nowReq, dateReq, repeatReq)
	if err != nil {
		logger.Err.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(newDate))
	if err != nil {
		logger.Err.Fatalln(err)
	}
}
