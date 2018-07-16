package graph

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
)

// User represents user who has ID and Name.
type User struct {
	ID   string
	Name string
}

// Resolver implements Resolvers interface
// type Resolvers interface {
// 	Mutation_createUser(ctx context.Context, input NewUser) (User, error)
// 	Query_user(ctx context.Context, id string) (*User, error)
// 	Query_users(ctx context.Context) ([]User, error)
// }
type Resolver struct {
	users []User
}

// Mutation_createUser creates a new user and add user to Resolver.
func (r *Resolver) Mutation_createUser(ctx context.Context, input NewUser) (User, error) {
	user := User{
		ID:   fmt.Sprintf("%d", rand.Int()),
		Name: input.Name,
	}
	r.users = append(r.users, user)
	return user, nil
}

// Query_user returns a user specified by the id.
func (r *Resolver) Query_user(ctx context.Context, id string) (*User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Query_users returns all the users who Resolver knows.
func (r *Resolver) Query_users(ctx context.Context) ([]User, error) {
	return r.users, nil
}
