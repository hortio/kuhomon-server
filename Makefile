PID=/tmp/kuhomon-server.pid
APP=./kuhomon
TAG=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

.PHONY: serve
serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

.PHONY: kill
kill:
	@kill `cat $(PID)` || true

.PHONY: before
before:
	@echo "actually do nothing"

$(APP): 
	@go build ./server -ldflags "-X main.version=$(TAG)" -o $@

.PHONY: pack
pack:
	GOOS=linux make $(APP)
	docker build -t kumekay/kuhomon-server:$(TAG) .

.PHONY: upload
upload: pack
	docker push kumekay/kuhomon-server:$(TAG)

.PHONY: restart
restart: kill before $(APP)
	$(APP) & echo $$! > $(PID)

.PHONY: test
test:
	@go test -v ./...

.PHONY: cover
cover:
	@go test -coverprofile=c.out -v ./... 