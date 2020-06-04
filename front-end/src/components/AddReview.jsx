import React, {useState} from "react";

import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import {reviewsService} from "../services/reviews";

export default function AddReview(props) {
  const [validated, setValidated] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
    event.stopPropagation();

    const form = event.currentTarget;

    if (form.checkValidity() === true) {
      const review = {
        "restaurant_id": props.restaurant_id,
        "rating": parseInt(form.rating.value, 10),
        "comment": form.comment.value,
      };

      reviewsService.add(review)
        .then(() => {
          props.showAlert(`Your review was successfully added`, true);
          props.handleClose();
          location.reload();
        })
        .catch((err) => {
          props.showAlert(err, false)
          props.handleClose();
        })
    } else {
      setValidated(true);
    }
  }

  const handleCommentValidation = (event) => {
    const comment = event.currentTarget.comment.value;
    event.currentTarget.comment.setCustomValidity(comment.length < 30 || comment.length > 300 ? "Invalid comment" : "")
  }

  return (
    <Modal show={props.show} onHide={props.handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Add new review</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form noValidate validated={validated} onInput={handleCommentValidation} onSubmit={handleSubmit}>
          <Form.Group controlId="formRating">
            <Form.Label>Rating</Form.Label>
            <Form.Control type="range" min="1" max="5" step="1" name="rating" />
          </Form.Group>

          <Form.Group controlId="formComment">
            <Form.Label>Comment</Form.Label>
            <Form.Control as="textarea" placeholder="Comment" name="comment" required />
            <Form.Control.Feedback type="invalid">Please provide a valid comment (30-300 characters)!</Form.Control.Feedback>
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