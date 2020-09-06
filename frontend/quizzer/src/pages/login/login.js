import {Component} from "react";
import React from "react";
import { Link } from 'react-router-dom';
import Swal from 'sweetalert2'

import apis from '../../apis/apis'

import 'bootstrap/dist/css/bootstrap.min.css';
import './login.css'


class Login extends Component {
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            password: ""
        }
    }

    handleInput = event => {
        this.setState({
            [event.target.name]: event.target.value
        })
    }

    handleSubmit = async event => {
        event.preventDefault()

        const {username, password} = this.state
        apis.initialize('')
        try{
            let token = await apis.auth().login({
                username,
                password
            })
            console.log(token)
            apis.initialize(token)
            window.location.pathname = "/categories"
        }catch (e) {
            Swal.fire({
                icon: 'error',
                title: 'error',
                text: e.message
            })
        }
    }

    render() {
        let {username,password} = this.state
        return (
            <div className={"main-container row h-100 justify-content-center align-items-center"}>
                <div className={"col-md-3 col-sm-6"}>
                    <div className={"text-center"}>
                        <span className={"h3"}>Sign in to Quizzer</span>
                    </div>
                    <form className={"login-form m-2 p-2"} onSubmit={this.handleSubmit}>
                        <div className={"form-group"}>
                            <label>Username</label>
                            <input name="username" value={username} onChange={this.handleInput} type="text" className={"form-control"}/>
                        </div>
                        <div className={"form-group"}>
                            <label>Password</label>
                            <input name="password"  value={password} onChange={this.handleInput} type="password" className={"form-control"}/>
                        </div>
                        <div>
                            <input type="submit" className={"btn btn-success btn-block"} value="Login"/>
                        </div>
                        <div className={"mt-2"}>
                            <span>Don't have an account? <Link to="/signup">Sign up</Link></span>
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}

export default Login