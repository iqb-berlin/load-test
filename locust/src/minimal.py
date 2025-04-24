import json
from jsonschema import validate
from locust import FastHttpUser, task, run_single_user
import  login_util

config = {}

def load_config_file():
    global config

    with open('config.json', 'r') as file:
        config = json.load(file)
        validate(config, login_util.config_schema)

class MinimalUser(FastHttpUser):
    host = 'http://localhost'

    @task
    def test_fun(self):
        self.client.get('/api/version')

if __name__ == "__main__":
    run_single_user(MinimalUser)