package templates

type PPUControlRegister struct {
	NMI              bool  `bitfield:"1"`
	PPUMode          bool  `bitfield:"1"`
	SpriteSize       uint8 `bitfield:"1"`
	PatternTable     uint8 `bitfield:"1"`
	SpriteTable      uint8 `bitfield:"1"`
	VRAMIncrement    bool  `bitfield:"1"`
	NameTableAddress uint8 `bitfield:"2"`
}

type PPUMaskRegister struct {
	Emphasize          uint8 `bitfield:"3"`
	ShowSprites        bool  `bitfield:"1"`
	ShowBackground     bool  `bitfield:"1"`
	SpritesLeftmost    bool  `bitfield:"1"`
	BackgroundLeftmost bool  `bitfield:"1"`
	Greyscale          bool  `bitfield:"1"`
}

type PPUStatusRegister struct {
	VerticalBlank  bool  `bitfield:"1"`
	SpriteZeroHit  bool  `bitfield:"1"`
	SpriteOverflow bool  `bitfield:"1"`
	_              uint8 `bitfield:"5"`
}

type PPUAddressRegister struct {
	CoarseXScroll uint8 `bitfield:"5"`
	CoarseYScroll uint8 `bitfield:"5"`
	NameTable     uint8 `bitfield:"2"`
	FineYScroll   uint8 `bitfield:"3"`
	_             bool  `bitfield:"1"`
}
