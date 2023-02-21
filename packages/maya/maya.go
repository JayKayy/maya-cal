package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type Request struct {
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       *Maya             `json:"body,omitempty"`
}

type Maya struct {
	Date            string `json:"date"`
	LongCount       string `json:"longCount"`
	Haab            string `json:"haab"`
	Tzolkin         string `json:"tzolkin"`
	Pronounce       string `json:"pronounce"`
	LordOfNight     string `json:"lordOfNights"`
	JulianDayNumber int    `json:"julianDayNumber"`
}

var (
	haabMonths    = map[int]string{0: "Pop", 1: "Wo", 2: "Sip", 3: "Sotz'", 4: "Sek", 5: "Xul", 6: "Yaxk'in", 7: "Mol", 8: "Ch'en", 9: "Yax", 10: "Sak", 11: "Kej", 12: "Mak", 13: "K'ank'in", 14: "Muwan", 15: "Pax", 16: "Kâ€™ayabâ€™", 17: "Kâ€™umkuâ€™", 18: "Wayeb'"}
	tzolkinMonths = map[int]string{0: "Imix", 1: "Ik'", 2: "Ak'b'al", 3: "K'an", 4: "Chikchan", 5: "Kimi", 6: "Manik'", 7: "Lamat", 8: "Muluk", 9: "Ok", 10: "Chuwen", 11: "Eb'", 12: "B'en", 13: "Ix", 14: "Men", 15: "K'ib'", 16: "Kab'an", 17: "Etz'nab'", 18: "Kawak", 19: "Ajaw"}
	daemon        bool
)

func main() {
	flag.BoolVar(&daemon, "d", false, "Run as a daemon")
	day := flag.String("day", "", "Day of the month")
	month := flag.String("month", "", "Month of the year")
	year := flag.String("year", "", "Year")
	flag.Parse()

	if !daemon {
		res, err := Main(Request{Day: *day, Month: *month, Year: *year})
		if err != nil {
			fmt.Println("errored: ", err)
		} else {
			output := fmt.Sprintf("Date: %s\nLongCount:%s\nHaab: %s\nTzolkin: %s\nLord of Nights: %s\n", res.Body.Date, res.Body.LongCount, res.Body.Haab, res.Body.Tzolkin, res.Body.LordOfNight)
			fmt.Printf(output)
		}
		return
	}
	// Otherwise Main used as DO Function call
}

// Help from https://www.omnicalculator.com/other/mayan-calendar#how-to-convert-a-date-to-the-long-count-calendar

// Calculators tools
// https://utahgeology.com/bin/maya-calendar-converter/
// https://maya.nmai.si.edu/calendar/maya-calendar-converter

func Main(req Request) (*Response, error) {
	//hash of month and day lengths for validation
	//	monthLengths := map[int]int{0: 31, 1: 29, 2: 31, 3: 30, 4: 31, 5: 30, 6: 31, 7: 31, 8: 30, 9: 31, 10: 30, 11: 31}
	//	months := map[int]string{0: "January", 1: "February", 2: "March", 3: "April", 4: "May", 5: "June", 6: "July", 7: "August", 8: "September", 9: "October", 10: "November", 11: "December"}

	// Used for local testing
	if req.Day == "" || req.Month == "" || req.Year == "" {
		return nil, errors.New("'day', 'month', and 'year' variables are required")
	}

	d, errDay := strconv.Atoi(req.Day)
	m, errMonth := strconv.Atoi(req.Month)
	y, errYear := strconv.Atoi(req.Year)

	if errDay != nil || errMonth != nil || errYear != nil {
		return nil, fmt.Errorf("error parsing day, month, and year as integers")
	}
	// d := 19
	// m := 10
	// y := 1991

	// The algorithm is valid for all Gregorian calendar dates starting on March 1, 4801 BC (astronomical year -4800) at noon UT.[14]
	if y < -4800 || y > 4000 {
		return nil, fmt.Errorf("year must be between -4800 and 4000")
	}

	maya := &Maya{}
	maya.Date = fmt.Sprintf("%d/%d/%d", d, m, y)
	// JDN of Mayan creation date = 584,238
	maya.JulianDayNumber = julianDayNumber(d, m, y, 584238)
	maya.LongCount, maya.Pronounce = longCountCalc(maya.JulianDayNumber)
	maya.Haab = haabCalc(maya.JulianDayNumber)
	maya.Tzolkin = tzolkinCalc(maya.JulianDayNumber)
	maya.LordOfNight = lordOfNightsCalc(maya.JulianDayNumber)

	return &Response{StatusCode: 200, Body: maya}, nil
}

func lordOfNightsCalc(jdn int) string {
	return fmt.Sprintf("G%d", (jdn%9)+1)
}

func tzolkinCalc(jdn int) string {
	// Tzolk'in date calculation
	tzNum := (jdn + 4) % 13
	tzDay := (jdn + 19) % 20
	return fmt.Sprintf("%d %s", tzNum, tzolkinMonths[tzDay])
}

func haabCalc(jdn int) string {
	haabMo := ((jdn - 17) % 365) / 20
	haabDay := ((jdn - 17) % 365) % 20

	return fmt.Sprintf("%d %s", haabDay, haabMonths[haabMo])
}

func longCountCalc(jdn int) (string, string) {
	// Long count calculation
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
	pronunciation := fmt.Sprintf("%d b'ak'tun %d ka'tun %d tun %d uinal %d k'in", bak, kat, tun, uin, kin)

	return longCount, pronunciation
}

func julianDayNumber(d, m, y, correlationConstant int) (jdn int) {
	// The Julian day number (JDN) is the number of days between a specific date and the 1ˢᵗ of January, 4713 B.C.E.
	// https://calendars.fandom.com/wiki/Julian_day#Converting_Gregorian_calendar_date_to_Julian_Day_Number
	// yy and mm are year and month
	alpha := (14 - m) / 12
	yy := y + 4800 - alpha
	mm := m + (12 * alpha) - 3

	jdn = d + (153*mm+2)/5 + 365*yy + (yy / 4) - (yy / 100) + (yy / 400) - 32045
	return jdn - correlationConstant
}

// TODO Maybe setup a Maya.Parse() to compute all given the date?
// TODO run computations in parallel when applicable
// TODO may Maya struct the receiver of the methods?
