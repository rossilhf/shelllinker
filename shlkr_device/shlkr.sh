#!/bin/sh

progs_debug="./shlkr"
#log_file=./task_monitor.log

detect_deamons_debug()
{
    all_progs="$@"

    for prog in $all_progs; do
        
        pid=$(pgrep -f "${prog}")
    
        if [ -z $pid ]; then
            echo `date "+%Y-%m-%d %H:%M:%S"` ':         resart '${prog} >> ${log_file}
            case ${prog} in
                ./shlkr)
                    #python3 run_device.py &
                    #nohup ./shlkr 1>/dev/null 2>&1 & 
                    ./shlkr & 
                    ;;
                *)
                    echo `date "+%Y-%m-%d %H:%M:%S"` ':         cannt resart '${prog} >> ${log_file}
                    ;;
            esac
        fi
    done
}

while [ true ]; do
    detect_deamons_debug $progs_debug

    sleep 15  
 
done

