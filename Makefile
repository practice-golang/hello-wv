build:
	go build -ldflags "-w -s -H windowsgui" -o bin/

clean:
	rm -rf ./bin