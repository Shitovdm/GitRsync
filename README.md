# GitRsync

![Release](https://github.com/Shitovdm/GitRsync/workflows/Release/badge.svg)
![Windows](https://github.com/Shitovdm/GitRsync/workflows/Windows/badge.svg)
![Linux](https://github.com/Shitovdm/GitRsync/workflows/Linux/badge.svg)
![Macos](https://github.com/Shitovdm/GitRsync/workflows/Macos/badge.svg)

### Build the project:  
`go build -o ./GitRsync.exe -ldflags -H=windowsgui`
`go build -o ./GitRsync.exe`

### Build app manifest:  
`rsrc -manifest GitRsync.exe.manifest -o rsrc.syso -ico="./public/assets/src/icon/iconwin.ico"`