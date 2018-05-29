import requests
import json

r = requests.post("http://127.0.0.1:12345", data = json.dumps({"op": 1, "key": "sb", "value" : "nb"}))

print(r.text)
