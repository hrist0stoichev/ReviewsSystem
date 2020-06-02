import React, {useEffect, useState} from "react";
import {restaurantsService} from "../services/restaurants";
import CardDeck from "react-bootstrap/CardDeck";
import Card from "react-bootstrap/Card";

export default function RestaurantList(props) {
  const [restaurants, setRestaurants] = useState([])

  useEffect(() => {
    restaurantsService.get( 20, 0, 1, 5, "average_rating")
      .then(res => {
        setRestaurants(res)
      })
      .catch(err => {
        props.showAlert(err, false)
      })
  }, []);

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
            <Card.Img style={{height:"50%"}} variant="top" src={restaurants[j].img} />
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
    <>{getDecks()}</>
  )
}