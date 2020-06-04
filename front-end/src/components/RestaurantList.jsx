import React, {useEffect, useState, useRef} from "react";
import {restaurantsService} from "../services/restaurants";
import CardDeck from "react-bootstrap/CardDeck";
import Card from "react-bootstrap/Card";
import InputRange from 'react-input-range';
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Badge from "react-bootstrap/Badge";

export default function RestaurantList(props) {
  const defaultPageSize = 9;
  const defaultOrderBy = "average_rating";

  const [restaurants, setRestaurants] = useState([]);
  const [ratingRange, setRatingRange] = useState({min: 0, max: 5});

  const page = useRef(0);
  const hasMoreResults = useRef(true);
  const isUpdating = useRef(false);

  useEffect(() => {
    updateRestaurantsList(false);
    window.addEventListener('scroll', listenToScroll, { passive: true });

    return () => {
      window.removeEventListener('scroll', listenToScroll);
    };
  }, []);

  const listenToScroll = () => {
    const winScroll = document.body.scrollTop || document.documentElement.scrollTop;
    const height = document.documentElement.scrollHeight - document.documentElement.clientHeight;

    if (!isUpdating.current && hasMoreResults.current && winScroll / height > 0.6) {
      isUpdating.current = true;
      page.current++;
      updateRestaurantsList(true);
    }
  };

  const updateRestaurantsList = (append) => {
    restaurantsService.get(defaultPageSize, page.current * defaultPageSize, ratingRange.min, ratingRange.max, defaultOrderBy)
      .then(res => {
        setRestaurants(restaurants => {
          return append ? [...restaurants, ...res] : res
        })

        if (res.length < defaultPageSize) {
          hasMoreResults.current = false;
        }

        isUpdating.current = false;
      })
      .catch(err => {
        props.showAlert(err, false);
      });
  };

  const handleCardClick = (event) => {
    props.history.push("/restaurants/" + event.currentTarget.id);
  };

  const getDecks = () => {
    const decks = [];
    for (let i = 0; i < restaurants.length; i += 3) {
      const deck = [];
      for (let j = i; j < i + 3 && j < restaurants.length; j++) {
        deck.push(
          <Card onClick={handleCardClick} style={{height:"500px", cursor: "pointer"}} key={restaurants[j].id} id={restaurants[j].id}>
            <Card.Img style={{height:"50%"}} variant="top" src={restaurants[j].img || "https://www.opentable.com/img/restimages/150568.jpg"} />
            <Card.Body style={{height:"40%", overflow: "hidden"}}>
              <Card.Title><Badge pill variant="primary">{restaurants[j].average_rating} <span>&#9733;</span></Badge> {restaurants[j].name}</Card.Title>
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

  const handleFilterChanged = () => {
    page.current = 0;
    hasMoreResults.current = true;
    isUpdating.current = true;
    updateRestaurantsList(false)
  };

  return (
    <>
      <Row style={{width: "100%", marginBottom: "30px"}}>
        <Col lg={{ span: 3, offset: 9 }}>
          <InputRange
            formatLabel={value => `${value} \u2B50`}
            allowSameValues={true}
            step={0.5}
            maxValue={5}
            minValue={0}
            value={ratingRange}
            onChangeComplete={handleFilterChanged}
            onChange={(value) => {setRatingRange(value)}} />
        </Col>
      </Row>
      {getDecks()}
    </>
  )
}