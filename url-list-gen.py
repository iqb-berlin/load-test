#!/usr/bin/env python3

# Helper script to generate list of URLs for the selected Booklet.
# This list can then be used to feed load-test tools.

import logging
from typing import List
import requests
import xml.etree.ElementTree as ET


SERVERNAME = 'http://172.28.33.165'
USERNAME = 'super'
PASSWORD = 'user123'
WORKSPACE_ID = '3'
# BOOKLET_FILE_NAME = 'SAMPLE_BOOKLET.XML'
BOOKLET_FILE_NAME = 'booklet1.xml'
TARGET_FILE_NAME = 'file_list.txt'


def get_token() -> str:
    response = requests.put(SERVERNAME + '/api/session/admin',
                            data=f'{{"name": "{USERNAME}", "password": "{PASSWORD}"}}')
    return response.json()['token']


def get_booklet(token: str) -> str:
    headers = {'AuthToken': token}
    response = requests.get(SERVERNAME + f'/api/workspace/{WORKSPACE_ID}/file/Booklet/{BOOKLET_FILE_NAME}',
                            headers=headers)
    logging.debug('get_booklet response: ' + response.text)
    return response.text


def get_unit_ids(booklet_xml: str) -> List[str]:
    try:
        tree = ET.fromstring(booklet_xml)
        unit_elements = tree.findall('.//Unit')
        return [unit.attrib['id'] for unit in unit_elements]
    except ET.ParseError:
        print('faulty XML: ' + booklet_xml)


def get_unit_metadata_list(token: str) -> List:
    headers = {'AuthToken': token}
    response = requests.get(SERVERNAME + f'/api/workspace/{WORKSPACE_ID}/files',
                            headers=headers)
    return response.json()['Unit']


def get_unit_file_name(unit_id: str, unit_metadata: List) -> str:
    print('Looking for ' + unit_id + ' in ' + str(unit_metadata))
    print([unit_meta['name'] for unit_meta in unit_metadata if unit_meta['id'] == unit_id.upper()])
    return [unit_meta['name'] for unit_meta in unit_metadata if unit_meta['id'] == unit_id.upper()][0]


def load_unit_file(filename: str, token: str) -> str:
    headers = {'AuthToken': token}
    response = requests.get(SERVERNAME + f'/api/workspace/{WORKSPACE_ID}/file/Unit/{filename}',
                            headers=headers)
    return response.text


def write_to_file(line: str):
    f = open(TARGET_FILE_NAME, 'a')
    f.write("'" + line + "'\n")
    f.close()


# logging.basicConfig(level=logging.DEBUG)
logging.info('geht los')

player_list = []

token = get_token()
booklet = get_booklet(token)
unit_ids = get_unit_ids(booklet)
for unit_id in unit_ids:
    write_to_file(f'/unit/{unit_id}/alias/{unit_id}')

unit_metadata = get_unit_metadata_list(token)
for unit_id in unit_ids:
    unit_file_name = get_unit_file_name(unit_id, unit_metadata)
    unit_xml = load_unit_file(unit_file_name, token)
    tree = ET.fromstring(unit_xml)
    def_ref = tree.findall('.//DefinitionRef')
    if def_ref:
        write_to_file(f'/resource/{def_ref[0].text}?v=f')
        player_list.append(def_ref[0].attrib['player'])
    definition = tree.findall('.//Definition')
    if definition:
        player_list.append(definition[0].attrib['player'])


player_list = list(set(player_list))
for player in player_list:
    write_to_file(f'/resource/{player.upper()}.html?v=1')
