package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"

	"github.com/go-chi/chi/v5"
)

type WalkTypeHandler struct {
	repo *repository.WalkTypeRepository
}

func NewWalkTypeHandler(repo *repository.WalkTypeRepository) *WalkTypeHandler {
	return &WalkTypeHandler{repo: repo}
}

func (h *WalkTypeHandler) ListByInstructor(w http.ResponseWriter, r *http.Request) {
	instructorID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid instructor ID", http.StatusBadRequest)
		return
	}

	walkTypes, err := h.repo.GetByInstructorID(r.Context(), instructorID)
	if err != nil {
		log.Printf("Error getting walk types: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		h.renderList(w, walkTypes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(walkTypes)
}

func (h *WalkTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	instructorID, _ := strconv.Atoi(r.FormValue("instructor_id"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	maxPeople, _ := strconv.Atoi(r.FormValue("max_people"))

	if instructorID < 1 || r.FormValue("name") == "" || price < 1 || maxPeople < 1 {
		http.Error(w, "Заполните все поля корректно", http.StatusBadRequest)
		return
	}

	walkType := &models.WalkType{
		InstructorID: instructorID,
		Name:         r.FormValue("name"),
		Price:        price,
		MaxPeople:    maxPeople,
	}

	if err := h.repo.Create(r.Context(), walkType); err != nil {
		log.Printf("Error creating walk type: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	walkTypes, err := h.repo.GetByInstructorID(r.Context(), instructorID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.renderList(w, walkTypes)
}

func (h *WalkTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	instructorID, _ := strconv.Atoi(r.URL.Query().Get("instructor_id"))

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	walkTypes, err := h.repo.GetByInstructorID(r.Context(), instructorID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.renderList(w, walkTypes)
}

func (h *WalkTypeHandler) renderList(w http.ResponseWriter, walkTypes []*models.WalkType) {
	if len(walkTypes) == 0 {
		fmt.Fprint(w, `<p class="text-sm text-gray-500">Типы прогулок не добавлены</p>`)
		return
	}
	for _, wt := range walkTypes {
		fmt.Fprintf(w, `<div class="flex items-center justify-between gap-2 border rounded p-2 mb-2">
			<div>
				<p class="font-medium">%s</p>
				<p class="text-sm text-gray-600">%d ₽ • до %d чел.</p>
			</div>
			<button hx-delete="/admin/api/walk-types/%d?instructor_id=%d" hx-target="#walk-types-%d" hx-swap="innerHTML" class="text-red-600 text-sm">Удалить</button>
		</div>`, wt.Name, wt.Price, wt.MaxPeople, wt.ID, wt.InstructorID, wt.InstructorID)
	}
}
