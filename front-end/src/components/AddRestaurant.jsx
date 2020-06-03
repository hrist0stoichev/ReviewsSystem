import React, {useState} from "react";

import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import {restaurantsService} from "../services/restaurants";

export default function AddRestaurant(props) {
  const [validated, setValidated] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
    event.stopPropagation();

    const form = event.currentTarget;

    if (form.checkValidity() === true) {
      const restaurant = {
        "name": form.name.value,
        "city": form.city.value,
        "address": form.address.value,
        "img": form.img.value,
        "description": form.description.value,
      };

      restaurantsService.add(restaurant)
        .then((res) => {
          props.showAlert(`Restaurant ${res.name} successfully added`, true);
          props.history.push(`/restaurants/${res.id}`);
          props.handleClose();
        })
        .catch((err) => {
          props.showAlert(err, false)
        })
    } else {
      setValidated(true);
    }
  }

  const handleDescriptionValidation = (event) => {
    const description = event.currentTarget.description.value;
    event.currentTarget.description.setCustomValidity(description.length < 30 || description.length > 500 ? "Invalid description" : "")
  }

  return (
    <Modal
      show={props.show}
      onHide={props.handleClose}
      backdrop="static"
      keyboard={false}
    >
      <Modal.Header closeButton>
        <Modal.Title>Create new restaurant</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form noValidate validated={validated} onInput={handleDescriptionValidation} onSubmit={handleSubmit}>
          <Form.Group controlId="formName">
            <Form.Label>Name</Form.Label>
            <Form.Control type="text" placeholder="Name" name="name" required pattern="^[\w\s']{5,60}$" />
            <Form.Control.Feedback type="invalid">Please provide a valid name (5-60 characters)!</Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="formCity">
            <Form.Label>City</Form.Label>
            <Form.Control type="text" placeholder="City" name="city" required pattern="^[\w\s]{5,30}$" />
            <Form.Control.Feedback type="invalid">Please provide a valid city (5-30 characters)!</Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="formAddress">
            <Form.Label>Address</Form.Label>
            <Form.Control type="text" placeholder="Address" name="address" required pattern="^[\w\s]{5,100}$" />
            <Form.Control.Feedback type="invalid">Please provide a valid address (5-100 characters)!</Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="formImg">
            <Form.Label>Image</Form.Label>
            <Form.Control type="url" placeholder="Image URL" name="img" required />
            <Form.Control.Feedback type="invalid">Please provide a valid URL!</Form.Control.Feedback>
          </Form.Group>

          <Form.Group controlId="formDescription">
            <Form.Label>Description</Form.Label>
            <Form.Control as="textarea" placeholder="Description" name="description" required />
            <Form.Control.Feedback type="invalid">Please provide a valid description (30-500 characters)!</Form.Control.Feedback>
          </Form.Group>

          <Modal.Footer>
            <Button variant="secondary" onClick={props.handleClose}>Cancel</Button>
            <Button variant="primary" type="submit">Create</Button>
          </Modal.Footer>
        </Form>
      </Modal.Body>

    </Modal>
  );
}