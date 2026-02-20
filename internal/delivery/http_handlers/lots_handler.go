package http_handlers

import (
	"encoding/json"
	"log/slog"
	"lots-service/internal/domain"
	"lots-service/internal/lib/responseHTTP"
	"lots-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LotsHandler struct {
	service *service.LotsService
}

func NewLotsHandler(service *service.LotsService) *LotsHandler {
	return &LotsHandler{service: service}
}

func (h *LotsHandler) GetLotsCount(w http.ResponseWriter, r *http.Request) {
	lotsCount, err := h.service.GetLotsCount()
	if err != nil {
		slog.Debug("Кількість 0, лоти не знайдені", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusNotFound, "Лоти не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, lotsCount)
}

func (h *LotsHandler) GetLotsByParamsCount(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	brand := params.Get("brand")
	model := params.Get("model")
	minPrice := params.Get("minPrice")
	maxPrice := params.Get("maxPrice")
	minYear := params.Get("minYear")
	maxYear := params.Get("maxYear")

	lotsCount, err := h.service.GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusNotFound, "Лоти за вказаними параметрами не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, lotsCount)
}

// @Summary		Get Lot by ID
// @Description	Get Lot by ID
// @Tags			Lots
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Lot ID"
// @Success		200	{object}	domain.Lot
// @Failure		400	{object}	responseHTTP.ErrorResponse
// @Failure		404	{object}	responseHTTP.ErrorResponse
// @Failure		500	{object}	responseHTTP.ErrorResponse
// @Router			/api/lots/id/{lot_id} [get]
func (h *LotsHandler) GetLotByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	lot, err := h.service.GetLotByID(userID, lotID)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusNotFound, "Лот не знайдено")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, lot)
}

func (h *LotsHandler) GetLotsPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	lots, err := h.service.GetPageLots(userID, page, limit)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusNotFound, "Лоти не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, lots)
}

func (h *LotsHandler) GetLotsPageByParams(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		userID = 0
	}

	params := r.URL.Query()

	brand := params.Get("brand")
	model := params.Get("model")
	minPrice := params.Get("minPrice")
	maxPrice := params.Get("maxPrice")
	minYear := params.Get("minYear")
	maxYear := params.Get("maxYear")

	pageStr := params.Get("page")
	limitStr := params.Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	lots, total, err := h.service.GetLotsByParams(userID, page, limit, brand, model, minPrice, maxPrice, minYear, maxYear)
	if err != nil {
		slog.Debug("Помилка при отриманні лотів з БД", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	if total == 0 {
		slog.Debug("Лоти за параметрами не знайдені", "brand", brand, "model", model, "minPrice", minPrice,
			"maxPrice", maxPrice, "minYear", minYear, "maxYear", maxYear)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.LotsResponse{Lots: []domain.Lot{}, Total: 0})
		return
	}

	response := domain.LotsResponse{
		Lots:  *lots,
		Total: total,
	}

	if response.Lots == nil {
		response.Lots = []domain.Lot{}
	}

	slog.Debug("Запит " + r.RequestURI)

	responseHTTP.JSONResp(w, http.StatusOK, response)
}

func (h *LotsHandler) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.service.GetBrands()
	if err != nil {
		slog.Debug("Бренди не знайдені", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusNotFound, "Бренди не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, brands)
}

func (h *LotsHandler) GetModels(w http.ResponseWriter, r *http.Request) {
	brandName := r.URL.Query().Get("brand")

	models, err := h.service.GetModels(brandName)
	if err != nil {
		slog.Debug("Моделі не знайдені", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusNotFound, "Моделі не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, models)
}

func (h *LotsHandler) CreateLot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		slog.Debug("Помилка парсингу форми", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusUnprocessableEntity, "Помилка парсингу форми")
		return
	}

	lot, err := ParseLotFromRequest(r)
	if err != nil {
		slog.Debug("Помилка валідації лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusUnprocessableEntity, "Помилка валідації форми")
		return
	}
	lot.SellerID = userID

	files := r.MultipartForm.File["NewImages"]

	if err := h.service.CreateLot(r.Context(), &lot, files); err != nil {
		slog.Debug("Помилка збереження лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	slog.Debug("Створено лот")

	responseHTTP.JSONRespMessage(w, http.StatusCreated, "Лот створено")
}

func (h *LotsHandler) UpdateLot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		slog.Debug("Некоректний ID лота", "lotID", lotID, "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Debug("Помилка парсингу форми", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusUnprocessableEntity, "Помилка парсингу форми")
		return
	}

	lot, err := ParseLotFromRequest(r)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Помилка парсингу лота з форми")
		return
	}
	lot.LotID = lotID
	lot.SellerID = userID

	files := r.MultipartForm.File["NewImages"]
	deleteImages := r.MultipartForm.Value["DeleteImagesNames"]
	oldImagesStr := r.MultipartForm.Value["OldImagesNames"]

	slog.Debug("Оновлення лота", "lotID", lotID, "userID", userID, "deleteImages", deleteImages, "oldImagesStr", oldImagesStr, "files", files)

	if err := h.service.UpdateLot(r.Context(), &lot, files, deleteImages, oldImagesStr); err != nil {
		slog.Debug("Помилка при оновленні лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	slog.Debug("Оновлено лот", "lotID", lotID)

	responseHTTP.JSONRespMessage(w, http.StatusOK, "Лот оновлено")
}

func (h *LotsHandler) DeleteLot(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		slog.Debug("Некоректний ID лота", "lot_id", lotID, "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	err = h.service.DeleteLot(r.Context(), lotID, userID)
	if err != nil {
		slog.Debug("Помилка видалення лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	slog.Debug("Видалено лот користувачем", "userID", userID, "lotID", lotID)

	responseHTTP.JSONRespMessage(w, http.StatusOK, "Лот видалено")
}
