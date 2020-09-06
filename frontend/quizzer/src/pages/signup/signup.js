import {Component} from "react";
import React from "react";
import { Link } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import './signup.css'

class Signup extends Component {
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            password: ""
        }
    }

    render() {
        return (
            <div className={"main-container row h-100 justify-content-center align-items-center"}>
                <div className={"col-md-3 col-sm-6"}>
                    <div className={"text-center"}>
                        <span className={"h3"}>Sign up to Quizzer</span>
                    </div>
                    <div className={"login-form m-2 p-2"}>
                        <div className={"form-group"}>
                            <label>Username</label>
                            <input type="text" className={"form-control"}/>
                        </div>
                        <div className={"form-group"}>
                            <label>Password</label>
                            <input type="password" className={"form-control"}/>
                        </div>
                        <div className={"form-group"}>
                            <label>Password</label>
                            <input type="password" className={"form-control"}/>
                        </div>
                        <div>
                            <button className={"btn btn-success btn-block"}>Sign up</button>
                        </div>
                        <div className={"mt-2"}>
                            <span>Already have an account? <Link to="/login">Login</Link></span>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

export default Signup