package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"
	"sup-anapa/internal/services"

	"github.com/gorilla/sessions"
)

var (
	store          *sessions.CookieStore
	authSvc        *services.AuthService
	bookingRepo    *repository.BookingRepository
	instructorRepo *repository.InstructorRepository
)

func Init(sessionSecret string, authService *services.AuthService) {
	store = sessions.NewCookieStore([]byte(sessionSecret))
	authSvc = authService
}

func SetRepositories(booking *repository.BookingRepository, instructor *repository.InstructorRepository) {
	bookingRepo = booking
	instructorRepo = instructor
}

func GetStore() *sessions.CookieStore {
	return store
}

func renderTemplate(w http.ResponseWriter, layoutFiles []string, data interface{}) {
	tmpl, err := template.ParseFiles(layoutFiles...)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTemplateData(r *http.Request, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Получить username из сессии
	session, err := store.Get(r, "admin-session")
	if err == nil {
		if username, ok := session.Values["username"].(string); ok {
			data["Username"] = username
		}
	}

	return data
}

func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/home.html",
	}, getTemplateData(r, nil))
}

func BookingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/booking.html",
	}, getTemplateData(r, nil))
}

func InstructorsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/instructors.html",
	}, getTemplateData(r, nil))
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var bookingData struct {
		SlotID      int    `json:"slot_id"`
		ClientName  string `json:"client_name"`
		ClientPhone string `json:"client_phone"`
		ClientEmail string `json:"client_email"`
		PeopleCount int    `json:"people_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bookingData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if bookingData.ClientName == "" {
		http.Error(w, "Имя клиента обязательно", http.StatusBadRequest)
		return
	}
	if bookingData.ClientPhone == "" {
		http.Error(w, "Телефон обязателен", http.StatusBadRequest)
		return
	}
	if bookingData.PeopleCount < 1 {
		http.Error(w, "Количество человек должно быть больше 0", http.StatusBadRequest)
		return
	}

	// Создать бронирование
	booking := &models.Booking{
		SlotID:      bookingData.SlotID,
		ClientName:  bookingData.ClientName,
		ClientPhone: bookingData.ClientPhone,
		ClientEmail: bookingData.ClientEmail,
		PeopleCount: bookingData.PeopleCount,
		Status:      "pending",
	}

	if err := bookingRepo.Create(r.Context(), booking); err != nil {
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"ID":     booking.ID,
		"status": booking.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/admin/admin-login.html",
	}, getTemplateData(r, nil))
}

func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := authSvc.Authenticate(r.Context(), username, password)
	if err != nil {
		// Показать ошибку на странице логина
		data := getTemplateData(r, map[string]interface{}{
			"Error": "Неверный логин или пароль",
		})
		renderTemplate(w, []string{
			"web/templates/layouts/base.html",
			"web/templates/admin/admin-login.html",
		}, data)
		return
	}

	// Создать сессию
	session, _ := store.Get(r, "admin-session")
	session.Values["admin_id"] = admin.ID
	session.Values["username"] = admin.Username
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/admin/admin-dashboard.html",
	}, getTemplateData(r, nil))
}

func AdminInstructors(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/admin/admin-instructors.html",
	}, getTemplateData(r, nil))
}

func AdminSlots(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/admin/admin-slots.html",
	}, getTemplateData(r, nil))
}

func AdminBookings(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/admin/admin-bookings.html",
	}, getTemplateData(r, nil))
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "admin-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
