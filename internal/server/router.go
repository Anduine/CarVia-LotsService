package server

import (
	"log/slog"
	"lots-service/internal/delivery/http_handlers"
	"lots-service/internal/lib/responseHTTP"
	"lots-service/pkg/auth"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(lotsHandler *http_handlers.LotsHandler) http.Handler {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3011/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	router.Handle("/api/lots/id/{lot_id}", auth.OptionalAuthMiddleware(lotsHandler.GetLotByID)).Methods("GET")
	router.Handle("/api/lots/filtered", auth.OptionalAuthMiddleware(lotsHandler.GetLotsPageByParams)).Methods("GET")

	router.Handle("/api/lots/sell_lots", auth.OptionalAuthMiddleware(lotsHandler.GetLotsPage)).Methods("GET")
	router.HandleFunc("/api/lots/sell_lots_count", lotsHandler.GetLotsCount).Methods("GET")
	router.HandleFunc("/api/lots/sell_lots_filtered_count", lotsHandler.GetLotsByParamsCount).Methods("GET")

	router.HandleFunc("/api/lots/brands", lotsHandler.GetBrands).Methods("GET")
	router.HandleFunc("/api/lots/models", lotsHandler.GetModels).Methods("GET")

	router.Handle("/api/lots/user_posted_lots", auth.AuthMiddleware(lotsHandler.GetUserPostedLots)).Methods("GET")
	router.Handle("/api/lots/user_liked_lots", auth.AuthMiddleware(lotsHandler.GetUserLikedLots)).Methods("GET")

	router.Handle("/api/lots/create_lot", auth.AuthMiddleware(lotsHandler.CreateLot)).Methods("POST")
	router.Handle("/api/lots/update_lot/{lot_id}", auth.AuthMiddleware(lotsHandler.UpdateLot)).Methods("PUT")
	router.Handle("/api/lots/delete_lot/{lot_id}", auth.AuthMiddleware(lotsHandler.DeleteLot)).Methods("DELETE")

	router.Handle("/api/lots/likes/{lot_id}", auth.AuthMiddleware(lotsHandler.LikeLot)).Methods("POST")
	router.Handle("/api/lots/likes/{lot_id}", auth.AuthMiddleware(lotsHandler.UnlikeLot)).Methods("DELETE")

	router.Handle("/api/lots/buy_lot/{lot_id}", auth.AuthMiddleware(lotsHandler.BuyLotHandler)).Methods("PUT")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Маршрут не знайдено", "method", r.Method, "path", r.URL.Path)
		responseHTTP.JSONError(w, http.StatusNotFound, "Маршрут не знайдено")
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Заборонений метод", "method", r.Method, "path", r.URL.Path)
		responseHTTP.JSONError(w, http.StatusMethodNotAllowed, "Заборонений метод")
	})

	return router
}
