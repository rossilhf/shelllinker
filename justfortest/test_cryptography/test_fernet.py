from cryptography.fernet import Fernet
import time

cipher_key = Fernet.generate_key()
print("cipher_key: ", cipher_key)

cipher = Fernet(cipher_key)
text = b'{"name":"rossi", "class":"car", context:"hello, i am Charmy."}'

start = time.time()
encrypted_text = cipher.encrypt(text)
end = time.time()
print("encry time cost: ", end-start)
print("encrypted_text: ", encrypted_text)

start = time.time()
decrypted_text = cipher.decrypt(encrypted_text)
end = time.time()
print("decry time cost: ", end-start)
print("decrypted_text: ", decrypted_text)

