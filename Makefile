BUILD_DIR=builds
VERSION=0.2.5
GIT_SHA=$(shell git rev-parse --short HEAD)
PASSWORD=root
MSPASS=zebra-IT-32

DOCKER_MYSQL=mysql/mysql-server:5.7
DOCKER_POSTGRES=postgres:9.6
DOCKER_MSSQL=mcr.microsoft.com/mssql/server:2017-latest
DB_NAME=testdb

BUILD_PREFIX=builds/sql-dataset
LDFLAGS="-X main.version=$(VERSION) -X main.gitSHA=$(GIT_SHA)"

build-darwin:
	rm builds/* || true
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${BUILD_PREFIX}-darwin-amd64 -ldflags=${LDFLAGS}
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ${BUILD_PREFIX}-darwin-arm64 -ldflags=${LDFLAGS}
	
build-unix:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${BUILD_PREFIX}-linux-amd64 -ldflags=${LDFLAGS}
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o ${BUILD_PREFIX}-linux-386 -ldflags=${LDFLAGS}

build-win:
	set GOARCH=amd64; go build -o ${BUILD_PREFIX}-windows-10.0-amd64.exe -ldflags=${LDFLAGS}
	set GOARCH=386; go build -o ${BUILD_PREFIX}-windows-10.0-386.exe -ldflags=${LDFLAGS}

pull-docker-images:
	docker pull ${DOCKER_MYSQL}
	docker pull ${DOCKER_POSTGRES}
	docker pull ${DOCKER_MSSQL}

run-containers:
	docker rm -f sd-mysql sd-postgres sd-mssql || true
	docker run --name sd-mssql -e ACCEPT_EULA=Y -e SA_PASSWORD=${MSPASS} -p 1433:1433 -d ${DOCKER_MSSQL} || true
	# MySQL
	docker run --name sd-mysql -e MYSQL_ROOT_PASSWORD=${PASSWORD} -p 3307:3306 -d ${DOCKER_MYSQL} || true
	scripts/wait_for_mysql sd-mysql
	# Postgres
	docker run --name sd-postgres -e POSTGRES_PASSWORD=${PASSWORD} -p 5433:5432 -d ${DOCKER_POSTGRES} || true
	scripts/wait_for_postgres sd-postgres

setup-db:
	# Mysql ensure root can access from anywhere
	docker exec -it sd-mysql mysql -uroot -proot -e "CREATE USER 'root'@'%' IDENTIFIED BY '${PASSWORD}'" || true
	docker exec -it sd-mysql mysql -uroot -proot -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%'" || true
	docker exec -it sd-mysql mysql -uroot -proot -e "CREATE DATABASE ${DB_NAME};" || true
	# Postgres db creation
	docker exec --user postgres -it sd-postgres psql -c "CREATE DATABASE ${DB_NAME}" || true
	# MSSQL db creation
	docker exec -it sd-mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P ${MSPASS} -Q "CREATE DATABASE testdb" || true


#
# Running the more full integration tests requires docker containers to be
# running database servers. The setup target helps with call the other targets
# which download docker images, run-containers and setup-db
#
setup: pull-docker-images run-containers setup-db

#
# Once you run setup target the docker containers will continously
# run the tests should handle different states already with SQL fixtures
# using REPLACE for example instead of INSERT, but eventually you will want
# to stop the containers running this target supports that
#
teardown:
	@docker stop sd-mysql > /dev/null 2>&1 || true
	@docker rm -v sd-mysql > /dev/null 2>&1 || true
	@docker stop sd-mssql > /dev/null 2>&1 || true
	@docker rm -v sd-mssql > /dev/null 2>&1 || true
	@docker stop sd-postgres > /dev/null 2>&1 || true
	@docker rm -v sd-postgres > /dev/null 2>&1 || true

test:
	MYSQL_URL="root:${PASSWORD}@tcp(localhost:3307)/testdb?parseTime=true" \
	POSTGRES_URL=postgres://postgres:${PASSWORD}@localhost:5433/testdb?sslmode=disable \
	MSSQL_URL="odbc:server=localhost;port=1433;user id=sa;password=${MSPASS};database=${DB_NAME}" \
	go test ./... -race -covermode=atomic -coverprofile=coverage.txt
