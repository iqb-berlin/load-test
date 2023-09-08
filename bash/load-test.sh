#!/bin/bash

### Author: Richard Henck
### Email: richard.henck@iqb.hu-berlin.de

### Required packages for scripts:
### - jq
### - /usr/bin/time (this differs slightly from the system built-in 'time')

if [ $# -eq 0 ]
  then
    echo "No arguments supplied. Using default of 1."
fi

#ulimit -n 4096

### Cleanup
rm -r results/1

LOGDIR="results/$(date +%FT%T)"
mkdir -p $LOGDIR

### NORMAL ###
seq $1 | parallel -j0 --progress --results results "/usr/bin/time -f \"%e\" ./testscript.sh 2>&1" 1> $LOGDIR/output.log

### FS ###
#parallel -j0 --progress --results results "/usr/bin/time -f \"%e\" ./testscript-fs.sh 2>&1" ::: $(seq 1 $1) 1> $LOGDIR/output.log

### MULTI HOST ###
#parallel -j0 --progress --results results --controlmaster -S 172.28.37.42,: "/usr/bin/time -f \"%e\" ./testscript.sh 2>&1" ::: $(seq 1 $1) 1> $LOGDIR/output.log

./summary.sh $1 $LOGDIR
