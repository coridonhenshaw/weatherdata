package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

type ObservationSummary struct {
	MinTemperature        float64
	MaxTemperature        float64
	MinHumidity           float64
	MaxHumidity           float64
	MinPressure           float64
	MaxPressure           float64
	MinWetBulbTemperature float64
	MaxWetBulbTemperature float64
	MinDewPoint           float64
	MaxDewPoint           float64
	TotalPrecipitation    float64
	PeakPrecipitation     float64
	PeakWindSpeed         float64
}

func (s *ObservationSummary) fill_defaults() {
	s.MaxDewPoint = -9999999
	s.MaxHumidity = -9999999
	s.MaxPressure = -9999999
	s.MaxTemperature = -9999999
	s.MaxDewPoint = -9999999
	s.MaxWetBulbTemperature = -9999999
	s.MinDewPoint = 9999999
	s.MinHumidity = 9999999
	s.MinPressure = 9999999
	s.MinTemperature = 9999999
	s.MinWetBulbTemperature = 9999999
	s.TotalPrecipitation = 0
	s.PeakWindSpeed = -9999999
}

func SetTotal(Total *float64, Input string) {
	v, err := strconv.ParseFloat(Input, 64)
	if err != nil {
		return
	}

	*Total += v

}

func SetMinMax(Min *float64, Max *float64, Input string) {
	v, err := strconv.ParseFloat(Input, 64)
	if err != nil {
		return
	}

	if v < *Min {
		*Min = v
	}
	if v > *Max {
		*Max = v
	}

}

func GetTotals(Station string, StartTime time.Time, EndTime time.Time) error {

	var err error

	BaseURL, err := MakeBaseURL(Station)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var Observation Observation
	var Obs ObservationSummary
	Obs.fill_defaults()

	var Hours = int(EndTime.Sub(StartTime).Hours())

	UTCLoc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(`Failed to load location "UTC"`)
	}

	LocalLoc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(`Failed to load location "Local"`)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Min\n°C", "Avg\n°C", "Max\n°C", "Hum\n%", "Pressure\nhPA", "Wet Bulb\n°C", "Dew Point\n°C", "Precip\nmm/hr", "Wind\nkm/h", "Gusts\nkm/h", "Wind Dir\n°"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("")
	table.SetNoWhiteSpace(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	ST := StartTime.In(UTCLoc).Format("2006-01-02 15 MST") + " (" + StartTime.In(LocalLoc).Format("2006-01-02 15 MST") + ")"
	ET := EndTime.In(UTCLoc).Format("2006-01-02 15 MST") + " (" + EndTime.In(LocalLoc).Format("2006-01-02 15 MST") + ")"

	fmt.Printf("\nTotalizing station %s from %s to %s (%d hours):\n\n", Station, ST, ET, Hours)

	for i := 1; i <= Hours; i++ {

		ObservationTime := StartTime.Add(time.Hour * time.Duration(i))
		ObservationTime = ObservationTime.In(UTCLoc)

		FinalURL, err := InjectURLDateTime(BaseURL, ObservationTime)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// fmt.Println(FinalURL)

		s, err := HTTPSGet(FinalURL)
		if err != nil {
			fmt.Println(err)
		}

		Observation, _ = ParseObservation(s)

		//		fmt.Println(Observation)

		Row := []string{ObservationTime.In(LocalLoc).Format("2006-01-02 15 MST"), Observation.MinTemperature, Observation.Temperature, Observation.MaxTemperature, Observation.Humidity, Observation.Pressure, Observation.WetBulbTemperature, Observation.DewPoint, Observation.Precipitation, Observation.AverageWindSpeed, Observation.PeakWindSpeed, Observation.AverageWindDirection}
		table.Append(Row)

		var Junk float64

		SetMinMax(&Obs.MinTemperature, &Junk, Observation.MinTemperature)
		SetMinMax(&Junk, &Obs.MaxTemperature, Observation.MaxTemperature)
		SetMinMax(&Obs.MinHumidity, &Obs.MaxHumidity, Observation.Humidity)
		SetMinMax(&Obs.MinPressure, &Obs.MaxPressure, Observation.Pressure)
		SetMinMax(&Obs.MinWetBulbTemperature, &Obs.MaxWetBulbTemperature, Observation.WetBulbTemperature)
		SetMinMax(&Obs.MinDewPoint, &Obs.MaxDewPoint, Observation.DewPoint)
		SetTotal(&Obs.TotalPrecipitation, Observation.Precipitation)
		SetMinMax(&Junk, &Obs.PeakPrecipitation, Observation.Precipitation)
		SetMinMax(&Junk, &Obs.PeakWindSpeed, Observation.PeakWindSpeed)
	}

	table.Render()

	fmt.Println()
	fmt.Println("       Station name:", Observation.Station)

	if Obs.MinTemperature != 9999999 && Obs.MaxTemperature != -9999999 {
		fmt.Println("  Temperature range:", Obs.MinTemperature, "-", Obs.MaxTemperature, "°C")
	} else {
		fmt.Println("  Temperature range: <not valid>")
	}

	if Obs.MinHumidity >= 0 && Obs.MinHumidity < 100 && Obs.MaxHumidity > 0 && Obs.MaxHumidity <= 100 {
		fmt.Println("     Humidity range:", Obs.MinHumidity, "-", Obs.MaxHumidity, "percent")
	} else {
		fmt.Println("     Humidity range: <not valid>")
	}

	if Obs.MinPressure != 9999999 && Obs.MaxPressure != -9999999 {
		fmt.Println("     Pressure range:", Obs.MinPressure, "-", Obs.MaxPressure, "hPa")
	} else {
		fmt.Println("     Pressure range: <not valid>")
	}

	if Obs.MinWetBulbTemperature != 9999999 && Obs.MaxWetBulbTemperature != -9999999 {
		fmt.Println("     Wet bulb range:", Obs.MinWetBulbTemperature, "-", Obs.MaxWetBulbTemperature, "°C")
	} else {
		fmt.Println("     Wet bulb range: <not valid>")
	}

	if Obs.MinDewPoint != 9999999 && Obs.MaxDewPoint != -9999999 {
		fmt.Println("     Dewpoint range:", Obs.MinDewPoint, "-", Obs.MaxDewPoint, "°C")
	} else {
		fmt.Println("     Dewpoint range: <not valid>")
	}

	if Obs.TotalPrecipitation >= 0 {
		fmt.Printf("Total precipitation: %.1f mm\n", Obs.TotalPrecipitation)
		fmt.Printf(" Mean precipitation: %.1f mm/hr\n", Obs.TotalPrecipitation/float64(Hours))
	} else {
		fmt.Println("Total precipitation: <not valid>")
		fmt.Println(" Mean precipitation: <not valid>")
	}

	if Obs.PeakPrecipitation >= 0 {
		fmt.Println(" Peak precipitation:", Obs.PeakPrecipitation, "mm/hr")
	} else {
		fmt.Println(" Peak precipitation: <not valid>")
	}

	if Obs.PeakWindSpeed >= 0 {
		fmt.Println("    Peak wind speed:", Obs.PeakWindSpeed, "km/h")
	} else {
		fmt.Println("    Peak wind speed: <not valid>")
	}

	fmt.Println()

	return nil
}
