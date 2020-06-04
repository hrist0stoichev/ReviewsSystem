import {authenticationService} from "./auth";
import config from "config";
import {handleResponse} from "./common";

export const reviewsService = {
  getForRestaurant,
  add,
  addAnswer
}
function getForRestaurant(id, top, skip, unanswered) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/reviews?restaurantId=${id}&top=${top}&skip=${skip}&unanswered=${unanswered ? "true" : "false"}`, requestOptions)
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

function addAnswer(id, answer) {
  const requestOptions = {
    method: 'PUT',
    headers: {
      ...authenticationService.authHeader(),
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({answer: answer}),
  };

  return fetch(`${config.apiUrl}/api/v1/reviews/${id}`, requestOptions)
    .then(handleResponse)
}