### General

This program consists of 3 scripts:
- `testscript.sh` loads the test
- `load-test.sh` runs multiple instances (simulated users) of the testscript
- `summary.sh` is called in the end to parse the results folder and write aggregated values to the console

### Preparation

#### Required software
- jq (for parsing JSON responses)
- /usr/bin/time (this may differ from system built-in `time`)

#### Config file
- Configure test parameters in `load-test.cfg`
  - There are 3 pre-configured booklets available - comment in the one you want
    - Check the resource folder for details of used URLs
  - The `WORKSPACE` parameter is only used in the *file-service* version of the Testcenter
- `testscript.sh` has a retry mechanism for failed requests. You can change or turn it off, by manipulating the
 following line in the script:
`DEFAULT_REQUEST_PARAMS="-sSf --retry 3"`

#### File service version

When testing the Testcenter version with file service, you need to change the used testscript in `load-test.sh`.
Swap the commented line in the script to use `testscript-fs.sh`.

#### ulimit
When running lots of users, the script may run into the OS-limit of open file descriptors.
You may need to raise the allowed value before running the script. You can do that with `ulimit`:
```
ulimit -n 8192
```
For larger values you may have to change the allowed value for your user in `/etc/security/limits.conf`


### Usage
`testscript.sh` can be run on it's own to load a booklet just once
```
rhenck@load-test-master:~/load-test/bash$ ./testscript.sh 
StartTime: 2023-06-05T09:48:42
EndTime: 2023-06-05T09:48:45

```
You can also measure the time it took with `time`:
```
rhenck@load-test-master:~/load-test/bash$ time ./testscript.sh 
StartTime: 2023-06-05T09:48:51
EndTime: 2023-06-05T09:48:54

real	0m2,948s
user	0m0,345s
sys	0m0,158s
```

To test with multiple users call `load-test.sh` in a similar way. Pass the target user count as parameter.
```
rhenck@load-test-master:~/load-test/bash$ ./load-test.sh 10

Computers / CPU cores / Max jobs to run
1:local / 8 / 10

Computer:jobs running/jobs completed/%of started jobs/Average seconds to complete
local:0/10/100%/0.7s 

Run users: 10
Error Count: 0
Successful: 10
```

Measure the total duration with `time`:
```
rhenck@load-test-master:~/load-test/bash$ time ./load-test.sh 10

Computers / CPU cores / Max jobs to run
1:local / 8 / 10

Computer:jobs running/jobs completed/%of started jobs/Average seconds to complete
local:0/10/100%/0.7s 

Run users: 10
Error Count: 0
Successful: 10

real	0m6,445s
user	0m3,880s
sys	0m3,639s
```

### Results
The main results are shown on the console after the run, listing the most common response-codes. Details can be found in the `results` directory.
- The output (and potentially errors messages) of every execution of `testscript.sh` can be found
 under `results/1` with a subfolder for every sub-process (simulated user)
  - If you need to dig into the results for debugging, you may do that here. Mind that the folder content is parsed by `summary.sh`
   so you usually don't need to.
  - This folder is cleared and overwritten on every call of `load-test.sh`! 
- Stdout for every user and collected timings are kept in a timestamped folder

```
rhenck@load-test-master:~/load-test/bash$ ls -la results/2023-05-23T17:37:23
insgesamt 16
drwxr-xr-x   2 rhenck rhenck 4096 23. Mai 17:37 .
drwxr-xr-x 140 rhenck rhenck 4096  5. Jun 09:53 ..
-rw-r--r--   1 rhenck rhenck  650 23. Mai 17:37 output.log
-rw-r--r--   1 rhenck rhenck   50 23. Mai 17:37 timings.txt
```
