#Declarations
# FLAG is used to notify the shell scripts that they are being invoked from a Makefile
FLAG="TRUE"
# BUILDNAME is the standard name tag used throughout, change it as required for the specific project
BUILDNAME="go-boilerplate"
# ENV is the standard env name passed for running in docker
ENV="docker"

#Commands
build: . # Creates binary
	go build -o ${BUILDNAME}

clean: . # Cleans unecessary junk
	rm -Rf ${BUILDNAME} _reports

re-init: . # Completely cleans current repo (vendor, gomod, gosum and junk files) and re initializes it
	rm -Rf ${BUILDNAME} _reports vendor go.mod go.sum
	go mod init ${BUILDNAME}
	go mod vendor

vendor: . # Initializes dependencies only
	go mod vendor

lint: # Runs go linter
	./scripts/lint.sh ${FLAG}

run: main.go # Runs app
	TIER=development go run main.go
		
test: # Runs test scripts
	./scripts/test.sh ${FLAG}

# Docker based command , needs docker to be installed on the system to work
docker-build: # Builds docker image
	docker image rm -f ${BUILDNAME}
	./docker/build ${BUILDNAME} ${FLAG}

docker-run: # Runs already built docker image (will fail if image has not been built)
	./docker/run ${BUILDNAME} ${ENV} ${FLAG}

docker-clean: # Cleans existing the docker image using the buildname
	docker image rm -f ${BUILDNAME}

