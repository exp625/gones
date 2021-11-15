package cartridge

import "log"

type Cartridge struct {
	PrgRomSize uint8
	PrgRom     []uint8
	ChrRomSize uint8
	ChrRom     []uint8
	MirrorBit  bool
	Mapper

}

//An iNES file consists of the following sections, in order:
//
//Header (16 bytes)
//Trainer, if present (0 or 512 bytes)
//PRG ROM data (16384 * x bytes)
//CHR ROM data, if present (8192 * y bytes)
//PlayChoice INST-ROM, if present (0 or 8192 bytes)
//PlayChoice PROM, if present (16 bytes Data, 16 bytes CounterOut) (this is often missing, see PC10 ROM-Images for details)

//The format of the header is as follows:
//
//0-3: Constant $4E $45 $53 $1A ("NES" followed by MS-DOS end-of-file)
//4: Size of PRG ROM in 16 KB units
//5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
//6: Flags 6 - Mapper, mirroring, battery, trainer
//7: Flags 7 - Mapper, VS/Playchoice, NES 2.0
//8: Flags 8 - PRG-RAM size (rarely used extension)
//9: Flags 9 - TV system (rarely used extension)
//10: Flags 10 - TV system, PRG-RAM presence (unofficial, rarely used extension)
//11-15: Unused padding (should be filled with zero, but some rippers put their name across bytes 7-15)

func LoadCartridge (rom []byte) *Cartridge {
	prgRomSize := rom[4]
	chrRomSize := rom[5]
	mapperNumber := rom[6] & 0xF0 >> 4
	trainerPresent := rom[6] & 0b00000100 >> 2 == 1
	mirrowBit := rom[6] & 0b00000001 == 1

	prgRom :=  make([]uint8, int(prgRomSize) * 0x4000)
	chrRom :=  make([]uint8, int(prgRomSize) * 0x2000)
	ptr := 0xF
	if trainerPresent {
		ptr += 0x200
	}
	for i := 0; i < int(prgRomSize) * 0x4000; i++ {
		prgRom[i] = rom[ptr]
		ptr++
	}
	for i := 0; i < int(chrRomSize) * 0x2000; i++ {
		chrRom[i] = rom[ptr]
		ptr++
	}

	cartridge := &Cartridge{
		PrgRomSize: prgRomSize,
		PrgRom: prgRom,
		ChrRomSize: chrRomSize,
		ChrRom: chrRom,
		Mapper:     nil,
		MirrorBit: mirrowBit,
	}

	switch mapperNumber {
	case 0:
		cartridge.Mapper = NewMapper000(cartridge)
		log.Println("Created Cartridge with Mapper 000")
	case 2:
		cartridge.Mapper = NewMapper002(cartridge)
		log.Println("Created Cartridge with Mapper 002")
	default:
		log.Panic("Unsupported ROM File")
	}


	return cartridge
}