package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
)

func GetValue(XMLDoc *xmlquery.Node, KeyStem string) string {

	KeyStem = "//element[@name='" + KeyStem + "']"

	Value := KeyStem + "/@value"
	Valid := KeyStem + "/qualifier/@value"

	OK := false

	c := xmlquery.FindOne(XMLDoc, Valid)
	if c != nil {
		if c.InnerText() == "100" {
			OK = true
		}
		if c.InnerText() == "1" {
			OK = true
		}
	}

	if OK == false {
		return ""
	}

	c = xmlquery.FindOne(XMLDoc, Value)
	if c != nil {
		return c.InnerText()
	}

	return ""
}

type Observation struct {
	Station              string
	Timestamp            string
	Temperature          string
	MinTemperature       string
	MaxTemperature       string
	Humidity             string
	Pressure             string
	WetBulbTemperature   string
	DewPoint             string
	Windchill            string
	Humidex              string
	Precipitation        string
	AverageWindSpeed     string
	PeakWindSpeed        string
	AverageWindDirection string
}

func RoundFloat(In float64) string {
	return fmt.Sprintf("%.0f", math.Round(In))
}

func ParseObservation(XMLString string) (Observation, error) {

	var O Observation

	XMLDoc, err := xmlquery.Parse(strings.NewReader(XMLString))
	if err != nil {
		return O, err
	}

	c := xmlquery.FindOne(XMLDoc, "//element[@name='stn_nam']/@value")
	if c != nil {
		O.Station = c.InnerText()
	}

	c = xmlquery.FindOne(XMLDoc, "//element[@name='date_tm']/@value")
	if c != nil {
		O.Timestamp = c.InnerText()
	}

	O.Temperature = GetValue(XMLDoc, "air_temp")
	O.MinTemperature = GetValue(XMLDoc, "min_air_temp_pst1hr")
	O.MaxTemperature = GetValue(XMLDoc, "max_air_temp_pst1hr")
	O.Humidity = GetValue(XMLDoc, "rel_hum")
	O.Pressure = GetValue(XMLDoc, "stn_pres")
	O.WetBulbTemperature = GetValue(XMLDoc, "wetblb_temp")
	O.DewPoint = GetValue(XMLDoc, "dwpt_temp")

	O.Precipitation = GetValue(XMLDoc, "pcpn_amt_pst1hr")
	if O.Precipitation == "" {
		O.Precipitation = GetValue(XMLDoc, "rnfl_snc_last_syno_hr")
	}

	O.AverageWindSpeed = GetValue(XMLDoc, "avg_wnd_spd_10m_pst1hr")
	if O.AverageWindSpeed == "" {
		O.AverageWindSpeed = GetValue(XMLDoc, "avg_wnd_spd_10m_pst10mts")
	}
	if O.AverageWindSpeed == "" {
		O.AverageWindSpeed = GetValue(XMLDoc, "avg_wnd_spd_10m_pst2mts")
	}

	O.PeakWindSpeed = GetValue(XMLDoc, "max_wnd_spd_10m_pst1hr")
	if O.PeakWindSpeed == "" {
		O.PeakWindSpeed = GetValue(XMLDoc, "max_wnd_gst_spd_10m_pst10mts")
	}

	O.AverageWindDirection = GetValue(XMLDoc, "avg_wnd_dir_10m_pst1hr")
	if O.AverageWindDirection == "" {
		O.AverageWindDirection = GetValue(XMLDoc, "avg_wnd_dir_10m_pst2mts")
	}

	var WindSpeed float64

	Temperature, err := strconv.ParseFloat(O.MinTemperature, 64)

	if err == nil {
		WindSpeed, err = strconv.ParseFloat(O.PeakWindSpeed, 64)

		if err != nil {
			WindSpeed, err = strconv.ParseFloat(O.AverageWindSpeed, 64)
		}

	}

	// Windchill formula per https://www.climate.weather.gc.ca/glossary_e.html
	if err == nil && Temperature < 0 {
		if WindSpeed >= 5 {
			O.Windchill = RoundFloat(13.12 + 0.6215*Temperature - 11.37*math.Pow(WindSpeed, 0.16) + 0.3965*Temperature*math.Pow(WindSpeed, 0.16))
		} else {
			O.Windchill = RoundFloat(Temperature + ((-1.59+0.1345*Temperature)/5)*WindSpeed)
		}

	}

	// Humidex formula per https://www.climate.weather.gc.ca/glossary_e.html
	if err == nil && Temperature >= 20 {
		var Dewpoint float64
		Dewpoint, err = strconv.ParseFloat(O.DewPoint, 64)

		if err != nil {

			var e float64 = 6.11 * math.Exp(5417.7530*((1/273.15)-(1/Dewpoint)))
			var h float64 = (0.5555) * (e - 10.0)
			if h >= 1 {
				O.Humidex = RoundFloat(Temperature + h)
			}

		}
	}

	return O, nil

}
