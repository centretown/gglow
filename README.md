# gglow
A graphical light effect editor

## Development Environment Ubuntu
ARM64 (pi4), ARM32(pi2), AMD64 

`sudo apt install build-essential git pkg-config libx11-dev xserver-xorg-dev xorg-dev`

[Install Go](https://go.dev/doc/install)

`mkdir src`

`git clone https://github.com/centretown/gglow.git`

`cd src`

`go mod tidy`

`go install .`

and wait...

`gglow -p glow.db`

or 

`go run . -p glow.db`

