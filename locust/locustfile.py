import json
from jsonschema import validate
from locust import HttpUser, task, run_single_user

config_schema = {
    "type": "object",
    "properties": {
        "hostname": {"type": "string"},
        "username": {"type": "string"},
        "password": {"type": "string"},
        "increment_user_id": {"type": "boolean"},
        "booklet_name": {"type": "string"},
        "workspace": {"type": "string"},
        "resource_dir": {"type": "string"},
        "retries": {"type": "number"},
        "timeout": {"type": "number"},
        "file_service_mode": {"type": "boolean"}
    }
}

config = {}
resource_list = []
unit_list = []
global_id = 1


def load_config_file():
    global config # Declare config as global to modify it in this scope, for reading not necessary
    global resource_list
    global unit_list

    with open('config.json', 'r') as file:
        config = json.load(file)
        validate(config, config_schema)
    with open(config['resource_dir'] + 'resources.txt', 'r') as file:
        for line in file:
            resource_list.append(line.rstrip())
    with open(config['resource_dir'] + 'units.txt', 'r') as file:
        for line in file:
            unit_list.append(line.rstrip())


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

        try:
            token, groupToken = self.get_tokens()
        except ValueError:
            return
        headers = {'AuthToken': token}

        with self.client.put('/api/test', data=json.dumps({'bookletName': config['booklet_name']}),
                             headers=headers, timeout=self.timeout, catch_response=True) as response:
            if response.status_code >= 400:
                response.failure("Error of some kind")
            else:
                test_number = response.content.decode("utf-8")
                response.success()

        with (self.client.get("/api/test/" + test_number, headers=headers, timeout=self.timeout, catch_response=True)
              as response):
            if response.status_code >= 400:
                response.failure("Error of some kind")
            else:
                response.success()

        for file in unit_list:
            with self.client.get('/api/test/' + test_number + '/unit/' + file, headers=headers, name='Unit: ' + file,
                                 timeout=self.timeout, catch_response=True) as response:
                if response.status_code >= 400:
                    response.failure("Error of some kind")
                else:
                    response.success()

        for file in resource_list:
            workspace = config['workspace']
            if config['file_service_mode'] is False:
                self.client.get('/api/test/' + test_number + '/resource/' + file, headers=headers, name=file,
                                timeout=self.timeout)
            else:
                with self.client.get(f"/fs/file/{groupToken}/{workspace}/Resource/{file}", headers=headers,
                                     name='Resource: ' + file, timeout=self.timeout, catch_response=True) as response:
                    if response.status_code >= 400:
                        response.failure("Error of some kind")
                    else:
                        response.success()

    def get_tokens(self) -> tuple[str, str]:
        username = config['username']
        if config['increment_user_id'] is True:
            username += str(self.id)

        data = json.dumps({'name': username, 'password': config['password']})

        with self.client.put('/api/session/login', data=data, name='putSessionLogin', timeout=self.timeout,
                             catch_response=True, max_retries=2) as response:
            if response.status_code >= 400:
                response.failure("Some error occurred")
            else:
                response.success()

        return response.json()['token'], response.json()['groupToken']


if __name__ == "__main__":
    run_single_user(QuickstartUser)
