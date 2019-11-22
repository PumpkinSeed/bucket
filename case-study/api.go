package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PumpkinSeed/bucket/case-study/models"
	"github.com/go-chi/chi"
)

/*
	Router section
 */

func Routes(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/profile", ProfilePost)
		r.Get("/profile/{id}", ProfileGet)
		r.Post("/profile/{id}", ProfileUpdate)

		r.Post("/event", EventPost)
		r.Get("/event/{id}", EventGet)
		r.Post("/event/{id}", EventUpdate)

		r.Post("/order", OrderPost)
		r.Get("/order/{id}", OrderGet)
		r.Post("/order/{id}", OrderUpdate)
	})
}

/*
	Profile section
 */

var profiles []string

func ProfilePost(w http.ResponseWriter, r *http.Request) {
	var req *models.Profile
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, id, err := han.Insert(context.Background(), ProfileType, "", req, 0)
	profiles = append(profiles, id)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Success")
}

func ProfileGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var resp *models.Profile
	if err := han.Get(context.Background(), ProfileType, id, resp); err != nil {
		respond(w, err)
	}

	respond(w, resp)
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req *models.Profile
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, _, err := han.Upsert(context.Background(), ProfileType, id, req, 0)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Successfully updated")
}

/*
	Event section
 */

var events []string

func EventPost(w http.ResponseWriter, r *http.Request) {
	var req *models.Event
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, id, err := han.Insert(context.Background(), EventType, "", req, 0)
	events = append(events, id)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Success")
}

func EventGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var resp *models.Event
	if err := han.Get(context.Background(), EventType, id, resp); err != nil {
		respond(w, err)
	}

	respond(w, resp)
}

func EventUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req *models.Event
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, _, err := han.Upsert(context.Background(), EventType, id, req, 0)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Successfully updated")
}

/*
	Order section
 */

var orders []string

func OrderPost(w http.ResponseWriter, r *http.Request) {
	var req *models.Order
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, id, err := han.Insert(context.Background(), OrderType, "", req, 0)
	orders = append(orders, id)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Success")
}

func OrderGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var resp *models.Order
	if err := han.Get(context.Background(), OrderType, id, resp); err != nil {
		respond(w, err)
	}

	respond(w, resp)
}

func OrderUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req *models.Order
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respond(w, err)
	}

	_, _, err := han.Upsert(context.Background(), OrderType, id, req, 0)
	if err != nil {
		respond(w, err)
	}

	respond(w, "Successfully updated")
}