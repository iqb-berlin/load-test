#!/bin/bash

#Param 1: number of users
#Param 2: logdir

echo
echo "SUMMARY"
echo "Run users: $1"

error_count=$(grep -r "Request-Error:" results/1/* | wc -l)
echo "Error Count: $error_count"

if [ ! $error_count -eq 0 ] ; then
    error_400=$(grep -r "400 Bad Request" results/1/* | wc -l)
    echo "Error Code 400: $error_400"
    error_404=$(grep -r "404 Not Found" results/1/* | wc -l)
    echo "Error Code 404: $error_404"
    error_500=$(grep -r "500 Internal Server Error" results/1/* | wc -l)
    echo "Error Code 500: $error_500"
    error_502=$(grep -r "502 Bad Gateway" results/1/* | wc -l)
    echo "Error Code 502: $error_502"
    error_504=$(grep -r "504 Gateway Timeout" results/1/* | wc -l)
    echo "Error Code 504: $error_504"
fi

error_count_transfer_closed=$(grep -ir "transfer closed" results/1/* | wc -l)
if [ ! $error_count_transfer_closed -eq 0 ] ; then
  echo "Error transfer closed with x bytes remaining to read: $error_count_transfer_closed"
fi

success_count=$(grep -r "EndTime" results/1/* | wc -l)
echo "Successful: $success_count"

success_files=$(grep -rl "EndTime" results/1/*)
while read -r line; do
    tail -1 $line >> "$2/timings.txt"
done <<< $success_files
