package user

import "userservice/internal/core/entities"

// TODO: complete
type Storage interface {
	Add(id entities.ID)
	Get(id entities.ID)
	Remove(id entities.ID)
	List() []entities.ID
}
