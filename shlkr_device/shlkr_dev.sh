#!/bin/sh

#progs_debug="shlkr_dev.py"
progs_debug="./bin/shlkr_dev"

#log_file=./task_monitor.log


detect_deamons_debug()
{
    all_progs="$@"

    for prog in $all_progs; do
        
        pid=$(pgrep -f "${prog}")
    
        if [ -z $pid ]; then
            echo `date "+%Y-%m-%d %H:%M:%S"` ':         resart '${prog} >> ${log_file}
            case ${prog} in
                #run_device.py)
                #shlkr_dev.py)
                ./bin/shlkr_dev)
                    #python3 run_device.py & # 重启进程的命令，请相应修改
                    ./bin/shlkr_dev & 
                    ;;
                *)
                    echo `date "+%Y-%m-%d %H:%M:%S"` ':         cannt resart '${prog} >> ${log_file}
                    ;;
            esac
        fi
    done
}


#检测监视进程
while [ true ]; do
    detect_deamons_debug $progs_debug

    sleep 15  
 
done

