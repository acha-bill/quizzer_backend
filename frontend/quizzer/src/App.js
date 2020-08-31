import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import './App.css';
import Login from "./pages/login/login";
import Signup from "./pages/signup/signup";
import Category from "./pages/admin/category/category";

class App extends Component {

  render(){
    return (
        <div className={"container-fluid main-wrapper"}>
          <Router>
            <Switch>
                <Route path="/categories" component={Category}/>
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
