#!/bin/bash

source load-test.cfg

DEFAULT_REQUEST_PARAMS="-sSf --retry 3"

echo "StartTime: $(date +%FT%T)"

### LOGIN
response=$(curl -X PUT $DEFAULT_REQUEST_PARAMS \
  "${HOSTNAME}/api/session/login" \
  --data-raw "{\"name\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}")
if [ ! $? -eq 0 ] ; then
  echo "Request-Error: /api/session/login"
  exit 0
else
  token=$(echo $response | jq -r '.token')
fi

### PUT Test
response=$(curl \
  -X PUT "${HOSTNAME}/api/test" $DEFAULT_REQUEST_PARAMS \
  --header "AuthToken: ${token}" \
  --data-raw "{\"bookletName\":\"${BOOKLET_NAME}\"}")
if [ ! $? -eq 0 ] ; then
  echo "Request-Error: $response_code for /test (PUT)"
  exit 0
else
  testID=$(echo $response)
fi

### GET Test
response_code=$(curl \
  -w '%{http_code}' $DEFAULT_REQUEST_PARAMS \
  -X GET "${HOSTNAME}/api/test/$testID" \
  --header "AuthToken: ${token}")
if [ ! $? -eq 0 ] ; then
  echo "Request-Error: $response_code for /test (GET)"
  exit 0
fi

### GET resources
while read file; do
  response_code=$(curl \
    -w '%{http_code}' $DEFAULT_REQUEST_PARAMS \
    -X GET "${HOSTNAME}/api/test/$testID/resource/${file}" \
    --header "AuthToken: ${token}")
  if [ ! $? -eq 0 ] ; then
    echo "Request-Error: $response_code for ${file}"
    exit 0
  fi
done < $RESOURCE_DIR/resources.txt

### GET units
while read file; do
  response_code=$(curl \
    -w '%{http_code}' $DEFAULT_REQUEST_PARAMS \
    -X GET "${HOSTNAME}/api/test/$testID/unit/${file}/alias/${file}" \
    --header "AuthToken: ${token}")
  if [ ! $? -eq 0 ] ; then
    echo "Request-Error: $response_code for ${file}"
    exit 0
  fi
done < $RESOURCE_DIR/units.txt

echo "EndTime: $(date +%FT%T)"
