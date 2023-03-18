package main

import (
	"fmt"
	"math"
	"strings"
)

func OutTableText(ObsList []ObservationStruct, ShowIdentifiers bool, RHAT string) error {

	var TZ string

	if UTC {
		TZ = "UTC"
	} else {
		TZ = LocalLoc.String()
	}

	var Columns = []ColumnStruct{
		{Header: []string{"Station", "Name", ""}, GenStatsTitles: true},
		{Header: []string{"Station", "Identifier", ""}, Suppress: !ShowIdentifiers},
		{Header: []string{"Observation", "Time", TZ}},
		{Header: []string{"Min", "Temp", "°C"}},
		{Header: []string{"Avg", "Temp", "°C"}},
		{Header: []string{"Max", "Temp", "°C"}},
		{Header: []string{"Rel", "Humid", "%"}},
		{Header: []string{"RH at", RHAT + " °C", "%"}, Suppress: len(RHAT) == 0},
		{Header: []string{"Barr", "Press", "hPa"}},
		{Header: []string{"Wet", "Bulb", "°C"}},
		{Header: []string{"Dew", "Point", "°C"}},
		{Header: []string{"Dew", "Point Δ", "°C"}},
		{Header: []string{"Humidex", "Max", "°C"}, HideIfEmpty: true},
		{Header: []string{"Windchill", "Max", "°C"}, HideIfEmpty: true},
		{Header: []string{"Precip", "Rate", "mm/hr"}, CalcTotal: true, HideIfEmpty: true},
		{Header: []string{"Wind", "Speed", "km/h"}},
		{Header: []string{"Gust", "Speed", "km/h"}},
		{Header: []string{"Wind", "Dir", "°"}},
		{Header: []string{"Remarks", "", ""}, HideIfEmpty: true, LeftAlign: true},
	}

	for _, Observation := range ObsList {
		DV := GenerateDerivedValues(&Observation, RHAT)

		Row := []ValueStruct{{String: Observation.Station},
			{String: Observation.Identifier},
			{String: DV.Timestamp},
			Observation.MinTemperature,
			Observation.Temperature,
			Observation.MaxTemperature,
			Observation.Humidity,
			DV.RHAT,
			Observation.Pressure,
			Observation.WetBulbTemperature,
			Observation.DewPoint,
			DV.DewPointDelta,
			Observation.Humidex,
			Observation.Windchill,
			Observation.Precipitation,
			Observation.AverageWindSpeed,
			Observation.PeakWindSpeed,
			DV.WindDirection,
			{String: DV.Remarks}}

		for i, v := range Row {
			Columns[i].Add(v)
		}
	}

	if len(ObsList) > 1 {
		for i := range Columns {
			Columns[i].GenStats = true
			Columns[i].Finalize()
		}
		Columns[len(Columns)-1].GenStats = false
	} else {
		for i := range Columns {
			Columns[i].Finalize()
		}
	}

	for r := range Columns[0].Cells {
		var sb strings.Builder
		for _, Col := range Columns {
			if Col.Suppress {
				continue
			}
			sb.WriteString(Col.Get(r))
		}
		fmt.Println(sb.String())
	}
	fmt.Println()
	return nil
}

type RowStruct struct {
	err  error
	Text []string
}

type ColumnStruct struct {
	Header []string

	GenStats       bool
	CalcTotal      bool
	GenStatsTitles bool

	Suppress    bool
	HideIfEmpty bool
	LeftAlign   bool

	HasData bool

	Min   float64
	Avg   float64
	Max   float64
	Total float64
	Count int
	Valid bool

	Width  int
	FmtStr string

	Cells []string
}

func (o *ColumnStruct) Add(Data ValueStruct) {
	if len(o.Cells) == 0 {

		o.Min = math.MaxInt64
		o.Max = math.MinInt64

		for _, v := range o.Header {
			if len(v) > o.Width {
				o.Width = len(v)
			}
			o.Cells = append(o.Cells, v)
		}
	}

	if len(Data.String) > o.Width {
		o.Width = len(Data.String)
	}
	o.Cells = append(o.Cells, Data.String)

	if len(Data.String) > 0 {
		o.HasData = true
	}

	if !Data.Valid {
		return
	}
	o.Valid = true
	if Data.Float < o.Min {
		o.Min = Data.Float
	}
	if Data.Float > o.Max {
		o.Max = Data.Float
	}
	o.Total += Data.Float
	o.Count++
}

func (o *ColumnStruct) Finalize() {

	o.Suppress = o.Suppress || (o.HideIfEmpty && !o.HasData)

	if o.Suppress {
		return
	}

	if o.LeftAlign {
		o.FmtStr = fmt.Sprintf(" %%-%ds ", o.Width)
	} else {
		o.FmtStr = fmt.Sprintf(" %%%ds ", o.Width)
	}

	if !o.GenStats {
		return
	}

	o.Avg = (float64)(o.Total) / (float64)(o.Count)

	o.Cells = append(o.Cells, "")

	if o.GenStatsTitles {
		o.Cells = append(o.Cells, "Minimum")
		o.Cells = append(o.Cells, "Average")
		o.Cells = append(o.Cells, "Maximum")
		o.Cells = append(o.Cells, "Total")
	} else if o.Valid {
		o.Cells = append(o.Cells, fmt.Sprintf("%.1f", o.Min))
		o.Cells = append(o.Cells, fmt.Sprintf("%.1f", o.Total/(float64)(o.Count)))
		o.Cells = append(o.Cells, fmt.Sprintf("%.1f", o.Max))
		if o.CalcTotal {
			o.Cells = append(o.Cells, fmt.Sprintf("%.1f", o.Total))
		} else {
			o.Cells = append(o.Cells, "")
		}
	} else {
		o.Cells = append(o.Cells, "")
		o.Cells = append(o.Cells, "")
		o.Cells = append(o.Cells, "")
		o.Cells = append(o.Cells, "")
	}
}

func (o *ColumnStruct) Get(Row int) string {

	rc := fmt.Sprintf(o.FmtStr, o.Cells[Row])

	return rc
}
