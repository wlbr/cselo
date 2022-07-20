LINKERFLAGS = -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DBPATH	=	$(PROJECTROOT)db/

all: clean build

.PHONY: clean release
clean:
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/ release/
	rm -f main eloudp elogql csscores cselo.tgz cselo.zip cselo.dmg

dep:
	go install github.com/99designs/gqlgen

generate:
	gqlgen

build: dep generate
	@echo Running build job...
	mkdir -p bin/linux bin/windows bin/mac/x64 bin/mac/arm
	GOOS=linux GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/linux ./...
	GOOS=windows GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin GOARCH=amd64 go build  -ldflags "$(LINKERFLAGS)" -o bin/mac/x64 ./...
	GOOS=darwin GOARCH=arm64 go build  -ldflags "$(LINKERFLAGS)" -o bin/mac/arm ./...

run: generate
#	go run -ldflags "$(LINKERFLAGS)" cmd/eloudp/main.go -cfg cselo-local.ini -import data/latest.log
	go run -ldflags "$(LINKERFLAGS)" cmd/elogql/server.go -cfg cselo-local.ini

test: recreatetables
	go run -ldflags "$(LINKERFLAGS)" cmd/eloudp/main.go -cfg cselo-local.ini -import data/test.log
	@echo Running test job...
	go test ./... -cover -coverprofile=coverage.txt -cfg $(PROJECTROOT)cselo-local.ini -import $(PROJECTROOT)data/test.log -loglevel Error -player Jagger

coverage: test
	@echo Running coverage job...
	go tool cover -html=coverage.txt

deploy:
	ssh cselo mkdir -p '~/cselo/bin'
	rsync -v --progress bin/linux/* cselo:~/cselo/bin

release:
	mkdir -p release
	rm -rf release/
	$(eval VER=$(shell sh -c "bin/mac/x64/eloudp -version |cut -f 2 -d ' '"))
	cd bin && tar -zcpv -s /linux/cselo-linux-$(VER)/ -f ../release/cselo-linux-$(VER).tgz linux/*   # OSX
	cd bin/windows &&  zip -r -9 ../../release/cselo-win-$(VER).tgz *

initelodb: resetdb recreatetables

wipe: initelodb clean

newpostgresdb:
	initdb -D $(DBPATH)

startdb:
	#postgres -D $(DBPATH)
	#osascript -e 'tell app "Terminal" to do script "postgres -D $(DBPATH)"'
	execInNewITerm "postgres -D $(DBPATH)"

stopdb:
	pg_ctl stop -D $(DBPATH) -m fast

resetdb: stopdb startdb
	sleep 3
	psql postgres -f scripts/create-db.sql

recreatetables:
	psql cselo -U cseloapp -f scripts/create-tables.sql

analysis:
	psql cselo -f scripts/analysis.sql
