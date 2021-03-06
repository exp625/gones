package cartridge

import (
	"crypto/md5"
	"github.com/exp625/gones/pkg/bus"
	"log"
)

type Cartridge struct {
	Mapper
	Bus        bus.Bus
	PrgRomSize uint8
	PrgRom     []uint8
	ChrRomSize uint8
	ChrRom     []uint8
	ChrRam     bool
	MirrorBit  bool
	Identifier [16]byte
}

// Load loads a Cartridge from an iNES file.
//
// An iNES file consists of the following sections, in order:
//
// Header (16 bytes)
// Trainer, if present (0 or 512 bytes)
// PRG ROM data (16384 * x bytes)
// CHR ROM data, if present (8192 * y bytes)
// PlayChoice INST-ROM, if present (0 or 8192 bytes)
// PlayChoice PROM, if present (16 bytes Data, 16 bytes CounterOut) (this is often missing, see PC10 ROM-Images for details)
//
// The format of the header is as follows:
//
// 0-3: Constant $4E $45 $53 $1A ("NES" followed by MS-DOS end-of-file)
// 4: Size of PRG ROM in 16 KB units
// 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
// 6: Flags 6 - Mapper, mirroring, battery, trainer
// 7: Flags 7 - Mapper, VS/Playchoice, NES 2.0
// 8: Flags 8 - PRG-RAM size (rarely used extension)
// 9: Flags 9 - TV system (rarely used extension)
// 10: Flags 10 - TV system, PRG-RAM presence (unofficial, rarely used extension)
// 11-15: Unused padding (should be filled with zero, but some rippers put their name across bytes 7-15)
func Load(rom []byte, bus bus.Bus) *Cartridge {
	prgRomSize := rom[4]
	chrRomSize := rom[5]
	chrRam := false
	if chrRomSize == 0 {
		chrRomSize = 1
		chrRam = true
	}
	mapperNumberLo := rom[6] >> 4
	mapperNumberHi := rom[7] >> 4
	mapperNumber := mapperNumberHi<<4 | mapperNumberLo

	trainerPresent := (rom[6]&0b0000_0100)>>2 == 1
	mirrorBit := rom[6]&0b0000_0001 == 1

	prgRom := make([]uint8, int(prgRomSize)*0x4000)
	chrRom := make([]uint8, int(chrRomSize)*0x2000)

	ptr := 0x10
	if trainerPresent {
		log.Println("Trainer present!")
		ptr += 0x200
	}
	for i := 0; i < int(prgRomSize)*0x4000; i++ {
		prgRom[i] = rom[ptr]
		ptr++
	}
	if !chrRam {
		for i := 0; i < int(chrRomSize)*0x2000; i++ {
			chrRom[i] = rom[ptr]
			ptr++
		}
	}

	c := &Cartridge{
		Bus:        bus,
		PrgRomSize: prgRomSize,
		PrgRom:     prgRom,
		ChrRomSize: chrRomSize,
		ChrRom:     chrRom,
		ChrRam:     chrRam,
		MirrorBit:  mirrorBit,
		Identifier: md5.Sum(rom),
	}

	switch mapperNumber {
	case 0:
		c.Mapper = NewMapper000(c)
		log.Println("Created Cartridge with Mapper 000")
	case 1:
		c.Mapper = NewMapper001(c)
		log.Println("Created Cartridge with Mapper 001")
	case 2:
		c.Mapper = NewMapper002(c)
		log.Println("Created Cartridge with Mapper 002")
	case 3:
		c.Mapper = NewMapper003(c)
		log.Println("Created Cartridge with Mapper 003")
	case 4:
		c.Mapper = NewMapper004(c)
		log.Println("Created Cartridge with Mapper 004")
	case 7:
		c.Mapper = NewMapper007(c)
		log.Println("Created Cartridge with Mapper 007")
	default:
		log.Printf("Unsupported ROM File with Mapper %00d", mapperNumber)
		return nil
	}

	return c
}
