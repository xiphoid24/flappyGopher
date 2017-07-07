package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	sdl.Main(func() {
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(2)
		}
	})
}

func run() error {
	var w *sdl.Window
	var r *sdl.Renderer
	var err error

	sdl.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			sdl.Quit()
		})
	}()

	sdl.Do(func() {
		err = ttf.Init()
	})
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			ttf.Quit()
		})
	}()

	sdl.Do(func() {
		w, r, err = sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	})
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			w.Destroy()
		})
	}()

	defer func() {
		sdl.Do(func() {
			r.Destroy()
		})
	}()

	if err = drawTitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case err := <-s.run(ctx, r):
		return err
	case <-time.After(5 * time.Second):
		return nil
	}
}

func drawTitle(r *sdl.Renderer) error {
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

	sdl.Do(func() {
		s, err = f.RenderUTF8_Solid("Flappy Gopher", sdl.Color{
			R: 255,
			G: 100,
			B: 0,
			A: 255,
		})
	})
	if err != nil {
		return fmt.Errorf("could not render solid: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			s.Free()
		})
	}()

	sdl.Do(func() {
		t, err = r.CreateTextureFromSurface(s)
	})
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			t.Destroy()
		})
	}()

	sdl.Do(func() {
		r.Clear()
		err = r.Copy(t, nil, nil)
	})
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	sdl.Do(func() {
		r.Present()
	})

	return nil
}
