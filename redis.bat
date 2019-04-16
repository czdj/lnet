cd %cd%
cd tools/redis
start redis /min redis.bat
start redis1 /min redis1.bat
start redis2 /min redis2.bat
python init.py
