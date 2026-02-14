package http

import (
	"echo/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UserPhotoHandler struct {
	photoService *service.UserPhotoService
}

func NewUserPhotoHandler(photoService *service.UserPhotoService) *UserPhotoHandler {
	return &UserPhotoHandler{
		photoService: photoService,
	}
}

// ReplaceUserPhotos godoc
// @Summary      Загрузка/замена фото пользователя
// @Description  Заменяет все фото пользователя
// @Tags         user-photos
// @Accept       multipart/form-data
// @Produce      json
// @Param        photos formData file true "Файлы фото" collectionFormat(multi)
// @Success      200  {array}  entity.UserPhoto
// @Failure      400  {object}  string "invalid files"
// @Failure      500  {object}  string "internal server error"
// @Router       /user/photos [put]
func (h *UserPhotoHandler) ReplaceUserPhotos(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111") // TODO: брать из токена

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid files", http.StatusBadRequest)
		return
	}

	form := r.MultipartForm
	files := form.File["photos"]
	if len(files) == 0 {
		http.Error(w, "no files provided", http.StatusBadRequest)
		return
	}

	photos, err := h.photoService.ReplaceUserPhotos(r.Context(), userID, files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photos)
}

// GetUserPhotos godoc
// @Summary      Получение фото пользователя
// @Description  Возвращает все фото пользователя
// @Tags         user-photos
// @Produce      json
// @Success      200  {array}  entity.UserPhoto
// @Failure      404  {object}  string "photos not found"
// @Router       /user/photos [get]
func (h *UserPhotoHandler) GetUserPhotos(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111") // TODO: брать из токена

	photos, err := h.photoService.GetUserPhotos(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(photos) == 0 {
		http.Error(w, "photos not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photos)
}
