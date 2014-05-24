from flask import Flask, request
import json
import os

app = Flask(__name__)

@app.route('/',methods=['POST'])
def foo():
   print 'triggerd'
   os.system('sh run.sh &')
   return "OK"

if __name__ == '__main__':
   app.run(host='0.0.0.0', port=7077)
