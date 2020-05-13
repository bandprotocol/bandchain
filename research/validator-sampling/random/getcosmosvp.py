
import requests

url = "https://rpc.cosmos.network/validators"
r = requests.get(url)

jsonData = r.json()["result"]["validators"]
l = []
ll = []
for x in jsonData:
    l.append(int(x["voting_power"]))

for x in jsonData:
    ll.append(x["address"])
print(l)
print(ll)
print(len(l))
