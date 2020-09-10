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
import loger
import getInfo
import readConfig
import encryption

localpath = readConfig.readProjPath()
os.chdir(localpath)
log = loger.Log("remotetool-server-db")
GotMsgContext = ""


def saveInfo(curtoolaccount, curmac, curcpu, curuser, curip, curos, curversion, heartbeattime):
    """
    save device info to db
    """
    dbpath = "./db/deviceList__"+curtoolaccount+".db"

    # create db file
    if not os.path.isfile(dbpath):
        con = sqlite3.connect(dbpath)
        cur = con.cursor()
        try:
            cur.execute('CREATE TABLE deviceList(mac TEXT PRIMARY KEY, cpu TEXT, os TEXT, ip TEXT, heartbeattime TEXT, user TEXT, version TEXT)')
            con.commit()
        except:
            log.warning("create db file failed.")
        con.close()

    # check if this device is signed
    con = sqlite3.connect(dbpath)
    cur = con.cursor()
    cur.execute("SELECT * FROM deviceList WHERE mac='%s'" %(curmac))
    output = cur.fetchall()
    if len(output) == 0:
        # add new device info
        try:
            cur.execute("INSERT INTO deviceList(mac, cpu, os, ip, heartbeattime, user, version) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')" % (curmac, curcpu, curos, curip, heartbeattime, curuser, curversion))
            con.commit()
        except:
            log.warning("insert to db failed.")
    else:
        # update old device info
        try:
            cur.execute("UPDATE deviceList SET cpu='%s', os='%s', ip='%s', heartbeattime='%s', user='%s', version='%s' WHERE mac='%s'" %(curcpu, curos, curip, heartbeattime, curuser, curversion, curmac))
            con.commit()
        except:
            log.warning("update db failed.")
    con.close()

    return


def on_connect0(client, userdata, flags, rc):
    """
    callback func
    """
    log.info("connected with result code "+str(rc)+"\n")
    client.subscribe("topic_dev2ser/dev_info/+")


def on_message0(client, userdata, msg):
    """
    callback func
    """
    global GotMsgContext
    GotMsgContext = str(msg.payload, encoding="utf-8")
    print("got a device info report: "+msg.topic+" "+GotMsgContext+"\n")


def process_main():
    '''
    1. get device info report
    2. save to db
    '''
    global GotMsgContext
    log.info("current version: "+getInfo.get_curVersion())
    log.info("start ...")
    mq_server, mq_user, mq_psw = readConfig.readMQinfo()

    mqttclient0 = mqtt.Client()
    mqttclient0.on_connect = on_connect0
    mqttclient0.on_message = on_message0
    mqttclient0.username_pw_set(mq_user, mq_psw)
    mqttclient0.connect(mq_server, 1883, 60)
    #mqttclient0.subscribe("topic_dev2ser/dev_info/+")
    mqttclient0.loop_start()

    while True:
        if GotMsgContext != "":
            gotMsgContext_decrypt = encryption.decrypt(11, GotMsgContext)
            print("got a device info report(decrypt): "+gotMsgContext_decrypt+"\n")
            reportdict = json.loads(gotMsgContext_decrypt)

            if "curtoolaccount" in reportdict:
                curtoolaccount = reportdict["curtoolaccount"]
            else:
                curtoolaccount = "unknown"

            if "curuser" in reportdict:
                curuser = reportdict["curuser"]
            else:
                curuser = "unknown"

            if "curmac" in reportdict:
                curmac = reportdict["curmac"]
            else:
                curmac = "unknown"

            if "curip" in reportdict:
                curip = reportdict["curip"]
            else:
                curip = "unknown"

            if "curcpu" in reportdict:
                curcpu = reportdict["curcpu"]
            else:
                curcpu = "unknown"

            if "curos" in reportdict:
                curos= reportdict["curos"]
            else:
                curos= "unknown"

            if "curversion" in reportdict:
                curversion= reportdict["curversion"]
            else:
                curversion= "unknown"

            if "heartbeattime" in reportdict:
                heartbeattime = reportdict["heartbeattime"]
            else:
                heartbeattime = "unknown"

            saveInfo(curtoolaccount, curmac, curcpu, curuser, curip, curos, curversion, heartbeattime)
            GotMsgContext = ""


if __name__ == "__main__":
    process_main()
