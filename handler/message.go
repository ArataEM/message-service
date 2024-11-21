package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArataEM/message-service/model"
	"github.com/ArataEM/message-service/repository/message"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Repo interface {
	Insert(ctx context.Context, message model.Message) error
	Get(ctx context.Context, id uuid.UUID) (model.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, message model.Message) error
	FindAll(ctx context.Context, page message.FindAllPage) (message.FindResult, error)
}

type Message struct {
	Repo Repo
}

func (m *Message) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId uuid.UUID `json:"user_id"`
		Text   string    `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	id, err := uuid.NewV7()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("error generating UUID for new message")
		return
	}
	message := model.Message{
		Id:        id,
		UserId:    body.UserId,
		Text:      body.Text,
		CreatedAt: &now,
	}

	err = m.Repo.Insert(r.Context(), message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("error inserting message")
		return
	}

	res, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("error marshaling message")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (m *Message) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := m.Repo.FindAll(r.Context(), message.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		slog.Error("failed to find all: ", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Items []model.Message `json:"items"`
		Next  uint64          `json:"next,omitempty"`
	}

	response.Items = res.Messages
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("failed to marshall: ", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (m *Message) GetById(w http.ResponseWriter, r *http.Request) {
	idParam, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		slog.Error("failed uuid parse")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messageRes, err := m.Repo.Get(r.Context(), idParam)
	if errors.Is(err, message.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to find by id", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(messageRes)
	if err != nil {
		slog.Error("failed to encode result")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *Message) UpdateById(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error decoding: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(body.Text)
	if body.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Field \"text\" cannot be empty"))
		return
	}
	idParam, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		slog.Error("failed uuid parse")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageRes, err := m.Repo.Get(r.Context(), idParam)
	if errors.Is(err, message.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to find by id", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	messageRes.UpdatedAt = &now
	messageRes.Text = body.Text

	err = m.Repo.Update(r.Context(), messageRes)
	if err != nil {
		slog.Error("failed to update record", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (m *Message) DeleteById(w http.ResponseWriter, r *http.Request) {
	idParam, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		slog.Error("failed uuid parse")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = m.Repo.Delete(r.Context(), idParam)
	if errors.Is(err, message.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to find by id", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
