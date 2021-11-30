package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

func GetOneObservation(Station string, ObsTime time.Time) error {

	UTCLoc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(`Failed to load location "UTC"`)
	}

	LocalLoc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(`Failed to load location "Local"`)
	}

	ST := ObsTime.In(UTCLoc).Format("2006-01-02 15 MST") + " (" + ObsTime.In(LocalLoc).Format("2006-01-02 15 MST") + ")"

	fmt.Printf("\nReports at %s from %s:\n", ST, Station)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Station\nIdentifier", "Min\n°C", "Avg\n°C", "Max\n°C", "Humidity\n%", "Pressure\nhPA", "Wet Bulb\n°C", "Dew Point\n°C", "Precipitation\nmm/hr", "Wind\nkm/h", "Gusts\nkm/h", "Wind Dir\n°"})
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

	StationArray := strings.Fields(Station)

	for _, Station := range StationArray {
		//		fmt.Println(Station)
		//		continue

		q, err := MakeBaseURL(Station)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		q, err = InjectURLDateTime(q, ObsTime.In(UTCLoc))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// fmt.Println(q)
		// os.Exit(0)

		s, err := HTTPSGet(q)
		if err != nil {
			fmt.Println("HTTP Error", err)
			os.Exit(1)
		}

		Observation, _ := ParseObservation(s)

		Row := []string{Observation.Station, Observation.MinTemperature, Observation.Temperature, Observation.MaxTemperature, Observation.Humidity, Observation.Pressure, Observation.WetBulbTemperature, Observation.DewPoint, Observation.Precipitation, Observation.AverageWindSpeed, Observation.PeakWindSpeed, Observation.AverageWindDirection}
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

	fmt.Println()
	table.Render()
	fmt.Println()

	return nil
}
