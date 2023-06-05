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

parallel -j0 --progress --results results "/usr/bin/time -f \"%e\" ./testscript.sh 2>&1" ::: $(seq 1 $1) 1> $LOGDIR/output.log
#parallel -j0 --progress --results results "/usr/bin/time -f \"%e\" ./testscript-fs.sh 2>&1" ::: $(seq 1 $1) 1> $LOGDIR/output.log

./summary.sh $1 $LOGDIR
