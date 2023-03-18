# Weatherdata

Weatherdata is a command-line tool to read Canadian weather observations from weather stations which report into the Environment Canada [Surface Weather Observation (SWOB)](https://dd.weather.gc.ca/observations/doc/) public data source. It is intended to access weather observations that Environment Canada does not include in either its [weather information site](https://weather.gc.ca/canada_e.html) or weather app.

In addition to extracting temperature/humidity/wind and other basic observational data, Weatherdata can extract:

 * Observations from non-EC weather stations that report into SWOB but are not shown by EC in public-facing weather station lists. 
 * Real-time (per minute) reporting, where available from some EC weather stations.
 * Precipitation rates.
 * Wet bulb temperatures.
 * Wind direction in degrees.
 * Per hour min-max temperatures.
 * 30 day observation history.

Weatherdata is intended for casual, interactive, use. In the interests of being polite to EC's SWOB servers, use of Weatherdata as a backend for large-scale scraping or to implement public-access services is prohibited. Tools for these use cases should connect with Environment Canada's HPFX server or AMQP feed.

## Sample Output
```
> weatherdata -u observation "CYYZ-MAN" -s -8

Weatherdata Release 1 -- https://github.com/coridonhenshaw/weatherdata

Reports from CYYZ-MAN over 2023-03-14 14:06 UTC to 2023-03-14 22:06 UTC:

                       Station  Observation   Min   Avg   Max    Rel   Barr   Wet    Dew       Dew  Windchill   Wind   Gust     Wind 
                          Name         Time  Temp  Temp  Temp  Humid  Press  Bulb  Point   Point Δ        Max  Speed  Speed      Dir 
                                        UTC    °C    °C    °C      %    hPA    °C     °C        °C         °C   km/h   km/h        ° 
 Toronto/Pearson International  230314 1500  -5.0  -4.4  -4.2     71  994.2  -8.4   -8.9       4.5        -14   36.0   49.7  335 NNW 
 Toronto/Pearson International  230314 1600  -4.5  -3.2  -3.2     66  994.4  -8.0   -8.7       5.5        -12   34.2   45.0  327 NNW 
 Toronto/Pearson International  230314 1700  -3.2  -2.5  -2.3     64  994.6  -7.7   -8.4       5.9        -10   33.5         337 NNW 
 Toronto/Pearson International  230314 1800  -2.5  -1.6  -1.4     63  994.3  -7.0   -7.8       6.2        -10   36.4   47.5  332 NNW 
 Toronto/Pearson International  230314 1900  -1.7  -1.3  -1.3     63  994.2  -6.7   -7.4       6.1        -10   35.3   47.9  336 NNW 
 Toronto/Pearson International  230314 2000  -2.4  -2.2  -1.0     69  994.7  -6.6   -7.1       4.9        -11   33.5   51.5  339 NNW 
 Toronto/Pearson International  230314 2100  -2.5  -2.5  -1.9     67  995.1  -7.2   -7.7       5.2        -12   40.7   57.6  328 NNW 
 Toronto/Pearson International  230314 2200  -2.7  -2.5  -2.1     69  995.8  -6.9   -7.5       5.0        -12   47.5   64.1  321  NW 
                                                                                                                                      
                       Minimum               -5.0  -4.4  -4.2   63.0  994.2         -8.9       4.5      -14.0   33.5   45.0          
                       Average               -3.1  -2.5  -2.2   66.5  994.7         -7.9       5.4      -11.4   37.1   51.9          
                       Maximum               -1.7  -1.3  -1.0   71.0  995.8         -7.1       6.2      -10.0   47.5   64.1          
                         Total


> weatherdata -u observation "CYYZ-MAN CYEG-MAN CYVR-MAN"

Weatherdata Release 1 -- https://github.com/coridonhenshaw/weatherdata

Reports from CYYZ-MAN CYEG-MAN CYVR-MAN at 2023-03-14 22:10 UTC:

                       Station  Observation   Min   Avg   Max    Rel    Barr   Wet    Dew       Dew  Windchill   Wind   Gust     Wind 
                          Name         Time  Temp  Temp  Temp  Humid   Press  Bulb  Point   Point Δ        Max  Speed  Speed      Dir 
                                        UTC    °C    °C    °C      %     hPA    °C     °C        °C         °C   km/h   km/h        ° 
 Toronto/Pearson International  230314 2200  -2.7  -2.5  -2.1     69   995.8  -6.9   -7.5       5.0        -12   47.5   64.1  321  NW 
        Edmonton International  230314 2200  -5.2  -4.3  -4.2     63   920.6  -9.6  -10.2       5.9         -8   10.1         219  SW 
       Vancouver International  230314 2200   8.7   9.2  10.1     56  1007.3   0.9    0.8       8.4              19.1   32.0  132  SE 
                                                                                                                                      
                       Minimum               -5.2  -4.3  -4.2   56.0   920.6  -9.6  -10.2       5.0      -12.0   10.1   32.0          
                       Average                0.3   0.8   1.3   62.7   974.6  -5.2   -5.6       6.4      -10.0   25.6   48.0          
                       Maximum                8.7   9.2  10.1   69.0  1007.3   0.9    0.8       8.4       -8.0   47.5   64.1          
                         Total

```
## Output Notes

With some exceptions, most output columns show raw values from the SWOB system. Values not provided by the originating weather station are represented as blanks.

Wet bulb temperature will be estimated (per the 2017 version of the ASHRAE Fundamentals Handbook) if no wet bulb value is provided in the station report. The appropriateness of using the ASHRAE formula to estimate outdoor wet bulb temperatures is not known. An 'E' will be shown in the web bulb column if the wet bulb temperature has been estimated.

The dew point difference column contains the difference between the average air temperature and the dew point. A larger difference between temperature and dew point usually indicates more comfortable weather. 

Humidex and windchill are computed internally using formulae provided by EC. Windchill is computed based on the lowest reported temperature and highest reported wind speed, even if these readings did not occur simultaneously, and may be colder than officially published figures from EC. Similarly, humidex is computed based on the highest reported temperature, and may be higher than officially published figures from EC.

Wind direction is given in degrees from true north: 0/360 = north, 90 = east, 180 = south, 270 = west, etc.

Some columns (such as humidex, windchill, and precipitation rate) are only shown when relevant data is available.

## Example Usage

#### Station List Mode

`weatherdata stations %toronto%`

Show all SWOB stations with names that contain the word Toronto. Uses SQL LIKE syntax, where `%` is a wildcard.

`weatherdata stations --kml stations.kml`

Exports the locations of known SWOB stations to `stations.kml` for use in Google Earth or similar tools.

#### Observation Mode

`weatherdata observation CXTO-AUTO`

Acquires and presents the most recent weather observations taken at CXTO-AUTO (Toronto downtown).

`weatherdata observation CXTO-AUTO -s -2`

Acquires and presents the weather observations taken at CXTO-AUTO (Toronto downtown) over the past two hours.

`weatherdata observation CXTO-AUTO -s "2021-11-26 12:00 EST"`

Acquires and presents the weather observations taken at CXTO-AUTO from 12 PM, 26 November 2021, Eastern Time, to the present. Timestamps must be enclosed in double quotes.

The SWOB system typically retains historical weather for 30 days; Weatherdata will return an error if no observations are available,

`weatherdata observation CVVR-AUTO --starttime "2021-11-25 00:00 PST" --endtime "2021-11-26 00:00 PST"`

Acquires and presents a summary of weather observations recorded by CVVR-AUTO from midnight, 25 November 2021, Pacific Time to midnight, 26 November 2021, Pacific Time. Timestamps must be enclosed in double quotes.

The SWOB system typically retains historical weather for 30 days; Weatherdata will return an error if no observations are available,

`weatherdata observation "CYYZ-MAN CWTQ-AUTO CYVR-MAN CYYC-MAN"`

Acquires and presents the most recent weather observations from the major airports in Toronto, Montreal, Vancouver, and Calgary.

When multiple stations are specified, the station list must be surrounded in double quotes as shown above.

## Notes

A map of Environment Canada stations (with IATA IDs) is available via [GeoMet](https://api.weather.gc.ca/collections/swob-realtime/items).

Data provided by Environment Canada are subject to assorted terms of use, as [made available by Environment Canada](https://eccc-msc.github.io/open-data/msc-data/obs_station/readme_obs_insitu_en/). These terms are separate from the terms that apply to Weatherdata.

## Build Instructions

Weatherdata will not build out-of-the box. Third-party code is required to be downloaded to the `psychrometrics` directory. See `psychrometrics/READ-THIS.md` for instructions.

Weatherdata should build and run on any platform supported by Go where Sqlite is available.

## Revision History

### Release 0

Initial release.

### Release 1

Major rewrite with breaking changes to the user interface.

Observation collection engine redesigned to simplify support for EC partner data providers, and to add support for finding observations nearest to a given time.

Consolidated totalize and observation subcommands into one component.

Observations are downloaded in parallel.

Station list is now cached for 24 hours to reduce startup times.

Increased Go version to 1.19 to access newer dependencies without known CVEs.

## License

Copyright 2021, 2023 Coridon Henshaw

Permission is granted to all natural persons to execute, distribute, and/or modify this software (including its documentation) subject to the following terms:

1. Subject to point \#2, below, **all commercial use and distribution is prohibited.** This software has been released for personal and academic use for the betterment of society through any purpose that does not create income or revenue. *It has not been made available for businesses to profit from unpaid labor.*

2. Re-distribution of this software on for-profit, public use, repository hosting sites (for example: Github) is permitted provided no fees are charged specifically to access this software.

3. **This software is provided on an as-is basis and may only be used at your own risk.** This software is the product of a single individual's recreational project. The author does not have the resources to perform the degree of code review, testing, or other verification required to extend any assurances that this software is suitable for any purpose, or to offer any assurances that it is safe to execute without causing data loss or other damage.

4. **This software is intended for experimental use in situations where data loss (or any other undesired behavior) will not cause unacceptable harm.** Users with critical data safety needs must not use this software and, instead, should use equivalent tools that have a proven track record.

5. If this software is redistributed, this copyright notice and license text must be included without modification.

6. Distribution of modified copies of this software is discouraged but is not prohibited. It is strongly encouraged that fixes, modifications, and additions be submitted for inclusion into the main release rather than distributed independently.

7. This software reverts to the public domain 10 years after its final update or immediately upon the death of its author, whichever happens first.

8. Use of Weatherdata as a backend for large-scale scraping or to implement any service accessible to the general public is prohibited.