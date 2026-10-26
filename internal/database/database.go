package database

import "httpServer/internal/domain"

type Repository struct {
	Users map[string]domain.User
}

func NewRepository() Repository {
	user := domain.User{Name: "Max", Id: "TEST"}
	users := map[string]domain.User{
		user.Id: user,
	}
	return Repository{
		Users: users,
	}
}

func (r *Repository) GetUsers() []domain.User {
	users := make([]domain.User, 0, len(r.Users))

	for _, event := range r.Users {
		users = append(users, event)
	}

	return users
}

func (r *Repository) CreateUser(name string, id string) (domain.User, error) {

	user := domain.User{
		Id:   id,
		Name: name,
	}
	r.Users[id] = user

	return user, nil
}
