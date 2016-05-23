# Points of Interest
Server that caches POI data in the background
in order to serve requests as fast as possible.

###Goals
- Provide faster access to POI data
- Protect against temporary downtime at the data provider
- Horizontal scalability
- Parsing POI data to a more workable format
- Filtering POI data according to client request parameters
- Use of data provider application token without exposing the token in client source files

###Future development
- Parse data in the background during cache refresh
- Automatic background cache refresh, to keep the cache fresh during downtime
- More types of POI data
- More filtering options

###My experience with Go
This was my first time coding in Go, hopefully it will not be my last :-)

###Links
Two other projects are related to this:
- Initialization of the server for Google App Engine: https://github.com/BoSunesen/pointsofinterestlauncher
- A minimal front-end: https://github.com/BoSunesen/FoodTrucks

The hosted application can be found here:
- Back-end: https://points-of-interest-1308.appspot.com/poi
- Front-end: https://points-of-interest-map.appspot.com

