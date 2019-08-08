echo off
echo "Start build worker..."
echo "build Env"
set worker_dir="../../deploy/worker/"
ver > %worker_dir%build_version.md
echo "build to " %worker_dir%
rd /S /Q %worker_dir%bin
echo "create dir"
mkdir %worker_dir%bin
copy worker.json %worker_dir%bin\worker.json

echo "build windows exe "
set GOOS=windows
go build -o %worker_dir%bin/worker.exe

echo "build linux bin  "
set GOOS=linux
go build -o %worker_dir%bin/worker
dir %worker_dir%bin
echo "build from git branch ,see bin\build_version.md"
git branch >> %worker_dir%build_version.md
echo "DONE"

echo on