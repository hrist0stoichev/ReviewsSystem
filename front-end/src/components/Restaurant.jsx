import React, {useEffect, useState} from "react";
import {restaurantsService} from "../services/restaurants";
import Image from "react-bootstrap/Image";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Badge from "react-bootstrap/Badge";
import Review from "./Review";
import Carousel from "react-bootstrap/Carousel";
import {reviewsService} from "../services/reviews";

export default function Restaurant(props) {
  const [restaurant, setRestaurant] = useState({});
  const [reviews, setReviews] = useState([]);

  useEffect(() => {
    restaurantsService.getSingle(props.match.params.id)
      .then(res => setRestaurant(res))
      .catch(err => props.showAlert(err, false));

    reviewsService.getForRestaurant(props.match.params.id, 9, 0)
      .then(res => setReviews(res))
      .catch(err => props.showAlert(err, false));
  }, []);

  return (
    <>
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
          <Image src={restaurant.img} fluid rounded />
        </Col>
      </Row>
      <hr />
      {/*If a min_review exists, it is guaranteed that a max_review also exists*/}
      {restaurant.min_review &&
      <Row style={{marginTop: "40px"}}>
        <Col lg={{ offset: 1, span: 4}}>
          <Review review={restaurant.min_review} />
        </Col>
        <Col lg={{ offset: 2, span: 4}}>
          <Review review={restaurant.max_review} />
        </Col>
      </Row>
      }
      {/*Hardcode the 9 most-recent reviews*/}
      {reviews.length > 0 &&
      <Row style={{borderRadius: "5px", marginTop: "80px"}}>
        <Carousel style={{width: "100%"}}>
          <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 0 && <Review review={reviews[0]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 1 && <Review review={reviews[1]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 2 && <Review review={reviews[2]}/>}
              </Col>
            </Row>
          </Carousel.Item>
          {reviews.length > 3 && <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 3 && <Review review={reviews[3]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 4 && <Review review={reviews[4]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 5 && <Review review={reviews[5]}/>}
              </Col>
            </Row>
          </Carousel.Item>}
          {reviews.length > 6 && <Carousel.Item>
            <Row>
              <Col lg={4}>
                {reviews.length > 6 && <Review review={reviews[6]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 7 && <Review review={reviews[7]}/>}
              </Col>
              <Col lg={4}>
                {reviews.length > 8 && <Review review={reviews[8]}/>}
              </Col>
            </Row>
          </Carousel.Item>}
        </Carousel>
      </Row>
      }
    </>
  )
}