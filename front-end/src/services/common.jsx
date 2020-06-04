import { authenticationService } from "./auth";

export function handleResponse(response) {
  return response.text().then(text => {
    const contentType = response.headers.get("Content-Type");
    const data = contentType === "application/json" ? JSON.parse(text) : text;

    if (!response.ok) {
      if (response.status === 401) {
        // The token has most probably expired, logout and redirect to login
        authenticationService.logout();
        location.reload();
      }

      return Promise.reject(data);
    }

    return data;
  });
}