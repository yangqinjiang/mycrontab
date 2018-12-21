echo off
echo "Start"
echo "build Env"
ver
rd /S /Q bin
echo "create dir"
mkdir bin\webroot
xcopy /y webroot bin\webroot
copy master.json bin\master.json

echo "build windows exe "
set GOOS=windows
go build -o bin\master.exe

echo "build linux bin  "
set GOOS=linux
go build -o bin/master
dir bin
echo "build from git branch ,see bin\build_version.md"
git branch >> bin\build_version.md
echo "DONE"

echo on