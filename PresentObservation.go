package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func PresentObservation(Stations string, ObsTime time.Time) error {

	var GOE GetObservationEngine

	ST := ObsTime.In(UTCLoc).Format("2006-01-02 15 MST") + " (" + ObsTime.In(LocalLoc).Format("2006-01-02 15 MST") + ")"

	fmt.Printf("\nReports at %s from %s:\n\n", ST, Stations)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Station\nName", "Min\n°C", "Avg\n°C", "Max\n°C", "RH\n%", "Barr\nhPA", "Wet Bulb\n°C", "Dew Point\n°C", "Perceived\n°C", "Precip\nmm/hr", "Wind\nkm/h", "Gusts\nkm/h", "W Dir\n°"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	StationArray := strings.Fields(Stations)

	for _, Station := range StationArray {
		//		fmt.Println(Station)
		//		continue

		// q, err := MakeBaseURL(Station)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		//
		// q, err = InjectURLDateTime(q, ObsTime.In(UTCLoc))
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		//
		// // fmt.Println(q)
		// // os.Exit(0)
		//
		// s, err := HTTPSGet(q)
		// if err != nil {
		// 	fmt.Println("HTTP Error", err, "acquiring", q)
		// 	continue
		// }
		//
		// Observation, _ := ParseObservation(s)

		Observation, err := GOE.Get(Station, ObsTime)
		if err != nil {
			fmt.Println("error", err, "acquiring observation from", Station)
			continue
		}

		Row := []string{Observation.Station, Observation.MinTemperature, Observation.Temperature, Observation.MaxTemperature, Observation.Humidity, Observation.Pressure, Observation.WetBulbTemperature, Observation.DewPoint, Observation.Windchill + Observation.Humidex, Observation.Precipitation, Observation.AverageWindSpeed, Observation.PeakWindSpeed, Observation.AverageWindDirection}
		table.Append(Row)

		// t1, err := time.Parse(time.RFC3339, Observation.Timestamp)
		//
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//
		// fmt.Println()
		// fmt.Printf("%s - %s (%s)\n", Observation.Station, t1.Format("2006-01-02 15:04 MST"), t1.In(LocalLoc).Format("2006-01-02 15:04 MST"))
		// fmt.Println()
		// fmt.Println(" Mean temperature:", Observation.Temperature, "°C")
		// fmt.Println("Temperature range:", Observation.MinTemperature, "-", Observation.MaxTemperature, "°C")
		// fmt.Println("         Humidity:", Observation.Humidity, "percent")
		// fmt.Println("         Pressure:", Observation.Pressure, "hPa")
		// fmt.Println("         Wet bulb:", Observation.WetBulbTemperature, "°C")
		// fmt.Println("        Dew point:", Observation.DewPoint, "°C")
		// fmt.Println("    Precipitation:", Observation.Precipitation, "mm/hr")
		// fmt.Println("  Mean wind speed:", Observation.AverageWindSpeed, "Km/h")
		// fmt.Println("  Peak wind speed:", Observation.PeakWindSpeed, "Km/h")
		// fmt.Println()
		// fmt.Println()

	}

	table.Render()
	fmt.Println()

	return nil
}
