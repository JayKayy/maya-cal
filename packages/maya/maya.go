package main

import (
	"errors"
	"fmt"
)

type Request struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// Help from https://www.omnicalculator.com/other/mayan-calendar#how-to-convert-a-date-to-the-long-count-calendar

// Calculators tools
// https://utahgeology.com/bin/maya-calendar-converter/
// https://maya.nmai.si.edu/calendar/maya-calendar-converter

func Main(req Request) (*Response, error) {
	//hash of month and day lengths for validation
	//	monthLengths := map[int]int{0: 31, 1: 29, 2: 31, 3: 30, 4: 31, 5: 30, 6: 31, 7: 31, 8: 30, 9: 31, 10: 30, 11: 31}
	//	months := map[int]string{0: "January", 1: "February", 2: "March", 3: "April", 4: "May", 5: "June", 6: "July", 7: "August", 8: "September", 9: "October", 10: "November", 11: "December"}

	haabMonths := map[int]string{0: "Pop", 1: "Wo", 2: "Sip", 3: "Sotz'", 4: "Sek", 5: "Xul", 6: "Yaxk'in", 7: "Mol", 8: "Ch'en", 9: "Yax", 10: "Sak", 11: "Kej", 12: "Mak", 13: "K'ank'in", 14: "Muwan", 15: "Pax", 16: "Kâ€™ayabâ€™", 17: "Kâ€™umkuâ€™", 18: "Wayeb'"}
	tzolkinMonths := map[int]string{0: "Imix", 1: "Ik'", 2: "Ak'b'al", 3: "K'an", 4: "Chikchan", 5: "Kimi", 6: "Manik'", 7: "Lamat", 8: "Muluk", 9: "Ok", 10: "Chuwen", 11: "Eb'", 12: "B'en", 13: "Ix", 14: "Men", 15: "K'ib'", 16: "Kab'an", 17: "Etz'nab'", 18: "Kawak", 19: "Ajaw"}

	// Used for local testing
	if req.Day == 0 || req.Month == 0 {
		return nil, errors.New("'day', 'month', and 'year' variables are required")
	}

	d := req.Day
	m := req.Month
	y := req.Year

	// d := 19
	// m := 10
	// y := 1991

	//	fmt.Println(d, m, y)
	// The algorithm is valid for all Gregorian calendar dates starting on March 1, 4801 BC (astronomical year -4800) at noon UT.[14]
	if y < -4800 || y > 4000 {
		return nil, fmt.Errorf("Bad Year")
	}

	// Calculate the Julian Day Number (JDN)
	// The Julian day number (JDN) is the number of days between a specific date and the 1ˢᵗ of January, 4713 B.C.E.
	// https://calendars.fandom.com/wiki/Julian_day#Converting_Gregorian_calendar_date_to_Julian_Day_Number
	// yy and mm are year and month but
	a := (14 - m) / 12
	yy := y + 4800 - a
	mm := m + (12 * a) - 3

	jdn := d + (153*mm+2)/5 + 365*yy + (yy / 4) - (yy / 100) + (yy / 400) - 32045

	// JDN of Mayan creation date = 584,238
	// Coorelation constant 584,238
	jdn -= 584283

	bak := jdn / 144000
	remain := jdn % 144000
	kat := remain / 7200
	remain = remain % 7200
	tun := remain / 360
	remain = remain % 360
	uin := remain / 20
	remain = remain % 20
	kin := remain

	longCount := fmt.Sprintf("%d.%d.%d.%d.%d", bak, kat, tun, uin, kin)
	p := fmt.Sprintf("%d b'ak'tun %d ka'tun %d tun %d uinal %d k'in", bak, kat, tun, uin, kin)

	// Haab date calculation
	haabMo := (((jdn - 17) % 365) / 20)
	haabDay := (((jdn - 17) % 365) % 20)

	// Tzolk'in date calculation
	tzNum := (jdn + 4) % 13
	tzDay := (jdn + 19) % 20

	haab := fmt.Sprintf("%d %s", haabDay, haabMonths[haabMo])
	tz := fmt.Sprintf("%d %s", tzNum, tzolkinMonths[tzDay])
	loN := fmt.Sprintf("G%d", (jdn%9)+1)

	return &Response{
		StatusCode: 200,
		Body:       fmt.Sprintf("{\"dd-mm-yyyy\": \"%d-%d-%d\", \"longCount\": \"%s\", \"pronounce\": \"%s\",\"haab\": \"%s\",\"tzolk'in\": \"%s\",\"Lord of the Night\": \"%s\"}", d, m, y, longCount, p, haab, tz, loN),
	}, nil
}
