package sdldebug

import (
	"io"

	"github.com/veandco/go-sdl2/sdl"
)

// the textures type keeps track of the textures used to render to the screen.
// textures are flipped every frame and the previous frame used as a ghosted
// image for the paused instance of sdeldebug.
//
// the overlay type looks after its own texture.
type textures struct {
	renderer *sdl.Renderer
	a        *sdl.Texture
	b        *sdl.Texture
	front    *sdl.Texture
	back     *sdl.Texture
}

// create a new instance of the textures type. called everytime the screen
// dimensions change.
func newTextures(renderer *sdl.Renderer, w, h int) (*textures, error) {
	txt := &textures{
		renderer: renderer,
	}

	var err error

	// texture is applied to the renderer to show the image. we copy the pixels
	// to it every NewFrame()
	txt.a, err = renderer.CreateTexture(uint32(sdl.PIXELFORMAT_ABGR8888),
		int(sdl.TEXTUREACCESS_STREAMING),
		int32(w), int32(h))
	if err != nil {
		return nil, err
	}

	err = txt.a.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return nil, err
	}

	txt.b, err = renderer.CreateTexture(uint32(sdl.PIXELFORMAT_ABGR8888),
		int(sdl.TEXTUREACCESS_STREAMING),
		int32(w), int32(h))
	if err != nil {
		return nil, err
	}

	err = txt.b.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return nil, err
	}

	txt.front = txt.a
	txt.back = txt.b

	return txt, err
}

// destroy texture resources
func (txt *textures) destroy(output io.Writer) {
	err := txt.a.Destroy()
	if err != nil {
		output.Write([]byte(err.Error()))
	}

	err = txt.b.Destroy()
	if err != nil {
		output.Write([]byte(err.Error()))
	}
}

// swap textures and set/reset alpha modifiers depending on if the texture is
// now the "front" or "back" texture
func (txt *textures) flip() {
	if txt.front == txt.a {
		txt.front = txt.b
		txt.back = txt.a
	} else {
		txt.front = txt.a
		txt.back = txt.b
	}

	txt.front.SetAlphaMod(255)
	txt.back.SetAlphaMod(150)
}

// update texture with pixels and render
func (txt *textures) render(cpyRect *sdl.Rect, pixels []byte, pitch int) error {
	// draw "back" texture to screen
	err := txt.renderer.Copy(txt.back, cpyRect, nil)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// update "front" texture
	err = txt.front.Update(nil, pixels, pitch)
	if err != nil {
		return err
	}

	// draw "front" texture to screen
	err = txt.renderer.Copy(txt.front, cpyRect, nil)
	if err != nil {
		return err
	}

	return nil
}