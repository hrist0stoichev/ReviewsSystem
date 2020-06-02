import React, { useState } from 'react';
import { HashRouter, Route, Switch } from "react-router-dom";

import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Alert from "react-bootstrap/Alert";

import Header from './Header';
import Login from "./Login";
import Register from "./Register";

export default function App() {
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
    <div >
      <Header showAlert={showAlert} />

      <Container>
        <Alert style={{position:"fixed", right: 30, width: "500px", zIndex: 999 }} show={alert.show} variant={alert.variant} onClose={closeAlert} dismissible>
          <Alert.Heading>{alert.heading}</Alert.Heading>
          <p>{alert.msg}</p>
        </Alert>

        <Row>
          <HashRouter>
            <Switch>
              <Route exact path="/login" render={(props) => <Login showAlert={showAlert} {...props} />} />
              <Route exact path="/register" render={(props) => <Register showAlert={showAlert} {...props} />} />
            </Switch>
          </HashRouter>
        </Row>
      </Container>
    </div>
  );
}