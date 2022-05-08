LINKERFLAGS = -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DBPATH	=	$(PROJECTROOT)db/

all: clean build

.PHONY: clean
clean:
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/ release/
	rm -f main eloudp elogql csscores cselo.tgz cselo.zip cselo.dmg

generate:
	gqlgen

dep:
	go install github.com/99designs/gqlgen

build: dep generate
	@echo Running build job...
	mkdir -p bin/linux bin/windows bin/mac
	GOOS=linux go build  -ldflags "$(LINKERFLAGS)" -o bin/linux ./...
	GOOS=windows go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin go build  -ldflags "$(LINKERFLAGS)" -o bin/mac ./...

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
