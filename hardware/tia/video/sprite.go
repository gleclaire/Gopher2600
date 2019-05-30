package video

import (
	"fmt"
	"gopher2600/hardware/tia/delay"
	"gopher2600/hardware/tia/delay/future"
	"gopher2600/hardware/tia/polycounter"
	"strings"
)

// the sprite type is used for those video elements that move about - players,
// missiles and the ball. the VCS doesn't really have anything called a sprite
// but we all know what it means
type sprite struct {
	// label is the name of a particular instance of a sprite (eg. player0 or
	// missile 1)
	label string

	// colorClock references the VCS wide color clock. we only use it to note
	// the Pixel() value of the color clock at the reset point of the sprite.
	colorClock *polycounter.Polycounter

	// position of the sprite as a polycounter value - the basic principle
	// behind VCS sprites is to begin drawing of the sprite when position
	// circulates to zero
	position polycounter.Polycounter

	// horizontal position of the sprite - may be affected by HMOVE
	resetPixel int

	// horizontal position after hmove has been applied
	currentPixel int

	// the draw signal controls which "bit" of the sprite is to be drawn next.
	// generally, the draw signal is activated when the position polycounter
	// matches the colorClock polycounter, but differenct sprite types handle
	// this differently in certain circumstances
	graphicsScanCounter int
	graphicsScanMax     int
	graphicsScanOff     int

	// the amount of horizontal movement for the sprite
	// -- as set by the 6502 - written into the HMP0/P1/M0/M1/BL register
	// -- normalised into the 0 to 15 range
	horizMovement int
	// -- whether HMOVE is still affecting this sprite
	moreMovementRequired bool

	// each type of sprite has slightly different spriteTick logic which needs
	// to be called from within the HMOVE logic common to all sprite types
	spriteTick func()

	// a note on whether the sprite is about to be reset its position
	resetFuture *future.Instance

	// 0 = force reset is off
	// 1 = force reset trigger
	// n = wait for trigger
	forceReset int
	// see comment in resolveHorizMovement()
}

func newSprite(label string, colorClock *polycounter.Polycounter, spriteTick func()) *sprite {
	sp := new(sprite)
	sp.label = label
	sp.colorClock = colorClock
	sp.spriteTick = spriteTick

	sp.position = *polycounter.New6Bit()
	sp.position.SetResetPoint(39) // "101101"

	// the direction of count and max is important - don't monkey with it
	sp.graphicsScanMax = 8
	sp.graphicsScanOff = sp.graphicsScanMax + 1
	sp.graphicsScanCounter = sp.graphicsScanOff

	return sp
}

// MachineInfoTerse returns the sprite information in terse format
func (sp sprite) MachineInfoTerse() string {
	s := strings.Builder{}
	s.WriteString(sp.label)
	s.WriteString(": ")
	s.WriteString(sp.position.String())
	s.WriteString(fmt.Sprintf(" pos=%d", sp.currentPixel))
	if sp.isDrawing() {
		s.WriteString(fmt.Sprintf(" drw=%d", sp.graphicsScanMax-sp.graphicsScanCounter))
	} else {
		s.WriteString(" drw=-")
	}
	if sp.resetFuture == nil {
		s.WriteString(" res=-")
	} else {
		s.WriteString(fmt.Sprintf(" res=%d", sp.resetFuture.RemainingCycles))
	}

	return s.String()
}

// MachineInfo returns the Video information in verbose format
func (sp sprite) MachineInfo() string {
	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("%s:\n", sp.label))
	s.WriteString(fmt.Sprintf("   polycounter: %s\n", sp.position))
	if sp.isDrawing() {
		s.WriteString(fmt.Sprintf("   drawing: %d\n", sp.graphicsScanMax-sp.graphicsScanCounter))
	} else {
		s.WriteString("   drawing: inactive\n")
	}
	if sp.resetFuture == nil {
		s.WriteString("   reset: none scheduled\n")
	} else {
		s.WriteString(fmt.Sprintf("   reset: %d cycles\n", sp.resetFuture.RemainingCycles))
	}

	// information about horizontal movement.
	// - horizMovement value normalised and inverted so that positive numbers
	// indicate movement to the right and negative numbers indicate movement to
	// the left
	// - value in square brackets is the value that was originally poked into
	// the move register
	// - the 4 bit binary number at the end is the representation of what the HMOVE
	// circuitry interacts with, bit-by-bit - see resolveHorizMovement()
	s.WriteString(fmt.Sprintf("   hmove: %d [%#02x] %04b\n", -sp.horizMovement+8, (sp.horizMovement<<4)^0x80, sp.horizMovement))

	s.WriteString(fmt.Sprintf("   reset pixel: %d\n", sp.resetPixel))
	s.WriteString(fmt.Sprintf("   current pixel: %d", sp.currentPixel))
	if sp.moreMovementRequired {
		s.WriteString(" *\n")
	} else {
		s.WriteString("\n")
	}

	return s.String()
}

// EmulatorInfo returns low state information about the type
func (sp sprite) EmulatorInfo() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%04b ", sp.horizMovement))
	if sp.moreMovementRequired {
		s.WriteString("*")
	} else {
		s.WriteString(" ")
	}
	s.WriteString(" ")
	s.WriteString(sp.label)
	return s.String()
}

func (sp *sprite) resetPosition() {
	sp.position.Reset()

	// note reset position of sprite, in pixels
	sp.resetPixel = sp.colorClock.Pixel()
	sp.currentPixel = sp.resetPixel
}

func (sp *sprite) checkForGfxStart(triggerList []int) (bool, bool) {
	if sp.position.Tick() {
		return true, false
	}

	// check for start positions of additional copies of the sprite
	for _, v := range triggerList {
		if v == sp.position.Count && sp.position.Phase == 0 {
			return true, true
		}
	}

	return false, false
}

func (sp *sprite) forceHMOVE(adjustment int) {
	hm := (sp.horizMovement - (15 - adjustment))
	for i := 0; i < hm; i++ {
		// adjust position information
		sp.currentPixel--
		if sp.currentPixel < 0 {
			sp.currentPixel = 159
		}

		// perform an additional tick of the sprite (different sprite types
		// have different tick logic)
		sp.spriteTick()
	}
}

func (sp *sprite) prepareForHMOVE() {
	// start horizontal movment of this sprite
	sp.moreMovementRequired = true

	// at beginning of hmove sequence, without knowing anything else, the final
	// position of the sprite will be the current position plus 8. the actual
	// value will be reduced depending on what happens during hmove ticking.
	// factors that effect the final position:
	//   o the value in the horizontal movement register (eg. HMP0)
	//   o whether the ticking is occuring during the hblank period
	// both these factors are considered in the resolveHorizMovement() function
	sp.currentPixel += 8
}

func compareBits(a, b uint8) bool {
	// return true if any corresponding bits in the lower nibble are the same
	return a&0x08 == b&0x08 || a&0x04 == b&0x04 || a&0x02 == b&0x02 || a&0x01 == b&0x01
}

func (sp *sprite) resolveHMOVE(count int) {
	sp.moreMovementRequired = sp.moreMovementRequired && compareBits(uint8(count), uint8(sp.horizMovement))

	if sp.moreMovementRequired {
		// this mental construct is designed to fix a problem in the Keystone
		// Kapers ROM. I don't believe for a moment that this is a perfect
		// solution but it makes sense in the context of that ROM.
		//
		// What seems to be happening in Keystone Kapers ROM is this:
		//
		//	o Ball is reset at end of scanline 95 ($f756); and other scanlines
		//  o HMOVE is tripped at beginning of line 96
		//  o but reset doesn't occur until we resume motion clocks, by which
		//		time HMOVE is finished
		//  o moreover, the game doesn't want the ball to appear at the
		//		beginning of the visible part of the screen; it wants the ball
		//		to appear in the HMOVE gutter on scanlines 97 and 98; so the
		//		move adjustments needs to happen such that the ball really
		//		appears at the end of the scanline
		//  o to cut a long story short, the game needs the ball to have been
		//		reset before the HMOVE has completed on line 96
		//
		// confusing huh?  this delay construct fixes the above issue while not
		// breaking other regression tests. I don't know if this is a generally
		// correct solution or if it's specific to the ball sprite but I'm
		// keeping it in for now.
		if sp.resetFuture != nil {
			if sp.forceReset == 1 {
				sp.resetFuture.Force()
				sp.forceReset = 0
			} else if sp.forceReset == 0 {
				sp.forceReset = delay.ForceReset
			} else {
				sp.forceReset--
			}
		}

		// adjust position information
		sp.currentPixel--
		if sp.currentPixel < 0 {
			sp.currentPixel = 159
		}

		// perform an additional tick of the sprite (different sprite types
		// have different tick logic)
		sp.spriteTick()
	}
}

func (sp *sprite) startDrawing() {
	sp.graphicsScanCounter = 0
}

func (sp *sprite) isDrawing() bool {
	return sp.graphicsScanCounter <= sp.graphicsScanMax
}

func (sp *sprite) tickGraphicsScan() {
	if sp.isDrawing() {
		sp.graphicsScanCounter++
	}
}
