all: help

APP=gallery-service
WORKDIR=/go/src/github.com/hyperjiang/gallery-service/app
NETWORK=docker-backend

GOLANG_IMG=golang:1.11
MARIADB_IMG=mariadb:10
MIGRATE_IMG=migrate/migrate
REDIS_IMG=redis:4.0-alpine

MARIADB=mariadb
MARIADB_DATA=${HOME}/data/mysql
MARIADB_ROOT_PW=hello123
REDIS=redis
REDIS_DATA=${HOME}/data/redis


#####   ####   ####  #    # ###### #####
#    # #    # #    # #   #  #      #    #
#    # #    # #      ####   #####  #    #
#    # #    # #      #  #   #      #####
#    # #    # #    # #   #  #      #   #
#####   ####   ####  #    # ###### #    #

help:
	###########################################################################################################
	# [DOCKER]
	# network                   - create docker bridge network
	# ps                        - docker ps -a (list all containers)
	# up                        - run docker-compose up (run up the container)
	# down                      - run docker-compose down (shutdown the container)
	# kill                      - run docker-compose rm -f (kill and rm the container)
	# restart                   - run docker-compose restart (restart the container)
	# logs                      - tail the container logs
	# clean                     - run: make rm-con, make rm-img
	# stats                     - show container stats (CPU%, memory, etc)
	# stats-all                 - show all containers stats (CPU%, memory, etc)
	# rm-con                    - remove all dead containers (non-zero Exited)
	# rm-img                    - remove all <none> images/layers
	#
	# [DOCKER PUBLIC]
	# mariadb                   - run up an MariaDB container
	# mariadb-down              - remove the MariaDB container
	# conndb                    - connect to MariaDB using root
	# conndb-app                - connect to MariaDB using APP-user
	# redis                     - run up an redis container
	# redis-down                - remove the redis container
	# connredis                 - connect to redis
	#
	# [PROJECT]
	# initdb                    - initialize DBs (current project need to provide ./db/ and ./db/migrations/)
	# migration                 - run db migration, you can specify v={version} to migrate to the {version} you want
	# gomod-init                - run "go mod init"
	# gomod-verify              - run "go mod verify"
	# gomod-vendor              - run "go mod vendor"
	# gotest                    - run "go test" in container
	# gofmt                     - format golang source code (change the codes)
	# godoc                     - serve godoc for source code and open in browser (Mac)
	###########################################################################################################
	@echo "Enjoy!"

network:
	docker network create -d bridge ${NETWORK} || true

ps:
	docker ps -a

up: network
	docker-compose -f docker-compose.yml up --build -d

down:
	docker-compose -f docker-compose.yml down

kill:
	docker-compose -f docker-compose.yml kill && \
	docker-compose -f docker-compose.yml rm -f

restart:
	docker-compose -f docker-compose.yml restart

logs:
	docker-compose -f docker-compose.yml logs -f --tail=10

clean-vendor:
	rm -rf ./app/vendor/*

clean: rm-img clean-vendor

stats:
	docker stats ${APP}

stats-all:
	docker stats `docker ps -a | sed 1d | awk '{print $$NF}'`

rm-con:
	deads=$$(docker ps -a | sed 1d | grep "Exited " | grep -v "Exited (0)" | awk '{print $$1}'); if [ "$$deads" != "" ]; then docker rm -f $$deads; fi

rm-img: rm-con
	none=$$(docker images | sed 1d | grep "^<none>" | awk '{print $$3}'); if [ "$$none" != "" ]; then docker rmi $$none; fi

gomod-init:
	# will create go.mod
	docker run --rm -t -v "${PWD}/app:${WORKDIR}" -w "${WORKDIR}" -e GO111MODULE=on "${GOLANG_IMG}" go mod init

gomod-verify:
	docker run --rm -t -v "${PWD}/app:${WORKDIR}" -w "${WORKDIR}" -e GO111MODULE=on "${GOLANG_IMG}" go mod verify

gomod-vendor:
	# will create go.sum
	docker run --rm -t -v "${PWD}/app:${WORKDIR}" -w "${WORKDIR}" -e GO111MODULE=on "${GOLANG_IMG}" go mod vendor -v

gotest: network
	docker run --rm -it --net=host -v "${PWD}/app:${WORKDIR}" -v "${PWD}/runtime:/runtime" -w "${WORKDIR}" "${GOLANG_IMG}" bash -c "go test -cover ./..."

gofmt:
	docker run --rm -t -v "${PWD}/app:${WORKDIR}" -w "${WORKDIR}" "${GOLANG_IMG}" gofmt -w .

godoc:
	@docker rm -f go-$(APP)-godoc || true
	@while [ true ]; do \
		PORT=$$(( ( RANDOM % 60000 )  + 1025 )); \
		nc -z -vv localhost $$PORT >/go/null 2>/go/null || break; \
	done; \
	docker run --rm -it --net=${NETWORK} -v "${PWD}/app:${WORKDIR}" -w "${WORKDIR}" \
		--expose=80 -p "$$PORT:80" \
		-e ENV=go -d --name=go-$(APP)-godoc "${GOLANG_IMG}" bash -c \
		"godoc -http :80";\
	docUrl="http://localhost:$$PORT/pkg/$(APP)"; \
	sleep 3s ; os=$$(uname -s); if [ "$$os" == Darwin ]; then echo "Opening browser..."; open "$$docUrl"; fi;

mariadb: network
	@mkdir -p ${MARIADB_DATA}
	docker rm -f ${MARIADB} || true; \
	docker run -d --name=${MARIADB} --hostname=${MARIADB} --restart=always \
		--net=${NETWORK} \
		-e MYSQL_ROOT_PASSWORD=${MARIADB_ROOT_PW} \
		-v ${MARIADB_DATA}:/var/lib/mysql \
		${MARIADB_IMG}
	# add "-p 3306:3306" if you want to publish the port

mariadb-down:
	docker rm -f ${MARIADB} || true

initdb:
	@a=$(eval TMP := "${HOME}/tmp")
	@mkdir -p ${TMP}
	@DB_PATH="${PWD}/db"; \
	if ! [ -d "$${DB_PATH}/" ]; then echo "ERR: $${DB_PATH}/ doesn't exist"; exit 1; fi; \
	if ! [ -f "$${DB_PATH}/db.rc" ]; then echo "ERR: $${DB_PATH}/db.rc doesn't exist"; exit 1; fi; \
	if ! [ -d "$${DB_PATH}/migrations/" ]; then echo "ERR: $${DB_PATH}/migrations/ doesn't exist"; exit 1; fi; \
	source "$${DB_PATH}/db.rc"; echo "DB: $$DB"; \
	a=$(eval TMP := $(shell mktemp -d ${TMP}/XXXXXX)); trap 'rm -rf $(TMP)' EXIT; \
	echo ""; echo "--------------------Drop & Create DB-----------------------------------"; \
	echo "DROP DATABASE IF EXISTS \`$$DB\`;" > $(TMP)/db.sql; \
	echo "CREATE DATABASE \`$$DB\` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;" >> $(TMP)/db.sql; \
	docker run -it --rm --network=${NETWORK} -v $(TMP)/db.sql:/file.sql ${MARIADB_IMG} bash -c "mysql -A -h${MARIADB} -uroot -p${MARIADB_ROOT_PW} < /file.sql"; \
	echo "DB root was created to manage schema change, you can connect with: make conndb"; \
	echo ""; echo "--------------------Creating schema------------------------------------"; \
	docker run -it --rm --network=${NETWORK} -v $${DB_PATH}/migrations:/migrations "${MIGRATE_IMG}" -path /migrations -database "mysql://root:${MARIADB_ROOT_PW}@tcp(${MARIADB}:3306)/$$DB" up; \
	echo "App user was created to simulate the connection from the app, you can connect with: make conndb-app"; \
	echo ""; echo "--------------------Creating App user(s)-------------------------------"; \
	echo "DROP USER IF EXISTS '$$APP_USER'@'%';" > $(TMP)/user.sql; \
	echo "CREATE USER '$$APP_USER'@'%' IDENTIFIED BY '$$APP_PASSWORD';" >> $(TMP)/user.sql; \
	docker run -it --rm --network=${NETWORK} -v $(TMP)/user.sql:/file.sql ${MARIADB_IMG} bash -c "mysql -A -h${MARIADB} -uroot -p${MARIADB_ROOT_PW} < /file.sql"; \
	echo ""; echo "--------------------Running customized SQL files----------------------"; \
	for sql in `ls $${DB_PATH}/*.sql | sort` ; do \
		echo "$$sql..."; \
		time docker run -it --rm --network=${NETWORK} -v "$$sql":/file.sql ${MARIADB_IMG} bash -c "mysql -A -h${MARIADB} -uroot -p${MARIADB_ROOT_PW} $$DB < /file.sql" ; \
		echo "$$sql...[OK]"; \
	done; \
	echo "initdb completed";

migration:
	@DB_PATH="${PWD}/db"; \
	if ! [ -d "$${DB_PATH}/" ]; then echo "ERR: $${DB_PATH}/ doesn't exist"; exit 1; fi; \
	if ! [ -f "$${DB_PATH}/db.rc" ]; then echo "ERR: $${DB_PATH}/db.rc doesn't exist"; exit 1; fi; \
	if ! [ -d "$${DB_PATH}/migrations/" ]; then echo "ERR: $${DB_PATH}/migrations/ doesn't exist"; exit 1; fi; \
	source "$${DB_PATH}/db.rc"; echo "DB: $$DB"; \
	if [ "$(v)" == "" ]; then \
	    docker run -it --rm --network=${NETWORK} -v "$${DB_PATH}"/migrations:/migrations "${MIGRATE_IMG}" -path /migrations -database "mysql://root:${MARIADB_ROOT_PW}@tcp(${MARIADB}:3306)/$$DB" up; \
	else \
	    docker run -it --rm --network=${NETWORK} -v "$${DB_PATH}"/migrations:/migrations "${MIGRATE_IMG}" -path /migrations -database "mysql://root:${MARIADB_ROOT_PW}@tcp(${MARIADB}:3306)/$$DB" goto $(v); \
	fi; \
	echo "migration completed";

conndb:
	@if ! [ -d "./db/" ]; then echo "ERR: ./db/ doesn't exist"; exit 1; fi; \
	if [ -f "./db/db.rc" ]; then source "./db/db.rc"; fi; \
	docker run -it --rm --network=${NETWORK} ${MARIADB_IMG} bash -c "mysql -A --default-character-set=utf8 -h${MARIADB} -uroot -p${MARIADB_ROOT_PW}"

conndb-app:
	@if ! [ -d "./db/" ]; then echo "ERR: ./db/ doesn't exist"; exit 1; fi; \
	if [ -f "./db/db.rc" ]; then source "./db/db.rc"; fi; \
	docker run -it --rm --network=${NETWORK} ${MARIADB_IMG} bash -c "mysql -A --default-character-set=utf8 -h${MARIADB} -u$$APP_USER -p$$APP_PASSWORD $$DB"

redis: network
	@mkdir -p ${REDIS_DATA}
	# add "-p 6379:7379" if you want to publish the port
	docker rm -f ${REDIS} || true
	docker run -d --name=${REDIS} --hostname=${REDIS} --restart=always \
		--net=${NETWORK} \
		-v ${REDIS_DATA}:/data \
		${REDIS_IMG}

redis-down:
	docker rm -f ${REDIS} || true

connredis:
	docker run --rm -it --net=${NETWORK} ${REDIS_IMG} redis-cli -h ${REDIS}
