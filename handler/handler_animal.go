package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vault-app/domain"

	"github.com/gorilla/mux"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type handlerAnimal struct {
	IUsecaseAnimal domain.AnimalUseCase
}

func NewHandlerAnimal(r *mux.Router, u *domain.AnimalUseCase) {
	h := handlerAnimal{
		IUsecaseAnimal: *u,
	}

	r.HandleFunc("/animals", h.Store()).Methods("POST")
	r.HandleFunc("/animals/{id}", h.GetByID()).Methods("GET")
	r.HandleFunc("/ping", h.Ping()).Methods("GET")

}

func (h *handlerAnimal) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("content-type", "application/json")

		vars := mux.Vars(r)

		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "ID not found in path", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		response, err := h.IUsecaseAnimal.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Invalid to get response", http.StatusInternalServerError)
			return
		}

	}
}

func (h *handlerAnimal) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "ping",
		})

	}
}

func (h *handlerAnimal) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("content-type", "application/json")

		var req *domain.Animal

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid Store", http.StatusInternalServerError)
			return
		}

		err = h.IUsecaseAnimal.Store(r.Context(), req)
		if err != nil {
			http.Error(w, "Invalid Store to Usecase", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		response := map[string]interface{}{
			"name": req.Name,
			"age":  req.Age,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Invalid to get response", http.StatusInternalServerError)
			return
		}

	}
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
