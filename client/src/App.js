import React, { useState } from "react";
import Navbar from "react-bootstrap/Navbar";
import "./App.css";
import Routes from "./Routes";
import Nav from "react-bootstrap/Nav";
import { NavHashLink } from 'react-router-hash-link';

function App() {

  const isAuthenticated = false;

  return (
    <div className="App container py-3">
      <Navbar collapseOnSelect bg="light" expand="md" className="mb-3">
        <Navbar.Brand href="/" className="font-weight-bold text-muted">
          Home
        </Navbar.Brand>
        <Navbar.Toggle />
	
	{ !isAuthenticated ? (
	<Navbar.Collapse className="justify-content-end">
          <Nav>
            <Nav.Link href="/signup">Signup</Nav.Link>
            <Nav.Link href="/login">Login</Nav.Link>
          </Nav>
        </Navbar.Collapse>
	) : (
	<Navbar.Collapse className="justify-content-end">
          <Nav>
            <Nav.Link><NavHashLink to="/#Spending">Spending</NavHashLink></Nav.Link>
            <Nav.Link><NavHashLink to="/#Income">Income</NavHashLink></Nav.Link>
            <Nav.Link><NavHashLink to="/#Goals">Goals</NavHashLink></Nav.Link>
	    <Nav.Link>Logged in: <a href="">Paul</a></Nav.Link>
          </Nav>
        </Navbar.Collapse>
	)}
      </Navbar>
      <Routes />

    </div>
  );
}

export default App;
