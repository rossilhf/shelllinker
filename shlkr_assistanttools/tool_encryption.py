#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import time

def encrypt(key, s):
    b = bytearray(str(s).encode("utf-8"))
    n = len(b)
    c = bytearray(n*2)
    j = 0
    for i in range(0, n):
        b1 = b[i]
        b2 = b1 ^ key
        c1 = b2 % 19
        c2 = b2 // 19
        c1 = c1 + 46
        c2 = c2 + 46
        c[j] = c1
        c[j + 1] = c2
        j = j + 2

    return c.decode("utf-8")


def decrypt(ksa, s):
    c = bytearray(str(s).encode("utf-8"))
    n = len(c)
    if n % 2 != 0:
        return ""

    n = n // 2
    b = bytearray(n)
    j = 0
    for i in range(0, n):
        c1 = c[j]
        c2 = c[j + 1]
        j = j + 2
        c1 = c1 - 46
        c2 = c2 - 46
        b2 = c2 * 19 + c1
        b1 = b2 ^ ksa
        b[i] = b1

    return b.decode("utf-8")


if __name__ == "__main__":
    #text = '{"name":"rossi", "class":"car", context:"hello, i am Charmy."}'
    text = "test"#"admin"#"117.50.109.189"#"rossi_lhf"
    start = time.time()
    s = encrypt(11, text)
    end = time.time()
    print("encry time cost: ", end-start)
    print("encrypted string: ", s)

    start = time.time()
    ds = decrypt(11, s)
    end = time.time()
    print("decry time cost: ", end-start)
    print("decrypted string: ", ds)
