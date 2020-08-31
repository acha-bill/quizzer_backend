import axios from 'axios';

var instance = null;
const baseURL = process.env.REACT_APP_BASE_URL ?  process.env.REACT_APP_BASE_URL : 'http://localhost:8081/api/v1'

export function getApi(token) {
  instance = axios.create({
    baseURL: baseURL,
    timeout: 0,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      Authorization: `Bear ${token}`,
    },
  });
  return instance;
}
