import json

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

def get_tokens(self, config) -> tuple[str, str]:
    username = config['username']
    if config['increment_user_id'] is True:
        username += str(self.id)

    data = json.dumps({'name': username, 'password': config['password']})

    with self.client.put('/api/session/login', data=data, name='putSessionLogin', timeout=self.timeout,
                         catch_response=True) as response:
        if response.status_code >= 400:
            response.failure("Some error occurred")
            raise Exception("das sollte nicht kommen")
        else:
            response.success()

    return response.json()['token'], response.json()['groupToken']

def put_test(self, headers, config):
    with self.client.put('/api/test', data=json.dumps({'bookletName': config['booklet_name']}),
                         headers=headers, timeout=self.timeout, catch_response=True) as response:
        if response.status_code >= 400:
            response.failure("Error of some kind")
        else:
            test_number = response.content.decode("utf-8")
            response.success()

    return test_number

def get_test(self, headers, test_number):
    with (self.client.get("/api/test/" + test_number, headers=headers, timeout=self.timeout, catch_response=True)
          as response):
        if response.status_code >= 400:
            response.failure("Error of some kind")
        else:
            response.success()


