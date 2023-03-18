package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (o *StationsStruct) Get(Name string) (StationDataStruct, error) {
	var SD StationDataStruct
	SQL := `SELECT Identifier, Name, Province, Latitude, Longitude, Elevation, URL FROM StationList WHERE Identifier = ?`
	err := o.SQLdb.QueryRow(SQL, Name).Scan(&SD.Identifier, &SD.Name, &SD.Province, &SD.Latitude, &SD.Longitude, &SD.Elevation, &SD.URL)

	if err == sql.ErrNoRows {
		err = errors.New("Station not found")
	}

	return SD, err
}

func (o *StationsStruct) ParseRuntimeTemplate(Station string, Timestamp time.Time) (string, error) {
	u := RuntimeTemplateStruct{
		Timestamp.Format("20060102"),
		Timestamp.Format("2006-01-02-1504"),
		Timestamp.Format("TimeFormat 20060102T1504Z")}

	ut, err := template.New("foo").Parse(Station)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = ut.Execute(&tpl, u)

	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (o *StationsStruct) GetObservationSpan(Station string, Start time.Time, End time.Time) []ObservationStruct {

	var RC []ObservationStruct
	var Obs = ObservationStruct{Station: Station, Identifier: Station}
	var err error

	SD, err := o.Get(Station)

	if err != nil {
		Obs.Station = Station
		Obs.err = err
		return append(RC, Obs)
	}

	var SC ScrapeStruct
	SC.Stations = o
	SC.URL = SD.URL
	SC.Start = Start
	SC.End = End

	StartCooked := Start.Truncate(time.Hour * 24)
	EndCooked := End.Truncate(time.Hour * 24)
	var Step = time.Hour * 24
	var Day = StartCooked
	for {
		//fmt.Println(Day.Format("2006-01-02 15:04:05"))
		err := SC.Scrape(Day)
		if err != nil {
			var Obse = ObservationStruct{Station: Station, Identifier: Station}
			Obse.err = err
			RC = append(RC, Obse)
		}
		Day = Day.Add(Step)
		if Day.After(EndCooked) {
			break
		}
	}

	for _, v := range SC.ObservationFiles {
		Obs.err = nil

		if strings.HasSuffix(v.URL, "-minute-swob.xml") {
			continue
		}

		ObservationXML, err := HTTPSGet(v.URL, time.Hour*30*24)
		if err != nil {
			Obs.err = err
		} else {
			Obs, err = ParseObservation(ObservationXML)
			Obs.Identifier = Station
			Obs.err = err
		}
		RC = append(RC, Obs)
	}

	return RC
}

func (o *StationsStruct) GetObservation(Station string, Timestamp time.Time, Direct bool) (ObservationStruct, error) {

	var Obs = ObservationStruct{Station: Station, Identifier: Station}
	var err error

	SD, err := o.Get(Station)
	if err != nil {
		Obs.err = err
		return Obs, err
	}

	var FinalURL string

	if Direct {

		if SD.TimeOffset != 0 {
			Timestamp.Add(time.Minute * time.Duration(SD.TimeOffset))
		}

		FinalURL, err = o.ParseRuntimeTemplate(SD.URL, Timestamp)
		if err != nil {
			Obs.err = err
			return Obs, err
		}
	} else {

		var SC *ScrapeStruct

		SC, err = o.Datespan(Timestamp, 4*time.Hour, 4*time.Hour, SD)

		if len(SC.ObservationFiles) > 0 {

			var BestTimeDelta = time.Hour * 24 * 30
			var BestTimeIndex int

			for i, v := range SC.ObservationFiles {
				var TimeDelta = Timestamp.Sub(v.Timestamp)
				if TimeDelta < BestTimeDelta {
					BestTimeDelta = TimeDelta
					BestTimeIndex = i
				}
			}

			FinalURL = SC.ObservationFiles[BestTimeIndex].URL
		}
	}

	if Verbose {
		fmt.Println("FinalURL:", FinalURL)
	}

	if len(FinalURL) == 0 {
		err = errors.New("No observation found")
		Obs.err = err
		return Obs, err
	}

	ObservationXML, err := HTTPSGet(FinalURL, time.Hour*30*24)
	if err != nil {
		Obs.err = err
		return Obs, err
	}

	Obs, err = ParseObservation(ObservationXML)
	Obs.Identifier = Station

	return Obs, err
}

type RuntimeTemplateStruct struct {
	Date          string
	FullTimestamp string
	DFOTimestamp  string
}

func (o *StationsStruct) Datespan(Center time.Time, NegativeRange time.Duration, PositiveRange time.Duration, SD StationDataStruct) (*ScrapeStruct, error) {

	var Start = Center.Add(-NegativeRange)
	var End = Center.Add(PositiveRange)
	if End.After(time.Now()) {
		End = time.Now().UTC()
	}

	var SC ScrapeStruct
	SC.Stations = o
	SC.URL = SD.URL
	SC.Start = Start
	SC.End = End

	Start = Start.Truncate(time.Hour * 24)
	End = End.Truncate(time.Hour * 24)
	var Step = time.Hour * 24
	var Day = Start
	for {
		//fmt.Println(Day.Format("2006-01-02 15:04:05"))
		err := SC.Scrape(Day)
		if err != nil {
			return &SC, err
		}
		Day = Day.Add(Step)
		if Day.After(End) {
			break
		}
	}
	return &SC, nil
}

type ObservationFileStruct struct {
	Timestamp time.Time
	URL       string
}

type ScrapeStruct struct {
	URL              string
	Start            time.Time
	End              time.Time
	Stations         *StationsStruct
	ObservationFiles []ObservationFileStruct
}

func (o *ScrapeStruct) Scrape(Day time.Time) error {

	var err error

	DirectoryURL, err := Stations.ParseRuntimeTemplate(o.URL, Day)
	if err != nil {
		return err
	}

	DirectoryURL = DirectoryURL[0:strings.LastIndex(DirectoryURL, "/")]

	if Verbose {
		fmt.Println("DirectoryURL:", DirectoryURL)
	}

	Listing, err := HTTPSGet(DirectoryURL, time.Minute*5)
	if err != nil {
		return err
	}

	Date := Day.Format("20060102")
	YearNumber := Date[0:4]
	MonthNumber := Date[4:6]
	DayNumber := Date[6:]

	ExpStr := `<a href=\"(` + YearNumber + `-?` + MonthNumber + `-?` + DayNumber + `[T-])([0-9]{4})([Z-].*\.xml)\">`

	Exp, err := regexp.Compile(ExpStr)
	if err != nil {
		return err
	}

	res := Exp.FindAllStringSubmatch(Listing, -1)
	for i := range res {

		HH, err := strconv.Atoi(res[i][2][0:2])
		if err != nil {
			return err
		}
		MM, err := strconv.Atoi(res[i][2][2:])
		if err != nil {
			return err
		}

		T := time.Date(Day.Year(), Day.Month(), Day.Day(), HH, MM, 0, 0, time.UTC)

		if !(T.Before(o.Start) || T.After(o.End)) {
			var O ObservationFileStruct
			O.Timestamp = T
			O.URL = DirectoryURL + "/" + res[i][1] + res[i][2] + res[i][3]

			o.ObservationFiles = append(o.ObservationFiles, O)
		}

	}

	return err

}
