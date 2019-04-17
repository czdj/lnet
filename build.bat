@echo off
echo input 1 to export lastest config 
echo input 2 build union for protobuf 2.5
echo input 3 build union for protobuf 3.6
echo input 4 to export lastest config and build union for protobuf 2.5
echo input 5 to export lastest config and build union for protobuf 3.6
python tools/build.py 3
pause
::set /p pb=chose opt type:
if %pb% == 1 (
echo please wait a moment
cd tools/conf
python exportConfig.py
cd ../..
)
if %pb% == 2 (
echo please wait a moment
echo working...
python tools/build.py 2
)
if %pb% == 3 (
echo please wait a moment
echo working...
python tools/build.py 3
)
if %pb% == 4 (
echo please wait a moment
echo working...
cd tools/conf
python exportConfig.py
cd ../..
python tools/build.py 2
)
if %pb% == 5 (
echo please wait a moment
echo working...
cd tools/conf
python exportConfig.py
cd ../..
python tools/build.py 3
)
echo 1 for runing server other for exit
set /p run=input 1 or other:
if %run% == 1 (
start run.bat
)
exit