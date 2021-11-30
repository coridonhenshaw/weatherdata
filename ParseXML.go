package main

import (
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
	Precipitation        string
	AverageWindSpeed     string
	PeakWindSpeed        string
	AverageWindDirection string
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

	return O, nil

}
