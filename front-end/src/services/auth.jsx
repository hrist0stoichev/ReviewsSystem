import { BehaviorSubject } from 'rxjs';

import config from 'config';
import {handleResponse} from "./common";

const currentUserSubject = new BehaviorSubject(JSON.parse(localStorage.getItem('currentUser')));

export const authenticationService = {
  login,
  facebookLogin,
  logout,
  register,
  authHeader,
  currentUser: currentUserSubject.asObservable(),
  get currentUserValue () { return currentUserSubject.value }
};

function authHeader() {
  const currentUser = currentUserSubject.value;
  if (currentUser && currentUser.token) {
    return { Authorization: `Bearer ${currentUser.token}` };
  }

  return {};
}

function login(email, password) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  };

  return handleLogin(`${config.apiUrl}/api/v1/token`, requestOptions);
}

function facebookLogin(state, code) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ state, code })
  };

  return handleLogin(`${config.apiUrl}/api/v1/token/facebook`, requestOptions);
}

function handleLogin(url, requestOptions) {
  return fetch(url, requestOptions)
    .then(handleResponse)
    .then(user => {
      localStorage.setItem('currentUser', JSON.stringify(user));
      currentUserSubject.next(user);

      return user;
    });
}

function register(user) {
  const requestOptions = {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(user),
  }

  return fetch(`${config.apiUrl}/api/v1/users`, requestOptions)
    .then(handleResponse)
}

function logout() {
  localStorage.removeItem('currentUser');
  currentUserSubject.next(null);
}