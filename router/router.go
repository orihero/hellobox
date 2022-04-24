package router

import (
	controller "hellobox/controllers"
	"hellobox/env"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRouter() {
	env.Router = mux.NewRouter()
	prefix := env.Router.PathPrefix("/api/").Subrouter()
	env.Router.PathPrefix("/manage").Handler(http.StripPrefix("/manage", http.FileServer(http.Dir("./admin"))))
	prefix.HandleFunc("/update-present-image", controller.UpdatePresentImage).Methods("POST")
	prefix.HandleFunc("/get-present-image", controller.GetPresentImage).Methods("GET")
	env.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./admin/static"))))

	//******CATEGORY
	prefix.HandleFunc("/category", controller.GetCategories).Methods("GET")
	prefix.HandleFunc("/category/{id}", controller.DeleteCategory).Methods("DELETE")
	prefix.HandleFunc("/category", controller.EditCategory).Methods("PUT")
	prefix.HandleFunc("/category", controller.CreateCategory).Methods("POST")

	//*******PRODUCT
	prefix.HandleFunc("/product", controller.GetProducts).Methods("GET")
	prefix.HandleFunc("/product/{id}", controller.DeleteProduct).Methods("DELETE")
	prefix.HandleFunc("/product", controller.EditProduct).Methods("PUT")
	prefix.HandleFunc("/product", controller.CreateProduct).Methods("POST")

	//******USER
	prefix.HandleFunc("/users", controller.GetUsers).Methods("GET")
	prefix.HandleFunc("/user/{id}", controller.DeleteUser).Methods("DELETE")
	prefix.HandleFunc("/user", controller.EditUser).Methods("PUT")
	prefix.HandleFunc("/users", controller.CreateUser).Methods("POST")

	//******ORDERS
	prefix.HandleFunc("/orders", controller.GetOrders).Methods("GET")

	//******NEWS
	prefix.HandleFunc("/news", controller.GetNews).Methods("GET")
	prefix.HandleFunc("/news/{id}", controller.DeleteNews).Methods("DELETE")
	prefix.HandleFunc("/news", controller.EditNews).Methods("PUT")
	prefix.HandleFunc("/news", controller.CreateNews).Methods("POST")

	//******PARTNER
	prefix.HandleFunc("/partner", controller.GetPartner).Methods("GET")
	prefix.HandleFunc("/partner/{id}", controller.DeletePartner).Methods("DELETE")
	prefix.HandleFunc("/partner", controller.EditPartner).Methods("PUT")
	prefix.HandleFunc("/partner", controller.CreatePartner).Methods("POST")

	//******SETTINGS
	prefix.HandleFunc("/settings", controller.GetSettings).Methods("GET")
	prefix.HandleFunc("/settings/{id}", controller.DeleteSettings).Methods("DELETE")
	prefix.HandleFunc("/settings", controller.EditSettings).Methods("PUT")
	prefix.HandleFunc("/settings", controller.CreateSettings).Methods("POST")

	prefix.HandleFunc("/profit-percent", controller.UpdateProfitPercent).Methods("PUT")
	prefix.HandleFunc("/profit-percent", controller.GetProfitPercent).Methods("GET")

	//*******UPLOAD/DOWNLOAD
	prefix.HandleFunc("/file-upload", controller.MultipleFileUpload).Methods("POST")
	prefix.HandleFunc("/file-download/{name}", controller.GetUploadedFiles).Methods("GET")

	env.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
}
