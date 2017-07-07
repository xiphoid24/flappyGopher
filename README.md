# Flappy Gopher

Flappy Gopher is a clone of the famous Flappy Bird game developed in Go with
bindings for SDL2.

This is my version of the [just for func](https://www.youtube.com/watch?v=aYkxFbd6luY&list=PL64wiCrrxh4Jisi7OcCJIUpguV_f5jGnZ&index=9) flappy gopher. the original code can be found [here](https://github.com/campoy/flappy-gopher).
I was having issues with go-sdl2 and goroutines on Ubuntu so I followed the example [here](https://github.com/veandco/go-sdl2/blob/master/examples/render_goroutines/render_goroutines.go). I changed my program to reflect that style and got it to work.
Hope it helps.

## Installation

You need to install first SDL2 and the SDL2 bindings for Go. To do so follow the instructions [here](https://github.com/veandco/go-sdl2).
It is quite easy to install on basically any platform.

You will also need to install [pkg-config](https://en.wikipedia.org/wiki/Pkg-config).

After that you should be able to simply run:

    go get github.com/gregpechiro/flappyGopher

And run the binary generated in `$GOPATH/bin`.

## Images, fonts, and licenses

All the images used in this game are CC0 and obtained from [Clipart](https://openclipart.org/tags/flapping).
You can find atlernative birds in there, so you can mod the game!

The fonts are copied from https://github.com/deano2390/OpenFlappyBird.
