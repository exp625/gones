# Gones - An attempt to program a NES emulator

The primary resources for the project is the great work on [wiki.nesdev.org]()
This project uses the two awesome libraries
- [https://github.com/faiface/pixel](https://github.com/faiface/pixel) for 2D display and user input
- [https://github.com/faiface/beep](https://github.com/faiface/beep) for audio streaming

## Building

### Windows (Powershell)
1. Build the Docker Image using `` docker build -t gones-builder .``
2. Compile the project: ``docker run -v "${PWD}":/usr/src/nes --rm builder``
3. Run the emulator: ``.\bin\nes.exe``