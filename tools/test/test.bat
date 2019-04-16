@echo off
setlocal EnableDelayedExpansion
echo 1 test_login
echo 2 test_social
set /p file=chose test file:
cd %cd%
if %file% == 1 (
set /p name=input name:
python test_login.py !name!
)
if %file% == 2 (
set /p name=input name:
python test_social.py !name!
)

pause