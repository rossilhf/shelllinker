#!/usr/bin/python3
# -*- coding: utf-8 -*-

import time
import paho.mqtt.client as mqtt
import os
import json
import time
import subprocess
import getInfo
import readConfig
import encryption

GotMsgContext = ""

def on_connect(client, userdata, flags, rc):
    """
    subscribe cmd from server to this device
    """
    print("connected with "+str(rc))
    cur_os = getInfo.get_curOs()
    cur_mac = getInfo.get_curMac(cur_os)
    client.subscribe("topic_ser2dev/exec_cmd/"+cur_mac)


def on_message(client, userdata, msg):
    """
    callback func
    """
    global GotMsgContext
    GotMsgContext = str(msg.payload, encoding="utf-8")
    print("got: "+msg.topic+" "+GotMsgContext)


def process_infoReport(client, cur_toolaccount):
    """
    when start up, report this device's basic info
    """
    # get basic info
    cur_os = getInfo.get_curOs()
    cur_user = getInfo.get_curUser(cur_os)
    cur_path = getInfo.get_curPath()
    cur_mac = getInfo.get_curMac(cur_os)
    cur_ip = getInfo.get_curIp(cur_os)
    cur_sn = getInfo.get_curSn(cur_os)
    cur_cpu = getInfo.get_curCpu(cur_os)
    cur_version = getInfo.get_curVersion()

    # send out msg
    resultdict = {}
    resultdict.update({"msgType": "info_report"})
    resultdict.update({"curtoolaccount": cur_toolaccount})
    resultdict.update({"curmac": cur_mac})
    resultdict.update({"curcpu": cur_cpu})
    resultdict.update({"curos": cur_os})
    resultdict.update({"curip": cur_ip})
    resultdict.update({"heartbeattime": time.strftime('%Y-%m-%d-%H', time.localtime(time.time())) })
    resultdict.update({"curversion": cur_version})
    resultdict.update({"curuser": cur_user})
    resultdict.update({"cursn": cur_sn})
    resultMsg_json = json.dumps(resultdict)
    resultMsg_json_encrypt = encryption.encrypt(11, resultMsg_json)

    topic = "topic_dev2ser/dev_info/" + cur_mac
    client.publish(topic, payload=resultMsg_json_encrypt)

    print("[x] Sent msg out to topic: " + topic)
    print("sent out msg: " + resultMsg_json)
    print("sent out msg(encrypt): " + resultMsg_json_encrypt)


def process_execCmd(client, cur_toolaccount):
    """
    recv cmd and exec and return result
    """
    # get basic info
    cur_os = getInfo.get_curOs()
    cur_user = getInfo.get_curUser(cur_os)
    cur_path = getInfo.get_curPath()
    cur_mac = getInfo.get_curMac(cur_os)

    # read cmd
    reportMsg = ""
    msgContext_decrypt = encryption.decrypt(11, GotMsgContext)
    msgdict = json.loads(msgContext_decrypt)
    cmd = msgdict["cmd"]

    # excute cmd
    if cmd == "pwd":
        reportMsg = os.getcwd()
    elif cmd.find("cd ") == 0:
        leng = len(cmd)
        cmd = cmd[3:leng]
        #print(cmd)
        try:
            os.chdir(cmd)
        except:
            reportMsg = "Error: execute " + cmd + "err!"
    else:
        try:
            proc = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        except:
            reportMsg = "Error: subprocess.Popen " + cmd + "err!"
    
        for line in proc.stdout.readlines():
            reportMsg = reportMsg + line.decode("UTF-8")

        for line in proc.stderr.readlines():
            reportMsg = reportMsg + line.decode("UTF-8")

    # send out msg
    resultdict = {}
    resultdict.update({"msgType": "exec_return"})
    resultdict.update({"curuser": cur_user})
    resultdict.update({"curtoolaccount": cur_toolaccount})
    resultdict.update({"curmac": cur_mac})
    resultdict.update({"curpath": cur_path})
    resultdict.update({"report": reportMsg})
    resultMsg_json = json.dumps(resultdict)
    resultMsg_json_encrypt = encryption.encrypt(11, resultMsg_json)

    topic = "topic_dev2ser/exec_result/" + cur_mac
    client.publish(topic, payload=resultMsg_json_encrypt)

    print("[x] Sent msg out to topic: ", topic)
    print("sent out msg: " + resultMsg_json)
    print("sent out msg(encrypt): " + resultMsg_json_encrypt)


def process_main():
    """
    main process
    """
    global GotMsgContext
    mq_server, mq_user, mq_psw = readConfig.readMQinfo()
    cur_toolaccount = readConfig.readToolAccount()

    # client: report dev info
    mqttclient0 = mqtt.Client()
    mqttclient0.username_pw_set(mq_user, mq_psw)
    mqttclient0.connect(mq_server, 1883, 60)
    mqttclient0.loop_start()

    # client: 1.receive cmd 2.execute 3.return result
    mqttclient = mqtt.Client()
    mqttclient.on_connect = on_connect
    mqttclient.on_message = on_message
    mqttclient.username_pw_set(mq_user, mq_psw)
    mqttclient.connect(mq_server, 1883, 60)
    mqttclient.loop_start()

    lastReportTime = time.time() - 5000
    while True:
        # report dev info
        time.sleep(0.2)
        now = time.time()
        if now - lastReportTime > 60 * 60 * 1:      # 1 hours
            process_infoReport(mqttclient0, cur_toolaccount)
            lastReportTime = now

        # recv cmd and exec and return result
        if GotMsgContext != "":
            process_execCmd(mqttclient, cur_toolaccount)
            GotMsgContext = ""


if __name__ == "__main__":
    process_main()
