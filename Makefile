
all: build

build:
	@./scripts/build.sh

up:
	@./test/server/bin/start.sh
	tail -f test/server/logs/server.log

down:
	@./test/server/bin/stop.sh
