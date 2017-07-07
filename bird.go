package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.25
	jumpSpeed = 5
)

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture
	r        *sdl.Renderer
	chunk    *mix.Chunk

	x, y  int32
	w, h  int32
	speed float64
	dead  bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	var err error

	chunk, err := mix.LoadWAV("resources/audio/boing.wav")
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 4; i++ {
		var texture *sdl.Texture
		path := fmt.Sprintf("resources/imgs/bird_frame_%d.png", i)
		sdl.Do(func() {
			texture, err = img.LoadTexture(r, path)
		})
		if err != nil {
			return nil, err
		}
		textures = append(textures, texture)
	}

	return &bird{
		textures: textures,
		chunk:    chunk,
		x:        10,
		y:        int32(windowH / 2),
		w:        int32(float64(windowH) * .07),
		h:        int32(float64(windowH) * .07),
		r:        r,
	}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y -= int32(b.speed)

	if b.y < 0 || b.y > 600+b.h {
		b.dead = true
	}

	/*if b.y > 0 {
		b.y -= int32(b.speed)
		b.speed += gravity
	}*/

	b.speed += gravity

}

func (b *bird) paint() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var err error

	rect := &sdl.Rect{X: b.x, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}
	i := b.time / 10 % len(b.textures)

	sdl.Do(func() {
		err = b.r.Copy(b.textures[i], nil, rect)
	})
	if err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}

	return nil
}

func (b *bird) restart() {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.chunk.Free()
	for _, t := range b.textures {
		sdl.Do(func() {
			t.Destroy()
		})
	}
}

func (b *bird) isDead() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dead
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if mix.Playing(-1) < 1 && !b.dead {
		b.chunk.Play(-1, 1)
	}

	b.speed = -jumpSpeed
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // too far left
		return
	}
	if p.x+p.w < b.x { // too far right
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 { // pipe is too low
		return
	}
	if p.inverted && (600-p.h) > b.y+b.h/2 { // inverted pipe is too high
		return
	}
	b.dead = true
}
