package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"weatherdata/psychrometrics"

	"github.com/antchfx/xmlquery"
)

func GetValue(XMLDoc *xmlquery.Node, KeyStem string) string {

	KeyStem = "//element[@name='" + KeyStem + "']"

	Value := KeyStem + "/@value"
	Valid := KeyStem + "/qualifier/@value"

	OK := false

	c := xmlquery.FindOne(XMLDoc, Valid)
	if c != nil {
		QualValue, err := strconv.Atoi(c.InnerText())
		if err == nil {
			if (QualValue > 0 && QualValue < 6) || QualValue == 100 {
				OK = true
			}
		}
	} else {
		OK = true
	}

	if OK == false {
		return ""
	}

	c = xmlquery.FindOne(XMLDoc, Value)
	if c != nil {
		rc := c.InnerText()
		if rc == "MSNG" {
			return ""
		}
		return rc
	}

	return ""
}

func GetValueFromList(XMLDoc *xmlquery.Node, KeyStem ...string) string {
	var Value string
	for _, v := range KeyStem {
		Value = GetValue(XMLDoc, v)
		if len(Value) > 0 {
			break
		}

	}
	return Value
}

type NVS struct {
	KeyStem       string
	Normalization float64
}

func GetNormalizedValueFromList(XMLDoc *xmlquery.Node, NV []NVS) ValueStruct {
	Value := ValueStruct{Valid: false}
	var Scale float64

	for _, v := range NV {
		Value.String = GetValue(XMLDoc, v.KeyStem)
		if len(Value.String) > 0 {
			Scale = v.Normalization
			break
		}

	}
	if len(Value.String) == 0 {
		return Value
	}

	q, err := strconv.ParseFloat(Value.String, 64)
	if err != nil {
		return Value
	}

	q /= Scale

	Value.LoadF(q)
	if Scale != 1 {
		Value.String = "E " + Value.String
	}

	return Value
}

type ObservationStruct struct {
	Station              string
	Identifier           string
	Timestamp            string
	Temperature          ValueStruct
	MinTemperature       ValueStruct
	MaxTemperature       ValueStruct
	Humidity             ValueStruct
	Pressure             ValueStruct
	WetBulbTemperature   ValueStruct
	DewPoint             ValueStruct
	Windchill            ValueStruct
	Humidex              ValueStruct
	Precipitation        ValueStruct
	AverageWindSpeed     ValueStruct
	PeakWindSpeed        ValueStruct
	AverageWindDirection ValueStruct
	err                  error
}

func RoundFloat(In float64) string {
	return fmt.Sprintf("%.0f", math.Round(In))
}

func RoundFloatTwoPlaces(In float64) string {
	return fmt.Sprintf("%.1f", In)
}

func ParseObservation(XMLString string) (ObservationStruct, error) {

	var O ObservationStruct

	XMLDoc, err := xmlquery.Parse(strings.NewReader(XMLString))
	if err != nil {
		O.err = err
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

	O.Temperature.Load(GetValueFromList(XMLDoc, "air_temp", "avg_air_temp_pst1hr", "air_temp_1", "air_temp_2",
		"avg_air_temp_pst2mts"))

	O.MinTemperature.Load(GetValue(XMLDoc, "min_air_temp_pst1hr"))
	O.MaxTemperature.Load(GetValue(XMLDoc, "max_air_temp_pst1hr"))

	O.Humidity.Load(GetValueFromList(XMLDoc, "rel_hum", "avg_rel_hum_pst1hr", "avg_rel_hum_pst5mts", "avg_rel_hum_pst2mts"))

	O.Pressure.Load(GetValue(XMLDoc, "stn_pres"))
	O.WetBulbTemperature.Load(GetValue(XMLDoc, "wetblb_temp"))

	O.DewPoint.Load(GetValue(XMLDoc, "dwpt_temp"))

	// O.Precipitation = GetValueFromList(XMLDoc, "pcpn_amt_pst1hr", "rnfl_amt_pst1hr", "rnfl_snc_last_syno_hr",
	// 	"pcpn_amt_pst30mts", "pcpn_amt_pst20mts", "pcpn_amt_pst15mts", "pcpn_amt_pst10mts", "pcpn_amt_pst5mts",
	// 	"pcpn_amt_pst1mt")

	O.Precipitation = GetNormalizedValueFromList(XMLDoc, []NVS{
		{"pcpn_amt_pst1hr", 1},
		{"rnfl_amt_pst1hr", 1},
		{"rnfl_snc_last_syno_hr", 1},
		{"pcpn_amt_pst30mts", 30 / 60},
		{"pcpn_amt_pst20mts", 20 / 60},
		{"pcpn_amt_pst15mts", 15 / 60},
		{"pcpn_amt_pst10mts", 10 / 60},
		{"pcpn_amt_pst5mts", 5 / 60},
		{"pcpn_amt_pst1mt", 1 / 60}})

	////
	O.AverageWindSpeed.Load(GetValueFromList(XMLDoc,
		"wnd_spd",
		"avg_wnd_spd_pst1hr",
		"avg_wnd_spd_10m_pst1hr",
		"avg_wnd_spd_sclr_pst1hr",
		"avg_wnd_spd_10m_pst15mts",
		"avg_wnd_spd_10m_pst10mts",
		"avg_wnd_spd_pcpn_gag_pst10mts",
		"avg_wnd_spd_pst10mts",
		"avg_wnd_spd_10m_pst2mts"))

	O.PeakWindSpeed.Load(GetValueFromList(XMLDoc, "max_wnd_spd_10m_pst1hr", "max_wnd_spd_pst1hr", "max_wnd_spd_10m_pst10mts",
		"max_wnd_gst_spd_10m_pst10mts"))

	O.AverageWindDirection.Load(GetValueFromList(XMLDoc, "avg_wnd_dir_10m_pst1hr", "avg_wnd_dir_pst1hr",
		"avg_wnd_dir_10m_pst2mts", "avg_wnd_dir_pst10mts"))

	O.AverageWindDirection.String = strings.Split(O.AverageWindDirection.String, ".")[0]

	O.CalcWetBulb()
	O.CalcWindchill()
	O.CalcHumidex()

	return O, nil
}

func (o *ObservationStruct) CalcWetBulb() {
	if !o.WetBulbTemperature.Valid && o.Temperature.Valid && o.Humidity.Valid {
		var Hum float64
		Hum = o.Humidity.Float / float64(100)
		o.WetBulbTemperature.LoadF(psychrometrics.GetTDewPointFromRelHum(o.Temperature.Float, Hum))
		o.WetBulbTemperature.String = "E " + o.WetBulbTemperature.String
	}
}

func GetFirstOf(Values ...ValueStruct) (Vs ValueStruct) {
	for _, v := range Values {
		if v.Valid {
			Vs = v
			break
		}
	}
	return
}

func (o *ObservationStruct) CalcWindchill() {

	Temperature := GetFirstOf(o.MinTemperature, o.Temperature)
	WindSpeed := GetFirstOf(o.PeakWindSpeed, o.AverageWindSpeed)

	if Temperature.Float >= 0 || !Temperature.Valid || !WindSpeed.Valid {
		return
	}

	// Windchill formula per https://www.climate.weather.gc.ca/glossary_e.html

	if WindSpeed.Float >= 5 {
		o.Windchill.LoadF(13.12 + 0.6215*Temperature.Float - 11.37*math.Pow(WindSpeed.Float, 0.16) + 0.3965*Temperature.Float*math.Pow(WindSpeed.Float, 0.16))
	} else {
		o.Windchill.LoadF(Temperature.Float + ((-1.59+0.1345*Temperature.Float)/5)*WindSpeed.Float)
	}

}

func (o *ObservationStruct) CalcHumidex() {
	Temperature := GetFirstOf(o.MaxTemperature, o.Temperature)

	if !o.Temperature.Valid || !o.DewPoint.Valid || Temperature.Float < 20 {
		return
	}

	// Humidex formula per https://www.climate.weather.gc.ca/glossary_e.html

	var e float64 = 6.11 * math.Exp(5417.7530*((1/273.15)-(1/(o.DewPoint.Float+273.15))))
	var h float64 = (0.5555) * (e - 10.0)
	if h >= 1 {
		o.Humidex.LoadF(Temperature.Float + h)
	}

}
