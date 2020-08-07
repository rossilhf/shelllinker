#!/usr/bin/env python3
import pika
import time
import json


if __name__ == '__main__':
    '''
    发布算法调用指令至rabbitMQ服务器
    '''
    # 建立连接
    credentials = pika.PlainCredentials('rossi', 'rossi1219')
    Parameter = pika.ConnectionParameters('10.39.251.182',5672,'/',credentials)
    connection = pika.BlockingConnection(Parameter)
    channel = connection.channel()

    # 声明/指定消息队列
    exchange_name = 'rossi_exchange_fanout'
    channel.exchange_declare(exchange=exchange_name, exchange_type='fanout')
    topic_key = "qdh_mytest"

    # 生成算法调用指令消息内容
    dict = {}
    dict.update({'time':str(time.time())})
    dict.update({'topic':topic_key})
    dict.update({'taskid':"sdsgfdsdf23345"})
    dict.update({'url':"http://img.mp.itc.cn/upload/20161210/7d24dd1648a94e6f9826f150d3f17f38_th.jpeg"})
    dict.update({'operate':"face_analysis"})
    commandmsg = json.dumps(dict)

    # 发布消息
    #channel.basic_publish(exchange=exchange_name, routing_key=topic_key, body=commandmsg)
    channel.basic_publish(exchange=exchange_name, routing_key="", body=commandmsg)
    print("msg: ", commandmsg)
    print(" [x] Sent msg out!")

    time.sleep(0.5)
    topic_key = "qdh_mytest_other"
    dict = {}
    dict.update({'time':str(time.time())})
    dict.update({'topic':topic_key})
    dict.update({'taskid':"sdsgfdsdf23345"})
    dict.update({'url':"http://img.mp.itc.cn/upload/20161210/7d24dd1648a94e6f9826f150d3f17f38_th.jpeg"})
    dict.update({'operate':"face_analysis"})
    commandmsg = json.dumps(dict)
    #channel.basic_publish(exchange = 'rossi_exchange_fanout', routing_key = topic_key, body = commandmsg)
    channel.basic_publish(exchange=exchange_name, routing_key="", body=commandmsg)
    print("msg: ", commandmsg)
    print(" [x] Sent msg out!")

    # 断开连接
    connection.close()
