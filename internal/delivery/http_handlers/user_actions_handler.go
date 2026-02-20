package http_handlers

import (
	"log/slog"
	"lots-service/internal/lib/responseHTTP"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *LotsHandler) GetUserPostedLots(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	postedLots, err := h.service.GetUserPostedLots(userID)
	if err != nil {
		slog.Debug("Опубликовані користувачем лоти не знайдені", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusNotFound, "Опубликовані користувачем лоти не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, postedLots)
}

func (h *LotsHandler) GetUserLikedLots(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	LikedLots, err := h.service.GetUserLikedLots(userID)
	if err != nil {
		slog.Debug("Лайкнуті користувачем лоти не знайдені", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusNotFound, "Лайкнуті користувачем лоти не знайдені")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, LikedLots)
}

func (h *LotsHandler) LikeLot(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		slog.Debug("Некоректний ID лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	err = h.service.LikeLot(userID, lotID)
	if err != nil {
		slog.Debug("Помилка встановлення лайку", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Не вдалося додати лайк")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LotsHandler) UnlikeLot(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		slog.Debug("Некоректний ID лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	err = h.service.UnlikeLot(userID, lotID)
	if err != nil {
		slog.Debug("Помилка прибирання лайку", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Не вдалося прибрати лайк")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LotsHandler) BuyLotHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		slog.Info("Некоректний ID лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Некоректний ID лота")
		return
	}

	err = h.service.BuyLot(userID, lotID)
	if err != nil {
		slog.Info("Помилка при купівлі лота", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Неможливо купити лот")
		return
	}

	slog.Debug("Куплено лот користувачем", "userID:", userID, "lotID:", lotID)

	responseHTTP.JSONRespMessage(w, http.StatusOK, "Лот успішно куплено")
}
