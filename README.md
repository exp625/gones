# GoNES - An attempt to program a NES emulator

The following resources were used for the project:

- The great work on [wiki.nesdev.org](https://wiki.nesdev.org)
- R650X and R651X
  Datasheet [http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf](http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf)
- Reset behaviour description [https://www.pagetable.com/?p=410](https://www.pagetable.com/?p=410)
- The [cc65 assembler ](https://cc65.github.io/index.html)
- [https://github.com/clbr/neslib](https://github.com/clbr/neslib)
- [NTSC NES palette generator](https://bisqwit.iki.fi/utils/nespalette.php)

This project the following awesome libraries

- [Ebiten](https://ebiten.org/) for 2D display, audio and input

## Status

![](./screenshots/VerboseLogging.PNG?raw=true)

- [x] CPU Working "_We've built a working 6502 emulator. Kinda cool_"
- [x] Nes Illegal Opcodes "_Look at the debug tool, it's amazing_"
- [ ] PPU Debug Working
- [ ] DMA Working
- [ ] PPU Background Rendering Clock Accurate
- [ ] PPU Sprite Rendering Working Clock Accurate
- [ ] Input Working
- [ ] APU Working
- Implemented iNES Mappers: Mapper000, Mapper002

## Usage

Start the emulator with ``nes romfile.rom``
The rom file should be valid rom file including iNES header. You can build your own rom file with the description below.

## Controls

* ``Space`` - Start or Stop auto mode
* ``Enter`` - Execute one CPU instruction
* ``Arrow Up`` - Execute one CPU Clock
* ``Arror Left`` - Execute one Master/PPU Clock
* ``R`` Reset
* ``F1`` Hide/Display CPU Debug display
* ``F2`` Hide/Display Pattern Tables
* ``F3`` Hide/Display Nametable Information
* ``F4`` Hide/Display Palette Information
* ``L`` Enable Logging
* ``Keypad`` Enter requested instruction that should be executed when pressing enter. 0 = Disabled,
* ``Esc`` Reset requested instructions

## Building

### For Windows

1. No CGO required. Just run ``go buid github.com/exp625/gones``

### For Linux

TODO

## Compiling C programs to 6502 using cc65

1. Install the [cc65 compiler ](https://github.com/cc65/cc65)
2. Write your C program inside the test folder ``test.c``
3. Assemble you C program to 6502 assembler ``cc65 -Os -T -t nes test.c``
4. Create object files for
    1. Your assembled program ``ca65 -t nes test.s``
    2. The startup code ``ca65 -t nes crt0.s``
    3. The default nes characters ``ca65 -t nes chars.s``
5. Create your rom file ``ld65 -C memory.cfg test.o crt0.o chars.o nes.lib -o $FILE.nes``

For a one liner (Set FILE
accordingly): ``FILE="test" && cc65 -Os -T -t nes $FILE.c && ca65 -t nes $FILE.s && ca65 -t nes crt0.s && ca65 -t nes chars.s && ld65 -C memory.cfg $FILE.o crt0.o chars.o nes.lib -o $FILE.nes``

## Testing

### CPU Testing

1. Download the awesome [nestest.rom](http://nickmass.com/images/nestest.nes)
2. Download the known good [cpu log](https://www.qmtpro.com/~nes/misc/nestest.log)
3. Start the emulator with the nestest.rom file ``nes nestest.rom``
4. Force the emulator to start execution on 0xC000 and set P flag to 0x24 using the ``Q`` Key (
   See [Note](https://wiki.nesdev.org/w/index.php?title=CPU_power_up_state#cite_note-reset-stack-push-3) for why)
5. Enable logging to file using ``L`` Key
6. Run the emulator step by step or in auto mode
7. Compare generated log file to the know good log file