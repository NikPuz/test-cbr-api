package cbr

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type cbrHandler struct {
	cbrService ICbrService
	logger     *zap.Logger
}

func RegisterCbrHandlers(r *mux.Router, logger *zap.Logger, service ICbrService) {
	cbrHandler := cbrHandler{
		cbrService: service,
		logger:     logger,
	}

	r.HandleFunc("/cbr/valute/curs", cbrHandler.GetValCurs).Methods("GET")
}

func (h cbrHandler) GetValCurs(w http.ResponseWriter, r *http.Request) {
	jsonValCurs, err := json.Marshal(h.cbrService.GetValCurs())
	if err != nil {
		h.logger.Error("Error Marshal ValCurs in handler", zap.Error(err))
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(jsonValCurs)
}
