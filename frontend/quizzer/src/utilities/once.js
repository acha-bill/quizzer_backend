export default class Once {
    constructor() {
        this.results = new Map()
    }
    do(func){
        let k = func.toString()
        let val = this.results.get(k)
        if(val){
            return val
        }
        val = func()
        this.results.set(k,val)
        return val
    }
}