package main

import (
	"bytes"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//SDL_True and SDL_False
	sdlFalse = 0
	sdlTrue  = 1
)

//LWindow is a wrapper for SDL_Window
type LWindow struct {
	//Window data
	mWindow   *sdl.Window
	mRenderer *sdl.Renderer
	mWindowID uint32

	//Window dimensions
	mWidth  int32
	mHeight int32

	//Window focus
	mMouseFocus    bool
	mKeyboardFocus bool
	mFullScreen    bool
	mMinimized     bool
	mShown         bool
}

//Init Creates window
func (w *LWindow) Init() error {
	//Local error declaration
	var err error

	//Create Window
	w.mWindow, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWitdh, screenHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}

	w.mMouseFocus = true
	w.mKeyboardFocus = true
	w.mWidth = screenWitdh
	w.mHeight = screenHeight

	//Create renderer for window
	w.mRenderer, err = sdl.CreateRenderer(w.mWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		dErr := w.mWindow.Destroy()
		if dErr != nil {
			return fmt.Errorf("could not destroy window after failing renderer creation: %v", err)
		}
		w.mWindow = nil

		return fmt.Errorf("could not create renderer for window: %v", err)
	}

	//Initialize renderer color
	w.mRenderer.SetDrawColor(255, 255, 255, 255)

	//Grab window identifier
	w.mWindowID, err = w.mWindow.GetID()
	if err != nil {
		return fmt.Errorf("could not grab window ID: %v", err)
	}

	//Flag as opened
	w.mShown = true

	return nil
}

//CreateRenderer creates renderer from internal window
func (w *LWindow) CreateRenderer() (*sdl.Renderer, error) {
	return sdl.CreateRenderer(w.mWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
}

//HandleEvent handles window events
func (w *LWindow) HandleEvent(e sdl.Event) error {
	//Window event occured
	if e.GetType() == sdl.WINDOWEVENT && (e.(*sdl.WindowEvent)).WindowID == w.mWindowID {
		//Caption update flag
		var updateCaption bool
		var wEvent = e.(*sdl.WindowEvent)

		switch wEvent.Event {
		//Window appeared
		case sdl.WINDOWEVENT_SHOWN:
			w.mShown = true
			break

		//Window disappeared
		case sdl.WINDOWEVENT_HIDDEN:
			w.mShown = false
			break

		//Get new dimensions and repaint on widow size change
		case sdl.WINDOWEVENT_SIZE_CHANGED:
			w.mWidth = wEvent.Data1
			w.mHeight = wEvent.Data2
			w.mRenderer.Present()
			break

		//Repaint on exposure
		case sdl.WINDOWEVENT_EXPOSED:
			w.mRenderer.Present()
			break

		//Mouse entered window
		case sdl.WINDOWEVENT_ENTER:
			w.mMouseFocus = true
			updateCaption = true
			break

		//Mouse left window
		case sdl.WINDOWEVENT_LEAVE:
			w.mMouseFocus = false
			updateCaption = true
			break

		//Window has keyboard focus
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			w.mKeyboardFocus = true
			updateCaption = true
			break

		//Window lost keyboard focus
		case sdl.WINDOWEVENT_FOCUS_LOST:
			w.mKeyboardFocus = false
			updateCaption = true
			break

		//Window minimized
		case sdl.WINDOWEVENT_MINIMIZED:
			w.mMinimized = true
			break

		//Window maximized
		case sdl.WINDOWEVENT_MAXIMIZED:
			w.mMinimized = false
			break

		//Window restored
		case sdl.WINDOWEVENT_RESTORED:
			w.mMinimized = false
			break

		case sdl.WINDOWEVENT_CLOSE:
			w.mWindow.Hide()
			break
		}

		//Update window caption with new data
		if updateCaption {

			mouseFocus := "Off"
			if w.mMouseFocus {
				mouseFocus = "On"
			}

			keyboardFocus := "Off"
			if w.mKeyboardFocus {
				keyboardFocus = "On"
			}

			var caption = bytes.NewBufferString("")
			fmt.Fprint(caption, "SDL Tutorial - MouseFocus:", mouseFocus, keyboardFocus)

			w.mWindow.SetTitle(caption.String())
		}
	} else if e.GetType() == sdl.KEYDOWN && (e.(*sdl.KeyboardEvent)).Keysym.Sym == sdl.K_RETURN { //Enter exit fullscreen on return key
		if w.mFullScreen {
			if err := w.mWindow.SetFullscreen(sdlFalse); err != nil {
				return fmt.Errorf("could not unset window from fullscreen: %v", err)
			}

			w.mFullScreen = false
		} else {
			if err := w.mWindow.SetFullscreen(sdlTrue); err != nil {
				return fmt.Errorf("could not set window to fullscreen: %v", err)
			}

			w.mFullScreen = true
			w.mMinimized = false
		}
	}

	return nil
}

//Focus focuses on window
func (w *LWindow) Focus() {
	//Restore window if needed
	if !w.mShown {
		w.mWindow.Show()
	}

	//Move window forward
	w.mWindow.Raise()
}

//Render shows window contents
func (w *LWindow) Render() error {
	if !w.mMinimized {
		//Clear screen
		if err := w.mRenderer.SetDrawColor(255, 255, 255, 255); err != nil {
			return fmt.Errorf("could not set draw color for renderer: %v", err)
		}
		if err := w.mRenderer.Clear(); err != nil {
			return fmt.Errorf("could not clear renderer")
		}

		//Update screen
		w.mRenderer.Present()
	}
	return nil
}

//MWidth returns window's width
func (w *LWindow) MWidth() int32 {
	return w.mWidth
}

//MHeight returns window's height
func (w *LWindow) MHeight() int32 {
	return w.mHeight
}

//HasMouseFocus returns mouse focus state
func (w *LWindow) HasMouseFocus() bool {
	return w.mMouseFocus
}

//HasKeyboardFocus returns keyboard focus state
func (w *LWindow) HasKeyboardFocus() bool {
	return w.mKeyboardFocus
}

//IsMinimized returns window minimization state
func (w *LWindow) IsMinimized() bool {
	return w.mMinimized
}

//IsShown return window show state
func (w *LWindow) IsShown() bool {
	return w.mShown
}

//Free dellocates internal
func (w *LWindow) Free() error {
	if w.mWindow != nil {
		if err := w.mWindow.Destroy(); err != nil {
			return fmt.Errorf("could not destroy window: %v", err)
		}
	}

	w.mMouseFocus = false
	w.mKeyboardFocus = false
	w.mWidth = 0
	w.mHeight = 0

	return nil
}
