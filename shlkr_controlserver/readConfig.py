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


def readUrlServer():
    """
    read update resource server ip-address
    """
    infolist = []
    with open("./configs/urlserver.dat", "r") as f:
        line = f.readline()
        urlserver = line.strip('\n')
        print("urlserver: ", urlserver)

    urlserver = encryption.decrypt(11, urlserver)

    return urlserver


if __name__ == "__main__":
    print("helo")
    #print(readProjPath())
    print(readMQinfo())
    print(readUrlServer())
