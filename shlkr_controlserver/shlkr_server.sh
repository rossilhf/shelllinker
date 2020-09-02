#!/bin/sh

progs_debug="run_server_db.py run_server_update.py"

log_file=./task_monitor.log

detect_deamons_debug()
{
    all_progs="$@"

    for prog in $all_progs; do
        
        pid=$(pgrep -f "${prog}")
    
        if [ -z $pid ]; then
            echo `date "+%Y-%m-%d %H:%M:%S"` ':         resart '${prog} >> ${log_file}
            case ${prog} in
                run_server_db.py)
                    #cd /home/aiot/workspace/qdh_eyecloud_aibox/qdh_remoteDebugTool/devicePart
                    python3 run_server_db.py & # 重启进程的命令，请相应修改
                    ;;
                run_server_update.py)
                    #cd /home/aiot/workspace/qdh_eyecloud_aibox/qdh_remoteDebugTool/devicePart
                    python3 run_server_update.py & # 重启进程的命令，请相应修改
                    ;;
                *)
                    echo `date "+%Y-%m-%d %H:%M:%S"` ':         cannt resart '${prog} >> ${log_file}
                    ;;
            esac
        fi
    done
}

while [ true ]; do
    #检测监视进程
    detect_deamons_debug $progs_debug

    sleep 15  
 
done

