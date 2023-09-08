#!/bin/bash

echo
echo "SUMMARY"
echo "Run users: $1"

error_count=$(grep -ir "error" $2 | wc -l)
echo "Error Count: $error_count"

if [ ! $error_count -eq 0 ] ; then
    error_400=$(grep -r "400 Bad Request" $2 | wc -l)
    echo "Error Code 400: $error_400"
    error_404=$(grep -r "404 Not Found" $2 | wc -l)
    echo "Error Code 404: $error_404"
    error_500=$(grep -r "500 Internal Server Error" $2 | wc -l)
    echo "Error Code 500: $error_500"
    error_502=$(grep -r "502 Bad Gateway" $2 | wc -l)
    echo "Error Code 502: $error_502"
    error_504=$(grep -r "504 Gateway Timeout" $2 | wc -l)
    echo "Error Code 504: $error_504"
fi

success_count=$(grep -r "success" $2 | wc -l)
echo "Successful: $success_count"

