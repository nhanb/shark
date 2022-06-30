# All builds use the `ebitensinglethread` tag to shave off some CPU usage.
# Details here: https://github.com/hajimehoshi/ebiten/issues/1367

build:
	go build -tags ebitensinglethread -o dist/shark

linux:
	GOOS=linux GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-linux

win:
	GOOS=windows GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-win.exe

mac:
	GOOS=darwin GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-mac

clean:
	rm -f dist/*
