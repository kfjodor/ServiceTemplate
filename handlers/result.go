package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ttm-keywallet/statserver/logger"
)

// Response interface
// swagger:response appJsonResult
type appJsonResult struct {
	// Can be a address
	Result interface{} `json:"result"`
}

func ReturnResult(ctx context.Context, w http.ResponseWriter, r interface{}) {
	log := logger.FromContext(ctx).WithField("m", "ReturnResult")
	log.Debugf("ReturnResult:: w: %v, r: %v", w, r)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respJson, err := json.Marshal(appJsonResult{r})
	if err == nil {
		log.Infof("Response:%v, %v", http.StatusOK, string(respJson))
	} else {
		log.Infof("Response error: %v", err)
	}
	json.NewEncoder(w).Encode(appJsonResult{r})
}

func ReturnResultWithCode(ctx context.Context, w http.ResponseWriter, s int, r interface{}) {
	log := logger.FromContext(ctx).WithField("m", "ReturnResultWithCode")
	log.Debugf("ReturnResultWithCode:: w: %v, s: %v, r: %v", w, s, r)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(s)
	respJson, err := json.Marshal(appJsonResult{r})
	if err == nil {
		log.Infof("Response: %v", s, string(respJson))
	} else {
		log.Infof("Response error: %v", err)
	}
	json.NewEncoder(w).Encode(appJsonResult{r})
}
