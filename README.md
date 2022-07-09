# What

![](https://user-images.githubusercontent.com/1446315/177188223-ad9759c9-4ef4-44e0-84d8-03cfd46129b8.png)

This is a PoC "desktop pet" à la [shimeji][1] using [ebitengine][2] that runs
on Windows, Linux, and macOS. It currently has only 5 animations:

- `Idle`
- `Dragging`
- `Right-click`
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

Simply run the provided binary for your OS. Mac & Linux users may need to first
make the file executable with `chmod +x <file-name>`.

If run from a terminal, use the `-h` argument to see available options.
Windows users can [create a shortcut][7] to save their desired options.

Here are the currently supported options:

```
-hungry int
      The number of seconds it takes for Gura to go hungry (default 3600)
-size int
      Size multiplier: make Gura as big as you want (default 1)
-x int
      X position on screen (default 9999)
-y int
      Y position on screen (default 9999)
```

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
[7]: https://superuser.com/questions/29569/how-to-add-command-line-options-to-shortcut

[srht]: https://builds.sr.ht/~nhanb/shark/commits/master
[gh]: https://github.com/nhanb/shark/actions/workflows/main.yml
