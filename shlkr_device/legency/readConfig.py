#!/usr/bin/python3
# -*- coding: UTF-8 -*-

import encryption

def readProjPath():
    """
    read proj path
    """
    with open("./configs/curpath.dat", "r") as f:
        line = f.readline()
        projpath = line.strip('\n')
        print("proj path: ", projpath)

    return projpath


def readMQinfo():
    """
    read rabbit mq info: ip, user, psw
    """
    infolist = []
    with open("./configs/mqserver.dat", "r") as f:
        for line in f.readlines():
            line = line.strip('\n')
            print("rabbitmq info: ", line)
            infolist.append(line)

    ip = infolist[0]
    user = infolist[1]
    psw = infolist[2]

    ip = encryption.decrypt(11, ip)
    user = encryption.decrypt(11, user)
    psw = encryption.decrypt(11, psw)

    return ip, user, psw


def readToolAccount():
    """
    read this remote tool account
    """
    infolist = []
    with open("./configs/toolaccount.dat", "r") as f:
        line = f.readline()
        toolaccount = line.strip('\n')
        print("tool account: ", toolaccount)

    toolaccount = encryption.decrypt(11, toolaccount)

    return toolaccount


if __name__ == "__main__":
    print("helo")
    #print(readProjPath())
    print(readMQinfo())
    print(readToolAccount())
