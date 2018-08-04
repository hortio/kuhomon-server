# Server for kuhomon - simple home air monitor: T, H, P & CO2 levels

[![CircleCI](https://circleci.com/gh/kumekay/kuhomon-server.svg?style=svg)](https://circleci.com/gh/kumekay/kuhomon-server)
[![Maintainability](https://api.codeclimate.com/v1/badges/1ccf0ab6df9087fd6c4f/maintainability)](https://codeclimate.com/github/kumekay/kuhomon-server/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/1ccf0ab6df9087fd6c4f/test_coverage)](https://codeclimate.com/github/kumekay/kuhomon-server/test_coverage)


## HTTP API

### `GET /measurements/:deviceID`  
Returns 10-minutes data for last 24 hours

#### Parameters
- `device_id` - **required** Device UUID

#### Headers
- `Device-Read-Token` -  **required** permanent token to read data stored on server for this device.

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

### `POST /measurements/:deviceID`
Creates a data point in DB

#### Parameters

- `device_id` - **required** Device UUID
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