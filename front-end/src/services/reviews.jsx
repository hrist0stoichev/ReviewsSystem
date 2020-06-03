import {authenticationService} from "./auth";
import config from "config";
import {handleResponse} from "./common";

export const reviewsService = {
  getForRestaurant
}
function getForRestaurant(id, top, skip) {
  const requestOptions = {
    method: 'GET',
    headers: authenticationService.authHeader()
  }

  return fetch(`${config.apiUrl}/api/v1/reviews?restaurantId=${id}&top=${top}&skip=${skip}`, requestOptions)
    .then(handleResponse)
}