package handlers

import (
	"Project/models"
	"Project/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SLHandler struct {
	service *services.SLService
}

func NewSLHandler(service *services.SLService) *SLHandler {
	return &SLHandler{service: service}
}

func (h *SLHandler) GetAllSLs(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Fetch DLs from the service
	dls, err := h.service.GetAllSLs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the response as JSON
	if err := json.NewEncoder(w).Encode(dls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *SLHandler) CreateSL(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Code  string `json:"code"`
		Title string `json:"title"`
		HasDL bool   `json:"hasdl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateSL(input.Code, input.Title, input.HasDL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (h *SLHandler) GetSLByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	sl, err := h.service.GetSLByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sl)
}

func (h *SLHandler) UpdateSL(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var input models.SL
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.UpdateSL(uint(id), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (h *SLHandler) DeleteSL(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := h.service.DeleteSL(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
