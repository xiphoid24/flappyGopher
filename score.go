package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type score struct {
	mu sync.RWMutex

	r          *sdl.Renderer
	high       int
	current    int
	multiplyer int
}

func newScore(r *sdl.Renderer) (*score, error) {
	b, err := ioutil.ReadFile("high.txt")
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not read high score %v", err)
	}
	high := 0
	if err == nil {
		high, err = strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("could not convert high score %v", err)
		}
	}

	return &score{
		r:          r,
		high:       high,
		current:    0,
		multiplyer: 10,
	}, nil
}

func (sc *score) increase() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.current += sc.multiplyer
	sc.multiplyer += 10
	if sc.current > sc.high {
		sc.high = sc.current
	}
}

func (sc *score) restart() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.current = 0
	sc.multiplyer = 10
}

func (sc *score) paintHigh() error {
	var f *ttf.Font
	var s *sdl.Surface
	var t *sdl.Texture
	var err error

	sdl.Do(func() {
		f, err = ttf.OpenFont("resources/fonts/Flappy.ttf", 20)
	})
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			f.Close()
		})
	}()

	c := sdl.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	sdl.Do(func() {
		s, err = f.RenderUTF8_Solid(fmt.Sprintf("High Score: %d", sc.high), c)
	})
	if err != nil {
		return fmt.Errorf("could not render surface: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			s.Free()
		})
	}()

	sdl.Do(func() {
		t, err = sc.r.CreateTextureFromSurface(s)
	})
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			t.Destroy()
		})
	}()

	rect := &sdl.Rect{X: 10, Y: 10, W: 200, H: 30}
	sdl.Do(func() {
		err = sc.r.Copy(t, nil, rect)
	})
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

func (sc *score) paintCurrent() error {
	var f *ttf.Font
	var s *sdl.Surface
	var t *sdl.Texture
	var err error

	sdl.Do(func() {
		f, err = ttf.OpenFont("resources/fonts/Flappy.ttf", 20)
	})
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			f.Close()
		})
	}()

	c := sdl.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	sdl.Do(func() {
		s, err = f.RenderUTF8_Solid(fmt.Sprintf("Score: %d", sc.current), c)
	})
	if err != nil {
		return fmt.Errorf("could not render surface: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			s.Free()
		})
	}()

	sdl.Do(func() {
		t, err = sc.r.CreateTextureFromSurface(s)
	})
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			t.Destroy()
		})
	}()

	rect := &sdl.Rect{X: 580, Y: 10, W: 200, H: 30}
	sdl.Do(func() {
		err = sc.r.Copy(t, nil, rect)
	})
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}
