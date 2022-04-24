import axios from "axios";

// export let url = "https://hellobox.uz/api";
export let url = "http://localhost:8081/api";
// export let url = 'http://192.168.31.87:8081/api';

export let requests = {
  news: {
    createAndSend: (credentails) => axios.post(`${url}/news`, credentails),
  },
  partners: {
    getPartners: () => axios.get(`${url}/partner`),
    addPartner: (credentials) => axios.post(`${url}/partner`, credentials),
    deletePartner: (id) => axios.delete(`${url}/partner/${id}`),
    editPartner: (credentials) => axios.put(`${url}/partner`, credentials),
  },
  products: {
    getProducts: () => axios.get(`${url}/product`),
    createProduct: (credentails) => axios.post(`${url}/product`, credentails),
    deleteProduct: (id) => axios.delete(`${url}/product/${id}`),
    editProduct: (credentails) => axios.put(`${url}/product`, credentails),
    updateProfitPercent: (percent) =>
      axios.put(`${url}/profit-percent`, { percent, id: 1 }),
    getProfitPercent: () => axios.get(`${url}/profit-percent`),
  },
  order: {
    getOrder: () => axios.get(`${url}/order`),
    createOrder: (credentails) => axios.post(`${url}/order`, credentails),
  },
  uploads: {
    uploadImage: (credentails) => axios.post(`${url}/file-upload`, credentails),
    updatePresentImage: (creds) =>
      axios.post(`${url}/update-present-image`, creds),
  },
  users: {
    getUsers: () => axios.get(`${url}/users`),
    getPartnerUsers: () => axios.get(`${url}/users?fromPartner=true`),
    editUser: (payload) => axios.put(`${url}/user`, payload),
  },
  category: {
    getCategories: () => axios.get(`${url}/category`),
    deleteCategory: (id) => axios.delete(`${url}/category/${id}`),
    createCategory: (credentails) => axios.post(`${url}/category`, credentails),
    editCategory: (credentails) => axios.put(`${url}/category`, credentails),
  },
  order: {
    getOrders: () => axios.get(`${url}/orders`),
  },
};
