import React, { useState } from "react";
import Navbar from "react-bootstrap/Navbar";
import "./App.css";
import Routes from "./Routes";
import Nav from "react-bootstrap/Nav";
import { NavHashLink } from 'react-router-hash-link';
import inMemoryJWT from './components/inMemoryJWT';
import { AppContext } from "./libs/contextLib";
import Button from "react-bootstrap/Button";
import {useAppContext} from "./libs/contextLib.js";
import jwt_decode from "jwt-decode";


export default function LoggedNav() {
  const redirect = "https://spend-it-gatekeeper.herokuapp.com/?redirect.html=spend-it.herokuapp.com/redirect"
  //const redirect = "http://localhost:8080/redirect.html?redirect=localhost:3001/redirect"

  const { tokenHolder } = useAppContext();

  const [loggedIn, setLoggedIn] = useState(tokenHolder.getToken());

  let token = "";

  if(!loggedIn) {
	tokenHolder.waitForTokenRefresh().then(() => {
		setLoggedIn(tokenHolder.getToken())
	});
  }else {
  	token = jwt_decode(loggedIn);
  	console.log("TOKEN DATA", token)
  }


  return (
	<>
	{ loggedIn === null  ? (
          <Nav>
            <Nav.Link href={redirect}>Login</Nav.Link>
          </Nav>
	) : (
		<>
	<Navbar.Collapse className="justify-content-end">
          <Nav>
            <Nav.Link><NavHashLink to="/#Spending">Spending</NavHashLink></Nav.Link>
            <Nav.Link><NavHashLink to="/#Income">Income</NavHashLink></Nav.Link>
            <Nav.Link><NavHashLink to="/#Goals">Goals</NavHashLink></Nav.Link>
	<Button variant="" onClick={() => {tokenHolder.ereaseToken()}}>Logout</Button>
          </Nav>
        </Navbar.Collapse>
	<Nav.Link>Logged in: <a href="">{token.username}</a></Nav.Link>
        <Navbar.Toggle />
		</>
	)}
	</>
  );
}
