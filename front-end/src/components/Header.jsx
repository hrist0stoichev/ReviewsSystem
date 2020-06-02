import React, {useEffect, useState} from "react";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import {authenticationService} from "../services/auth";

export default function Header(props) {
  const [user, setUser] = useState(null)

  useEffect(() => {
    authenticationService.currentUser.subscribe(x => setUser(x))
  }, [])

  const logOut = () => {
    authenticationService.logout()
    props.showAlert("You logged out", true)
  }

  return (
    <Navbar bg="dark" variant="dark" style={{ marginBottom: "2em"}}>
      <Navbar.Brand href="#home">ReviewsSystem</Navbar.Brand>
      <Nav className="mr-auto">
        {user && <Nav.Link href="#restaurants">Restaurants</Nav.Link>}
      </Nav>
      <Nav className="justify-content-end" activeKey="/home">
        {user === null && <Nav.Link href="#register">Register</Nav.Link>}
        {user === null && <Nav.Link href="#login">Login</Nav.Link>}
        {user && <Nav.Link>{user.email}</Nav.Link>}
        {user && <Nav.Link onClick={logOut} href="#login">Logout</Nav.Link>}
      </Nav>
    </Navbar>
  )
}