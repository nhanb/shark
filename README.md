# What

![](https://user-images.githubusercontent.com/1446315/176945811-18387e27-df3e-4303-ad17-e1637e547f04.jpg)

This is a PoC "desktop pet" à la [shimeji][1] using [ebitengine][2] that runs
on Windows, Linux, and macOS. It currently has only 3 animations: idle
(default), left-click dragging, and right-click.

Go [here][7] for a demo video.

Fair warning: I'm a Go noob who mostly has no idea what he's doing.
Read the source code at your own peril.

# Download

Official tagged releases are available on [GitHub Releases][5].

Newer builds straight from master branch are also available as artifacts on
GitHub and sourcehut:

[![github status](https://github.com/nhanb/shark/actions/workflows/main.yml/badge.svg)][gh]
[![builds.sr.ht status](https://builds.sr.ht/~nhanb/shark/commits/master.svg)][srht]

Sourcehut doesn't have macOS builds though.

# Compile from source

- Follow [ebitengine's install guide][6]
- Run: `go build -tags ebitensinglethread -o dist/`
- Your compiled binary should now be in `./dist/`

Apparently it should compile on FreeBSD too but I haven't tried that.
Do let me know if it works on your FreeBSD desktop!

# Artist

The sprites were graciously provided by Mee Way:

![](https://user-images.githubusercontent.com/1446315/176449384-7a06250d-7dfe-4371-b998-707ddbda66b1.jpg)  
**[Behance][3] | [Facebook][4]**

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
[7]: https://user-images.githubusercontent.com/1446315/176439983-091dec3d-bc36-4ae3-8b78-2a2a7f11e90d.mp4

[srht]: https://builds.sr.ht/~nhanb/shark/commits/master
[gh]: https://github.com/nhanb/shark/actions/workflows/main.yml
