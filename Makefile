
all: build

build:
	@./scripts/build.sh

up:
	@./run/server/start.sh
	tail -f run/server/logs/server.log

down:
	@./run/server/bin/stop.sh
