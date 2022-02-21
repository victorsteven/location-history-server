# location-history-server

Implementation of an in-memory server to store, retrieve and delete the several values of location based on an order input.

There are three endpoints implemented:

## API Endpoints

| Method | Description                                                           |   Endpoints                     |
| ------ | ----------------------------------------------------------------------| ------------------------------- |
| GET    | Get the locations based on order_id provided and max value provided   | `/location/:order_id?max=:max`  |
| POST   | Create a location based on order_id provided                          | `/location/:order_id/now`       |
| DELETE | Remove list of locations based on order_id provided                   | `/location/:order_id`           |


## Running the application
First build the application
```
go build
```
Then run the application
```
./location-history-server
```

## Tests
Run all test cases from the root directory:
```
go test -v ./...
```
