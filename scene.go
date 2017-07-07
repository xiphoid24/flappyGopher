package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var quitError = errors.New("quit")

type scene struct {
	bg       *sdl.Texture
	bird     *bird
	pipes    *pipes
	score    *score
	music    *mix.Music
	r        *sdl.Renderer
	quitMenu *sdl.MessageBoxData
}

func newScene(r *sdl.Renderer) (*scene, error) {

	var bg *sdl.Texture
	var err error

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, 100); err != nil {
		return nil, err
	}

	music, err := mix.LoadMUS("resources/audio/music.mp3")
	if err != nil {
		return nil, fmt.Errorf("could not load music: %v", err)
	}

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

	sc, err := newScore(r)
	if err != nil {
		return nil, err
	}

	return &scene{
		bg:       bg,
		bird:     b,
		pipes:    ps,
		score:    sc,
		music:    music,
		r:        r,
		quitMenu: newQuitMenu(),
	}, nil
}

func (s *scene) run(events chan sdl.Event) <-chan error {
	s.music.Play(-1)
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
				if s.bird.isDead() {
					continue
				}
				fmt.Println("tick")
				s.update()

				if s.bird.isDead() {
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>dead")
					err := s.gameOver()
					if err == quitError {
						return
					}
					if err != nil {
						errc <- err
					}
				} else {
					if err := s.paint(); err != nil {
						errc <- err
					}
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
		/*if s.bird.isDead() {
			s.restart()
		} else {
		}*/

	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.CommonEvent:
	default:
		log.Printf("unknown event %T\n", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update(s.score)
	s.touch()
}

func (s *scene) touch() {
	s.pipes.mu.RLock()
	defer s.pipes.mu.RUnlock()
	for _, p := range s.pipes.pipes {
		s.bird.touch(p)
	}
}

func (s *scene) gameOver() error {
	mix.PauseMusic()
	mix.RewindMusic()
	// drawTitle(s.r, "Game Over")
	s.score.mu.RLock()
	if err := ioutil.WriteFile("high.txt", []byte(strconv.Itoa(s.score.high)), 0666); err != nil {
		s.score.mu.RUnlock()
		return err
	}
	s.score.mu.RUnlock()

	err, button := sdl.ShowMessageBox(s.quitMenu)
	if err != nil {
		return err
	}
	fmt.Printf(">>>>>>>>button %d\n", button)
	if button == 0 {
		return quitError
	}

	// time.Sleep(2 * time.Second)
	s.restart()
	return nil
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
	s.score.restart()
	mix.ResumeMusic()
}

func (s *scene) paint() error {
	var err error

	sdl.Do(func() {
		s.r.Clear()
		err = s.r.Copy(s.bg, nil, nil)
	})
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err = s.bird.paint(); err != nil {
		return err
	}

	if err = s.pipes.paint(); err != nil {
		return err
	}

	if err = s.score.paintHigh(); err != nil {
		return err
	}

	if err = s.score.paintCurrent(); err != nil {
		return err
	}

	sdl.Do(func() {
		s.r.Present()
	})

	return nil
}

func (s *scene) destroy() {
	s.music.Free()

	sdl.Do(func() {
		s.bg.Destroy()
	})

	s.bird.destroy()
	s.pipes.destroy()
}

func newQuitMenu() *sdl.MessageBoxData {
	buttons := []sdl.MessageBoxButtonData{
		{
			Flags:    sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT,
			ButtonId: 0,
			Text:     "No",
		},
		{
			Flags:    sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT,
			ButtonId: 1,
			Text:     "Yes",
		},
	}
	color := &sdl.MessageBoxColorScheme{
		Colors: [5]sdl.MessageBoxColor{
			/* .colors (.r, .g, .b) */
			/* [SDL_MESSAGEBOX_COLOR_BACKGROUND] */
			{R: 255, G: 0, B: 0},
			/* [SDL_MESSAGEBOX_COLOR_TEXT] */
			{R: 0, G: 255, B: 0},
			/* [SDL_MESSAGEBOX_COLOR_BUTTON_BORDER] */
			{R: 255, G: 255, B: 0},
			/* [SDL_MESSAGEBOX_COLOR_BUTTON_BACKGROUND] */
			{R: 0, G: 0, B: 255},
			/* [SDL_MESSAGEBOX_COLOR_BUTTON_SELECTED] */
			{R: 255, G: 0, B: 255},
		},
	}

	return &sdl.MessageBoxData{
		Flags:       sdl.MESSAGEBOX_INFORMATION,
		Window:      nil,
		Title:       "Game Over",
		Message:     "Game Over! Would you like to try again?",
		NumButtons:  2,
		Buttons:     buttons,
		ColorScheme: color,
	}
}
