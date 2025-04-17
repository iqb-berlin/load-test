import json
from tenacity import retry, wait_chain, wait_fixed
from jsonschema import validate
from locust import HttpUser, task, run_single_user
import login_util

# Module Definitions
config = {}
global_id = 1
wait_retries_frontend=[wait_fixed(0.05), wait_fixed(0.5), wait_fixed(1), wait_fixed(2)]

def load_config_file():
    global config # Declare config as global to modify it in this scope, for reading not necessary

    with open('config.json', 'r') as file:
        config = json.load(file)
        validate(config, login_util.config_schema)

@retry(wait=wait_chain(*wait_retries_frontend))
def get_test(i):
    login_util.get_test(i, i.headers, i.test_number)

# Module initialisation
load_config_file()  # call this function on module level for class to attribute "host" to access it at class definition time

class NavigationUser(HttpUser):
    host = config['hostname']
    timeout = config['timeout']
    retries = config['retries']

    def __init__(self, *args, **kwargs):
        global global_id
        self.id = global_id
        global_id += 1
        super().__init__(*args, **kwargs)

    def on_start(self) -> None:
        print('global id ', self.id)

        token, groupToken = login_util.get_tokens(self, config)
        self.headers = {'AuthToken': token}

        self.test_number = login_util.put_test(self, self.headers, config)

        get_test(self)

    @task(1)
    @retry(wait=wait_chain(*wait_retries_frontend))
    def move_unit(self):
        newUnitState = {
            'newState': [{
                'content': 'LOADING',
                'key': 'PLAYER',
                'timeStamp': 1744900611009
            }],
            'originalUnitId': 'UNIT.SAMPLE'
        }

        newTestState = [{
            'content': 'HAS_NOT',
            'key': 'FOCUS',
            'timeStamp': 1744900611009
        }]

        self.client.patch(f"/api/test/{self.test_number}/unit/UNIT.SAMPLE/state", headers=self.headers,
                               data=json.dumps(newUnitState), name='moveToUnit', timeout=self.timeout)

        self.client.patch(f"/api/test/{self.test_number}/state", headers=self.headers,
                               data=json.dumps(newTestState), name='moveToTest', timeout=self.timeout)

    @task(3)
    @retry(wait=wait_chain(*wait_retries_frontend))
    def put_response(self):
        newResponse = {
            'dataParts': {
                'answers': '[{"id":"required-text-field","status":"VALUE_CHANGED","value":"eei"},{"id":"readonly-text-field","status":"VALUE_CHANGED","value":"read-only"},{"id":"text-field","status":"VALUE_CHANGED","value":"ei"},{"id":"number-field","status":"VALUE_CHANGED","value":""},{"id":"date-field","status":"VALUE_CHANGED","value":""},{"id":"email-field","status":"VALUE_CHANGED","value":"ei"},{"id":"range-field","status":"VALUE_CHANGED","value":"6"},{"id":"radio-group","status":"VALUE_CHANGED","value":""},{"id":"check-box","status":"VALUE_CHANGED","value":""},{"id":"select-box","status":"VALUE_CHANGED","value":"a"},{"id":"list-field","status":"VALUE_CHANGED","value":"eiei"},{"id":"text-area","status":"VALUE_CHANGED","value":"Type something...eiei"},{"id":"multi-select","status":"VALUE_CHANGED","value":"blue"},{"id":"","status":"VALUE_CHANGED","value":["a","b","c","d"]},{"id":"editable-named-required","status":"VALUE_CHANGED","value":""}]'
            },
            'originalUnitId': 'UNIT.SAMPLE',
            'responseType': 'iqb-standard@1.3',
            'timeStamp': 1744901285997
        }

        self.client.put(f"/api/test/{self.test_number}/unit/an_alias/response", headers=self.headers,
                             data=json.dumps(newResponse), name='putResponse', timeout=self.timeout)


if __name__ == "__main__":
    run_single_user(NavigationUser)
