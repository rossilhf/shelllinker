import paho.mqtt.client as mqtt
import time

def on_connect(client, userdata, flags, rc):
    """
    callback func
    """
    print("connected with result code "+str(rc))
    client.subscribe("iot")


def on_message(client, userdata, msg):
    """
    callback func
    """

    print("got: "+msg.topic+" "+str(msg.payload))
    tmpmsg = str(msg.payload) #"ssssssssss"
    client.publish("iotreturn", payload=str(time.time())+"xxxxxx"+tmpmsg+"__1")


if __name__ == "__main__":
    mqttclient = mqtt.Client()
    mqttclient.on_connect = on_connect
    mqttclient.on_message = on_message
    mqttclient.username_pw_set("admin", "public")
    mqttclient.connect("117.50.109.189", 1883, 60)
    #mqttclient.loop_start()
    mqttclient.loop_forever()

