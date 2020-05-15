import requests

url = "https://data.kava.io/staking/validators"
r = requests.get(url)

jsonData = r.json()["result"]
l = []
ll = []
for x in jsonData:
    l.append(int(x["tokens"]))

for x in jsonData:
    ll.append(x["operator_address"])
print(l)
print(ll)
print(len(l))
