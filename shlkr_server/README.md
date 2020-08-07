linux-devices manager/controllor/debugger tool (server part)
qding(sz) luohaifeng
2020.8.4

program 1:
1. nohup python3 run_server_db.py &
   recommend: nohup sh autoRestartPy_server_db.sh &
   this program should run for 7x24

program 2:
1. python3 run_server_echo.py
2. when logging in, if you don't have tool-account, could use "test"/"test".
3. after logging in, input ".help", you will get all info.
