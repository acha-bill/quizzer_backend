import { categoryActionType } from "../actions/category";

export const defaultState = {
    categories: []
}

const categories = (state = defaultState, action) => {
    switch(action.type){
        case categoryActionType.setCategories:
            state = { ...state, categories: action.categories };
            return state
        case categoryActionType.addCategory:
            let categories = [...state.categories].push(action.category)
            state = {...state, categories}
            return state
        case categoryActionType.updateCategory:
                let i = state.categories.findIndex(cat => cat.id === action.id)
                if(i >= 0){
                    state.categories[i] = action.category
                }
                state = {...state,categories: state.categories}
                return state
            case categoryActionType.deleteCategory:
                i = state.categories.findIndex(cat => cat.id === action.id)
                if (i >= 0){
                    categories = [...state.categories]
                    categories.splice(i,1)
                    state = {...state, categories}
                }
                return state
        default:
            return state
    }
}
export default categories