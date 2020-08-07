import paho.mqtt.client as mqtt

def on_connect(client, userdata, flags, rc):
    print("connected with result code "+str(rc))
    client.subscribe("iot")


def on_message(client, userdata, msg):
    print(msg.topic+" "+str(msg.payload))


if __name__ == "__main__":
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_message = on_message
    #client.connect("127.0.0.1", 1883, 60)
    client.username_pw_set("admin", "public")
    #client.connect("10.39.251.182", 1883, 60)
    client.connect("117.50.109.189", 1883, 60)
    client.loop_forever()
