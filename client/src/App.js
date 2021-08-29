import React, { useState } from "react";
import Navbar from "react-bootstrap/Navbar";
import "./App.css";
import Routes from "./Routes";
import Nav from "react-bootstrap/Nav";
import { NavHashLink } from 'react-router-hash-link';
import inMemoryJWT from './components/inMemoryJWT';
import { AppContext } from "./libs/contextLib";
import Button from "react-bootstrap/Button";
import LoggedNav from "./Nav";

function App() {

  const [tokenHolder] = useState(inMemoryJWT);
 
  if(!tokenHolder.getToken()) {
	tokenHolder.getRefreshedToken();
  }

  return (
    <div className="App container py-3">
      <Navbar collapseOnSelect bg="light" expand="md" className="mb-3">
        <Navbar.Brand href="/" className="font-weight-bold text-muted">
          Home
        </Navbar.Brand>
      <AppContext.Provider value={{ tokenHolder }}>
	<LoggedNav/>
      </AppContext.Provider>
      </Navbar>
      <AppContext.Provider value={{ tokenHolder }}>
      	<Routes/>
      </AppContext.Provider>
      


    </div>
  );
}

export default App;
