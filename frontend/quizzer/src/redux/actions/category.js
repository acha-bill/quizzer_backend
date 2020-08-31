import apis from "../../apis/apis";

export const categoryActionType = {
    setCategories: 'set_categories',
    updateCategory: 'update_category',
    addCategory: 'add_category',
    deleteCategory: 'delete_category',
    loading: 'loading'
}

function loading(bool) {
  return {
      type: categoryActionType.loading,
      isLoading: bool
  };
}

function setCategories(cats) {
  return {
    type: categoryActionType.setCategories,
    categories: cats
  }
}

function updateCategory(id, cat) {
  return {
    type: categoryActionType.updateCategory,
    category: cat,
    id: id
  }
}

function addCategory(cat) {
  return {
    type: categoryActionType.addCategory,
    category: cat,
  }
}

function _deleteCategory(id) {
  return {
    type: categoryActionType.addCategory,
    id,
  }
}

export function getAllCategories() {
  return async (dispatch) => {
    try{
      let res = await apis.categories().getAll()
      dispatch(setCategories(res))
    } catch (e) {
      throw e
    }finally {
      dispatch(loading(false))
    }    
  }
}

export function editCategory(id, cat) {
  return async(dispatch) => {
    try {
      let res = await apis.categories().edit(id, cat)
      dispatch(updateCategory(id, res))
    }catch(e){
      throw e
    }
  }
}

export function createCategory(cat) {
  return async(dispatch) => {
    try {
      let res = await apis.categories().create(cat)
      dispatch(addCategory(res))
    }catch(e){
      throw e
    }
  }
}

export function deleteCategory(id) {
  return async(dispatch) => {
    try {
      await apis.categories().delete(id)
      dispatch(_deleteCategory(id))
    }catch(e){
      throw e
    }
  }
}