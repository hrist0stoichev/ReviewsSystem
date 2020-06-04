import React, {useEffect, useState} from "react";
import {restaurantsService} from "../services/restaurants";
import Image from "react-bootstrap/Image";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Badge from "react-bootstrap/Badge";
import Review from "./Review";
import Carousel from "react-bootstrap/Carousel";
import {reviewsService} from "../services/reviews";
import Button from "react-bootstrap/Button";
import {authenticationService} from "../services/auth";
import AddReview from "./AddReview";
import FormCheck from "react-bootstrap/FormCheck";

export default function Restaurant(props) {
  const [restaurant, setRestaurant] = useState({});
  const [reviews, setReviews] = useState([]);
  const [addReviewVisible, setAddReviewVisible] = useState(false);
  const [unansweredSwitchToggled, setUnansweredSwitchToggled]  = useState(false);

  useEffect(() => {
    loadRestaurant();
    loadReviews(isOwner());
    setUnansweredSwitchToggled(isOwner());
  }, []);

  const handleUnansweredSwitch = () => {
    setUnansweredSwitchToggled(prev => {
      loadReviews(!prev)
      return !prev
    });
  }

  const loadRestaurantAndReviews = () => {
    loadRestaurant();
    loadReviews(unansweredSwitchToggled);
  }

  const loadRestaurant = () => {
    restaurantsService.getSingle(props.match.params.id)
      .then(res => setRestaurant(res))
      .catch(err => props.showAlert(err, false));
  }

  const loadReviews = (unansweredOnly) => {
    reviewsService.getForRestaurant(props.match.params.id, 9, 0, unansweredOnly)
      .then(res => setReviews(res))
      .catch(err => props.showAlert(err, false));
  }
  
  const isOwner = () => {
    return authenticationService.currentUserValue && authenticationService.currentUserValue.role === "owner"
  }

  return (
    <>
      <AddReview handleAddingNewReview={loadRestaurantAndReviews} restaurant_id={props.match.params.id} showAlert={props.showAlert} show={addReviewVisible} handleClose={() => {setAddReviewVisible(false)}}/>
      <Row>
        <Col lg={6}>
          <Row>
            <Col lg={9}>
              <h1>{restaurant.name}</h1>
              <small>{`${restaurant.address}, ${restaurant.city}`}</small>
            </Col>
            <Col lg={3}>
              <h1 style={{lineHeight:"2em"}}><Badge pill variant="primary">{restaurant.average_rating} <span>&#9733;</span></Badge> </h1>
            </Col>
          </Row>
          <hr />
          <p>{restaurant.description}</p>
        </Col>
        <Col lg={6}>
          {authenticationService.currentUserValue && authenticationService.currentUserValue.role !== "owner" &&<Row>
            <Col lg={{ offset: 8, span: 4}}>
              <div style={{padding: "10px", right: 0}}><Button onClick={() => {setAddReviewVisible(true)}} variant="primary">Leave a review</Button></div>
            </Col>
          </Row>}
          <Image src={restaurant.img} fluid rounded />
        </Col>
      </Row>
      <hr />
      {/*If a min_review exists, it is guaranteed that a max_review also exists*/}
      {restaurant.min_review &&
      <Row>
        <Col style={{marginTop: "40px"}} lg={{ offset: 1, span: 4}}>
          <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={restaurant.min_review} />
        </Col>
        <Col style={{marginTop: "40px"}} lg={{ offset: 2, span: 4}}>
          <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={restaurant.max_review} />
        </Col>
      </Row>
      }
      {isOwner() &&
      <Row style={{marginTop: "60px"}}>
        <FormCheck
          id="unansweredSwitch"
          type="switch"
          checked={unansweredSwitchToggled}
          onChange={handleUnansweredSwitch}
          label="Unanswered only"
        />
      </Row>}
      {/*Hardcode the 9 most-recent reviews*/}
      {reviews.length > 0 &&
      <Row style={{borderRadius: "5px", marginTop: isOwner() ? "20px" : "80px"}}>
        <Carousel style={{width: "100%"}}>
          <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 0 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[0]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 1 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[1]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 2 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[2]}/>}
              </Col>
            </Row>
          </Carousel.Item>
          {reviews.length > 3 && <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 3 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[3]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 4 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[4]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 5 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[5]}/>}
              </Col>
            </Row>
          </Carousel.Item>}
          {reviews.length > 6 && <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 6 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[6]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 7 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[7]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 8 && <Review handleAnsweringReview={loadRestaurantAndReviews} showAlert={props.showAlert} review={reviews[8]}/>}
              </Col>
            </Row>
          </Carousel.Item>}
        </Carousel>
      </Row>
      }
    </>
  )
}