import {authenticationService} from "./auth";
import config from "config";
import {handleResponse} from "./common";

export const restaurantsService = {
  get,
  getSingle,
  add,
  del
}

function add(restaurant) {
  const requestOptions = {
    method: 'POST',
    headers: {
      ...authenticationService.authHeader(),
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(restaurant),
  };

  return fetch(`${config.apiUrl}/api/v1/restaurants`, requestOptions)
    .then(handleResponse)
}

function get(top, skip, minRating, maxRating, orderBy) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/restaurants?top=${top}&skip=${skip}&minRating=${minRating}&maxRating=${maxRating}&orderBy=${orderBy}`, requestOptions)
    .then(handleResponse)
}

function getSingle(id) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/restaurants/${id}`, requestOptions)
    .then(handleResponse)
}

function del(id) {
  const requestOptions = {
    method: 'DELETE',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/restaurants/${id}`, requestOptions)
    .then(handleResponse)
}