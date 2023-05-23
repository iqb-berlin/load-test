import http from 'k6/http';
import { sleep, check } from 'k6';

// To output to Prometheus and Grafana:
// Run with: k6 run -o experimental-prometheus-rw script.js
// Otherwise remove the output parameter.

const hostname = 'http://localhost'
const booklet = {bookletName: 'BOOKLET.SAMPLE-1'}
// const booklet = {bookletName: 'BOOKLETV8HT2022ENGTHA'}
const auth = {
    name: 'test-no-pw',
    password: ''
}
// const auth = {
//     name: 'lasttest.gross01.mit',
//     password: 'last'
// }

export const options = {
    vus: 10,
    duration: '20s',
    iterations: 10,
    tags: { testid: 'ricrun6' }
};

export default function () {
    const token = getToken();
    const payload = JSON.stringify(booklet);
    const headers = {
        headers: { 'AuthToken': token }
    };
    const response = http.put(hostname + '/api/test', payload, headers)
    const testID = response.body;

    http.get(hostname + '/api/test/' + testID, headers)
    http.get(hostname + '/api/test/' + testID + '/unit/UNIT.SAMPLE/alias/UNIT.SAMPLE', headers);
    sleep(1);
}

function getToken() {
    const response = http.put(hostname + '/api/session/login', JSON.stringify(auth))
    check(response, {
        'is status 200': (r) => r.status === 200,
        'verify token': (r) => r.body.includes('token')
    });
    // console.log('response.body', response.body);
    // console.log('response.body token parsed', JSON.parse(response.body).token);
    return JSON.parse(response.body).token;
}
