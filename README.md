# AOSPA Flashboot Tool - Marble Edition

<img src="tool.png" alt="Image" width="640px">

## Getting started

AOSPA flashtool is a helper GUI tool to allow users to flash aospa roms on their Marble device. No more CLI invocation is needed to switch to experience AOSPA.

Features:
- Less errorprone
- Basic firmware validation
- Easy UI

> NOTE: USE AT YOUR OWN RISK <br>
Flashing wrong firmware or ROM files could brick your device. Be cautious about the action you perform.

## Installation
Download the binary from Releases specific to your operating system. 

## Usage
Download the required firmware and fastboot ROM file from aospa.co

> Step1: Check Flash Firmware and select the firmware zip file 

and/or

> Step2: Check Flash ROM and select the rom zip file

> Step3: Click the Flash button.
Once the  flashing is completed. Click Reboot button.

Flashing might take 5 to 10 minutes depending on the underlying hardware.

## Limitations
- No realtime output. Log output displayed as buffered streams which comes with some delay.

## Roadmap
- allow driver installation on windows
- Support kernel flashing from fastboot
- more validation logic to prevent user flashing incorrect files

## Authors and acknowledgment
This project is not possible without the engine created by ghost-rider reborn. 

## License
Licensed under GPL