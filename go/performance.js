import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 50 },  // Ramp up to 50 virtual users over 1 minute
    { duration: '3m', target: 50 },  // Stay at 50 virtual users for 3 minutes
    { duration: '1m', target: 0 },   // Ramp down to 0 virtual users over 1 minute
  ],
};

export default function () {
  // Replace these URLs with your application's endpoints
  const baseURL = 'http://localhost:1323';
  const stringinateURL = `${baseURL}/stringinate`;
  const transformURL = `${baseURL}/transform`;
  const statsURL = `${baseURL}/stats`;

  // Send requests to your endpoints
  const payload = JSON.stringify({ input: 'your-string-goes-here' });
  
  // Send a POST request to /stringinate
  let stringinateResponse = http.post(stringinateURL, payload, { headers: { 'Content-Type': 'application/json' } });
  check(stringinateResponse, { 'Stringinate status is 200': (r) => r.status === 200 });

  // Send a POST request to /transform
  let transformResponse = http.post(transformURL, payload, { headers: { 'Content-Type': 'application/json' } });
  check(transformResponse, { 'Transform status is 200': (r) => r.status === 200 });

  // Send a GET request to /stats
  let statsResponse = http.get(statsURL);
  check(statsResponse, { 'Stats status is 200': (r) => r.status === 200 });

  // Sleep for a short duration between requests (e.g., 0.1 seconds)
  sleep(0.1);
}
