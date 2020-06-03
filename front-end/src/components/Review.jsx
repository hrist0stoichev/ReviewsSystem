import React from "react";

export default function Review(props) {
  return (
    <div style={{textAlign: "left"}}>
      <h5>{props.review.reviewer}</h5>
      <div>My rating: {props.review.rating} <span>&#9733;</span></div>
      <div style={{marginTop: "10px"}}><i>{props.review.comment}</i></div>
      <div ><small>{new Intl.DateTimeFormat('en-US', {year: 'numeric', month: '2-digit',day: '2-digit', hour: '2-digit', minute: '2-digit'}).format(Date.parse(props.review.timestamp))}</small></div>
    </div>
  )
}