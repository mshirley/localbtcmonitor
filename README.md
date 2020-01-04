# localbtcmonitor
this tool generates png images of current localbitcoins.com advertisements converted to USD.

## install golang
* `sudo apt install golang`
## install/run go-swp
https://github.com/me-io/go-swap

quickest setup for go-swp is to use docker then modify monitor.go to point to the docker container.
## clone this repo
* `git clone https://github.com/mshirley/localbtcmonitor.git`
## build
* `cd localbtcmonitor`
* `go get .`
* `go build`
## run
* `./localbtcmonitor`
