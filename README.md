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

###Setup
Implementations of the following interfaces must be provided in order to initialize the back-end:
- `github.com/BoSunesen/pointsofinterest/webapi/logging.Logger`
- `github.com/BoSunesen/pointsofinterest/webapi/factories.ContextFactory`
- `github.com/BoSunesen/pointsofinterest/webapi/factories.ClientFactory`
- `github.com/BoSunesen/pointsofinterest/webapi/factories.WorkerFactory`

See https://godoc.org/github.com/BoSunesen/pointsofinterest/webapi/logging
for more information on the Logger interface and
see https://godoc.org/github.com/BoSunesen/pointsofinterest/webapi/factories
for more information on the factory interfaces.
The project includes a main package that starts the back-end using very simple implementations
of the interfaces. See [Links](#links) for a Google App Engine initialization project.

###Future development
- Remove dependency on `golang.org/x/net/context`
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

Me: https://www.linkedin.com/in/bosunesen
