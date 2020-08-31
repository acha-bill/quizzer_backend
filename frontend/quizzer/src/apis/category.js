
const baseUrl = '/category'
export default class CategoryApi {
    constructor(api) {
        this.api = api
    }
    async getAll(obj) {
        try{
            let res = await this.api.get(`${baseUrl}/`)
            return res.data
        }catch(e){
            throw  e
        }
    }
    async edit(id, obj) {
        try{
            let res = await this.api.put(`${baseUrl}/${id}`, obj)
            return res.data
        }catch(e){
            throw  e
        }
    }
    async create(obj) {
        try{
            let res = await this.api.post(`${baseUrl}/`, obj)
            return res.data
        }catch(e){
            throw  e
        }
    }

    async delete(id) {
        try{
            let res = await this.api.delete(`${baseUrl}/${id}`)
            return res.data
        }catch(e){
            throw  e
        }
    }
}