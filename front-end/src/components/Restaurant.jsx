import React, {useEffect, useState} from "react";
import {restaurantsService} from "../services/restaurants";
import Image from "react-bootstrap/Image";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Badge from "react-bootstrap/Badge";
import Review from "./Review";
import Container from "react-bootstrap/Container";

export default function Restaurant(props) {
  const [restaurant, setRestaurant] = useState({});

  useEffect(() => {
    restaurantsService.getSingle(props.match.params.id)
      .then(res => setRestaurant(res))
      .catch(err => props.showAlert(err, false))
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
          <p>{restaurant.description} asdoja sdj aiojsdi jaspdj klamsdpuas idka slkdn ajpsidu asjid amsdpo i</p>
        </Col>
        <Col lg={6}>
          <Image src={restaurant.img} fluid rounded />
        </Col>
      </Row>
      <hr />
      {restaurant.min_review &&
      <Row>
        <Col lg={{ offset: 1, span: 4}}>
          <Review review={restaurant.min_review} />
        </Col>
        <Col lg={{ offset: 2, span: 4}}>
          <Review review={restaurant.max_review} />
        </Col>
      </Row>
      }
    </>
  )
}