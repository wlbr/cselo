LINKERFLAGS = -X main.Version=`git describe --tags --always --long --dirty` -X main.BuildTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`
PROJECTROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DBPATH	=	$(PROJECTROOT)db/

all: clean build

.PHONY: clean
clean:
	@echo Running clean job...
	rm -f coverage.txt
	rm -rf bin/
	rm -f main eloudp

build: #generate
	@echo Running build job...
	mkdir -p bin/linux bin/windows bin/mac
	# GOOS=linux go build  -ldflags "$(LINKERFLAGS)" -o bin/linux ./...
	# GOOS=windows go build  -ldflags "$(LINKERFLAGS)" -o bin/windows ./...
	GOOS=darwin go build  -ldflags "$(LINKERFLAGS)" -o bin/mac ./...

run:
	go run -ldflags "$(LINKERFLAGS)" cmd/eloudp/main.go -cfg cselo.ini -cslog data/short.log

test: #generate
	@echo Running test job...
	go test ./... -cover -coverprofile=coverage.txt

coverage: test
	@echo Running coverage job...
	go tool cover -html=coverage.txt

deploy: clean
	mkdir -p bin/linux
	GOOS=linux go build  -ldflags "$(LINKERFLAGS)" -o bin/linux ./...
	rsync -v --progress bin/linux/* wlbr:~/bin/


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

resetdb:
	psql postgres -f create-db.sql

recreatetables:
	psql cselo -U cseloapp -f create-tables.sql

