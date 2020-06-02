import {authenticationService} from "./auth";
import config from "config";
import {handleResponse} from "./common";

export const restaurantsService = {
  get
}

function get(top, skip, minRating, maxRating, orderBy) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/restaurants?top=${top}&skip=${skip}&minRating=${minRating}&maxRating=${maxRating}&orderBy=${orderBy}`, requestOptions)
    .then(handleResponse)
}
