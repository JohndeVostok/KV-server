import requests
import json

url = "http://127.0.0.1:12345"

args = {"Op": 0, "Key": "", "Value": ""}

while True:
    content = input()
    req = content.split(" ")
    if (req[0] == "put"):
        args["Op"] = 1
        args["Key"] = req[1]
        args["Value"] = req[2]
    if (req[0] == "get"):
        args["Op"] = 2
        args["Key"] = req[1]
    r = requests.post(url, data = json.dumps(args))
    print(r.text)
