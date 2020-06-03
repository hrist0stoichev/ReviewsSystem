import {authenticationService} from "./auth";
import config from "config";
import {handleResponse} from "./common";

export const reviewsService = {
  getForRestaurant,
  add
}
function getForRestaurant(id, top, skip) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/reviews?restaurantId=${id}&top=${top}&skip=${skip}`, requestOptions)
    .then(handleResponse)
}

function add(review) {
  const requestOptions = {
    method: 'POST',
    headers: {
      ...authenticationService.authHeader(),
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(review),
  };

  return fetch(`${config.apiUrl}/api/v1/reviews`, requestOptions)
    .then(handleResponse)
}