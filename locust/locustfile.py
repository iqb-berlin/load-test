import json
from tenacity import retry, wait_chain, wait_fixed
from jsonschema import validate
from locust import HttpUser, task, run_single_user
import login_util

config = {}
resource_list = []
unit_list = []
global_id = 1
wait_retries_frontend=[wait_fixed(0.05), wait_fixed(0.5), wait_fixed(1), wait_fixed(2)]

def load_config_file():
    global config # Declare config as global to modify it in this scope, for reading not necessary
    global resource_list
    global unit_list

    with open('config.json', 'r') as file:
        config = json.load(file)
        validate(config, login_util.config_schema)
    with open(config['resource_dir'] + 'resources.txt', 'r') as file:
        for line in file:
            resource_list.append(line.rstrip())
    with open(config['resource_dir'] + 'units.txt', 'r') as file:
        for line in file:
            unit_list.append(line.rstrip())


@retry(wait=wait_chain(*wait_retries_frontend))
def get_test(i, headers, test_number):
    login_util.get_test(i, headers, test_number)

@retry(wait=wait_chain(*wait_retries_frontend))
def get_units(i, headers, test_number):
    for file in unit_list:
        with i.client.get('/api/test/' + test_number + '/unit/' + file, headers=headers, name='Unit: ' + file,
                             timeout=i.timeout, catch_response=True) as response:
            if response.status_code >= 400:
                response.failure("Error of some kind")
            else:
                response.success()

@retry(wait=wait_chain(*wait_retries_frontend))
def get_resources(i, groupToken, headers, test_number):
    for file in resource_list:
        workspace = config['workspace']
        if config['file_service_mode'] is False:
            i.client.get('/api/test/' + test_number + '/resource/' + file, headers=headers, name=file,
                            timeout=i.timeout)
        else:
            with i.client.get(f"/fs/file/{groupToken}/{workspace}/Resource/{file}", headers=headers,
                                 name='Resource: ' + file, timeout=i.timeout, catch_response=True) as response:
                if response.status_code >= 400:
                    response.failure("Error of some kind")
                else:
                    response.success()


load_config_file()  # call this function on module level for class to attribute "host" to access it at class definition time
# timing of this code:
# -> module level code definitions and function calls
# -> class definition (w/ class attribute initial values are calculated here)
# -> locust @events and their hooks
# -> class instantiation (w/ instance attribute values are calculated here, when changed in __init__)

class QuickstartUser(HttpUser):
    host = config['hostname']
    timeout = config['timeout']
    retries = config['retries']

    def __init__(self, *args, **kwargs):
        global global_id
        self.id = global_id
        global_id += 1
        super().__init__(*args, **kwargs)

    @task
    def load_test(self):
        print('global id ', self.id)

        token, groupToken = login_util.get_tokens(self, config)
        headers = {'AuthToken': token}

        test_number = login_util.put_test(self, headers, config)

        get_test(self, headers, test_number)

        get_units(self, headers, test_number)

        get_resources(self, groupToken, headers, test_number)

if __name__ == "__main__":
    run_single_user(QuickstartUser)
