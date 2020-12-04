package main

import (
	"encoding/json"
	"net/http"
)

type tickerService interface {
	Start() (started bool, error error)
	Stop() (stopped bool, error error)
}

type TickerController struct {
	TickerService tickerService
}

func (tc *TickerController) HandleStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	started, err := tc.TickerService.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewInternalServerError("unable to start ticker service"))
		return
	}
	if !started {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(NewMessage("ticker service is started already"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewMessage("ticker service started"))
	return
}

func (tc *TickerController) HandleStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	stopped, err := tc.TickerService.Stop()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NewInternalServerError("unable to stop ticker service"))
		return
	}
	if !stopped {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(NewMessage("ticker service is stopped already"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewMessage("ticker service stopped"))
	return
}
