package main

import (
	"github.com/exp625/gones/gen"
	"github.com/exp625/gones/gen/bitfield"
	"github.com/exp625/gones/gen/templates"
	"log"
)

type GenConf struct {
	structInstance interface{}
	fileName       string
	packageName    string
	structName     string
}

func main() {
	for _, entry := range []GenConf{
		{templates.CPUStatusRegister{}, "pkg/cpu/status_register.gen.go", "cpu", "StatusRegister"},
		{templates.PPUControlRegister{}, "pkg/ppu/control_register.gen.go", "ppu", "ControlRegister"},
		{templates.PPUMaskRegister{}, "pkg/ppu/mask_register.gen.go", "ppu", "MaskRegister"},
		{templates.PPUStatusRegister{}, "pkg/ppu/status_register.gen.go", "ppu", "StatusRegister"},
		{templates.PPUAddressRegister{}, "pkg/ppu/address_register.gen.go", "ppu", "AddressRegister"},
	} {
		if err := GenerateBitfield(entry); err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateBitfield(e GenConf) error {
	w := gen.NewCodeWriter()
	defer w.WriteGoFile(e.fileName, e.packageName)

	if err := bitfield.Gen(w, e.structInstance, e.structName, nil); err != nil {
		return err
	}
	return nil
}
