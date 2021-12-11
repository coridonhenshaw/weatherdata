# Weatherdata

Weatherdata is a command-line tool to read weather observations from weather stations which report into the Environment Canada [Surface Weather Observation (SWOB)](https://dd.weather.gc.ca/observations/doc/) public data source.

The intended use case is to access weather information that Environment Canada does not include in public-facing observations, such as precipitation rates, wet bulb temperatures, and intra-hour temperature extremes.

Weatherdata is also capable of totalizing observation data over multiple hours (limited only by data retention on the SWOB system) to report on cumulative precipitation and the ranges of observed values over time.

Weatherdata is intended for casual, interactive, use. Use of Weatherdata as a backend for large-scale scraping or to implement public-access services is prohibited. Tools for these use cases should connect with Environment Canada's HPFX server or AMQP push feeds.

## Example
```
> weatherdata -s CYEG-MAN --total --hours 24       

Totalizing station CYEG-MAN from 2021-12-04 04 UTC (2021-12-03 20 PST) to 2021-12-05 04 UTC (2021-12-04 20 PST) (24 hours):

  STATION             MIN     AVG     MAX     RH   BARR    WET BULB   DEW POINT   PERCEIVED   PRECIP   WIND   GUSTS   W DIR  
  IDENTIFIER          °C      °C      °C      %    HPA     °C         °C          °C          MM/HR    KM/H   KM/H    °      
  2021-12-03 21 PST   -12.0   -12.0   -10.3   73   935.6              -15.9       -18                  11.5           203    
  2021-12-03 22 PST   -12.5   -11.0   -9.9    69   935.3              -15.7       -18                  8.6            184    
  2021-12-03 23 PST   -13.9   -13.9   -10.8   73   935.0              -17.7       -20                  9.0            179    
  2021-12-04 00 PST   -14.8   -14.4   -11.0   78   934.7              -17.4       -21                  9.4            186    
  2021-12-04 01 PST   -14.6   -14.6   -13.2   76   934.1              -18.0       -21                  9.4            178    
  2021-12-04 02 PST   -15.4   -14.4   -12.8   75   933.3              -18.0       -22                  10.8           166    
  2021-12-04 03 PST   -14.5   -13.7   -11.6   74   932.5              -17.4       -18                  5.0            187    
  2021-12-04 04 PST   -15.4   -12.9   -12.9   75   931.6              -16.4       -17                  1.8            0      
  2021-12-04 05 PST   -16.2   -15.6   -12.5   82   930.9              -18.0       -22                  7.9            143    
  2021-12-04 06 PST   -17.3   -16.7   -15.0   80   930.2              -19.4       -21                  5.0            166    
  2021-12-04 07 PST   -17.0   -15.5   -14.8   82   929.1              -17.8       -21                  5.8            138    
  2021-12-04 08 PST   -18.0   -17.7   -15.6   83   928.8              -19.9       -25                  9.4            119    
  2021-12-04 09 PST   -17.7   -16.6   -15.7   83   928.0              -18.9       -21                  3.6            175    
  2021-12-04 10 PST   -16.7   -14.6   -14.2   80   927.2              -17.3       -18                  1.8            0      
  2021-12-04 11 PST   -14.5   -11.8   -11.8   80   926.2              -14.6       -19                  6.1            116    
  2021-12-04 12 PST   -11.7   -10.6   -10.3   76   925.3              -14.0       -16                  6.1            341    
  2021-12-04 13 PST   -10.7   -10.5   -10.3   76   924.4              -14.0       -16                  10.4           37     
  2021-12-04 14 PST   -10.4   -9.8    -9.5    76   924.5              -13.4       -13                  5.0            293    
  2021-12-04 15 PST   -10.8   -10.5   -9.8    83   924.7              -12.9       -14                  5.4            308    
  2021-12-04 16 PST   -11.2   -11.2   -10.4   81   925.3              -13.9       -13                  3.6            242    
  2021-12-04 17 PST   -11.7   -11.3   -11.1   83   926.2              -13.7       -16                  6.5            322    
  2021-12-04 18 PST   -11.4   -11.2   -10.5   77   927.0              -14.4       -18                  14.8           341    
  2021-12-04 19 PST   -12.4   -12.4   -11.2   80   927.9              -15.2       -21                  19.1           338    
  2021-12-04 20 PST   -13.7   -13.7   -12.5   84   929.1              -15.8       -22                  18.0           334    

       Station name: Edmonton International
  Temperature range: -18 - -9.5 °C
     Humidity range: 69 - 84 percent
     Pressure range: 924.4 - 935.6 hPa
     Wet bulb range: <not valid>
     Dewpoint range: -19.9 - -12.9 °C
    Windchill range: -25 - -13 °C
      Humidex range: <not valid>
Total precipitation: 0.0 mm
 Mean precipitation: 0.0 mm/hr
 Peak precipitation: 0 mm/hr
    Peak wind speed: 19.1 km/h
```
## Output Notes

Except for the perceived temperature column, all columns in the output table are passed through from the raw SWOB data without interpretation. Missing values mean that the underlying value was not provided in the station observation report.

Weather stations operated by NavCanada (-MAN type stations) do not generally report wet bulb temperatures. In addition, many automatic weather stations (-AUTO type stations) do not report precipitation amounts.

Wind direction is given in degrees from true north: 0 = north, 90 = east, 180 = south, 270 = west, etc.

The perceived temperature column contains humidex (positive) or windchill (negative) temperature values computed internally by Weatherdata. These values are computed using formulae published by Environment Canada, however the windchill is computed based on worst-case conditions (lowest reported temperature and highest reported wind speed, even if these readings did not occur simultaneously) and may be colder than officially published figures from EC.

## Usage
#### Single-shot Mode

`weatherdata -s <identifier>`

Acquires and presents the weather observations from the station specified for the most recent hour if available.

`weatherdata -s <identifier> --hours 2`

Acquires and presents the weather observations from the station specified two hours into the past.

`weatherdata -s <identifier> --starttime "2021-11-26 12 PST"`

Acquires and presents the weather observations from the station specified at 12 PM, 26 November 2021, Pacific Time.

`weatherdata -s "CYYZ-MAN CWTQ-AUTO CYVR-MAN CYYC-MAN"`

Acquires and presents the weather observations from the major airports in Toronto, Monteal, Vancouver, and Calgary for the most recent hour.

#### Totalizing Mode

`weatherdata -s <identifier> --total --hours 12`

Acquires and presents a summary of weather observations recorded by the specified station over the past 12 hours

`weatherdata -s <identifier> --total --starttime "2021-11-25 00 PST"`

Acquires and presents a summary of weather observations recorded by the specified station from midnight, 25 November 2021, Pacific Time to the present time.

`weatherdata -s <identifier> --total --starttime "2021-11-25 00 PST" --endtime "2021-11-26 00 PST"`

Acquires and presents a summary of weather observations recorded by the specified station from midnight, 25 November 2021, Pacific Time to midnight, 26 November 2021, Pacific Time.

`weatherdata -s <identifier> --total --hours 12 --endtime "2021-11-26 00 PST"`

Acquires and presents a summary of weather observations recorded by the specified station for twelve hours up to midnight, 26 November 2021, Pacific Time.

## Station Identifiers

The Environment Canada SWOB API provides three separate endpoints for different types of weather stations.

Most weather stations are included in the [SWOB station list CSV](https://dd.weather.gc.ca/observations/doc/swob-xml_station_list.csv) and are identified to Weatherdata in the form `Cxxx-yyyy`, where `Cxxx` is the IATA four letter station code (found in CSV column A) and `yyyy` indicates whether the station is AUTOmatic or MANual (found in CSV column J).

For example, the automatic weather station in Flin Flon, MB, is identified as `CWFO-AUTO` while the manual weather station at Toronto Pearson Airport is identified as `CYYZ-MAN`.

A map of Environment Canada stations (with IATA IDs) is available via [GeoMet](https://api.weather.gc.ca/collections/swob-realtime/items).

Weatherdata can also harvest data from Environment Canada partners (typically Provincial entities) included in the [partner station list](https://dd.weather.gc.ca/observations/doc/swob-xml_partner_station_list.csv).

Unfortunately, Environment Canada does not provide any way to automatically build data feed URLs for partner stations. To determine the identifier for a partner station, except for one operated by the Yukon government:

 1. Find a station of interest in the station list.
 2. Note the 'IATA ID' in column A.
 3. Note the 'data provider' in column N.
 4. Browse to https://dd.weather.gc.ca/observations/swob-ml/partners
 5. Make an educated guess as to which directory name most closely matches the name of the data provider. Note the result.

The Weatherdata station identifier is in the form `partners/<step5>/<step2>` For example: `partners/dfo-ccg-lighthouse/nootka`

To determine the identifier for a partner station operated by the Yukon government:

 1. Browse to https://dd.weather.gc.ca/observations/swob-ml/partners/yt-gov
 2. Enter any sub-directory.
 3. Note the location names in the directory listing.

The Weatherdata station identifier for YT government stations is in the form `/partners/yt-gov/<step3>`. For example: `partners/yt-gov/hasselberg`

Weatherdata does not currently read data provided by MSC-operated marine buoys under the https://dd.weather.gc.ca/observations/swob-ml/marine/moored-buoys/ hierarchy.

## Note

Data provided by Environment Canada are subject to assorted terms of use, as [made available by Environment Canada](https://eccc-msc.github.io/open-data/msc-data/obs_station/readme_obs_insitu_en/).

These terms are separate from the terms that apply to Weatherdata.

## License

Copyright 2021 Coridon Henshaw

Permission is granted to all natural persons to execute, distribute, and/or modify this software (including its documentation) subject to the following terms:

1. Subject to point \#2, below, **all commercial use and distribution is prohibited.** This software has been released for personal and academic use for the betterment of society through any purpose that does not create income or revenue. *It has not been made available for businesses to profit from unpaid labor.*

2. Re-distribution of this software on for-profit, public use, repository hosting sites (for example: Github) is permitted provided no fees are charged specifically to access this software.

3. **This software is provided on an as-is basis and may only be used at your own risk.** This software is the product of a single individual's recreational project. The author does not have the resources to perform the degree of code review, testing, or other verification required to extend any assurances that this software is suitable for any purpose, or to offer any assurances that it is safe to execute without causing data loss or other damage.

4. **This software is intended for experimental use in situations where data loss (or any other undesired behavior) will not cause unacceptable harm.** Users with critical data safety needs must not use this software and, instead, should use equivalent tools that have a proven track record.

5. If this software is redistributed, this copyright notice and license text must be included without modification.

6. Distribution of modified copies of this software is discouraged but is not prohibited. It is strongly encouraged that fixes, modifications, and additions be submitted for inclusion into the main release rather than distributed independently.

7. This software reverts to the public domain 10 years after its final update or immediately upon the death of its author, whichever happens first.

8. Use of Weatherdata as a backend for large-scale scraping or to implement any service accessible to the general public is prohibited.
