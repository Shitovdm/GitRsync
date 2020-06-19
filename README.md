# GitRsync

![Release](https://github.com/Shitovdm/GitRsync/workflows/Release/badge.svg)
![Windows](https://github.com/Shitovdm/GitRsync/workflows/Windows/badge.svg)
![Linux](https://github.com/Shitovdm/GitRsync/workflows/Linux/badge.svg)
![Macos](https://github.com/Shitovdm/GitRsync/workflows/Macos/badge.svg)

### Build the project:  
`go build -o ./GitRsync.exe -ldflags -H=windowsgui`
`go build -o ./GitRsync-x64.exe`

### Embedding binary resources:  
`rsrc -manifest x64.manifest -o rsrc.syso -ico="./public/assets/src/icon/iconwin.ico" -arch amd64`
`rsrc -manifest x32.manifest -o rsrc.syso -ico="./public/assets/src/icon/iconwin.ico" -arch 386`