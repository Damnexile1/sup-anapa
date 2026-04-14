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

type BookingHandler struct {
	repo           *repository.BookingRepository
	slotRepo       *repository.SlotRepository
	instructorRepo *repository.InstructorRepository
}

func NewBookingHandler(repo *repository.BookingRepository, slotRepo *repository.SlotRepository, instructorRepo *repository.InstructorRepository) *BookingHandler {
	return &BookingHandler{
		repo:           repo,
		slotRepo:       slotRepo,
		instructorRepo: instructorRepo,
	}
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	// Получить фильтр по статусу из query параметров
	status := r.URL.Query().Get("status")

	var bookings []*models.Booking
	var err error

	if status != "" && status != "all" {
		bookings, err = h.repo.GetByStatus(r.Context(), status)
	} else {
		bookings, err = h.repo.GetAll(r.Context())
	}

	if err != nil {
		log.Printf("Error getting bookings: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Проверить, нужен ли HTML ответ (для HTMX)
	if r.Header.Get("HX-Request") == "true" {
		h.renderBookingsList(w, bookings)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) renderBookingsList(w http.ResponseWriter, bookings []*models.Booking) {
	if len(bookings) == 0 {
		fmt.Fprint(w, `<tr><td colspan="7" class="px-6 py-4 text-center text-gray-500">Бронирований нет</td></tr>`)
		return
	}

	for _, booking := range bookings {
		// Получить информацию о слоте
		slotInfo := fmt.Sprintf("Слот #%d", booking.SlotID)

		statusClass := "bg-yellow-100 text-yellow-800"
		statusText := "Ожидает"
		switch booking.Status {
		case "confirmed":
			statusClass = "bg-green-100 text-green-800"
			statusText = "Подтверждено"
		case "cancelled":
			statusClass = "bg-red-100 text-red-800"
			statusText = "Отменено"
		}

		fmt.Fprintf(w, `
		<tr class="border-b hover:bg-gray-50">
			<td class="px-6 py-4">#%d</td>
			<td class="px-6 py-4">%s</td>
			<td class="px-6 py-4">%s</td>
			<td class="px-6 py-4">%s</td>
			<td class="px-6 py-4">%d чел.</td>
			<td class="px-6 py-4">
				<span class="px-2 py-1 rounded text-xs font-semibold %s">%s</span>
			</td>
			<td class="px-6 py-4">
				<div class="flex gap-2">
					<button onclick="updateBookingStatus(%d, 'confirmed')" class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm">
						Подтвердить
					</button>
					<button onclick="updateBookingStatus(%d, 'cancelled')" class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm">
						Отменить
					</button>
				</div>
			</td>
		</tr>`, booking.ID, booking.ClientName, booking.ClientPhone, slotInfo, booking.PeopleCount, statusClass, statusText, booking.ID, booking.ID)
	}
}

func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, data.Status); err != nil {
		log.Printf("Error updating booking status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Вернуть обновленный список
	bookings, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Error getting bookings: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderBookingsList(w, bookings)
}
