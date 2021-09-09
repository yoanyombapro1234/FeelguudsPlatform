# Makefile for releasing feelguuds platform
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=feelguuds-platform
DOCKER_REPOSITORY:=feelguuds
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')
EXTRA_RUN_ARGS?=
# TMP_BASE is the base directory used for TMP.
# Use TMP and not TMP_BASE as the temporary directory.
TMP_BASE := .tmp
# TMP_COVERAGE is where we store code coverage files.
TMP_COVERAGE := $(TMP_BASE)/coverage

# runs an instance of the service locally
.PHONY: run
run:
	go run -ldflags "-s -w -X github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version.REVISION=$(GIT_COMMIT)" cmd/feelguuds_platform/* \
	--level=debug --grpc-port=9999 --backend-url=https://httpbin.org/status/401 --backend-url=https://httpbin.org/status/500 \
	--ui-logo=https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif $(EXTRA_RUN_ARGS)

# builds the service as an executable
.PHONY: build
build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/feelguuds_platform ./cmd/feelguuds_platform/*
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/feelguuds_platform_cli ./cmd/feelguuds_platform_cli/*

# formats the service's codebase
.PHONY: fmt
fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

# builds various associated helm charts
.PHONY: build-charts
build-charts:
	helm lint charts/*
	helm package charts/*

# builds a docker container in which the service's executable will run
.PHONY: build-container
build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

# builds the final part of the docker build
.PHONY: push-base
build-base:
	docker build -f Dockerfile.base -t $(DOCKER_REPOSITORY)/feelguuds_platform-base:latest .

# push the base part of the docker build
.PHONY: push-base
push-base: build-base
	docker push $(DOCKER_REPOSITORY)/feelguuds_platform-base:latest

# test the docker container (endpoint test) TODO: expand this -- perform a suite of operations against the container
.PHONY: test-container
test-container:
	@docker rm -f feelguuds_platform || true
	@docker run -dp 9898:9898 --name=feelguuds_platform $(DOCKER_IMAGE_NAME):$(VERSION)
	@docker ps
	@TOKEN=$$(curl -sd 'test' localhost:9898/token | jq -r .token) && \
	curl -sH "Authorization: Bearer $${TOKEN}" localhost:9898/token/validate | grep test

# push the container to some docker registry
.PHONY: push-container
push-container:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):latest
	docker push quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker push quay.io/$(DOCKER_IMAGE_NAME):latest

# set the version of the service
.PHONY: version-set
version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	sed -i '' "s/tag: $$current/tag: $$next/g" charts/feelguuds_platform/values.yaml && \
	sed -i '' "s/tag: $$current/tag: $$next/g" charts/feelguuds_platform/values-prod.yaml && \
	sed -i '' "s/appVersion: $$current/appVersion: $$next/g" charts/feelguuds_platform/Chart.yaml && \
	sed -i '' "s/version: $$current/version: $$next/g" charts/feelguuds_platform/Chart.yaml && \
	sed -i '' "s/feelguuds_platform:$$current/feelguuds_platform:$$next/g" kustomize/deployment.yaml && \
	sed -i '' "s/feelguuds_platform:$$current/feelguuds_platform:$$next/g" deploy/webapp/frontend/deployment.yaml && \
	sed -i '' "s/feelguuds_platform:$$current/feelguuds_platform:$$next/g" deploy/webapp/backend/deployment.yaml && \
	sed -i '' "s/feelguuds_platform:$$current/feelguuds_platform:$$next/g" deploy/bases/frontend/deployment.yaml && \
	sed -i '' "s/feelguuds_platform:$$current/feelguuds_platform:$$next/g" deploy/bases/backend/deployment.yaml && \
	echo "Version $$next set in code, deployment, chart and kustomize"

# define a release of the current code base
.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)

# generate swagger docs for the service
.PHONY: swagger
swagger:
	go get github.com/swaggo/swag/cmd/swag
	cd pkg/api && $$(go env GOPATH)/bin/swag init -g server.go

# terminate the current set of docker containers
.PHONY: kill-containers
kill-containers:
	docker-compose -f docker-compose.yaml -f \
					  docker-compose.authn.yaml -f \
					  docker-compose.jaeger.yaml -f \
					  docker-compose.merchant.dep.yaml -f \
					  docker-compose.shopper.dep.yaml down

setup-authn-deps:
	./scripts/run_authn.sh

start-deps: setup-authn-deps
	docker-compose -f docker-compose.yaml -f \
				   	  docker-compose.jaeger.yaml -f \
				   	  docker-compose.merchant.dep.yaml -f \
				   	  docker-compose.shopper.dep.yaml up --remove-orphans --detach

# start docker containers in the backgound
.PHONY: start-containers
start-containers:
	docker-compose -f docker-compose.yaml -f \
					  docker-compose.authn.yaml -f \
					  docker-compose.jaeger.yaml -f \
					  docker-compose.merchant.dep.yaml -f \
					  docker-compose.shopper.dep.yaml config
	docker-compose -f docker-compose.yaml -f \
					  docker-compose.authn.yaml -f \
				   	  docker-compose.jaeger.yaml -f \
				   	  docker-compose.merchant.dep.yaml -f \
				   	  docker-compose.shopper.dep.yaml up --remove-orphans --detach

# start docker containers with logs running in the foreground
.PHONY: start-containers-live
start-containers-live:
	docker-compose -f docker-compose.yaml -f \
					  docker-compose.authn.yaml -f \
					  docker-compose.jaeger.yaml -f \
					  docker-compose.merchant.dep.yaml -f \
					  docker-compose.shopper.dep.yaml config
	docker-compose -f docker-compose.yaml -f \
				  	  docker-compose.authn.yaml -f \
				   	  docker-compose.jaeger.yaml -f \
				   	  docker-compose.merchant.dep.yaml -f \
				   	  docker-compose.shopper.dep.yaml up --remove-orphans

# Cover runs go_test on GO_PKGS and produces code coverage in multiple formats.
# A coverage.html file for human viewing will be at $(TMP_COVERAGE)/coverage.html
# This target will echo "open $(TMP_COVERAGE)/coverage.html" with TMP_COVERAGE
# expanded  so that you can easily copy "open $(TMP_COVERAGE)/coverage.html" into
# your terminal as a command to run, and then see the code coverage output locally.
.PHONY: cover
cover:
	$(AT) rm -rf $(TMP_COVERAGE)
	$(AT) mkdir -p $(TMP_COVERAGE)
	go test $(GO_TEST_FLAGS) -json -cover -coverprofile=$(TMP_COVERAGE)/coverage.txt $(GO_PKGS) | tparse
	$(AT) go tool cover -html=$(TMP_COVERAGE)/coverage.txt -o $(TMP_COVERAGE)/coverage.html
	$(AT) echo
	$(AT) go tool cover -func=$(TMP_COVERAGE)/coverage.txt | grep total
	$(AT) echo
	$(AT) echo Open the coverage report:
	$(AT) echo open $(TMP_COVERAGE)/coverage.html
	$(AT) if [ "$(OPEN_COVERAGE_HTML)" == "1" ]; then open $(TMP_COVERAGE)/coverage.html; fi

.PHONY: go-mod
go-mod:
	go list -m -u all

ci-test: start-deps
	docker ps -a
	docker logs authentication_service
	go get github.com/mfridman/tparse
	go test -v -race ./... -json -cover  -coverprofile cover.out | tparse -all -top
	go tool cover -html=cover.out

.PHONY: test
test: start-containers
	./scripts/run_authn.sh
	echo "starting unit tests and integration tests"
	docker ps -a
	docker logs authentication_service
	go get github.com/mfridman/tparse
	go test -v -race ./... -json -cover  -coverprofile cover.out | tparse -all -top
	go tool cover -html=cover.out

# runs service load tests
.PHONY: load-test
load-test: start-containers
	cd ./load_test && ./load.sh
	cd ../

# profile the serivice
.PHONY: install-pprof
install-pprof:
	go get -u github.com/google/pprof

## Profiling (https://blog.golang.org/pprof)
# profiles cpu usage
.PHONY: profile-cpu
profile-cpu: install-pprof start-containers
	go tool pprof http://localhost:9898/debug/pprof/profile\?seconds\=20

# profile heap allocations
.PHONY: profile-heap
profile-heap: install-pprof start-containers
	go tool pprof http://localhost:9898/debug/pprof/heap

# profile block go routines
.PHONY: install-pprof profile-goroutines
profile-goroutines: start-containers
	go tool pprof http://localhost:9898/debug/pprof/block

# start minikube cluster
start-minikube:
	minikube config set memory 16384
	minikube start

# deploy artifacts to minikube cluster
kube-deploy: start-minikube
	kubectl apply -f ./kubernetes/auth-service/auth-service-db-deployment.yaml \
				  -f ./kubernetes/auth-service/auth-service-db-service.yaml \
				  -f ./kubernetes/auth-service/auth-service-deployment.yaml \
				  -f ./kubernetes/auth-service/auth-service-env-configmap.yaml \
				  -f ./kubernetes/auth-service/auth-service-redis-deployment.yaml \
				  -f ./kubernetes/auth-service/auth-service-redis-service.yaml \
				  -f ./kubernetes/auth-service/auth-service-service.yaml
	kubectl apply -f ./kubernetes/feelguuds-platform/feelguuds-platform-claim0-persistentvolumeclaim.yaml \
				  -f ./kubernetes/feelguuds-platform/feelguuds-platform-deployment.yaml \
				  -f ./kubernetes/feelguuds-platform/feelguuds-platform-env-configmap.yaml \
				  -f ./kubernetes/feelguuds-platform/feelguuds-platform-service.yaml
	kubectl apply -f ./kubernetes/merchant-component-db/merchant-component-db-deployment.yaml \
				  -f ./kubernetes/merchant-component-db/merchant-component-db-service.yaml \
				  -f ./kubernetes/merchant-component-db/merchant-component-db-persistentvolumeclaim.yaml
	kubectl apply -f ./kubernetes/shopper-component-db/shopper-component-db-deployment.yaml \
				  -f ./kubernetes/shopper-component-db/shopper-component-db-service.yaml \
				  -f ./kubernetes/shopper-component-db/shopper-component-db-persistentvolumeclaim.yaml
	minikube dashboard

# kubectl convert -f ./my-deployment.yaml --output-version apps/v1
