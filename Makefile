# All builds use the `ebitensinglethread` tag to shave off some CPU usage.
# Details here: https://github.com/hajimehoshi/ebiten/issues/1367

build:
	go build -tags ebitensinglethread -o dist/

linux:
	GOOS=linux GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-linux

windows:
	GOOS=windows GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-windows.exe -ldflags -H=windowsgui

macos:
	GOOS=darwin GOARCH=amd64 go build -tags ebitensinglethread -o dist/shark-macos

clean:
	rm -f dist/*

# https://ebiten.org/documents/install.html#Debian_/_Ubuntu
deps-debian:
	sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
