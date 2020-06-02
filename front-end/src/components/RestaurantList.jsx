import React, {useEffect, useState} from "react";
import {restaurantsService} from "../services/restaurants";
import CardDeck from "react-bootstrap/CardDeck";
import Card from "react-bootstrap/Card";
import InputRange from 'react-input-range';
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";

export default function RestaurantList(props) {
  const defaultPageSize = 21;
  const defaultOrderBy = "average_rating"

  const [restaurants, setRestaurants] = useState([])
  const [ratingRange, setRatingRange] = useState({min: 1, max: 5})
  const [page, setPage] = useState(1)

  useEffect(() => {
    updateRestaurantsList()
  }, []);

  const updateRestaurantsList = () => {
    restaurantsService.get(defaultPageSize, (page - 1) * defaultPageSize, ratingRange.min, ratingRange.max, defaultOrderBy)
      .then(res => {
        setRestaurants(res)
      })
      .catch(err => {
        props.showAlert(err, false)
      })
  }

  const handleCardClick = (event) => {
    props.history.push("/restaurants/" + event.currentTarget.id)
  }

  const getDecks = () => {
    const decks = [];
    for (let i = 0; i < restaurants.length; i += 3) {
      const deck = [];
      for (let j = i; j < i + 3 && j < restaurants.length; j++) {
        deck.push(
          <Card onClick={handleCardClick} style={{height:"500px", cursor: "pointer"}} key={restaurants[j].id} id={restaurants[j].id}>
            <Card.Img style={{height:"50%"}} variant="top" src={restaurants[j].img || "https://www.opentable.com/img/restimages/150568.jpg"} />
            <Card.Body style={{height:"40%", overflow: "hidden"}}>
              <Card.Title>{restaurants[j].name}</Card.Title>
              <Card.Text>{restaurants[j].description}</Card.Text>
            </Card.Body>
            <Card.Footer style={{height:"10%"}}>
              <small className="text-muted">{restaurants[j].city + ", " + restaurants[j].address}</small>
            </Card.Footer>
          </Card>
        );
      }

      decks.push(<CardDeck key={i} style={{width: deck.length * 33.33 + "%", marginBottom: "20px", marginTop: "20px"}}>{deck[0]}{deck[1]}{deck[2]}</CardDeck>)
    }

    return decks
  }

  return (
    <>
      <Row style={{width: "100%"}}>
        <Col lg={{ span: 3, offset: 9 }}>
          <InputRange
            formatLabel={value => `${value} stars`}
            step={0.5}
            maxValue={5}
            minValue={1}
            value={ratingRange}
            onChangeComplete={(value) => {updateRestaurantsList}}
            onChange={(value) => {setRatingRange(value)}} />
        </Col>
      </Row>
      {getDecks()}
    </>
  )
}