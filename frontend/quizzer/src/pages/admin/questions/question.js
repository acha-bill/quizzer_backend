import { Component } from "react";
import React from "react";
import { connect } from 'react-redux';

import { getAllCategories } from '../../../redux/actions/category'


class Question extends Component {
    constructor(props) {
        super(props)
        this.state = {

        }
    }

    render() {
        let { categories } = this.props

        return (
            <div className={"row h-100"}>
                <div className={"col-md-4"}>
                    <h3>Categories</h3>
                    <table className={"table table-bordered"}>
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Name</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {categories.map((cat, i) => (
                                <tr key={i}>
                                    <td>{i + 1}</td>
                                    <td>cat.name</td>
                                    <td>
                                        <button className={"btn btn-primary"}>Edit</button>
                                        <button className={"btn btn-danger"}>Delete</button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        categories: state.categories,
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        getAllCategories
    }
};

export default connect(mapStateToProps, mapDispatchToProps)(Question);