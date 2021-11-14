# Gones - An attempt to program a NES emulator

The following resources were used for the project: 
- The great work on [wiki.nesdev.org](wiki.nesdev.org)
- R650X and R651X Datasheet [http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf](http://www.6502.org/documents/datasheets/rockwell/rockwell_r650x_r651x.pdf)
- Reset behaviour description [https://www.pagetable.com/?p=410](https://www.pagetable.com/?p=410)
- The [cc65 assembler ](https://cc65.github.io/index.html)
- [https://github.com/clbr/neslib](https://github.com/clbr/neslib)

This project uses the two awesome libraries
- [https://github.com/faiface/pixel](https://github.com/faiface/pixel) for 2D display and user input
- [https://github.com/faiface/beep](https://github.com/faiface/beep) for audio streaming

## Building

### Windows (Powershell)
1. Build the Docker Image using `` docker build -t gones-builder .``
2. Compile the project: ``docker run -v "${PWD}":/usr/src/nes --rm builder``
3. Run the emulator: ``.\bin\nes.exe``

## Compiling C programs to 6502 using cc65
1. Write your c test program