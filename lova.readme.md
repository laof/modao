rsrc -manifest test.manifest --ico icon.ico -o rsrc.syso

go build -ldflags="-w -s"
upx -9 *.exe