package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	var bg *sdl.Texture
	var err error

	sdl.Do(func() {
		bg, err = img.LoadTexture(r, "resources/imgs/background.png")
	})
	if err != nil {
		return nil, err
	}

	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		var bird *sdl.Texture
		path := fmt.Sprintf("resources/imgs/bird_frame_%d.png", i)
		sdl.Do(func() {
			bird, err = img.LoadTexture(r, path)
		})
		if err != nil {
			return nil, err
		}
		birds = append(birds, bird)
	}

	return &scene{bg: bg, birds: birds}, nil
}

func (s *scene) run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	var err error
	s.time++
	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	i := s.time / 10 % len(s.birds)

	sdl.Do(func() {
		r.Clear()
		err = r.Copy(s.bg, nil, nil)
	})
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	sdl.Do(func() {
		err = r.Copy(s.birds[i], nil, rect)
	})
	if err != nil {
		return err
	}
	sdl.Do(func() {
		r.Present()
	})
	return nil
}

func (s *scene) destroy() {
	sdl.Do(func() {
		s.bg.Destroy()
	})
}
