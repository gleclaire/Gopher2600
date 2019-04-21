package sdl

import (
	"gopher2600/gui"
	"gopher2600/television"

	"github.com/veandco/go-sdl2/sdl"
)

// GUI is the SDL implementation of a the gui/television
type GUI struct {
	television.HeadlessTV

	// much of the sdl magic happens in the screen object
	scr *screen

	// regulates how often the screen is updated
	fpsLimiter *fpsLimiter

	// connects SDL guiLoop with the parent process
	eventChannel chan gui.Event

	// whether the emulation is currently paused. if paused is true then
	// as much of the current frame is displayed as possible; the previous
	// frame will take up the remainder of the screen.
	paused bool

	// ther's a small bug significant performance boost if we disable certain
	// code paths with this allowDebugging flag
	allowDebugging bool
}

// NewGUI initiliases a new instance of an SDL based display for the VCS
func NewGUI(tvType string, scale float32) (*GUI, error) {
	var err error

	tv := new(GUI)

	tv.fpsLimiter, err = newFPSLimiter(50)
	if err != nil {
		return nil, err
	}

	err = television.InitHeadlessTV(&tv.HeadlessTV, tvType)
	if err != nil {
		return nil, err
	}

	// set up sdl
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, err
	}

	// initialise the screens we'll be using
	tv.scr, err = newScreen(tv)

	// set window size and scaling
	err = tv.scr.setScaling(scale)
	if err != nil {
		return nil, err
	}

	// register headlesstv callbacks
	// --leave SignalNewScanline() hook at its default
	tv.HookNewFrame = func() error {
		defer tv.scr.clearPixels()
		err := tv.scr.stb.checkStableFrame()
		if err != nil {
			return err
		}
		return tv.update()
	}
	tv.HookSetPixel = tv.scr.setRegPixel

	// update tv (with a black image)
	err = tv.update()
	if err != nil {
		return nil, err
	}

	// gui events are serviced by a separate loop
	go tv.guiLoop()

	// note that we've elected not to show the window on startup
	// window is instead opened on a ReqSetVisibility request

	return tv, nil
}

// update the gui so that it reflects changes to buffered data in the tv struct
func (tv *GUI) update() error {
	tv.fpsLimiter.wait()

	// abbrogate most of the updating to the screen instance
	err := tv.scr.update(tv.paused)
	if err != nil {
		return err
	}

	tv.scr.renderer.Present()

	return nil
}

func (tv *GUI) setDebugging(allow bool) {
	tv.allowDebugging = allow
	if allow {
		tv.HookSetAltPixel = tv.scr.setAltPixel
	} else {
		tv.HookSetAltPixel = nil
	}
}