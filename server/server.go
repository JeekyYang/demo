package server

import (
	"demo/controller"
	"demo/model"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
)

type Server struct {
	m      *model.Model
	router *mux.Router
}

//register router
func (s *Server) register() {
	s.router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	s.router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/:user_id/relationships", controller.GetAllRelationShipByUid).Methods("GET")
	s.router.HandleFunc("/users/:user_id/relationships/:other_user_id", controller.UpdateRelationShip).Methods("PUT")
}

func NewServer() (*Server, error) {
	m := model.NewDB()
	r := mux.NewRouter()

	s := &Server{m, r}

	s.register()
	return s, nil
}

func (s *Server) Start() {
	//todo: handler panic
	fmt.Println("begin to start server")
	log.Fatal(http.ListenAndServe(":8088", s.router))
}

func (s *Server) Close() {
	//todo: graceful shutdown??
}
