#!/bin/bash

echo
echo "SUMMARY"
echo
echo "Run users: $1"

error_count=$(grep -ir "error" $2 | wc -l)
echo "Error Count: $error_count"

if [ ! $error_count -eq 0 ] ; then
    error_400=$(grep -r "400 Bad Request" $2 | wc -l)
    echo " Error Code 400: $error_400"
    error_404=$(grep -r "404 Not Found" $2 | wc -l)
    echo " Error Code 404: $error_404"
    error_500=$(grep -r "500 Internal Server Error" $2 | wc -l)
    echo " Error Code 500: $error_500"
    error_502=$(grep -r "502 Bad Gateway" $2 | wc -l)
    echo " Error Code 502: $error_502"
    error_504=$(grep -r "504 Gateway Timeout" $2 | wc -l)
    echo " Error Code 504: $error_504"
fi

success_count=$(grep -r "success" $2 | wc -l)
echo "Successful: $success_count"

echo
echo "TIMINGS"

total_time=$(cat $2 | grep Total | awk '{print $5}')
echo " total time: $total_time"

median=$(cat $2 | grep success| sort -n -k5 | awk '{all[NR] = $5} END{print all[int(NR*0.5 - 0.5)]}')
echo " median: $median"
percentile90=$(cat $2 | grep success| sort -n -k5 | awk '{all[NR] = $5} END{print all[int(NR*0.9 - 0.5)]}')
echo " 90th percentile: $percentile90"
