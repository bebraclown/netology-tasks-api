package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct {
	store storage.Storage
}

func New(store storage.Storage) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	switch {
	case r.URL.Path == "/tasks":
		h.TasksCollection(w, r)
	case strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.TaskItem(w, r)
	default:
		writeError(w, http.StatusNotFound, "маршрут не найден")
	}
}

func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.store.List())
	case http.MethodPost:
		task, ok := readTask(w, r)
		if !ok {
			return
		}
		created, err := h.store.Create(task)
		if err != nil {
			writeStorageError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	idPart := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if idPart == "" || strings.Contains(idPart, "/") {
		writeError(w, http.StatusNotFound, "маршрут не найден")
		return
	}
	id, err := strconv.Atoi(idPart)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, found := h.store.Get(id)
		if !found {
			writeStorageError(w, storage.ErrNotFound)
			return
		}
		writeJSON(w, http.StatusOK, task)
	case http.MethodPut:
		task, ok := readTask(w, r)
		if !ok {
			return
		}
		updated, err := h.store.Update(id, task)
		if err != nil {
			writeStorageError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := h.store.Delete(id); err != nil {
			writeStorageError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w, "GET, PUT, DELETE")
	}
}

func readTask(w http.ResponseWriter, r *http.Request) (models.Task, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusBadRequest, "ожидается Content-Type application/json")
		return models.Task{}, false
	}
	var input struct {
		Title string          `json:"title"`
		Done  json.RawMessage `json:"done"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return models.Task{}, false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return models.Task{}, false
	}
	task := models.Task{Title: strings.TrimSpace(input.Title)}
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "поле title обязательно")
		return models.Task{}, false
	}
	switch string(input.Done) {
	case "", "false":
	case "true":
		task.Done = true
	default:
		writeError(w, http.StatusBadRequest, "поле done должно быть логическим значением")
		return models.Task{}, false
	}
	return task, true
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
}

func writeStorageError(w http.ResponseWriter, err error) {
	if errors.Is(err, storage.ErrNotFound) {
		writeError(w, http.StatusNotFound, storage.ErrNotFound.Error())
		return
	}
	log.Printf("ошибка хранилища: %v", err)
	writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	body, err := json.Marshal(value)
	if err != nil {
		log.Printf("ошибка кодирования ответа: %v", err)
		status = http.StatusInternalServerError
		body = []byte(`{"error":"внутренняя ошибка сервера"}`)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Printf("ошибка записи ответа: %v", err)
	}
}
