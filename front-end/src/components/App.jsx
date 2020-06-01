import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';

export default class App extends Component {
  constructor() {
    super();
  }

  render() {
    return <Button variant="danger">It works</Button>
  }
}