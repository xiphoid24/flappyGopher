package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var windowW, windowH int

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
		w, r, err = sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN_DESKTOP)
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

	windowW, windowH = w.GetSize()

	if err = drawTitle(r, "Flappy Gopher"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(events)

	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
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
		s, err = f.RenderUTF8_Solid(text, sdl.Color{
			R: 255,
			G: 100,
			B: 0,
			A: 255,
		})
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
