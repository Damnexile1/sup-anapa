package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"
	"time"

	"github.com/go-chi/chi/v5"
)

type SlotHandler struct {
	repo           *repository.SlotRepository
	instructorRepo *repository.InstructorRepository
}

func NewSlotHandler(repo *repository.SlotRepository, instructorRepo *repository.InstructorRepository) *SlotHandler {
	return &SlotHandler{
		repo:           repo,
		instructorRepo: instructorRepo,
	}
}

func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request) {
	slots, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting slots: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Проверить, нужен ли HTML ответ (для HTMX)
	if r.Header.Get("HX-Request") == "true" {
		h.renderSlotsList(w, slots)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

func (h *SlotHandler) renderSlotsList(w http.ResponseWriter, slots []*models.Slot) {
	if len(slots) == 0 {
		fmt.Fprint(w, `<tr><td colspan="7" class="px-6 py-4 text-center text-gray-500">Слоты не добавлены</td></tr>`)
		return
	}

	for _, slot := range slots {
		// Получить имя инструктора
		instructorName := fmt.Sprintf("Инструктор #%d", slot.InstructorID)
		instructor, err := h.instructorRepo.GetByID(context.Background(), slot.InstructorID)
		if err == nil {
			instructorName = instructor.Name
		}

		// Убрать секунды из времени
		startTime := slot.StartTime
		if len(startTime) > 5 {
			startTime = startTime[:5]
		}
		endTime := slot.EndTime
		if len(endTime) > 5 {
			endTime = endTime[:5]
		}

		fmt.Fprintf(w, `
		<tr class="border-b hover:bg-gray-50">
			<td class="px-6 py-4">%s</td>
			<td class="px-6 py-4">%s - %s</td>
			<td class="px-6 py-4">%d ₽</td>
			<td class="px-6 py-4">%d чел.</td>
			<td class="px-6 py-4">%s</td>
			<td class="px-6 py-4">
				<div class="flex gap-2">
					<button onclick="editSlot(%d)" class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-sm">
						Редактировать
					</button>
					<button hx-delete="/admin/api/slots/%d" hx-confirm="Удалить слот?" hx-target="#slots-table-body" hx-swap="innerHTML" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm">
						Удалить
					</button>
				</div>
			</td>
		</tr>`, slot.Date.Format("02.01.2006"), startTime, endTime, slot.Price, slot.MaxPeople, instructorName, slot.ID, slot.ID)
	}
}

func (h *SlotHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	slot, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting slot: %v", err)
		http.Error(w, "Slot not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slot)
}

func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	price, _ := strconv.Atoi(r.FormValue("price"))
	maxPeople, _ := strconv.Atoi(r.FormValue("max_people"))
	instructorID, _ := strconv.Atoi(r.FormValue("instructor_id"))

	// Парсинг даты
	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	slot := models.Slot{
		Date:         date,
		StartTime:    r.FormValue("start_time") + ":00",
		EndTime:      r.FormValue("end_time") + ":00",
		Price:        price,
		MaxPeople:    maxPeople,
		InstructorID: instructorID,
	}

	if err := h.repo.Create(r.Context(), &slot); err != nil {
		log.Printf("Error creating slot: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	slots, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting slots: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderSlotsList(w, slots)
}

func (h *SlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Date         string `json:"date"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		Price        int    `json:"price"`
		MaxPeople    int    `json:"max_people"`
		InstructorID int    `json:"instructor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Парсинг даты
	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	slot := models.Slot{
		ID:           id,
		Date:         date,
		StartTime:    data.StartTime,
		EndTime:      data.EndTime,
		Price:        data.Price,
		MaxPeople:    data.MaxPeople,
		InstructorID: data.InstructorID,
	}

	if err := h.repo.Update(r.Context(), &slot); err != nil {
		log.Printf("Error updating slot: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	slots, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting slots: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderSlotsList(w, slots)
}

func (h *SlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		log.Printf("Error deleting slot: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	slots, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting slots: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderSlotsList(w, slots)
}
