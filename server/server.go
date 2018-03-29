package server

import (
	"demo/controller"
	"demo/model"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	m      *model.Model
	router *mux.Router
}

//register router
func (s *Server) register() {
	s.router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	s.router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/{user_id}/relationships", controller.GetAllRelationShipByUid).Methods("GET")
	s.router.HandleFunc("/users/{user_id}/relationships/{other_user_id}", controller.UpdateRelationShip).Methods("PUT")
}

func NewServer() (*Server, error) {
	m := model.NewDB()
	r := mux.NewRouter()

	s := &Server{m, r}

	s.register()
	return s, nil
}

func (s *Server) Start() {
	fmt.Println("begin to start server")
	log.Fatal(http.ListenAndServe(":80", s.router))
}

//todo: graceful shutdown??
func (s *Server) Close() {
}
