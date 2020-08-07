import paho.mqtt.client as mqtt
import time


if __name__ == "__main__":
    client = mqtt.Client()
    #client.connect("127.0.0.1", 1883, 60)
    client.username_pw_set("admin", "public")
    #client.connect("10.39.251.182", 1883, 60)
    client.connect("117.50.109.189", 1883, 60)

    for _ in range(5):
        client.publish("test", payload=str(time.time())+" hello, world!")
        time.sleep(5)

    client.publish("iot111", payload=str(time.time())+" hello, world!")
