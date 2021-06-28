DIR?="."
PARALLEL?=false

.PHONY: run
## run: запустить программу. Можно указать значение переменной DIR - каталог, в котором будет производиться поиск. По умолчанию '.'. И значение переменной PARALLEL - запуск в параллельном режиме. По умолчанию 'false'.
run:
	go run cmd/finddup/main.go --dir=$(DIR) --parallel=$(PARALLEL)

.PHONY: build
## build: создаёт исполяемый файл.
build:
	CGO_ENABLED=0 && go build -o bin/finddup -ldflags="-s -w" cmd/finddup/main.go

.PHONY: clean
## clean: удалить содеримое папки bin.
clean:
	rm -f bin/*

.PHONY: help
help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
