#!/bin/sh

DEBUG=true #false

#progs="run_device.pyc"
progs="shlkr_dev.pyc"

#progs_debug="run_device.py"
progs_debug="shlkr_dev.py"

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
                shlkr_dev.py)
                    #python3 run_device.py & # 重启进程的命令，请相应修改
                    python3 shlkr_dev.py & 
                    ;;
                *)
                    echo `date "+%Y-%m-%d %H:%M:%S"` ':         cannt resart '${prog} >> ${log_file}
                    ;;
            esac
        fi
    done
}


detect_deamons()
{
    all_progs="$@"

    for prog in $all_progs; do
        
        pid=$(pgrep -f "${prog}")
    
        if [ -z $pid ]; then
            echo `date "+%Y-%m-%d %H:%M:%S"` ':         resart '${prog} >> ${log_file}
            case ${prog} in
                #run_device.pyc)
                shlkr_dev.pyc)
                    #python3 run_device.pyc & # 重启进程的命令，请相应修改
                    python3 shlkr_dev.pyc & 
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
    if $DEBUG
    then
        detect_deamons_debug $progs_debug
    else
        detect_deamons $progs
    fi

    sleep 15  
 
done

