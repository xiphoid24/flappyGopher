package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
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

	x, y  int32
	w, h  int32
	speed float64
	dead  bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	var err error

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

	return &bird{textures: textures, x: 10, y: 300, w: 50, h: 43}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity

}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var err error

	rect := &sdl.Rect{X: b.x, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}
	i := b.time / 10 % len(b.textures)

	sdl.Do(func() {
		err = r.Copy(b.textures[i], nil, rect)
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

	b.speed = -jumpSpeed
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()
	p.mu.RLock()
	p.mu.RUnlock()

	if p.x > b.x+b.w { // too far left
		return
	}
	if p.x+p.w < b.x { // too far right
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 { // pipe is too low
		return
	}
	if p.inverted && (600-p.h) > b.y-b.h/2 { // inverted pipe is too high
		return
	}
	b.dead = true
}
