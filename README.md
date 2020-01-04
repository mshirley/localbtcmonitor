# localbtcmonitor
this tool generates png images of current localbitcoins.com advertisements converted to USD.

the default configuration is to use the heroku app associated with https://github.com/me-io/go-swap for currency conversion.  if this server isn't available or you want to use your own copy of go-swap just change the value in monitor.go and recompile 

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
