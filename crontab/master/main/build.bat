echo off
echo "Start build master..."
echo "build Env"
set master_dir="../../deploy/master/"
ver > %master_dir%build_version.md
echo "build to " %master_dir%
rd /S /Q %master_dir%bin
echo "create dir"
mkdir %master_dir%bin\webroot
xcopy /y webroot %master_dir%bin\webroot
copy master.json %master_dir%bin\master.json

echo "build windows exe "
set GOOS=windows
go build -o %master_dir%bin\master.exe

echo "build linux bin  "
set GOOS=linux
go build -o %master_dir%bin/master
dir bin
echo "build from git branch ,see bin\build_version.md"
git branch >> %master_dir%build_version.md
echo "DONE"

echo on