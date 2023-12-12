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
		{templates.APUControlRegister{}, "pkg/apu/control_register.gen.go", "apu", "ControlRegister"},
		{templates.APUStatusRegister{}, "pkg/apu/status_register.gen.go", "apu", "StatusRegister"},
		{templates.APUFrameCounterRegister{}, "pkg/apu/frame_counter_register.gen.go", "apu", "FrameCounterRegister"},
		{templates.APUPulseChannelGlobalRegister{}, "pkg/apu/pulse_channel_global_register.gen.go", "apu", "PulseChannelGlobalRegister"},
		{templates.APUPulseChannelSweepRegister{}, "pkg/apu/pulse_channel_sweep_register.gen.go", "apu", "PulseChannelSweepRegister"},
		{templates.APUPulseChannelTimerLowRegister{}, "pkg/apu/pulse_channel_timer_low_register.gen.go", "apu", "PulseChannelTimerLowRegister"},
		{templates.APUPulseChannelTimerHighRegister{}, "pkg/apu/pulse_channel_timer_high_register.gen.go", "apu", "PulseChannelTimerHighRegister"},
		{templates.APUTriangleChannelGlobalRegister{}, "pkg/apu/triangle_channel_global_register.gen.go", "apu", "TriangleChannelGlobalRegister"},
		{templates.APUTriangleChannelTimerLowRegister{}, "pkg/apu/triangle_channel_timer_low_register.gen.go", "apu", "TriangleChannelTimerLowRegister"},
		{templates.APUTriangleChannelTimerHighRegister{}, "pkg/apu/triangle_channel_timer_high_register.gen.go", "apu", "TriangleChannelTimerHighRegister"},
		{templates.APUNoiseChannelGlobalRegister{}, "pkg/apu/noise_channel_global_register.gen.go", "apu", "NoiseChannelGlobalRegister"},
		{templates.APUNoiseChannelPeriodRegister{}, "pkg/apu/noise_channel_period_register.gen.go", "apu", "NoiseChannelPeriodRegister"},
		{templates.APUNoiseChannelLengthRegister{}, "pkg/apu/noise_channel_length_register.gen.go", "apu", "NoiseChannelLengthRegister"},
		{templates.APUDMCChannelGlobalRegister{}, "pkg/apu/dmc_channel_global_register.gen.go", "apu", "DMCChannelGlobalRegister"},
		{templates.APUDMCChannelDirectLoadRegister{}, "pkg/apu/dmc_channel_direct_load_register.gen.go", "apu", "DMCChannelDirectLoadRegister"},
		{templates.APUDMCChannelSampleAddressRegister{}, "pkg/apu/dmc_channel_sample_address_register.gen.go", "apu", "DMCChannelSampleAddressRegister"},
		{templates.APUDMCChannelSampleLengthRegister{}, "pkg/apu/dmc_channel_sample_length_register.gen.go", "apu", "DMCChannelSampleLengthRegister"},
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
