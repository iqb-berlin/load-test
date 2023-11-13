#!/bin/bash

# Color variables
red='\033[0;91m'
green='\033[0;32m'
yellow='\033[0;33m'
# Clear the color after that
clear='\033[0m'

echo
echo "SUMMARY"
echo
echo "Run users: $1"

error_count=$(grep "ERROR" $2 | wc -l)
if [ ! $error_count -eq 0 ] ; then
    echo -e "${red}Error Count: $error_count${clear}"
    error_400=$(grep "WARNING" $2 | grep "400 Bad Request" | wc -l)
    echo " Error Code 400: $error_400"
    error_404=$(grep "WARNING" $2 | grep "404 Not Found" | wc -l)
    echo " Error Code 404: $error_404"
    error_500=$(grep "WARNING" $2 | grep "500 Internal Server Error" | wc -l)
    echo " Error Code 500: $error_500"
    error_502=$(grep "WARNING" $2 | grep "502 Bad Gateway" | wc -l)
    echo " Error Code 502: $error_502"
    error_504=$(grep "WARNING" $2 | grep "504 Gateway Timeout" | wc -l)
    echo " Error Code 504: $error_504"
    error_timeout=$(grep "WARNING" $2 | grep "Timeout; Retries left: 0" | wc -l)
    echo " Error Timeout: $error_timeout"
fi

success_count=$(grep "SUCCESS" $2 | wc -l)
echo -e "${green}Successful: $success_count${clear}"

retries_needed=$(grep "WARNING" $2 | grep "Retries left" | wc -l)
if [ ! $retries_needed -eq 0 ] ; then
  echo -e "${yellow}Retries needed: $retries_needed${clear}"
fi

echo
echo "TIMINGS"

total_time=$(cat $2 | grep Total | awk '{print $5}')
echo " total time: $total_time"

median=$(cat $2 | grep SUCCESS | sort -n -k5 | awk '{all[NR] = $7} END{print all[int(NR*0.5 - 0.5)]}')
echo " median: $median"
percentile90=$(cat $2 | grep SUCCESS | sort -n -k5 | awk '{all[NR] = $7} END{print all[int(NR*0.9 - 0.5)]}')
echo " 90th percentile: $percentile90"
