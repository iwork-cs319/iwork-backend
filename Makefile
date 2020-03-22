run:
	DATABASE_URL=$$(heroku config:get DATABASE_URL -a icbc-go-api) \
	REDIS_URL=$$(heroku config:get REDIS_URL -a icbc-go-api) \
	G_DRIVE_CREDENTIALS=$$(heroku config:get G_DRIVE_CREDENTIALS -a icbc-go-api) \
	MICROSOFT_CLIENT_ID=$$(heroku config:get MICROSOFT_CLIENT_ID -a icbc-go-api) \
	MICROSOFT_GRAPH_SCOPE=$$(heroku config:get MICROSOFT_GRAPH_SCOPE -a icbc-go-api) \
	MICROSOFT_CLIENT_SECRET=$$(heroku config:get MICROSOFT_CLIENT_SECRET -a icbc-go-api) \
	ADMIN_USER_ID=$$(heroku config:get ADMIN_USER_ID -a icbc-go-api) \
	PORT=8000 \
    go run .

test: test-routes test-db

test-routes:
	TEST_DB_URL=$$(heroku config:get HEROKU_POSTGRESQL_AQUA_URL -a icbc-go-api) \
	G_DRIVE_CREDENTIALS=$$(heroku config:get G_DRIVE_CREDENTIALS -a icbc-go-api) \
 	go test -count=1 -cover -race ./routes/...

test-db:
	TEST_DB_URL=$$(heroku config:get HEROKU_POSTGRESQL_AQUA_URL -a icbc-go-api) \
	REDIS_URL=$$(heroku config:get REDIS_URL -a icbc-go-api) \
	G_DRIVE_CREDENTIALS=$$(heroku config:get G_DRIVE_CREDENTIALS -a icbc-go-api) \
 	go test -count=1 -cover -race ./db/...

.PHONY: run