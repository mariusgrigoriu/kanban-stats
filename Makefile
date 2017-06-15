app_name=kanban-stats
account=mariusgrigoriu
release=latest
target_os=linux
target_arch=amd64

container: Dockerfile build
	docker build -t $(account)/$(app_name):$(release) .
	touch build/container

.PHONY: release
release: container
	docker push $(account)/$(app_name):$(release)

.PHONY: clean
clean:
	rm -rf build
	rm -rf dist

build: *.go
	mkdir -p build
	GOOS=$(target_os) GOARCH=$(target_arch) go build -o build/$(app_name)

dist:
	mkdir -p dist
