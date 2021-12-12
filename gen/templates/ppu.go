package templates

type PPUStatusRegister struct {
	CoarseXScroll uint8 `bitfield:"5"`
	CoarseYScroll uint8 `bitfield:"5"`
	NameTable     uint8 `bitfield:"2"`
	FineYScroll   uint8 `bitfield:"3"`
	_             bool  `bitfield:"1"`
}
