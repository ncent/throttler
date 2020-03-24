HANDLER=handler
PACKAGE=package
ifdef DOTENV
	DOTENV_TARGET=dotenv
else
	DOTENV_TARGET=./.env
endif

.PHONY: build clean deploy


build: clean #test # generate_graphql test
	export $(grep -v '^#' ./.env | xargs)
	env
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -a -tags netgo -installsuffix netgo -o bin/throttler/throttle/email handlers/throttler/throttle/email/main.go
	chmod +x bin/throttler/throttle/email
	zip -j bin/throttler/throttle/email.zip bin/throttler/throttle/email

clean:
	-rm -rf ./bin

test: build
	go test github.com/ncent/arber-core/services/kinesis/client

deploy: build
	sls deploy --force --verbose

publish: deploy
	sls invoke --function publisher --path event.json
	ping -c 11 127.0.0.1 > nul
	sls logs -f publisher
	sls invoke --function consumer 
	ping -c 11 127.0.0.1 > nul
	sls logs -f consumer
	sls logs -f archiver
	rm nul
