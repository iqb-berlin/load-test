#!/bin/bash

# Color variables
red='\033[0;91m'
green='\033[0;32m'
yellow='\033[0;33m'
grey='\033[0;37m'
# Clear the color after that
clear='\033[0m'

echo
echo "SUMMARY"
echo
echo "Run users: $1"

error_count=$(grep "ERROR" $2 | wc -l)
if [ ! $error_count -eq 0 ] ; then
    echo -e "${red}Error Count: $error_count${clear}"
    error_400=$(grep "ERROR" $2 | grep "Response: 400" | wc -l)
    echo " Error Code 400: $error_400"
    error_401=$(grep "ERROR" $2 | grep "Response: 401" | wc -l)
    echo " Error Code 401: $error_401"
    error_403=$(grep "ERROR" $2 | grep "Response: 403" | wc -l)
    echo " Error Code 403: $error_403"
    error_404=$(grep "ERROR" $2 | grep "Response: 404" | wc -l)
    echo " Error Code 404: $error_404"
    error_429=$(grep "ERROR" $2 | grep "Response: 429" | wc -l)
    echo " Error Code 429: $error_429"
    error_500=$(grep "ERROR" $2 | grep "Response: 500" | wc -l)
    echo " Error Code 500: $error_500"
    error_502=$(grep "ERROR" $2 | grep "Response: 502" | wc -l)
    echo " Error Code 502: $error_502"
    error_503=$(grep "ERROR" $2 | grep "Response: 503" | wc -l)
    echo " Error Code 503: $error_503"
    error_504=$(grep "ERROR" $2 | grep "Response: 504" | wc -l)
    echo " Error Code 504: $error_504"
    error_timeout=$(grep "ERROR" $2 | grep "Client Timeout" | wc -l)
    echo " Error Client Timeout: $error_timeout"
fi

success_count=$(grep "SUCCESS" $2 | wc -l)
echo -e "${green}Successful: $success_count${clear}"

automatic_retries_needed=$(grep "WARNING" $2 | grep "Automatic retries left" | wc -l)
if [ ! $automatic_retries_needed -eq 0 ] ; then
  echo -e "${grey}Automatic Retries nedded for ressource loading: $automatic_retries_needed${clear}"
fi
manual_retries_needed=$(grep "WARNING" $2 | grep "Manual retries left" | wc -l)
if [ ! $manual_retries_needed -eq 0 ] ; then
  echo -e "${yellow}Manual Retries needed: $manual_retries_needed${clear}"
fi

echo
echo "TIMINGS"

total_time=$(cat $2 | grep Total | awk '{print $5}')
echo " total time: $total_time"

median=$(cat $2 | grep SUCCESS | sort -n -k5 | awk '{all[NR] = $7} END{print all[int(NR*0.5 - 0.5)]}')
echo " median: $median"
percentile90=$(cat $2 | grep SUCCESS | sort -n -k5 | awk '{all[NR] = $7} END{print all[int(NR*0.9 - 0.5)]}')
echo " 90th percentile: $percentile90"
