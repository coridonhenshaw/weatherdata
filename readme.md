# Weatherdata

Weatherdata is a command-line tool to read weather observations from weather stations which report into the Environment Canada [Surface Weather Observation (SWOB)](https://dd.weather.gc.ca/observations/doc/) public data source.

The intended use case is to access weather information that Environment Canada does not include in public-facing observations, such as precipitation rates, wet bulb temperatures, and intra-hour temperature extremes.

Weatherdata is also capable of totalizing observation data over multiple hours (limited only by data retention on the SWOB system) to report on cumulative precipitation and the ranges of observed values over time.

Weatherdata is intended for casual, interactive, use. Use of Weatherdata as a backend for large-scale scraping or to implement public-access services is prohibited. Tools for these use cases should use Environment Canada's HPFX server or AMQP push feeds.

## Example
```
> weatherdata -s CWAS-AUTO --total --starttime "2021-11-25 00 PST"

Totalizing station CWAS-AUTO from 2021-11-25 08 UTC (2021-11-25 00 PST) to 2021-11-26 21 UTC (2021-11-26 13 PST) (37 hours):

  TIME                MIN   AVG   MAX   HUM   PRESSURE   WET BULB   DEW POINT   PRECIP   WIND   GUSTS  
                      °C    °C    °C    %     HPA        °C         °C          MM/HR    KM/H   KM/H   
  2021-11-25 01 PST   5.7   5.7   6.3   92    1024.3     5.2        4.5         0.4      16.6   25.8   
  2021-11-25 02 PST   5.7   5.8   6.1   93    1023.8     5.3        4.8         1        16.6   24.1   
  2021-11-25 03 PST   5.4   5.7   5.8   92    1023.6     5.2        4.6         1.2      14.9   23.4   
  2021-11-25 04 PST   4.9   5.0   5.7   92    1022.6     4.5        3.8         2.4      17.7   27.2   
  2021-11-25 05 PST   4.3   4.5   5.1   93    1022.0     4.1        3.5         2        14.2   23.3   
  2021-11-25 06 PST   4.5   4.8   4.9   91    1021.4     4.1        3.3         1        16.3   23.5   
  2021-11-25 07 PST   4.6   4.7   4.9   94    1020.5     4.3        3.8         1.6      16.8   24.2   
  2021-11-25 08 PST   4.5   4.9   4.9   95    1020.1     4.5        4.1         2        19.0   25.3   
  2021-11-25 09 PST   4.7   4.8   5.1   94    1019.9     4.4        3.9         4.2      19.1   24.9   
  2021-11-25 10 PST   4.7   4.9   5.3   92    1020.2     4.4        3.8         3.8      20.5   25.8   
  2021-11-25 11 PST   4.8   5.0   5.2   94    1019.5     4.6        4.1         3.2      20.1   27.9   
  2021-11-25 12 PST   4.9   5.3   5.4   95    1019.2     4.9        4.5         2.8      18.6   26.0   
  2021-11-25 13 PST   5.0   5.3   5.4   94    1018.6     4.9        4.4         3.2      24.2   33.2   
  2021-11-25 14 PST   5.2   5.3   5.8   94    1017.8     4.9        4.5         2.2      21.0   35.8   
  2021-11-25 15 PST   5.2   5.6   5.9   96    1017.8     5.2        4.9         2.8      17.3   25.0   
  2021-11-25 16 PST   5.2   5.3   5.8   96    1017.7     5.0        4.7         1.8      19.8   27.6   
  2021-11-25 17 PST   5.3   5.5   5.8   95    1018.0     5.2        4.8         1.8      22.5   29.9   
  2021-11-25 18 PST   5.4   5.7   5.8   95    1018.2     5.3        4.9         1.2      24.5   32.3   
  2021-11-25 19 PST   5.6   5.7   6.1   93    1017.9     5.2        4.6         0.4      26.6   37.7   
  2021-11-25 20 PST   5.7   5.9   6.1   94    1017.6     5.4        4.9         0.8      17.0   26.5   
  2021-11-25 21 PST   5.7   5.8   6.1   95    1017.7     5.4        5           0.2      14.7   23.1   
  2021-11-25 22 PST   5.8   6.0   6.2   96    1017.4     5.7        5.4         0.6      7.7    13.3   
  2021-11-25 23 PST   5.8   6.1   6.2   96    1016.9     5.9        5.6         1        10.5   16.1   
  2021-11-26 00 PST   5.9   6.0   6.2   97    1016.9     5.8        5.6         3        12.3   18.2   
  2021-11-26 01 PST   5.8   6.1   6.2   97    1017.2     5.9        5.7         3.4      20.7   26.7   
  2021-11-26 02 PST   5.8   6.6   6.6   97    1017.1     6.3        6.1         1.6      17.4   25.3   
  2021-11-26 03 PST   6.0   6.4   6.5   97    1017.2     6.2        5.9         0.4      8.3    17.5   
  2021-11-26 04 PST   6.4   6.6   6.7   96    1017.1     6.4        6.1         0        4.0    10.9   
  2021-11-26 05 PST   6.6   6.6   6.6   96    1017.5     6.4        6.1         0        1.8    5.1    
  2021-11-26 06 PST   6.3   6.6   6.6   97    1017.8     6.4        6.1         0        7.1    12.7   
  2021-11-26 07 PST   6.6   6.7   6.7   97    1018.4     6.4        6.2         0        2.0    7.3    
  2021-11-26 08 PST   6.5   6.7   6.7   97    1018.8     6.4        6.2         0        2.7    7.9    
  2021-11-26 09 PST   6.6   6.8   6.8   97    1019.8     6.5        6.3         0        2.0    6.5    
  2021-11-26 10 PST   6.7   8.3   8.3   95    1020.4     7.9        7.6         0        2.2    6.6    
  2021-11-26 11 PST   7.5   7.6   8.6   94    1021.1     7.1        6.6         0        2.1    5.9    
  2021-11-26 12 PST   7.6   8.3   9.8   92    1021.2     7.6        7           0        1.7    6.0    
  2021-11-26 13 PST   8.0   8.4   9.6   91    1021.4     7.7        7           0        3.1    6.6    

       Station name: HOWE SOUND - PAM ROCKS
  Temperature range: 4.3 - 9.8 °C
     Humidity range: 91 - 97 percent
     Pressure range: 1016.9 - 1024.3 hPa
     Wet bulb range: 4.1 - 7.9 °C
     Dewpoint range: 3.3 - 7.6 °C
Total precipitation: 50.0 mm
 Mean precipitation: 1.4 mm/hr
 Peak precipitation: 4.2 mm/hr
    Peak wind speed: 37.7 km/h

```

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

`weatherdata -s <identifier> --total --hours 36 --endtime "2021-11-26 00 PST"`

Acquires and presents a summary of weather observations recorded by the specified station for twelve hours up to midnight, 26 November 2021, Pacific Time.

## Station Identifiers

The Environment Canada SWOB API provides three separate endpoints for different types of weather stations.

Most weather stations are included in the [SWOB station list CSV](https://dd.weather.gc.ca/observations/doc/swob-xml_station_list.csv) and are identified to Weatherdata in the form `Cxxx-yyyy`, where `Cxxx` is the IATA four letter station code (found in CSV column A) and `yyyy` indicates whether the station is AUTOmatic or MANual (found in CSV column J).

For example, the automatic weather station in Flin Flon, MB, is identified as `CWFO-AUTO` while the manual weather station at Toronto Pearson Airport is identified as `CYYZ-MAN`.

A map of Environment Canada stations (with IATA ids) is  available via [GeoMet](https://api.weather.gc.ca/collections/swob-realtime/items).

Weatherdata can also harvest data from Environment Canada partners (typically Provincial entities) included in the [partner station list](https://dd.weather.gc.ca/observations/doc/swob-xml_partner_station_list.csv).

Unfortunately, Environment Canada does not provide any way to automatically create data feed URLs for partner stations. To determine the identifier for a partner station:

 1. Find a station of interest in the station list.
 2. Note the 'IATA ID' in column A.
 3. Note the 'data provider' in column N.
 4. Browse to https://dd.weather.gc.ca/observations/swob-ml/partners
 5. Make an educated guess as to which directory name most closely matches the name of the data provider. Note the result.

The Weatherdata station identifier is in the form `partners/<step5>/<step2>` For example: `partners/dfo-ccg-lighthouse/nootka`

Note: this technique does not work for partner stations operated by the Yukon government. To determine the identifier for a YT partner station:

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
