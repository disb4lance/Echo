package http

import (
	"echo/internal/service"
	"echo/internal/service/dto"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UserProfileHandler struct {
	profileService *service.UserProfileService
}

func NewUserProfileHandler(profileService *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		profileService: profileService,
	}
}

// CreateProfile godoc
// @Summary      Создание профиля пользователя
// @Description  Создает новый профиль пользователя
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Param        request body dto.UserProfileRequest true "Данные профиля"
// @Success      201  {object}  dto.UserProfileResponse
// @Failure      400  {object}  string "invalid body"
// @Failure      409  {object}  string "profile already exists"
// @Router       /profiles [post]
func (h *UserProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.UserProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.profileService.CreateProfile(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetProfile godoc
// @Summary      Получение профиля пользователя
// @Description  Возвращает профиль пользователя по ID
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Param        id path string true "ID пользователя" Format(uuid)
// @Success      200  {object}  dto.UserProfileResponse
// @Failure      400  {object}  string "invalid user id"
// @Failure      404  {object}  string "profile not found"
// @Router       /profiles/{id} [get]
func (h *UserProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := uuid.New().String() // TODO брать из токена

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	resp, err := h.profileService.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateProfile godoc
// @Summary      Обновление профиля пользователя
// @Description  Обновляет существующий профиль пользователя
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Param        id path string true "ID пользователя" Format(uuid)
// @Param        request body dto.UserProfileRequest true "Данные для обновления"
// @Success      200  {object}  dto.UserProfileResponse
// @Failure      400  {object}  string "invalid body or user id"
// @Failure      404  {object}  string "profile not found"
// @Router       /profiles/{id} [put]
func (h *UserProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := uuid.New().String() // TODO брать из токена

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req dto.UserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.profileService.UpdateProfile(userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
