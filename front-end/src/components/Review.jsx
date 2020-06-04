import React, {useState} from "react";
import Button from "react-bootstrap/Button";
import {authenticationService} from "../services/auth";
import Popover from "react-bootstrap/Popover";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Form from "react-bootstrap/Form";
import {reviewsService} from "../services/reviews";

export default function Review(props) {
  const [validated, setValidated] = useState(false);

  const handleAnswerValidation = (event) => {
    const answer = event.currentTarget.answer.value;
    event.currentTarget.answer.setCustomValidity(answer.length < 30 || answer.length > 300 ? "Invalid answer" : "")
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    event.stopPropagation();

    const form = event.currentTarget;

    if (form.checkValidity() === true) {
      reviewsService.addAnswer(props.review.id, form.answer.value)
        .then(() => {
          props.showAlert(`You answered successfully!`, true);
          location.reload();
        })
        .catch((err) => {
          props.showAlert(err, false)
        })
    } else {
      setValidated(true);
    }
  }

  const popover = (
    <Popover>
      <Popover.Title as="h3">Answer</Popover.Title>
      <Popover.Content>{props.review.answer}</Popover.Content>
    </Popover>
  );

  const answerPopover = (
    <Popover style={{textAlign: "right", width:"35vw"}} >
      <Popover.Content>
        <Form noValidate validated={validated} onInput={handleAnswerValidation} onSubmit={handleSubmit}>
          <Form.Group controlId="formAnswer">
            <Form.Control as="textarea" placeholder="Answer" name="answer" required />
            <Form.Control.Feedback type="invalid">Please provide a valid answer (30-300 characters)!</Form.Control.Feedback>
          </Form.Group>

          <Button size="sm" variant="primary" type="submit">Submit</Button>
        </Form>
      </Popover.Content>
    </Popover>
  )

  return (
    <div style={{textAlign: "center"}}>
      <h5>{props.review.reviewer}</h5>
      <div>My rating: {props.review.rating} <span>&#9733;</span></div>
      <div style={{marginTop: "10px"}}><i>{props.review.comment}</i></div>
      <div>
        <small>{new Intl.DateTimeFormat('en-US', {year: 'numeric', month: '2-digit',day: '2-digit', hour: '2-digit', minute: '2-digit'}).format(Date.parse(props.review.timestamp))}</small>

        {authenticationService.currentUserValue && authenticationService.currentUserValue.role === "owner" && props.review.answer === null &&
        <OverlayTrigger trigger="click" placement="right" overlay={answerPopover}>
          <a style={{marginLeft: "1em"}} onClick={(e) => {e.preventDefault()}} href="">Answer</a>
        </OverlayTrigger>}

        {authenticationService.currentUserValue && props.review.answer !== null &&
        <OverlayTrigger trigger="click" placement="right" overlay={popover}>
          <a style={{marginLeft: "1em", color: "green"}} onClick={(e) => {e.preventDefault()}} href="">View Answer</a>
        </OverlayTrigger>}
      </div>
    </div>
  )
}