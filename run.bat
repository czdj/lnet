@echo off
set notfind=0
tasklist|find "redis-server">nul|find "Services">nul||set notfind=1
if %notfind% == 1 (
echo not found redis now auto staret
cd %cd%
cd tools/redis
::start redis /min redis.bat
::start redis1 /min redis1.bat
::python init.py
cd ..
cd ..
)
echo 0 start server group except test server
echo 1 auth_1
echo 101 logic_101
echo 102 logic_102
echo 2001 match_2001
echo 3000 battle_3000
echo 4000 league_4000
echo 4001 league_4001
echo 100001 test_100001.yml

:LOOP
set /p file=chose conf file:
cd %cd%
if %file% == 0 (
start "auth_1" /min xserver.exe -f conf/auth_1.yml -n 1
echo start auth_1
timeout /T 1
start "league_4000" /min xserver.exe -f conf/league_4000.yml -n 1
echo start league_4000
timeout /T 1
start "league_4001" /min xserver.exe -f conf/league_4001.yml -n 1
echo start league_4001
timeout /T 1
start "logic_101" /min xserver.exe -f conf/logic_101.yml -n 1
echo start logic_101
timeout /T 1
start "logic_102" /min xserver.exe -f conf/logic_102.yml -n 1
echo start logic_102
timeout /T 1
start "match_2001" /min xserver.exe -f conf/match_2001.yml -n 1
echo start match_2001
timeout /T 1
start "battle_3000" /min xserver.exe -f conf/battle_3000.yml -n 1
echo start battle_3000
echo start server group end
goto end
)
if %file% == 1 (
start "auth_1" /min xserver.exe -f conf/auth_1.yml -n 1
)
if %file% == 101 (
start "logic_101" /min xserver.exe -f conf/logic_101.yml -n 1
)
if %file% == 102 (
start "logic_102" /min xserver.exe -f conf/logic_102.yml -n 1
)
if %file% == 2001 (
start "match_2001" /min xserver.exe -f conf/match_2001.yml -n 1
)
if %file% == 3000 (
start "battle_3000" /min xserver.exe -f conf/battle_3000.yml -n 1
)
if %file% == 4000 (
start "league_4000" /min xserver.exe -f conf/league_4000.yml -n 1
)
if %file% == 4001 (
start "league_4001" /min xserver.exe -f conf/league_4001.yml -n 1
)
if %file% == 100001 (
start "test_100001" /min xserver.exe -f conf/test_100001.yml -n 1
)
goto LOOP

:end

pause

