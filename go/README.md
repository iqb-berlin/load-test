### General

### Preparation

#### Required software
- go (https://go.dev/doc/install)

#### Config and Resources
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
  - `retries` and `timeout` are the number of retries and the timeout in seconds for each request before deemed a failure.
  - `file_service_mode` adjusts the resource URL, depending on weather the file service is active. Default to `true` 
    starting with Testcenter v15.

### Usage
- build the executable with `go build load-test.go`
- execute the program with `./load-test 10 config-fs-large.json`
  - first parameter is the number of simulated users
  - second parameter is the config file to be used

### Results
The script `summary.sh` is automatically called when all simulated users are finished. It parses the logfile and handily shows
several values.
Log files can be accessed in the `results` directory. Files are named by Unix-timestamp.

### Troubleshooting

#### ulimit
When running lots of users, the script may run into the OS-limit of open file descriptors.
You may need to raise the allowed value before running the script. You can do that with `ulimit`. For example:
```
ulimit -n 10000
```
For larger values you may have to change the allowed value for your user in `/etc/security/limits.conf`
