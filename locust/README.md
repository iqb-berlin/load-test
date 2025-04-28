# Preparation
- Copy your own version of `config.json_template` and `locust.conf_template` as `config.json` and `locust.conf` and configure them to your liking
- `config.json` handles arguments for the load scripts. `locust.conf` handle the locust framework settings

## Config and Resources
Configure test parameters in `config.json` (file name may differ; passed as command line argument)
- `hostname` is the hostname of the Testcenter. Specify with 'http://' or 'https://', even for localhost.
- `username` and `password` are the credentials for the Testcenter.
- `increment_user_id` makes the username differ for each simulated user. An incremented number is appended to the username.
    - If set to false all users use the same login, meaning that the user MUST be of type `run-hot-restart` in the Testtakers.xml
    - If set to true, the tested users must have the same name-prefix as given in `username` with a number appended, starting from 1.
      E.g. if `username` is `testuser`, the users must be named `testuser1`, `testuser2`, etc.
      Here the mode can be set to `run-hot-restart` or `run-hot-reload` in the Testtakers.xml.
- `booklet_name` is the Id of the booklet to be used. The booklet MUST NOT be assigned an attribute `code` in the Testtakers.xml.
  This booklet must exist in the workspace.
- `workspace` is the workspace to be used.
- `resource_dir` expects a directory, with 2 files present:
    - resources.txt and units.txt
    - These contain all the files that the program needs to download, each file on a new line. These files must exist in the workspace.
    - resources.txt and units.txt may be empty to skip loading the type of resource.
- `retries` and `timeout` are the number of retries and the timeout in seconds for each request before deemed a failure. The standard timeout is 5 minutes as in most modern browsers
- `file_service_mode` adjusts the resource URL, depending on weather the file service is active. Default to `true`
  starting with Testcenter v15.

# Running Locust
- Install all dependencies with `pip install -r requirements.txt`
- Run `locust -f <filename>`. The test will run in headless mode as default (can be configured in `locust.conf`)

# Testing Recipes
## Quick Testing
- If you want to quickly test a new user/spawn-rate/iteration setup, use the CLI arguments `--users (-u)`, `spawn-rate (-r)` to override the settings in `locust.conf` 

## Spike Testing
- When you want to declare a number of users that should be generated and run at once, set `users` and `spawn-rate` to the same number and run the test with the `--iteration <number>` or `-i <number>` flag with the same number as argument, e.g. `locust -f src/minimal -i 1000`

## Distribute the load generation
- When testing a lot of users, you may want to run several Locust-workers. Those run on the same machine. Start and put them in the background with `locust --worker <number>`. If you want to use every available core of your machine use `locust --worker -1`

