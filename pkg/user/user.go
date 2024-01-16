package user

import (
	"go-boilerplate-api/apm"
	"go-boilerplate-api/config"
	"go-boilerplate-api/pkg/user/favourite"
	"go-boilerplate-api/pkg/user/rating"
	"go-boilerplate-api/pkg/user/repo"

	"github.com/ralstan-vaz/go-errors"
)

// UsersInterface ...
type UsersInterface interface {
	Get(query string) ([]*User, error)
	GetOne(id string) (*User, error)
	GetAll() ([]*User, error)
	Insert(u User) error
	GetWithInfo(id string) (*User, error)
}

// NewUser creates an instance of Users using the dependencies passed
// The dependency params can be moved to an interface to reduce to make it clean
func NewUser(conf config.IConfig, user repo.UserRepoInterface, apm apm.HandlerInterface, rating rating.Rater, favourite favourite.FavouriteInterface) UsersInterface {
	return &Users{config: conf, user: user, rating: rating, favourite: favourite, apm: apm}
}

// Users provides a way to perform operations on a user
type Users struct {
	config    config.IConfig
	user      repo.UserRepoInterface
	rating    rating.Rater
	favourite favourite.FavouriteInterface
	apm       apm.HandlerInterface
}

// Get gets users from the store using the query passed
func (pkg *Users) Get(query string) ([]*User, error) {
	repoUsers, err := pkg.user.Get(query)
	if err != nil {
		return nil, err
	}
	users := bindToUsers(repoUsers)
	return users, nil
}

// GetOne gets a user from the store using the query
func (pkg *Users) GetOne(id string) (*User, error) {
	repoUser, err := pkg.user.GetOne(id)
	if err != nil {
		return nil, err
	}
	user := bindToUser(repoUser)
	return user, nil
}

// GetAll gets all the users
func (pkg *Users) GetAll() ([]*User, error) {
	repoUsers, err := pkg.user.GetAll()
	if err != nil {
		return nil, err
	}
	users := bindToUsers(repoUsers)
	return users, nil
}

// Insert stores a user
func (pkg *Users) Insert(u User) error {

	user := repo.User{
		ID:   u.ID,
		Name: u.Name,
	}
	err := pkg.user.Insert(user)
	if err != nil {
		return err
	}

	return nil
}

// GetWithInfo get a user from the store along with the ratings and favourites
func (pkg *Users) GetWithInfo(id string) (*User, error) {
	repoUser, err := pkg.user.GetOne(id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	rating, err := pkg.rating.Get(rating.GetRequest{ID: id})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	favourite, err := pkg.favourite.Get(favourite.GetRequest{ID: id})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	fav := Favourite{Beers: favourite.Beers}
	user := &User{
		ID:        repoUser.ID,
		Name:      repoUser.Name,
		Stars:     rating.Stars,
		Favourite: fav,
	}

	return user, nil
}

func bindToUsers(u []*repo.User) []*User {
	user := []*User{}
	for i := 0; i < len(u); i++ {
		user = append(user, bindToUser(u[i]))
	}
	return user
}

func bindToUser(u *repo.User) *User {
	user := &User{ID: u.ID, Name: u.Name}
	return user
}
