import os
import json
import sys
import requests
import subprocess

url = "https://api.pinata.cloud/pinning/pinFileToIPFS"
dir = os.path.abspath(__file__)
data_source_dir = os.path.abspath(os.path.join(dir, '../../datasources'))
data_sources = []
owner = 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs'

f = open("add_os_ds.sh", "w")

for filename in sorted(os.listdir(data_source_dir)):
    file_path = os.path.join(data_source_dir, filename)
    data_sources.append({'name': filename, 'description': filename,
                         'path': os.path.join(data_source_dir, filename)})

for value in data_sources:
    name = value['name']
    description = value['description']
    path = value['path']
    command_line = f'bandd add-data-source "{name}" \ \n\t"{description}" \ \n\t{owner} \ \n\t{path}\n\n'
    # print(command_line)
    f.write(command_line)


os_source_dir = os.path.abspath(
    os.path.join(dir, '../../../owasm/chaintests'))

oracle_scripts = []
for filename in sorted(os.listdir(os_source_dir)):
    file_path = os.path.join(os_source_dir, filename, 'src', 'lib.rs')
    payload = {'pinataMetadata': f'{{"name": "{filename}"}}',
               'pinataOptions': '{"cidVersion":0}'}
    files = [
        ('file', open(file_path, 'rb')),
    ]
    headers = {
        'pinata_api_key': '5f7169a396725c53075b',
        'pinata_secret_api_key': '2fdca43a889df5602fd79dcc65a310de8a6ec85ec0aca6ab55a78aa1d45e7ce7'
    }

    r = requests.request(
        "POST", url, headers=headers, data=payload, files=files)

    # os.system(f"cd {os.path.join(os_source_dir, filename)}")
    os.chdir(os.path.join(os_source_dir, filename))
    print (os.getcwd())

    x = subprocess.check_output(
        'cargo test -- --nocapture | grep -v "^test" | grep -v "^running" | grep -v "^$" > out.txt',
        stderr=subprocess.STDOUT,
        shell=True)

    schema = ""
    ff = open("out.txt", "r")
    schema = ff.read().strip()
    ff.close()
    print(schema.strip())

    oracle_scripts.append({
        'name': filename, 'description': filename, 'schema': schema,
        'url': 'https://ipfs.io/ipfs/' + json.loads(r.content)['IpfsHash'],
        'path': os.path.join(os.path.abspath(os.path.join(dir, '../../')), 'pkg/owasm/res', filename + '.wasm')})

for value in oracle_scripts:
    name = value['name']
    description = value['description']
    schema = value['schema']
    url = value['url']
    path = value['path']
    command_line = f'bandd add-oracle-script "{name}" \ \n\t"{description}" \ \n\t"{schema}" \ \n\t"{url}" \ \n\t{owner} \ \n\t{path}\n\n'
    # print(command_line)
    f.write(command_line)

f.close()
