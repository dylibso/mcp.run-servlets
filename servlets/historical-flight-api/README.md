# Historical Flight API Servlet

Queries flight APIs for info about past flights, 
using the following services:

- https://opensky-network.org
- https://adsbdb.com
- https://airport-data.com

Note:

> In OpenSky, Flights are updated by a batch process at night, 
> i.e., only flights from the previous day or earlier are available using this endpoint.

OpenSky requires a username, password pair, you can easily get
with [a free account](https://opensky-network.org/login?view=registration).

## Configuration

- The servlet expects the config keys `username` `password` to be provided.
- It requires network access to the following domains:
    - opensky-network.org
    - api.adsbdb.com
    - airport-data.com

## Usage

We expose the [OpenSky `arrival` and `departure` endpoints](https://openskynetwork.github.io/opensky-api/rest.html)
following their structure, except we also require a `requestType` field:

```json
{
  "requestType": "departure",
  "airport": "LIMC",
  "begin": "1701428400",
  "end": "1701435600"
}
```

Where `airport` is the ICAO code for the airport, and `begin`, `end` are UNIX timestamps
at UTC.

The return value is the contents of the return value for such endpoints.

We also expose the `aircraft` endpoint from [adsbdb.com](https://www.adsbdb.com).
This requires the `icao24` identifier of the aircraft, and its `callsign`,
which are always returned as part of the `arrival` and `departure` responses. The request looks like:

```json
{
  "icao24": "440170",
  "callsign": "EJU73BJ"
}
```

The result contains the returned value from the [adsbdb.com](https://www.adsbdb.com)
endpoint; we also automatically fetch and return the aircraft picture returned
in the [adsbdb.com](https://www.adsbdb.com) response in the `url_photo` field.
