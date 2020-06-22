import serial
ser = serial.Serial('/dev/tty20')
print(ser.name)
ser.write(b'hello')
ser.close()
