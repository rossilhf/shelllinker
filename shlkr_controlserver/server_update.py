#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys
import os
import platform
import paho.mqtt.client as mqtt
#import pika
import json
import subprocess
import time
import datetime
import sqlite3
from dateutil.parser import parse
import loger
import getInfo
import readConfig
import encryption

localpath = readConfig.readProjPath()
os.chdir(localpath)
log = loger.Log("remotetool-server-update")


def isAlive(lasttime):
    """
    input: string, form like "2020-09-02 16:39:57"
    judge if the input time is active 
    """
    nowtime = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())

    start = parse(lasttime)
    end = parse(nowtime)
    timediff = (end - start).total_seconds()
    #print(timediff)

    isalive = 0
    if timediff < 60 * 60 * 1.5:
        isalive = 1
    else:
        isalive = 0

    return isalive


def process_updateCmd(mqttclient, urlserver, curversion):
    """
    find every active device and it's version
    send update cmd 
    """
    file_names = os.listdir("./db/")    
    file_list = [os.path.join("./db/", file) for file in file_names]

    # find all device db (different db belongs to different account)
    for filepath in file_list:
        if "deviceList__" in filepath:
            #print(filepath)

            # check all devices 
            con = sqlite3.connect(filepath)
            cur = con.cursor()
            cur.execute("SELECT * FROM deviceList")
            output = cur.fetchall()
            #print(output)
            for node in output:
                node_mac = node[0]
                node_arch = node[1]
                node_os = node[2]
                node_time = node[4]
                node_ver = node[6]
                if isAlive(node_time) and node_ver != curversion:
                    topic = "topic_ser2dev/exec_cmd/" + node_mac
                    binfile = "device_" + node_os + "_" + node_arch + "_" + node_ver
                    content = os.path.join(urlserver, binfile)

                    # send out msg
                    resultdict = {}
                    resultdict.update({"type": "update"}) # type update: this cmd tell device to update version
                    resultdict.update({"cmd": content})
                    resultMsg_json = json.dumps(resultdict)
                    print(resultMsg_json)
                    resultMsg_json_encrypt = encryption.encrypt(11, resultMsg_json)
                    mqttclient.publish(topic, payload=resultMsg_json_encrypt)

                    time.sleep(10) # avoid all device download update file at same time

            con.close()

    return


def process_main():
    '''
    1. get newest version num 
    2. check all alive devices which version is old
    3. send out update cmd
    '''
    curversion = getInfo.get_curVersion()
    urlserver = readConfig.readUrlServer()
    log.info("current version: "+curversion)
    log.info("update url server: "+urlserver)
    log.info("start ...")
    mq_server, mq_user, mq_psw = readConfig.readMQinfo()

    mqttclient = mqtt.Client()
    mqttclient.username_pw_set(mq_user, mq_psw)
    mqttclient.connect(mq_server, 1883, 60)
    mqttclient.loop_start()

    while True:
        process_updateCmd(mqttclient, urlserver, curversion)
        time.sleep(60 * 60)


if __name__ == "__main__":
    process_main()
    #process_updateCmd()
