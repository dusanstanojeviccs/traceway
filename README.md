# Temporary documentation

## Project structure
There are 3 projects:

backend
clients/go-client
frontend

## Backend 

Represents the main backend for traceway. Stores stack traces, endpoint info and metrics into a clickhouse DB. Has an api that can be accessed from the frontend project. 

### How to run:

To run the backend you need to:
1 - install golang
2 - install clickhosue
3 - create a clickhouse database (eg: traceway)
4 - create a clickhouse user (eg: if one does not exist default/empty password)
5 - create the .env file (instructions below)
6 - open the backend folder and run `go run .`

### Backend .env file
To run it you need to create a .env file in the backend project locally with your own clickhouse credentials.

Example .env:

```
TOKEN="demotoken"
APP_TOKEN="nice"
CLICKHOUSE_SERVER=localhost:9000
CLICKHOUSE_DATABASE=traceway
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=
CLICKHOUSE_TLS=false
```

The TOKEN value is what the clients/go-client will use to report while APP_TOKEN is what the frontend uses to access the backend.

## Frontend

This is a work in progress, it's in the wireframing/designing stages. It's a sveltekit app that is expected to run in the SPA mode (client running only). 

If you really want to run it:
1 - install node/npm
2 - run `npm install`
3 - run `npm run dev`

## clients/go-client

This is a client that users would include in their app. Right now it only supports the gin framework.

To run the basic devtesting app you can cd into `./devtesting` and run `go run .` this will start the server located in devtesting.go the keycode here is:
```
router.Use(tracewaygin.New(
		"tracewaydemo",
		"demotoken@http://localhost:8082/api/report",
		traceway.WithDebug(true),
	))
```

Which sets up the code to upload to localhost:8082 with the token 'demotoken' (this has to match your .env from your backend). After that the code will start uploading your metrics/stracktraces to the backend, you should be able to see it in tableplus or any other clickhouse client you like using.
