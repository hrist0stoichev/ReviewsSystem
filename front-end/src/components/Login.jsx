import React, {useState} from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import {authenticationService} from "../services/auth";
import config from 'config';

export default function Login(props) {
  const [validated, setValidated] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
    event.stopPropagation();

    const form = event.currentTarget;

    if (form.checkValidity() === true) {
      const form = event.currentTarget;
      authenticationService.login(form.email.value, form.password.value)
        .then((res) => {
          props.showAlert("Hello, " + res.email, true)
          props.history.push("/restaurants")
        })
        .catch((err) => {
          props.showAlert(err, false)
        })
    } else {
      setValidated(true)
    }
  }

  return (
    <Col style={{ marginTop: "3em" }} lg={{ span: 4, offset: 4 }}>
      <h2 style={{textAlign:"center"}}>Login</h2>
      <Form noValidate validated={validated} onSubmit={handleSubmit}>
        <Form.Group controlId="formEmail">
          <Form.Label>Email address</Form.Label>
          <Form.Control type="email" placeholder="Email" name="email" required />
          <Form.Control.Feedback type="invalid">Please provide a valid email!</Form.Control.Feedback>
        </Form.Group>

        <Form.Group controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control type="password" placeholder="Password" name="password" required />
          <Form.Control.Feedback type="invalid">Please provide a password!</Form.Control.Feedback>
        </Form.Group>


        <Button variant="success" type="submit" >Login</Button>
        <Button href={`${config.apiUrl}/api/v1/facebookauth`} style={{position: "absolute", right: "1em"}} variant="primary">Login with facebook</Button>
      </Form>
    </Col>
  );
}
