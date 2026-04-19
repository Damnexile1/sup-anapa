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

type InstructorHandler struct {
	repo *repository.InstructorRepository
}

func NewInstructorHandler(repo *repository.InstructorRepository) *InstructorHandler {
	return &InstructorHandler{repo: repo}
}

func (h *InstructorHandler) List(w http.ResponseWriter, r *http.Request) {
	instructors, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting instructors: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Проверить, нужен ли HTML ответ (для HTMX)
	if r.Header.Get("HX-Request") == "true" {
		h.renderInstructorsList(w, instructors)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructors)
}

func (h *InstructorHandler) renderInstructorsList(w http.ResponseWriter, instructors []*models.Instructor) {
	if len(instructors) == 0 {
		fmt.Fprint(w, `<div class="bg-white rounded-lg shadow-lg p-6 text-center col-span-full">
			<p class="text-gray-500">Инструкторы не добавлены</p>
		</div>`)
		return
	}

	for _, instructor := range instructors {
		photoURL := instructor.Photo
		if photoURL == "" {
			photoURL = "https://via.placeholder.com/300x300?text=Instructor"
		}

		fmt.Fprintf(w, `
		<div class="bg-white rounded-lg shadow-lg overflow-hidden">
			<img src="%s" alt="%s" class="w-full h-48 object-cover" onerror="this.src='https://via.placeholder.com/300x300?text=Instructor'">
			<div class="p-6">
				<h3 class="text-xl font-semibold mb-2">%s</h3>
				<p class="text-gray-600 mb-2">%s</p>
				<p class="text-gray-500 text-sm mb-4">%s</p>

				<div class="mb-4">
					<div class="flex items-center justify-between mb-2">
						<p class="font-semibold text-sm">Типы прогулок</p>
						<button onclick="addWalkType(%d)" class="text-blue-600 text-sm">+ Добавить</button>
					</div>
					<div id="walk-types-%d" hx-get="/admin/api/instructors/%d/walk-types" hx-trigger="load" hx-swap="innerHTML">
						<p class="text-sm text-gray-500">Загрузка...</p>
					</div>
				</div>

				<div class="flex gap-2">
					<button onclick="editInstructor(%d)" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded">Редактировать</button>
					<button hx-delete="/admin/api/instructors/%d" hx-confirm="Удалить инструктора?" hx-target="#instructors-list" hx-swap="innerHTML" class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded">Удалить</button>
				</div>
			</div>
		</div>`, photoURL, instructor.Name, instructor.Name, instructor.Phone, instructor.Description, instructor.ID, instructor.ID, instructor.ID, instructor.ID, instructor.ID)
	}
}

func (h *InstructorHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	instructor, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting instructor: %v", err)
		http.Error(w, "Instructor not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructor)
}

func (h *InstructorHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	instructor := models.Instructor{
		Name:        r.FormValue("name"),
		Phone:       r.FormValue("phone"),
		Description: r.FormValue("description"),
		Photo:       r.FormValue("photo_url"),
	}

	if err := h.repo.Create(r.Context(), &instructor); err != nil {
		log.Printf("Error creating instructor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	instructors, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting instructors: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderInstructorsList(w, instructors)
}

func (h *InstructorHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var instructor models.Instructor
	if err := json.NewDecoder(r.Body).Decode(&instructor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	instructor.ID = id
	if err := h.repo.Update(r.Context(), &instructor); err != nil {
		log.Printf("Error updating instructor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Если это HTMX запрос, вернуть обновленный список
	if r.Header.Get("HX-Request") == "true" {
		instructors, err := h.repo.GetAll(r.Context())
		if err != nil {
			log.Printf("Error getting instructors: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.renderInstructorsList(w, instructors)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructor)
}

func (h *InstructorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		log.Printf("Error deleting instructor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	instructors, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting instructors: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderInstructorsList(w, instructors)
}
