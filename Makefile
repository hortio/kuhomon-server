PID=/tmp/kuhomon-server.pid
APP=./kuhomon
TAG=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG
DB_SERVICE=$(shell kubectl get services  --no-headers=true --output='name' | grep --color=never 'cockroachdb-public' | cut -c9- )
export DB_SERVICE

.PHONY: serve
serve: $(APP) restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

.PHONY: kill
kill:
	@kill `cat $(PID)` || true

.PHONY: before
before:
	@echo "actually do nothing"

$(APP): 
	@go build -ldflags "-X main.version=$(TAG)" -o $@ ./server

.PHONY: pack
pack:
	GOOS=linux make $(APP)
	docker build -t kumekay/kuhomon-server:$(TAG) .

.PHONY: upload
upload: 
	docker push kumekay/kuhomon-server:$(TAG)

.PHONY: dev-deploy
dev-deploy:
	envsubst < ./k8s/dev-deployment.yaml | kubectl apply -f -

.PHONY: clean
clean:
	rm $(APP)

.PHONY: ship
ship: test clean pack upload dev-deploy

.PHONY: restart
restart: kill before $(APP)
	$(APP) & echo $$! > $(PID)

.PHONY: test
test:
	@go test -v ./...

.PHONY: cover
cover:
	@go test -coverprofile=c.out -v ./... 