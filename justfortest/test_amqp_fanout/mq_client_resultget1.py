#!/usr/bin/env python3
import pika
import time

def callback(ch, method, properties, body):
    print(" [x] Received %r" % body)
    ch.basic_ack(delivery_tag = method.delivery_tag)


if __name__ == '__main__':
    '''
    从rabbitMQ服务器消息队列中取消息
    注意：pika版本为1.x
    '''
    # 连接至rabbitMQ服务器
    credentials = pika.PlainCredentials('rossi', 'rossi1219')
    Parameter = pika.ConnectionParameters('10.39.251.182',5672,'/',credentials)
    connection = pika.BlockingConnection(Parameter)
    channel = connection.channel()

    # 声明/指定消息队列
    exchange_name = 'rossi_exchange_fanout'
    queue_name = "myqueue1_"+str(time.time())

    channel.exchange_declare(exchange=exchange_name, exchange_type='fanout')
    channel.queue_declare(queue=queue_name, auto_delete=True)#, durable=True)

    #topic_key = "qdh_mytest"
    channel.queue_bind(exchange='rossi_exchange_fanout', queue=queue_name)#, routing_key=topic_key)

    # 配置从rabbitMQ服务器消息队列中取消息的参数
    channel.basic_consume(queue=queue_name, on_message_callback=callback)

    # 永久持续运行
    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()
