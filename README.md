# Server for kuhomon - simple home air monitor: T, H, P & CO2 levels

## HTTP API

### `GET /measurements`  
Returns 10-minutes data for last 24 hours

#### Parameters
none

#### Headers
`Device-Read-Token` -  **required** permanent token to read data stored on server for this device.

#### Response
`200 OK` with list of measurements

```json
{
  "measurements":[
    {
      "h":29.2,
      "t": 300.1,
      "p":100000,
      "co2":435,
      "at":"2012-04-23T18:25:43.511Z"
    },
    ...
  ]
}
```

### `POST /measurements`
Creates a data point in DB

#### Parameters

- `h` - relative humidity, float 0-100
- `t` - temperature, in Kelvins, float  0 - 400
- `p` - pressure, in Pascals, integer 0 - 200000
- `co2` - CO2 level, ppm, integer 0 - 10000

```json
{
  "h":29.2,
  "t": 300.1,
  "p":100000,
  "co2":435
}
```

#### Headers
`Device-Write-Token` -  **required** permanent token to read data stored on server for this device.

#### Response

`201 Created` with empty body