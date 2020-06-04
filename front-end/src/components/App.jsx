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
            <Route exact path="/" render={(props) => {props.history.push("login")}} />
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