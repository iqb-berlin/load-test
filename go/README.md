### General

### Preparation

#### Required software
- go (https://go.dev/doc/install)

#### Config and Resources
- Configure test parameters in `config.json` (file name my differ; passed as command line argument)
  - `increment_user_id` makes the username differ for each simulated user. An incremented number is appended to the username.
    If set to false all users use the same login. Make sure the test-takers configuration allows this.
  - `resource_dir` expects a directory with 2 files present:
    - resources.txt and units.txt
    - These contain all the files that the program needs to download, line by line.
    - Files may be empty to skip loading this type of resource.
  - `file_service_mode` adjusts the resource URL for Testcenter v15

### Usage
Example call:
`./load-test 10 config-fs-large.json`
- first parameter is the number of simulated users
- second the config file to be used

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
