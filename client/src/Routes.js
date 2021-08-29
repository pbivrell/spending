import React from "react";
import { Route, Switch, Redirect, useLocation } from "react-router-dom";
import Home from "./containers/Home";
import Login from "./containers/Login";
import NotFound from "./containers/NotFound";
import { useAppContext } from "./libs/contextLib";

export default function Routes() {

  const { tokenHolder } = useAppContext();

  var loc = new URLSearchParams(useLocation().search).get("token");
  return (
    <Switch>
      <Route exact path="/">
        <Home />
      </Route>
      <Route exact path="/redirect" render={() => {
	     if (loc !== null) {
		tokenHolder.setToken(loc);
	     }
	     return <Redirect to="/"/>
      }}/>
      <Route>
        <NotFound />
      </Route>
    </Switch>
  );
}

