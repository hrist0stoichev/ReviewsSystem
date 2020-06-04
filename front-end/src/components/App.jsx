import React, { useState } from 'react';
import { HashRouter, Route, Switch } from "react-router-dom";

import Container from "react-bootstrap/Container";
import Alert from "react-bootstrap/Alert";

import Header from './Header';
import Login from "./Login";
import Register from "./Register";
import Restaurant from "./Restaurant";
import RestaurantList from "./RestaurantList";
import AddRestaurant from "./AddRestaurant";

import queryString from 'query-string';
import {authenticationService} from "../services/auth";

export default function App() {
  const [addRestaurantModalVisible, setAddRestaurantModalVisible] = useState(false)
  const [alert, setAlert] = useState({
    show: false,
    variant: "",
    heading: "",
    msg: "",
  })

  const closeAlert = () => {
    setAlert({...alert, show: false})
  }

  const showAlert = (msg, success) => {
    setAlert({
      show: true,
      variant: success ? "success" : "danger",
      heading: success ? "Success!" : "Failed!",
      msg: msg,
    })
  }

  const handleRedirections = (props) => {
    props.history.push("/login");

    const params = queryString.parse(props.location.search);
    if (params.confirmation_successful === "true") {
      showAlert("Email confirmation successful", true);
    }
  }

  const handleFacebookCallback = (props) => {
    const fullUrl = window.location.href;
    const queryParameters = fullUrl.substring(fullUrl.indexOf("?"), fullUrl.length - 5);
    const params = queryString.parse(queryParameters);

    if (params.state && params.code) {
      window.history.replaceState({}, '', location.pathname);
      authenticationService.facebookLogin(params.state, params.code)
        .then((res) => {
          props.history.push("/restaurants");
          showAlert("Hello, " + res.email, true);
        })
        .catch((err) => {
          props.history.push("/login")
          showAlert(err, false);
        })
    }
  }

  return (
    <div>
      <Header showAlert={showAlert} showAddRestaurantModal={() => setAddRestaurantModalVisible(true)} />

      <Container>
        <Alert style={{position:"fixed", right: 30, width: "500px", zIndex: 999 }} show={alert.show} variant={alert.variant} onClose={closeAlert} dismissible>
          <Alert.Heading>{alert.heading}</Alert.Heading>
          <p>{alert.msg}</p>
        </Alert>

        <HashRouter>
          <Switch>
            <Route exact path="/_=_" render={(props) => {handleFacebookCallback(props)}} /> {/* Hack for handling facebook callback using a hashrouter */}
            <Route exact path="/" render={(props) => {handleRedirections(props)}} />
            <Route exact path="/login" render={(props) => <Login showAlert={showAlert} {...props} />} />
            <Route exact path="/register" render={(props) => <Register showAlert={showAlert} {...props} />} />
            <Route exact path="/restaurants" render={(props) => <RestaurantList showAlert={showAlert} {...props} />} />
            <Route exact path="/restaurants/:id" render={(props) => <Restaurant showAlert={showAlert} {...props} />} />
            <Route render={() => <h1 style={{textAlign: "center"}}>Page Not Found!</h1>} />
          </Switch>
          <Route render={(props) => <AddRestaurant show={addRestaurantModalVisible} showAlert={showAlert} handleClose={() => setAddRestaurantModalVisible(false)} {...props} />} />
        </HashRouter>
      </Container>
    </div>
  );
}