#!/usr/bin/python3
# -*- coding: utf-8 -*-

import os
import platform
import subprocess

def get_curVersion():
    """
    get program version
    """
    return "V20200806.1"


def get_curOs():
    """
    get os info
    """
    osInfo = "unknown"

    if platform.system().lower() == "windows":
        osInfo = "windows"

    if platform.system().lower() == "linux":
        osInfo = "linux"

    return osInfo


def get_curCpu(cur_os):
    """
    get cpu info
    """
    curCpu = "unknown"

    if cur_os == "linux":
        cmd = "lscpu"
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        for line in p.stdout.readlines():
            context = line.decode("UTF-8").strip()
            tmplist = context.split("    ")
            break

        if len(tmplist) > 2:
            curCpu = tmplist[2].strip()
        elif len(tmplist) > 1:
            curCpu = tmplist[1].strip()

    if cur_os == "windows":
        cmd = "wmic cpu list brief"
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        i = 0
        for line in p.stdout.readlines():
            i += 1
            #print(line)
            if i == 2:
                context = line.decode("UTF-8").strip()
                tmplist = context.split(" ")
                curCpu = tmplist[0].strip()


    return curCpu
    

def get_curMac(cur_os):
    """
    get current mac address
    """
    macAddr = "unknown"
    if cur_os == "linux":
        curip = get_curIp(cur_os)
        cmd = "ifconfig"
        net_card = "eth0"

        # find net-card
        startlineidx = 0
        endlineidx = 0
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        i = 0
        for line in p.stdout.readlines():
            context = line.decode("UTF-8").strip()
            if context == "":
                startlineidx = i
            if curip in context:
                endlineidx = i
                #print(context)
                break
            i += 1
        #print("startlineidx", startlineidx)
        #print("endlineidx", endlineidx)

        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        if startlineidx == 0:
            for line in p.stdout.readlines():
                context = line.decode("UTF-8").strip()
                tmplist = context.split(":")
                #print(tmplist)
                net_card = tmplist[0]
                break
            #print("net_card", net_card)
        else:
            i = 0
            for line in p.stdout.readlines():
                if (i-1) == startlineidx:
                    context = line.decode("UTF-8").strip()
                    tmplist = context.split(":")
                    #print(tmplist)
                    net_card = tmplist[0]
                i += 1
            #print("net_card", net_card)

        # find mac address
        macAddr = "unknown"
        cmd = "cat /sys/class/net/" + net_card + "/address"
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        for line in p.stdout.readlines():
            macAddr = line.decode("UTF-8").strip()

    if cur_os == "windows":
        cmd = "getmac"
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        for line in p.stdout.readlines():
            context = line.decode("gbk").strip()
            # print(context)
            if "Tcpip_{" in context:
                tmplist = context.split(" ")
                macAddr = tmplist[0].strip()

    return macAddr


def get_curUser(cur_os):
    """
    get current user
    """
    #if cur_os == "linux":
    line_list = []

    #cmd = "w -hs"
    cmd = "whoami"
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
        #print(line)
        line_list = line.decode("UTF-8").split()
        #print(line_list)
        break

    if len(line_list) >= 1:
        return line_list[0]
    else:
        return "unknown"

    #if cur_os == "windows":
    #    return "unknown"

    return "unknown"
   

def get_curIp(cur_os):
    """
    get this device's ip, no matter is local ip or not
    """
    curip = "unknown"
    iplist = []
    boardcastlist = []

    if cur_os == "linux":
        # find all ip
        cmd = "ifconfig"
        ignoreflag = False
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        for line in p.stdout.readlines():
            context = line.decode("UTF-8").strip()
            #print(context)
            if "docker" in context:
                ignoreflag = True

            if context == "":
                ignoreflag = False

            if not ignoreflag:
                if "inet" in context and "netmask" in context:
                    tmplist = context.split("  ")
                    if len(tmplist) == 3:
                        ipcontext = tmplist[0].strip()
                        boardcastcontext = tmplist[2].strip()
                        iptmplist = ipcontext.split(" ")
                        boardcasttmplist = boardcastcontext.split(" ")
                        if len(iptmplist) > 1:
                            ip = iptmplist[1]
                            boardcast = boardcasttmplist[1]
                            iplist.append(ip)
                            boardcastlist.append(boardcast)
        # exclude 127.0.0.1 and 0.0.0.0 and localhost
        for i in range(len(iplist)):
            ip = iplist[i]
            boardcast = boardcastlist[i]
            if (ip != "127.0.0.1") and (ip != "0.0.0.0") and (ip != "localhost") and (boardcast != "0.0.0.0"):
                curip = ip
                break

    if cur_os == "windows":
        # find all ip
        cmd = "ipconfig"
        ignoreflag = False
        p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        for line in p.stdout.readlines():
            #print(line)
            #print(type(line))
            context = line.decode("gbk").strip()
            # print(context)
            if "WLAN" in context:
                ignoreflag = False

            if not ignoreflag:
                if "IPv4 " in context:
                    tmplist = context.split(": ")
                    ignoreflag = True
                    if len(tmplist) == 2:
                        ipcontext = tmplist[1].strip()
                        iplist.append(ipcontext)

        # exclude 127.0.0.1 and 0.0.0.0 and localhost
        for i in range(len(iplist)):
            ip = iplist[i]
            if (ip != "127.0.0.1") and (ip != "0.0.0.0") and (ip != "localhost"):
                curip = ip
                break
    
    return curip


def get_curPath():
    """
    get current work path
    """
    return os.getcwd()


def get_curSn(cur_os):
    """
    get current product serial number
    """
    sn = "unknown"

    if cur_os == "linux":
        path = "/etc/qd/.devSerialNo"
        sn = "unknown"
        if os.path.isfile(path):
            cmd = "cat /etc/qd/.devSerialNo"
            p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
            for line in p.stdout.readlines():
                sn = line.decode("UTF-8").strip()
                break

    if cur_os == "windows":
        pass

    return sn


if __name__ == "__main__":
    print(get_curPath())
    cur_os = get_curOs()
    print(get_curIp(cur_os))
