run:
	PORT=8000 go run .
run-db:
	DATABASE_URL=$(heroku config:get DATABASE_URL -a icbc-go-api) PORT=8000 go run .

.PHONY: run run-db