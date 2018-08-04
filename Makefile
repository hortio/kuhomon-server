PID      = /tmp/kuhomon-server.pid
APP      = ./kuhomon

serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

kill:
	@kill `cat $(PID)` || true

before:
	@echo "actually do nothing"

$(APP): 
	@go build ./server -o $@

restart: kill before $(APP)
	$(APP) & echo $$! > $(PID)

test:
	@go test -v ./...

.PHONY: serve restart kill before test