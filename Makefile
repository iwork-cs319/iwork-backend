run:
	DATABASE_URL=$$(heroku config:get DATABASE_URL -a icbc-go-api) \
	G_DRIVE_CREDENTIALS=$$(heroku config:get G_DRIVE_CREDENTIALS -a icbc-go-api) \
	PORT=8000 \
    go run .

test:
	TEST_DB_URL=$$(heroku config:get HEROKU_POSTGRESQL_AQUA_URL -a icbc-go-api) \
	G_DRIVE_CREDENTIALS=$$(heroku config:get G_DRIVE_CREDENTIALS -a icbc-go-api) \
 	go test -cover -race ./...

.PHONY: run run-db