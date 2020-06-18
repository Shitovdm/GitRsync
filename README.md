# GitRsync

![Go](https://github.com/Shitovdm/GitRsync/workflows/Go/badge.svg)

### Build the project:  
`go build -o ./GitRsync.exe -ldflags -H=windowsgui`
`go build -o ./GitRsync.exe`

### Build app manifest:  
`rsrc -manifest GitRsync.exe.manifest -o rsrc.syso -ico="./public/assets/src/icon/iconwin.ico"`