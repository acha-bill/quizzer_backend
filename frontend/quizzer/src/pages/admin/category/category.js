import { Component } from "react";
import React from "react";
import { connect } from 'react-redux';
import Modal from 'react-modal';

import { getAllCategories, createCategory, deleteCategory, editCategory } from '../../../redux/actions/category'
import Swal from "sweetalert2";

const customStyles = {
    content: {
        top: '50%',
        left: '50%',
        right: 'auto',
        bottom: 'auto',
        marginRight: '-50%',
        transform: 'translate(-50%, -50%)'
    }
};

Modal.setAppElement('#root')

class Category extends Component {
    constructor(props) {
        super(props)
        this.state = {
            isModalOpen: false,
            editingCategory: null,
            name: '',
            modalTitle: ''
        }
    }

    openModal = () => {
        this.setState({
            isModalOpen: true
        })
    }
    closeModal = () => {
        this.setState({
            isModalOpen: false
        })
    }

    onEditClicked = (cat) => {
        this.setState({
            modalTitle: 'Edit',
            editingCategory: cat
        })
        this.openModal()
    }

    onCancleClick = (e) => {
        e.preventDefault()
        this.closeModal()
        this.setState({
            name: ''
        })
    }
    handleInput = (e) => {
        this.setState({
            [e.target.name] : e.target.value
        })
    }

    handleSubmit = (e) => {
        e.preventDefault()
        const {name} = this.state
        const {createCategory, editCategory} = this.props

        let cat = this.state.editingCategory
        if(cat){
            cat.name = name
            editCategory(cat.id, cat)
            return 
        }
        cat = {
            name
        }
        createCategory(cat)
        Swal.fire({
            icon: 'success',
            title: 'success'
        })
    }

    render() {
        let { categories } = this.props
        let { name, modalTitle } = this.state

        categories = [
            { name: 'programing' },
            { name: 'action' }
        ]

        return (
            <div className={"container mt-5"}>
                <Modal
                    isOpen={this.state.isModalOpen}
                    onRequestClose={this.closeModal}
                    style={customStyles}
                    contentLabel="Edit category"
                >
                    <h5>{modalTitle}</h5>
                    <form className={"m-10"} onSubmit={this.handleSubmit}>
                        <div className={"form-group"}>
                            <label className="text-center">Name</label>
                            <input className={"form-control"} type="text" name="name" value={name} onChange={this.handleInput}/>
                        </div>
                        <div className={"d-flex justify-content-center align-items-center"}>
                            <input type="submit" value="submit" className={"btn btn-success"}/> &nbsp;
                            <button className={"btn btn-danger"} onClick={this.onCancleClick}>Cancel</button>
                        </div>
                    </form>
                </Modal>
                <h3 className="text-center">Categories</h3>
                <button className={"btn btn-primary align-items-flex-end justify-content-flex-end"}>Add</button>
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
                                <td>{cat.name}</td>
                                <td>
                                    <button className={"btn btn-sm btn-warning"} onClick={this.onEditClicked}>Edit</button> &nbsp;
                                    <button className={"btn btn-sm btn-danger"}>Delete</button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        categories: state.categories.categories
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        getAllCategories,
        createCategory,
        editCategory,
        deleteCategory
    }
};

export default connect(mapStateToProps, mapDispatchToProps)(Category);