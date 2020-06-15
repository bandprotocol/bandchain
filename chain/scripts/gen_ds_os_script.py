import os
import json
import sys
import requests
import subprocess
from data_sources import data_source_info
from oracle_scripts import oracle_scripts_info


def gen_data_source_oracle_script(ds_info, os_info, owner):
    dir = os.path.abspath(__file__)
    data_source_dir = os.path.abspath(os.path.join(dir, '../../datasources'))
    os_source_dir = os.path.abspath(
        os.path.join(dir, '../../../owasm/chaintests'))

    file = open("add_os_ds.sh", "w")
    file.write('DIR=$(dirname "$0")\n')
    if len(ds_info) != len(os.listdir(data_source_dir)):
        raise Exception(
            "length data source info does not match with amount of data sources in directory")
    if len(os_info) != len(os.listdir(os_source_dir)):
        raise Exception(
            "length oracle script info does not match with amount of oracle scripts in directory")

    gen_data_source(file, data_source_dir, ds_info, owner)
    gen_oracle_script(file, os_source_dir, os_info, owner)

    file.close()


def gen_data_source(file, data_source_dir, ds_info, owner):
    data_sources = []
    for (name, des), filename in zip(ds_info, sorted(os.listdir(data_source_dir))):
        file_path = os.path.join('$DIR/../datasources', filename)
        data_sources.append({'name': name, 'description': des,
                             'path': file_path})

    write_add_data_source(file, data_sources)


def write_add_data_source(file, data_sources):
    for value in data_sources:
        name = value['name']
        description = value['description']
        path = value['path']
        command_line = f'bandd add-data-source \\\n\t"{name}" \\\n\t"{description}" \\\n\t{owner} \\\n\t{path}\n\n'
        file.write(command_line)


def gen_oracle_script(file, os_source_dir, os_info, owner):
    oracle_scripts = []
    for (name, des), filename in zip(os_info, sorted(os.listdir(os_source_dir))):
        file_path = os.path.join(os_source_dir, filename, 'src', 'lib.rs')
        oracle_scripts.append({
            'name': name, 'description': des, 'schema': get_schema(os.path.join(os_source_dir, filename)),
            'url': get_url(file_path, filename),
            'path': get_oracle_script_path(filename)})

    write_add_oracle_script(file, oracle_scripts)


def get_schema(path):
    pwd = os.getcwd()
    os.chdir(path)
    subprocess.check_output(
        'cargo test -- --nocapture | grep -v "^test" | grep -v "^running" | grep -v "^$" > schema.txt',
        stderr=subprocess.STDOUT,
        shell=True)

    with open("schema.txt", 'r') as file:
        schema = file.read().strip()
        os.remove("schema.txt")
        os.chdir(pwd)

        return schema


def get_url(file_path, filename):
    url = "https://api.pinata.cloud/pinning/pinFileToIPFS"
    headers = {
        'pinata_api_key': os.getenv('PINATA_API_KEY'),
        'pinata_secret_api_key': os.getenv('PINATA_API_SECRET')
    }
    payload = {'pinataMetadata': f'{{"name": "{filename}"}}',
               'pinataOptions': '{"cidVersion":0}'}
    files = [
        ('file', open(file_path, 'rb')),
    ]

    r = requests.request(
        "POST", url, headers=headers, data=payload, files=files)
    r.raise_for_status()
    return 'https://ipfs.io/ipfs/' + json.loads(r.content)['IpfsHash']


def get_oracle_script_path(filename):
    return os.path.join('$DIR/..', 'pkg/owasm/res', filename + '.wasm')


def write_add_oracle_script(file, oracle_scripts):
    for value in oracle_scripts:
        name = value['name']
        description = value['description']
        schema = value['schema']
        url = value['url']
        path = value['path']
        command_line = f'bandd add-oracle-script \\\n\t"{name}" \\\n\t"{description}" \\\n\t"{schema}" \\\n\t"{url}" \\\n\t{owner} \\\n\t{path}\n\n'
        file.write(command_line)


if __name__ == "__main__":
    owner = 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs'
    gen_data_source_oracle_script(data_source_info, oracle_scripts_info, owner)
