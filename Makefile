LINKERFLAGS = -X main.Version=`git describe --tags --always --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DBPATH	=	$(PROJECTROOT)db/

all: clean build

.PHONY: clean
clean:
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/ release/
	rm -f main eloudp csscores cselo.tgz cselo.zip cselo.dmg


build: #generate
	@echo Running build job...
	mkdir -p bin/linux bin/windows bin/mac
	GOOS=linux go build  -ldflags "$(LINKERFLAGS)" -o bin/linux ./...
	GOOS=windows go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin go build  -ldflags "$(LINKERFLAGS)" -o bin/mac ./...

run:
	go run -ldflags "$(LINKERFLAGS)" cmd/eloudp/main.go -cfg cselo-local.ini -cslog data/latest.log

test: recreatetables
	go run -ldflags "$(LINKERFLAGS)" cmd/eloudp/main.go -cfg cselo-local.ini -cslog data/test.log
	@echo Running test job...
	go test ./... -cover -coverprofile=coverage.txt -cfg $(PROJECTROOT)cselo-local.ini -cslog $(PROJECTROOT)data/test.log -loglevel Error -player Jagger

analysis:
	psql cselo -f scripts/analysis.sql

coverage: test
	@echo Running coverage job...
	go tool cover -html=coverage.txt

deploy:
	mkdir -p release
	$(eval VER=$(shell sh -c "bin/mac/eloudp -version |cut -f 2 -d ' '"))
	cd bin && tar -zcpv -s /linux/cselo-linux-$(VER)/ -f ../release/cselo-linux-$(VER).tgz linux/*
	cd bin/windows &&  zip -r -9 ../../release/cselo-win-$(VER).zip *

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

