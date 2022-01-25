package main

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type GetObservationEngine struct {
	Initialized          bool
	SQLdb                *sql.DB
	reECStation          *regexp.Regexp
	rePartnersBC         *regexp.Regexp
	rePartnersBCEnvSnow  *regexp.Regexp
	rePartnersBCForestry *regexp.Regexp
	rePartnersNLWater    *regexp.Regexp
	rePartnersNTForestry *regexp.Regexp
	rePartnersNTWater    *regexp.Regexp
	rePartnersSKForestry *regexp.Regexp
	rePartnersYTGov      *regexp.Regexp
	rePartnersYTWater    *regexp.Regexp
}

func (o *GetObservationEngine) Initialize() error {
	o.reECStation = regexp.MustCompile(`C...-((AUTO)|(MAN))`)
	o.rePartnersBC = regexp.MustCompile(`partners/((bc-env-aq)|(bc-tran))/(.*)`)

	// bc-env-snow :: 2021-12-07-0100-bc-env-asw-1a01p-AUTO-swob.xml
	o.rePartnersBCEnvSnow = regexp.MustCompile(`partners/bc-env-snow/(.*)`)

	// bc-forestry :: 2021-12-07-0000-bc-wmb-1002-AUTO-swob.xml
	o.rePartnersBCForestry = regexp.MustCompile(`partners/bc-forestry/(.*)`)

	// nl-water :: https://dd.weather.gc.ca/observations/swob-ml/partners/nl-water/20211207/2021-12-07-0030-nl-deccm-wrmd-nlencl0001-nlencl0001-AUTO-swob.xml
	o.rePartnersNLWater = regexp.MustCompile(`partners/nl-water/(.*)`)
	// fails due to time being in the past

	// nt-forestry :: https://dd.weather.gc.ca/observations/swob-ml/partners/nt-forestry/20211207/aca016d8/2021-12-07-0005-nwt-enr-aca016d8-AUTO-swob.xml
	o.rePartnersNTForestry = regexp.MustCompile(`partners/nt-forestry/(.*)`)

	// nt-water :: https://dd.weather.gc.ca/observations/swob-ml/partners/nt-water/20211207/aca0180a/2021-12-07-0004-nwt-enr-aca0180a-AUTO-swob.xml
	o.rePartnersNTWater = regexp.MustCompile(`partners/nt-water/(.*)`)

	// sk-forestry :: https://dd.weather.gc.ca/observations/swob-ml/partners/sk-forestry/20211207/bagwa/2021-12-07-0000-sk-spsa-wmb-bagwa-bagwa-AUTO-swob.xml
	o.rePartnersSKForestry = regexp.MustCompile(`partners/sk-forestry/(.*)`)

	// yt-gov :: https://dd.weather.gc.ca/observations/swob-ml/partners/yt-gov/20211207/antimony_creek/2021-12-07-0000-ytg-antimonycreek-antimony_creek-AUTO-swob.xml
	o.rePartnersYTGov = regexp.MustCompile(`partners/yt-gov/(.*)`)

	// yt-water :: https://dd.weather.gc.ca/observations/swob-ml/partners/yt-water/20211207/2021-12-07-0000-yt-de-wrb-09aa-m1-09aa-m1-AUTO-swob.xml
	o.rePartnersYTWater = regexp.MustCompile(`partners/yt-water/(.*)`)

	// https://dd.weather.gc.ca/observations/swob-ml/partners/yt-gov/20211208/willow_creek/2021-12-08-0400-ytg-willowcreek-willow_creek-AUTO-swob.xml
	// https://dd.weather.gc.ca/observations/swob-ml/partners/yt-gov/20211208/willow_creek/2021-12-08-0300-ytg-willowcreek-willow_creek-AUTO-swob.xml

	o.SQLdb = Cache.GetSQLConnection()

	o.Initialized = true

	return nil
}

func (o *GetObservationEngine) Get(Station string, Timestamp time.Time) (Observation, error) {

	var Dummy Observation
	var err error

	if Stations.Validate(Station) == false {
		return Dummy, errors.New("invalid station identifier " + Station)
	}

	if o.Initialized != true {
		err = o.Initialize()
		if err != nil {
			return Dummy, err
		}

	}

	var OutString string

	if o.reECStation.MatchString(Station) == true {
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/{{ .Date }}/" + Station[0:4] + "/{{ .FullTimestamp }}00-" + Station + "-swob.xml"
	} else if o.rePartnersBC.MatchString(Station) == true {
		Parts := o.rePartnersBC.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/" + Parts[1] + "/{{ .Date }}/" + Parts[4] + "/{{ .FullTimestamp }}00-" + Parts[1] + "-" + Parts[4] + "-AUTO-swob.xml"
	} else if o.rePartnersBCEnvSnow.MatchString(Station) == true {
		Parts := o.rePartnersBCEnvSnow.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/bc-env-snow/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}00-bc-env-asw-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersBCForestry.MatchString(Station) == true {
		Parts := o.rePartnersBCForestry.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/bc-forestry/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}00-bc-wmb-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersNLWater.MatchString(Station) == true {
		Parts := o.rePartnersNLWater.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/nl-water/{{ .Date }}/{{ .FullTimestamp }}30-nl-deccm-wrmd-" + Parts[1] + "-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersNTForestry.MatchString(Station) == true {
		Parts := o.rePartnersNTForestry.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/nt-forestry/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}05-nwt-enr-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersNTWater.MatchString(Station) == true {
		Parts := o.rePartnersNTWater.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/nt-water/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}04-nwt-enr-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersSKForestry.MatchString(Station) == true {
		Parts := o.rePartnersSKForestry.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/sk-forestry/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}00-sk-spsa-wmb-" + Parts[1] + "-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersYTGov.MatchString(Station) == true {
		Parts := o.rePartnersYTGov.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/yt-gov/{{ .Date }}/" + Parts[1] + "/{{ .FullTimestamp }}00-ytg-" + strings.ReplaceAll(Parts[1], "_", "") + "-" + Parts[1] + "-AUTO-swob.xml"
	} else if o.rePartnersYTWater.MatchString(Station) == true {
		Parts := o.rePartnersYTWater.FindStringSubmatch(Station)
		OutString = "https://dd.weather.gc.ca/observations/swob-ml/partners/yt-water/{{ .Date }}/{{ .FullTimestamp }}00-yt-de-wrb-" + Parts[1] + "-" + Parts[1] + "-AUTO-swob.xml"
	} else {
		return Dummy, errors.New("Invalid station name")
	}

	u := URLDateTime{Timestamp.Format("20060102"), Timestamp.Format("2006-01-02-15")}

	ut, err := template.New("foo").Parse(OutString)
	if err != nil {
		return Dummy, err
	}

	var tpl bytes.Buffer
	err = ut.Execute(&tpl, u)

	if err != nil {
		return Dummy, err
	}

	FinalURL := tpl.String()

	//	fmt.Println(FinalURL)

	ObservationXML, err := HTTPSGet(FinalURL)
	if err != nil {
		return Dummy, err
	}

	return ParseObservation(ObservationXML)

}

type URLDateTime struct {
	Date          string
	FullTimestamp string
}
