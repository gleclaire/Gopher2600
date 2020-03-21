// This file is part of Gopher2600.
//
// Gopher2600 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gopher2600 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gopher2600.  If not, see <https://www.gnu.org/licenses/>.
//
// *** NOTE: all historical versions of this file, as found in any
// git repository, are also covered by the licence, even when this
// notice is not present ***

package colors

// colors used when alternative colors are selected. these colors mirror the
// so called debug colors used by the Stella emulator
var alt32bit = []uint32{
	0x111111, 0x84c8fc, 0x9246c0, 0x901c00, 0xe8e84a, 0xd5824a, 0x328432,
}

// the raw color values are the component values expressed as a single 32 bit
// number. we'll use these raw values in the init() function below to create
// the real palette
var ntsc32bit = []uint32{
	0x000000, 0x404040, 0x6c6c6c, 0x909090, 0xb0b0b0, 0xc8c8c8, 0xdcdcdc, 0xececec,
	0x444400, 0x646410, 0x848424, 0xa0a034, 0xb8b840, 0xd0d050, 0xe8e85c, 0xfcfc68,
	0x702800, 0x844414, 0x985c28, 0xac783c, 0xbc8c4c, 0xcca05c, 0xdcb468, 0xecc878,
	0x841800, 0x983418, 0xac5030, 0xc06848, 0xd0805c, 0xe09470, 0xeca880, 0xfcbc94,
	0x880000, 0x9c2020, 0xb03c3c, 0xc05858, 0xd07070, 0xe08888, 0xeca0a0, 0xfcb4b4,
	0x78005c, 0x8c2074, 0xa03c88, 0xb0589c, 0xc070b0, 0xd084c0, 0xdc9cd0, 0xecb0e0,
	0x480078, 0x602090, 0x783ca4, 0x8c58b8, 0xa070cc, 0xb484dc, 0xc49cec, 0xd4b0fc,
	0x140084, 0x302098, 0x4c3cac, 0x6858c0, 0x7c70d0, 0x9488e0, 0xa8a0ec, 0xbcb4fc,
	0x000088, 0x1c209c, 0x3840b0, 0x505cc0, 0x6874d0, 0x7c8ce0, 0x90a4ec, 0xa4b8fc,
	0x00187c, 0x1c3890, 0x3854a8, 0x5070bc, 0x6888cc, 0x7c9cdc, 0x90b4ec, 0xa4c8fc,
	0x002c5c, 0x1c4c78, 0x386890, 0x5084ac, 0x689cc0, 0x7cb4d4, 0x90cce8, 0xa4e0fc,
	0x003c2c, 0x1c5c48, 0x387c64, 0x509c80, 0x68b494, 0x7cd0ac, 0x90e4c0, 0xa4fcd4,
	0x003c00, 0x205c20, 0x407c40, 0x5c9c5c, 0x74b474, 0x8cd08c, 0xa4e4a4, 0xb8fcb8,
	0x143800, 0x345c1c, 0x507c38, 0x6c9850, 0x84b468, 0x9ccc7c, 0xb4e490, 0xc8fca4,
	0x2c3000, 0x4c501c, 0x687034, 0x848c4c, 0x9ca864, 0xb4c078, 0xccd488, 0xe0ec9c,
	0x442800, 0x644818, 0x846830, 0xa08444, 0xb89c58, 0xd0b46c, 0xe8cc7c, 0xfce08c,
}

var pal32bit = []uint32{
	0x000000, 0x282828, 0x505050, 0x747474, 0x949494, 0xb4b4b4, 0xd0d0d0, 0xececec,
	0x000000, 0x282828, 0x505050, 0x747474, 0x949494, 0xb4b4b4, 0xd0d0d0, 0xececec,
	0x805800, 0x947020, 0xa8843c, 0xbc9c58, 0xccac70, 0xdcc084, 0xecd09c, 0xfce0b0,
	0x445c00, 0x5c7820, 0x74903c, 0x8cac58, 0xa0c070, 0xb0d484, 0xc4e89c, 0xd4fcb0,
	0x703400, 0x885020, 0xa0683c, 0xb48458, 0xc89870, 0xdcac84, 0xecc09c, 0xfcd4b0,
	0x006414, 0x208034, 0x3c9850, 0x58b06c, 0x70c484, 0x84d89c, 0x9ce8b4, 0xb0fcc8,
	0x700014, 0x882034, 0xa03c50, 0xb4586c, 0xc87084, 0xdc849c, 0xec9cb4, 0xfcb0c8,
	0x005c5c, 0x207474, 0x3c8c8c, 0x58a4a4, 0x70b8b8, 0x84c8c8, 0x9cdcdc, 0xb0ecec,
	0x70005c, 0x842074, 0x943c88, 0xa8589c, 0xb470b0, 0xc484c0, 0xd09cd0, 0xe0b0e0,
	0x003c70, 0x1c5888, 0x3874a0, 0x508cb4, 0x68a4c8, 0x7cb8dc, 0x90ccec, 0xa4e0fc,
	0x580070, 0x6c2088, 0x803ca0, 0x9458b4, 0xa470c8, 0xb484dc, 0xc49cec, 0xd4b0fc,
	0x002070, 0x1c3c88, 0x3858a0, 0x5074b4, 0x6888c8, 0x7ca0dc, 0x90b4ec, 0xa4c8fc,
	0x3c0080, 0x542094, 0x6c3ca8, 0x8058bd, 0x9470cc, 0xa884dc, 0xb89cec, 0xc8b0fc,
	0x000088, 0x20209c, 0x3c3cb0, 0x5858c0, 0x7070d0, 0x8484e0, 0x9c9cec, 0xb0b0fc,
	0x000000, 0x282828, 0x505050, 0x747474, 0x949494, 0xb4b4b4, 0xd0d0d0, 0xececec,
	0x000000, 0x282828, 0x505050, 0x747474, 0x949494, 0xb4b4b4, 0xd0d0d0, 0xececec,
}

// this init() function converts the "raw" color values to the RGB components
func init() {
	for _, col := range ntsc32bit {
		red, green, blue := byte((col&0xff0000)>>16), byte((col&0xff00)>>8), byte(col&0xff)

		// repeat color twice in palette
		PaletteNTSC = append(PaletteNTSC, RGB{red, green, blue})
		PaletteNTSC = append(PaletteNTSC, RGB{red, green, blue})
	}

	for _, col := range pal32bit {
		red, green, blue := byte((col&0xff0000)>>16), byte((col&0xff00)>>8), byte(col&0xff)

		// repeat color twice in palette
		PalettePAL = append(PalettePAL, RGB{red, green, blue})
		PalettePAL = append(PalettePAL, RGB{red, green, blue})
	}

	for _, col := range alt32bit {
		red, green, blue := byte((col&0xff0000)>>16), byte((col&0xff00)>>8), byte(col&0xff)
		PaletteAlt = append(PaletteAlt, RGB{red, green, blue})
	}
}
