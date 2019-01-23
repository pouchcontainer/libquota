build: build-setquota

build-setquota:
	mkdir -p bin
	GOOS=linux go build -o bin/setquota cmd/setquota/main.go

clean:
	rm -rf bin/