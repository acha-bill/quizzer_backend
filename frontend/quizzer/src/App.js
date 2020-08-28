import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import './App.css';
import Login from "./pages/login/login";
import Signup from "./pages/signup/signup";

class App extends Component {
  constructor(props) {
    super(props);
  }

  render(){
    return (
        <div className={"container-fluid main-wrapper"}>
          <Router>
            <Switch>
                <Route path="/login" component={Login}/>
                <Route path="/signup" component={Signup}/>
                <Route path="/" component={Login}/>
            </Switch>
          </Router>
        </div>
    )
  }
}

export default App;
