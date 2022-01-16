rsrc -manifest test.manifest --ico icon.ico -o rsrc.syso
go build -ldflags="-H windowsgui -w -s"  -o modao.exe
upx -9 *.exe