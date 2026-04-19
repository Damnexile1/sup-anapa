package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sup-anapa/internal/middleware"
	bookinguc "sup-anapa/internal/usecase/booking"
	"time"
)

func UserRegisterPage(w http.ResponseWriter, r *http.Request) {
	next := safeNextURL(r.URL.Query().Get("next"), "/booking")
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/user-register.html",
	}, getTemplateData(r, map[string]interface{}{"Next": next}))
}

func UserLoginPage(w http.ResponseWriter, r *http.Request) {
	next := safeNextURL(r.URL.Query().Get("next"), "/booking")
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/user-login.html",
	}, getTemplateData(r, map[string]interface{}{"Next": next}))
}

func UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	next := safeNextURL(r.FormValue("next"), "/booking")
	user, err := userAuthSvc.Register(r.Context(), r.FormValue("username"), r.FormValue("password"), r.FormValue("phone"))
	if err != nil {
		data := getTemplateData(r, map[string]interface{}{
			"Next":  next,
			"Error": "Не удалось зарегистрироваться. Возможно, логин уже занят.",
		})
		renderTemplate(w, []string{
			"web/templates/layouts/base.html",
			"web/templates/public/user-register.html",
		}, data)
		return
	}
	session, _ := store.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func UserLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	next := safeNextURL(r.FormValue("next"), "/booking")
	user, err := userAuthSvc.Login(r.Context(), r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		data := getTemplateData(r, map[string]interface{}{
			"Next":  next,
			"Error": "Неверный логин или пароль",
		})
		renderTemplate(w, []string{
			"web/templates/layouts/base.html",
			"web/templates/public/user-login.html",
		}, data)
		return
	}
	session, _ := store.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserCabinet(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(int)
	if !ok || userID < 1 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	bookings, err := bookingRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Ошибка загрузки бронирований", http.StatusInternalServerError)
		return
	}
	data := getTemplateData(r, map[string]interface{}{"Bookings": bookings})
	renderTemplate(w, []string{
		"web/templates/layouts/base.html",
		"web/templates/public/user-cabinet.html",
	}, data)
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	userSession, _ := store.Get(r, "user-session")
	userID, ok := userSession.Values["user_id"].(int)
	if !ok || userID < 1 {
		http.Error(w, "Для бронирования нужно войти в аккаунт", http.StatusUnauthorized)
		return
	}

	var bookingData struct {
		SlotID      int    `json:"slot_id"`
		ClientName  string `json:"client_name"`
		ClientPhone string `json:"client_phone"`
		ClientEmail string `json:"client_email"`
		PeopleCount int    `json:"people_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&bookingData); err != nil {
		log.Printf("CreateBooking: invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	correlationID := middleware.GetCorrelationID(r.Context())
	log.Printf("CreateBooking: correlation_id=%s incoming request slot_id=%d people=%d client=%q", correlationID, bookingData.SlotID, bookingData.PeopleCount, bookingData.ClientName)

	booking, holdExpires, err := createBooking.Execute(r.Context(), bookinguc.CreateBookingInput{
		UserID:      userID,
		SlotID:      bookingData.SlotID,
		PeopleCount: bookingData.PeopleCount,
		ClientEmail: bookingData.ClientEmail,
	})
	if err != nil {
		switch err.Error() {
		case "slot_required":
			http.Error(w, "Выберите слот для бронирования", http.StatusBadRequest)
		case "invalid_people_count":
			http.Error(w, "Количество человек должно быть больше 0", http.StatusBadRequest)
		case "user_not_found":
			http.Error(w, "Не удалось получить данные пользователя", http.StatusUnauthorized)
		case "slot_not_found":
			log.Printf("CreateBooking: correlation_id=%s slot not found slot_id=%d", correlationID, bookingData.SlotID)
			http.Error(w, "Слот не найден", http.StatusNotFound)
		case "slot_unavailable":
			http.Error(w, "Слот уже занят или недоступен", http.StatusConflict)
		case "too_many_people":
			http.Error(w, "Слишком много человек для выбранной прогулки", http.StatusBadRequest)
		default:
			log.Printf("CreateBooking: correlation_id=%s usecase error: %v", correlationID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("CreateBooking: correlation_id=%s created booking_id=%d slot_id=%d status=%s", correlationID, booking.ID, booking.SlotID, booking.Status)

	response := map[string]interface{}{
		"ID":           booking.ID,
		"status":       booking.Status,
		"hold_expires": holdExpires.Format(time.RFC3339),
		"hold_minutes": 20,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func safeNextURL(next, fallback string) string {
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return fallback
	}
	return next
}
