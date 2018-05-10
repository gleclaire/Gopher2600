package video

import (
	"gopher2600/hardware/tia/colorclock"
)

// Video contains all the components of the video sub-system of the VCS TIA chip
type Video struct {
	colorClock *colorclock.ColorClock
	hblank     *bool

	// sprite objects
	player0  *sprite
	player1  *sprite
	missile0 *sprite
	missile1 *sprite
	ball     *sprite

	// colors
	colup0 Color
	colup1 Color
	colupf Color
	colubk Color

	// TODO: player sprite data
	// TODO: playfield

	// playfield control
	// -- including ball size
	ctrlpfReflection bool
	ctrlpfPriority   bool
	ctrlpfScoremode  bool
	ctrlpfBallSize   int

	// TODO: player/missile number & spacing
	// TODO: trigger lists
	// TODO: missile/ball size

	// TODO: player reflection

	// missile/ball enabling
	enam0      bool
	enam1      bool
	enabl      bool
	enam0Delay *delayCounter
	enam1Delay *delayCounter
	enablDelay *delayCounter
	enam0Prev  bool
	enam1Prev  bool
	enablPrev  bool

	// vertical delay
	vdelp0 bool
	vdelp1 bool
	vdelbl bool

	// TODO: horizontal movement
}

// New is the preferred method of initialisation for the Video structure
func New(colorClock *colorclock.ColorClock, hblank *bool) *Video {
	vd := new(Video)
	if vd == nil {
		return nil
	}

	vd.colorClock = colorClock
	vd.hblank = hblank

	// TODO: remove this tempory setting of color
	vd.colupf = 0xff

	// sprite objects
	vd.player0 = newSprite("player0")
	if vd.player0 == nil {
		return nil
	}
	vd.player1 = newSprite("player1")
	if vd.player1 == nil {
		return nil
	}
	vd.missile0 = newSprite("missile0")
	if vd.missile0 == nil {
		return nil
	}
	vd.missile1 = newSprite("missile1")
	if vd.missile1 == nil {
		return nil
	}
	vd.ball = newSprite("ball")
	if vd.ball == nil {
		return nil
	}

	// missile/ball enabling
	vd.enam0Delay = newDelayCounter("(dis/en)abling")
	if vd.enam0Delay == nil {
		return nil
	}
	vd.enam1Delay = newDelayCounter("(dis/en)abling")
	if vd.enam1Delay == nil {
		return nil
	}
	vd.enablDelay = newDelayCounter("(dis/en)abling")
	if vd.enablDelay == nil {
		return nil
	}

	return vd
}

// TickPlayfield moves playfield on one video cycle
func (vd *Video) TickPlayfield() {
	// TODO: tick playfield
}

// TickSprites moves sprite elements on one video cycle
func (vd *Video) TickSprites() {
	// TODO: tick other sprites
	vd.TickBall()
}

// PixelColor returns the color of the pixel at the current time
func (vd Video) PixelColor() Color {
	col := vd.colubk
	if vd.ctrlpfPriority {
		// TODO: complete priority pixel ordering
		col = vd.PixelBall()
	} else {
		// TODO: complete non-priority pixel ordering
		col = vd.PixelBall()
	}
	return col
}

// ServiceTIAMemory checks the TIA memory for changes to registers that are
// interesting to the video sub-system
func (vd *Video) ServiceTIAMemory(register string, value uint8) bool {
	switch register {
	case "NUSIZ0":
	case "NUSIZ1":
	case "COLUP0":
	case "COLUP1":
	case "COLUPF":
	case "COLUBK":
	case "CTRLPF":
	case "REFP0":
	case "REFP1":
	case "PF0":
	case "PF1":
	case "PF2":
	case "RESP0":
	case "RESP1":
	case "RESM0":
	case "RESM1":
	case "RESBL":
	case "GRP0":
	case "GRP1":
	case "ENAM0":
	case "ENAM1":
	case "ENABL":
		vd.enablDelay.set(1, value&0x20 == 0x20)
	case "HMP0":
	case "HMP1":
	case "HMM0":
	case "HMM1":
	case "HMBL":
	case "VDELP0":
	case "VDELP1":
	case "VDELBL":
	case "RESMP0":
	case "RESMP1":
	case "HMCLR":
	case "CXCLR":
	}

	return false
}
