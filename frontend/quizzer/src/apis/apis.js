
import AuthApi from "./auth";
import {getApi} from './axios'
import Once from '../utilities/once'
import CategoryApi from "./category";

let once = new Once()
class Apis {
    initialize(token) {
        this.token = token
        this.api = getApi(this.token)
    }
    auth(){
        return once.do(() => new AuthApi(this.api))
    }
    categories(){
        return once.do(() => new CategoryApi(this.api))
    }
}

let apis = new Apis()
export default apis