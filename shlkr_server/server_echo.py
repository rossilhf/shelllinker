#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys
import os
import paho.mqtt.client as mqtt
import json
import time
import datetime
import sqlite3
import loger
import getInfo
import readConfig
import encryption

localpath = readConfig.readProjPath()
os.chdir(localpath)
log = loger.Log("remotetool-sever-echo")
GotMsgContext = ""

helpInfo = """
"NOTICE: there are 2 modes: local-cmd-mode, and remote-cmd-mode"
"=============================================================="

"now is in local-cmd-mode."
"=============================================================="

".help: show help info"
".list: list all devices"
".echo xxxx: into remote-cmd-mode, echo with device which mac address is xxxx"
".exit: finish remoete-cmd-mode, return to local-cmd-mode"
"=============================================================="

"NOTIC: high-level usage!!!:"
"step1:"
"in remote-cmd-mode, input:"
"sshpass -p xxxx ssh -R nnnn:localhost:22 uuuu@111.222.33.44 -tt &"
"xxxx is ssh password of publicIPserver"
"nnnn is port num of server, should >= 10000"
"uuuu is ssh user of publicIPserver"
"111.222.33.44 is publicIPserver ip"
"example: sshpass -p qding@123 ssh -R 10000:localhost:22 root@117.50.109.189 -tt &"
"now server process will be blocked."
"======================================"
"step2:"
"open a new console window."
"input 'sudo rm ~/.ssh/known_hosts'"
"input 'ssh -p nnnn devuuuu@localhost'"
"nnnn is port num, same with upper nnnn"
"devuuuu is ssh user of device"
"example: ssh -p 10000 aiot@localhost"
"now you can control device by ssh"
"====================================="
"step3:"
"input 'ps -aux | grep ssh'"
"input 'kill -9 xxxxxx', kill process about ssh"
"notice: killing process 'sshpass -p zzzz.....' is useless! Kill process 'ssh -R ....'"
"then you will return to remote-cmd-mode"
"""

privateCmdList = [".help", ".list", ".echo", ".exit"]


def process_sendCmd(client, msg, topic):
    """
    send cmd to device
    """
    # send out msg
    resultdict = {}
    resultdict.update({"cmd": msg})
    resultMsg_json = json.dumps(resultdict)
    resultMsg_json_encrypt = encryption.encrypt(11, resultMsg_json)

    client.publish(topic, payload=resultMsg_json_encrypt)

    return


def process_logIn():
    """
    log in this tool by tool-account and password
    """
    print("")
    print("you should log in system first.")
    print("if no siged tool-account, could use test-account: 'test', password: 'test'.")
    toolaccount = ""
    tmppwd = ""

    #log in 
    ifaccount = False
    while True:
        toolaccount = input("Enter tool-account: ")
        tmplist = toolaccount.split()
        if len(tmplist) == 0:
            continue
        else:
            dbpath = "./db/toolaccount.db"
            con = sqlite3.connect(dbpath)
            cur = con.cursor()
            cur.execute("SELECT * FROM toolaccount WHERE toolaccount='%s'" %(toolaccount))
            output = cur.fetchall()
            if len(output) == 0:
                print("This tool-account is not signed yet.")
                continue
            else:
                ifaccount = True
                tmppwd = output[0][1]
                tmppwd = encryption.decrypt(11, tmppwd)
            con.close()
        if ifaccount:
            break

    iflogin = False
    for _ in range(3):
        pwd = input("Enter tool-account-password: ")
        if pwd == tmppwd:
            iflogin = True
            print("logged in qdh-remote-tool.")
            break
        else:
            print("password err. try again.")
       
    if iflogin == False:
        print("wrong password. log in err.")
        while True:
            pass

    #check useful db
    dbpath = "./db/deviceList__"+toolaccount+".db"
    if os.path.isfile(dbpath) == False:
        print("This tool-account has no useful devices yet. Please exit.")
        while True:
            pass

    return toolaccount


def process_localExcute(toolaccount):
    """
    excuet private cmd locally
    .help: show help
    .list: list all devices
    .echo: link to device
    """
    mac = ""

    # input and local excute cmd
    print("user guide:")
    print("support basic cmd: .help, .list, .echo, .exit")
    while True:
        cmd = input("Enter cmd: ")
        cmdlist = cmd.split()
        if len(cmdlist) == 0:
            continue

        keycmd = cmdlist[0]
        if keycmd in privateCmdList:
            if keycmd == privateCmdList[0]: # .help
                print(helpInfo)

            if keycmd == privateCmdList[1]: # .list
                dbpath = "./db/deviceList__"+toolaccount+".db"
                con = sqlite3.connect(dbpath)
                cur = con.cursor()
                try:
                    cur.execute("SELECT * FROM deviceList")
                    output = cur.fetchall()
                    for i in range(len(output)):
                        context = str(i) + "\t"
                        for item in output[i]:
                            item = str(item).strip()
                            context += (item + "\t")
                        print(context)
                except:
                    log.error("select db err.")
                con.close()

            if keycmd == privateCmdList[2]: # .echo
                macstr = cmdlist[1]
                macstr = macstr.strip()
                macstr = macstr.strip('"')
                macstr = macstr.strip("'")
                mac = macstr
                log.info("link to " + mac)
                break

    return mac


def on_connect0(client, userdata, flags, rc):
    """
    callback func
    """
    log.info("connected with result code "+str(rc)+"\n")
    client.subscribe("topic_dev2ser/exec_result/+")


def on_message0(client, userdata, msg):
    """
    callback func
    """
    global GotMsgContext
    GotMsgContext = str(msg.payload, encoding="utf-8")
    print("got dev returned msg: "+msg.topic+" "+GotMsgContext+"\n")
    

def cmdCheck(cmd, curuser):
    """
    check if cmd is supported.
    if not, cmd = ""
    """
    cmd_reform = cmd.replace("~", "/home/"+curuser)

    if cmd.find("su") == 0:
        print("su or sudo not support yet.")
        cmd_reform = ""

    if cmd.find("scp") == 0:
        print("scp not support yet.")
        cmd_reform = ""

    if cmd.find("vi") == 0:
        print("vi or vim not support yet.")
        cmd_reform = ""

    if cmd.find("top") == 0:
        print("top not support yet.")
        cmd_reform = ""
    
    if cmd.find("jtop") == 0:
        print("jtop not support yet.")
        cmd_reform = ""
    
    if cmd.find("htop") == 0:
        print("htop not support yet.")
        cmd_reform = ""
    
    if cmd.find("tail -f") == 0:
        print("tail -f not support yet.")
        cmd_reform = ""

    if cmd == "ll":
        print("ll change to ls -l")
        cmd_reform = "ls -l"

    return cmd_reform


def process_echo2Device(topic, mqttclient_send):
    """
    echo with device
    note: .exit is to quit echo with device
    """ 
    global GotMsgContext
    process_sendCmd(mqttclient_send, "uname -a", topic)

    while True:
        if GotMsgContext != "":
            # show cmd exec result
            gotMsgContext_decrypt = encryption.decrypt(11, GotMsgContext)
            reportdict = json.loads(gotMsgContext_decrypt)
            curuser = reportdict["curuser"]
            curmac = reportdict["curmac"]
            curpath = reportdict["curpath"]
            report = reportdict["report"]
            print(report)

            # send cmd to device
            cmd = input("echo@" + curuser + "@" + curmac + "@" + curpath + "@: ")
            if cmd == ".exit":
                break
            cmd_reform = cmdCheck(cmd, curuser)
            process_sendCmd(mqttclient_send, cmd_reform, topic)
            GotMsgContext = ""


def process_main():
    '''
    1. get cmd from keyboard
    2. send cmd to device
    3. get cmd exec result context, print
    '''
    log.info("current version: "+getInfo.get_curVersion())
    mq_server, mq_user, mq_psw = readConfig.readMQinfo()

    mqttclient0 = mqtt.Client()
    mqttclient0.on_connect = on_connect0
    mqttclient0.on_message = on_message0
    mqttclient0.username_pw_set(mq_user, mq_psw)
    mqttclient0.connect(mq_server, 1883, 60)
    mqttclient0.loop_start()

    mqttclient = mqtt.Client()
    mqttclient.username_pw_set(mq_user, mq_psw)
    mqttclient.connect(mq_server, 1883, 60)
    mqttclient.loop_start()
	
    #log in this tool by tool-account and password
    toolaccount = process_logIn()

    # get cmd from keyboard
    while True:
        linkmac = process_localExcute(toolaccount)

        topic = "topic_ser2dev/exec_cmd/" + linkmac
        process_echo2Device(topic, mqttclient)


if __name__ == "__main__":
    process_main()
