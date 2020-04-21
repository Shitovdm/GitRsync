# GitRsync

### Build the project:  
`go build -o ./GitRsync.exe -ldflags -H=windowsgui`

### Build app manifest:  
`rsrc -manifest GitRsync.exe.manifest -o rsrc.syso -ico="./public/assets/src/icon/iconwin.ico"`