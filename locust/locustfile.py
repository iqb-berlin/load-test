import json
import yaml
from locust import HttpUser, task, constant, events

config = {}


@events.init.add_listener
def on_locust_init(environment, **_kwargs):
    load_config_file()


def load_config_file():
    global config
    with open('load-test-config.yml', 'r') as file:
        config = yaml.safe_load(file)


class QuickstartUser(HttpUser):
    wait_time = constant(10)

    @task
    def load_test(self):
        try:
            token, groupToken = self.get_tokens()
        except ValueError:
            return
        headers = {'AuthToken': token}

        with self.client.put('/api/test', data=json.dumps({'bookletName': config['bookletName']}),
                             headers=headers, catch_response=True) as response:
            if response.status_code == 504:
                response.failure("Timeout on Test start (504)")
                return
            else:
                test_number = response.content.decode("utf-8")

        self.client.get("/api/test/" + test_number, headers=headers)

        for file in config['UNIT_URLS']:
            self.client.get('/api/test/' + test_number + '/unit/' + file, headers=headers, name='Unit: ' + file)

        for file in config['RESOURCE_FILES']:
#             self.client.get('/api/test/' + test_number + '/resource/' + file, headers=headers, name=file)
            self.client.get('/fs/file/' + groupToken + '/ws_1/Resource/' + file, headers=headers, name='Resource: ' + file)

    def get_tokens(self) -> tuple[str, str]:
        auth = json.dumps(config['auth'])
        with self.client.put('/api/session/login', data=auth, catch_response=True) as response:
            if response.status_code == 504:
                response.failure("Timeout: Getting token failed! (504)")
                raise ValueError
            else:
                return response.json()['token'], response.json()['groupToken']
