import paho.mqtt.client as mqtt
import time

"""def on_connect0(client, userdata, flags, rc):
    print("llllll")


def on_message0(client, userdata, msg):
    print("sssssss")"""


def on_connect(client, userdata, flags, rc):
    """
    callback func
    """
    print("connected with result code "+str(rc))
    client.subscribe("iot/+")


def on_message(client, userdata, msg):
    """
    callback func
    """
    print("got: "+msg.topic+" "+str(msg.payload))
    tmpmsg = str(msg.payload) #"ssssssssss"
    client.publish("iotreturn", payload=str(time.time())+"xxxxxx"+tmpmsg+"__0")


if __name__ == "__main__":
    mqttclient0 = mqtt.Client()
    #mqttclient0.on_connect = on_connect0
    #mqttclient0.on_message = on_message0
    mqttclient0.username_pw_set("admin", "public")
    mqttclient0.connect("117.50.109.189", 1883, 60)
    mqttclient0.loop_start()

    mqttclient = mqtt.Client()
    mqttclient.on_connect = on_connect
    mqttclient.on_message = on_message
    mqttclient.username_pw_set("admin", "public")
    mqttclient.connect("117.50.109.189", 1883, 60)
    mqttclient.loop_start()
    #mqttclient.loop_forever()

    while True:
        time.sleep(1.0)
        t = int(time.time())
        if t % 5 == 0:
            mqttclient0.publish("iotreport", payload=str(time.time())+"report"+"__0")

