.PHONY: build clean deploy

build:
	dep ensure -v

	env GOOS=linux go build -ldflags="-s -w" -o bin/game game/game.model.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wsclient wsclient/wsclient.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dbclient dbclient/dbclient.go

	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login login/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/auth auth/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/ws ws/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wsconnect wsconnect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wsauthenticate wsauthenticate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wsdisconnect wsdisconnect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/wsactionhandler wsactionhandler/main.go


clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
