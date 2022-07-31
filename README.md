# What

![](https://user-images.githubusercontent.com/1446315/177188223-ad9759c9-4ef4-44e0-84d8-03cfd46129b8.png)

This is a PoC "desktop pet" à la [shimeji][1] using [ebitengine][2] that runs
on Windows, Linux, and macOS. It currently has these animations:

- `Idle`
- `Dragging`
- `Right-click`
- Randomly `Walk` horizontally
- After some time has passed (1 hour by default), a `Hungry` animation will be
  activated, during which dragging is disabled.
- When `Hungry`, right-click to start `Feeding` animation and reset to the
  normal idle state.

Here's a [demo video](https://user-images.githubusercontent.com/1446315/178103169-006c2bc0-ebb9-4014-aba5-8a1fbc3d0733.mp4).

Fair warning: I'm a Go noob who mostly has no idea what he's doing.
Read the source code at your own peril.

# Download

[Download latest version][5] from GitHub Releases.

Newer builds straight from master branch are also available as artifacts on
GitHub and sourcehut:

[![github status](https://github.com/nhanb/shark/actions/workflows/main.yml/badge.svg)][gh]
[![builds.sr.ht status](https://builds.sr.ht/~nhanb/shark/commits/master.svg)][srht]

GitHub requires logging in to download artifacts, while Sourcehut doesn't have
macOS builds. Such is life.

# Usage

## Windows & Linux

Simply unzip then run the `shark-windows.exe` or `shark-linux` executable.

## macOS

Since I'm not participating in Apple's $99/yr [protection racket][pr], macOS
users will need to jump through some hoops to run this program:

- Double click on the downloaded zip file to get the `Shark` app bundle.
  (skip this step if you downloaded using Safari, which automatically unzips)
- Drag the `Shark` app bundle into your `Applications` folder.
- Right-click on `Shark` -> `Open`. You'll see a warning pop-up saying this
  application was created by an unverified developer (yours truly). Note: you
  must **right-click instead of double-clicking**, because double-clicking will
  open a different pop-up which hides the option to open the app.

![](https://user-images.githubusercontent.com/1446315/178136989-247b5d70-ee37-47a6-95b2-a726103b95f3.png)

- Click "Open" anyway.
- From now on you can launch the Shark application just like any other app,
  either from Spotlight or from the Applications folder.

There's also Apple's official guide [here][apple-guide].

In the future I might pay the $99 if I end up writing more macOS apps and this
becomes enough of a nuisance. Maybe.

## Options

If run from a terminal, use the `-h` argument to see available options.
Windows users can [create a shortcut][7] to save their desired options.

Here are the currently supported options:

```
  -hungry int
        The number of seconds it takes for Gura to go hungry (default 3600)
  -size int
        Size multiplier: make Gura as big as you want (default 1)
  -stop int
        chance to stop walking, in % (default 40)
  -walk int
        chance to start walking, in % (default 5)
  -x int
        X position on screen (default 9999)
  -y int
        Y position on screen (default 9999)
```

# Compile from source

- Clone this repo
- Follow [ebitengine's install guide][6]
- Run: `go build -tags ebitensinglethread -o dist/`
- Your compiled binary should now be in `./dist/`

Alternatively, if you already have Go, you can run `go run
go.imnhan.com/shark@latest` to compile and run the latest version without
manually cloning the repo. You still need to install ebiten's dependencies
first though.

Apparently it should compile on FreeBSD too but I haven't tried that.
Do let me know if it works on your FreeBSD desktop!

# Artist

The sprites were graciously provided by Mee Way:

![](https://user-images.githubusercontent.com/1446315/176449384-7a06250d-7dfe-4371-b998-707ddbda66b1.jpg)  
**[Twitter][8] | [Behance][3] | [Facebook][4]**

# License

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License version 3.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <https://www.gnu.org/licenses/>.

[1]: https://shimejis.xyz/
[2]: https://ebiten.org/
[3]: https://www.behance.net/meeway/projects
[4]: https://www.facebook.com/meexway
[5]: https://github.com/nhanb/shark/releases/latest
[6]: https://ebiten.org/documents/install.html
[7]: https://superuser.com/questions/29569/how-to-add-command-line-options-to-shortcut
[8]: https://twitter.com/mee_way

[srht]: https://builds.sr.ht/~nhanb/shark/commits/master
[gh]: https://github.com/nhanb/shark/actions/workflows/main.yml
[pr]: https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution
[apple-guide]: https://support.apple.com/en-vn/guide/mac-help/mh40616/mac
