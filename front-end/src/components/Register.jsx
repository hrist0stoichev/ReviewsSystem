import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import {authenticationService} from "../services/auth";

export default function Register(props) {
  const [validated, setValidated] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
    event.stopPropagation();

    const form = event.currentTarget;

    if (form.checkValidity() === true) {
      const user = {
        "email": form.email.value,
        "password": form.password.value,
        "confirm_password": form.confirmPassword.value,
        "is_owner": form.isOwner.value === "on",
      };

      authenticationService.register(user)
        .then(() => {
          props.showAlert("Please, confirm your email and log in.", true)
          props.history.push("/login")
        })
        .catch((err) => {
          props.showAlert(err, false)
        })
    } else {
      setValidated(true);
    }
  }

  const handlePasswordChange = (event) => {
    const password = event.currentTarget.password.value;
    const confirmPassword = event.currentTarget.confirmPassword.value;

    event.currentTarget.confirmPassword.setCustomValidity(password !== confirmPassword ? "Passwords do not match" : "")
  }

  return (
    <Col style={{ marginTop: "3em" }} lg={{ span: 4, offset: 4 }}>
      <h2 style={{textAlign:"center"}}>Register</h2>
      <Form noValidate validated={validated} onInput={handlePasswordChange} onSubmit={handleSubmit}>
        <Form.Group controlId="formEmail">
          <Form.Label>Email address</Form.Label>
          <Form.Control type="email" placeholder="Email" name="email" required />
          <Form.Control.Feedback type="invalid">Please provide a valid email!</Form.Control.Feedback>
        </Form.Group>

        <Form.Group controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control type="password" placeholder="Password" name="password" required pattern="(?=^.{8,64}$)(?=.*\d)(?=.*\W+)(?![.\n])(?=.*[A-Z])(?=.*[a-z]).*$" />
          <Form.Control.Feedback type="invalid">Please provide a valid password! (8+ characters, uppercase and lowercase letter, digit, and special symbol)</Form.Control.Feedback>
        </Form.Group>

        <Form.Group controlId="formConfirmPassword">
          <Form.Label>Confirm Password</Form.Label>
          <Form.Control type="password" placeholder="Confirm Password" name="confirmPassword" />
          <Form.Control.Feedback type="invalid">Confirmation password does not match!</Form.Control.Feedback>
        </Form.Group>

        <Form.Group controlId="formIsOwner">
          <Form.Check type="switch" label="Register as an owner" name="isOwner" />
        </Form.Group>

        <Button variant="primary" type="submit" >
          Register
        </Button>
      </Form>
    </Col>
  );
}
