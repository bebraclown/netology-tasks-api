package storage

import (
	"errors"
	"sort"
	"sync"

	"tasks-api/internal/models"
)

var ErrNotFound = errors.New("задача не найдена")

type Storage interface {
	List() []models.Task
	Create(task models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}

type Memory struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

func NewMemory() *Memory {
	return &Memory{tasks: make(map[int]models.Task), nextID: 1}
}

func (store *Memory) List() []models.Task {
	store.mu.RLock()
	tasks := make([]models.Task, 0, len(store.tasks))
	for _, task := range store.tasks {
		tasks = append(tasks, task)
	}
	store.mu.RUnlock()

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks
}

func (store *Memory) Create(task models.Task) (models.Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	task.ID = store.nextID
	store.nextID++
	store.tasks[task.ID] = task
	return task, nil
}

func (store *Memory) Get(id int) (models.Task, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	task, found := store.tasks[id]
	return task, found
}

func (store *Memory) Update(id int, task models.Task) (models.Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, found := store.tasks[id]; !found {
		return models.Task{}, ErrNotFound
	}
	task.ID = id
	store.tasks[id] = task
	return task, nil
}

func (store *Memory) Delete(id int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, found := store.tasks[id]; !found {
		return ErrNotFound
	}
	delete(store.tasks, id)
	return nil
}
