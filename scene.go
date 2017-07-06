package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	var bg *sdl.Texture
	var err error

	sdl.Do(func() {
		bg, err = img.LoadTexture(r, "resources/imgs/background.png")
	})
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	b, err := newBird(r)
	if err != nil {
		return nil, err
	}

	ps, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: b, pipes: ps}, nil
}

func (s *scene) run(events chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(1 * time.Second)
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.CommonEvent:
	default:
		log.Printf("unknown event %T\n", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.bird.touch(s.pipe)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	var err error

	sdl.Do(func() {
		r.Clear()
		err = r.Copy(s.bg, nil, nil)
	})
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err = s.bird.paint(r); err != nil {
		return err
	}

	if err = s.pipes.paint(r); err != nil {
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

	s.bird.destroy()
	s.pipes.destroy()
}
