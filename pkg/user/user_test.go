package user

import (
	"go-boilerplate-api/pkg/user/favourite"
	"go-boilerplate-api/pkg/user/rating"
	"go-boilerplate-api/pkg/user/repo"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create a MockStore struct with an embedded mock instance
type MockStoreRepo struct {
	mock.Mock
}
type MockStoreRating struct {
	mock.Mock
}

type MockStoreFavourite struct {
	mock.Mock
}

// REPO MOCKS
func (m *MockStoreRepo) Get(query string) ([]*repo.User, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(query)
	// return the values which we define
	return returnVals.Get(0).([]*repo.User), returnVals.Error(1)
}

func (m *MockStoreRepo) GetOne(id string) (*repo.User, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(id)
	// return the values which we define
	return returnVals.Get(0).(*repo.User), returnVals.Error(1)
}

func (m *MockStoreRepo) GetAll() ([]*repo.User, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called()
	// return the values which we define
	return returnVals.Get(0).([]*repo.User), returnVals.Error(1)
}

func (m *MockStoreRepo) Insert(u repo.User) error {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(u)
	// return the values which we define
	return returnVals.Error(0)
}

// RATING MOCKS
func (m *MockStoreRating) Get(req rating.GetRequest) (*rating.GetResponse, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(req)
	// return the values which we define
	return returnVals.Get(0).(*rating.GetResponse), returnVals.Error(1)
}

// FAVOURITE MOCKS
func (m *MockStoreFavourite) Get(req favourite.GetRequest) (*favourite.GetResponse, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(req)
	// return the values which we define
	return returnVals.Get(0).(*favourite.GetResponse), returnVals.Error(1)
}

////////

// TESTS /////////

// declarations (common)
var repoUsers []*repo.User
var repoUser *repo.User
var m *MockStoreRepo

// global inits
func setup() {
	m = new(MockStoreRepo)
	repoUsers = []*repo.User{{ID: "111", Name: "Shourie"}}
	repoUser = &repo.User{ID: "111", Name: "Shourie"}
}

// can be used if some prior setup is required , ideally this should be the point of invocation
func TestMain(m *testing.M) {
	// Do stuff BEFORE the tests!
	// if any initializations are required , make them here
	setup()
	t := m.Run()
	// Do stuff AFTER the tests!
	// teardown() remove any setups that requires manual removal
	os.Exit(t)
}

func TestGetSuccess(t *testing.T) {
	var query = "111"
	var users = []*User{{ID: "111", Name: "Shourie"}}

	// Defines input and return type
	m.On("Get", query).Return(repoUsers, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Users{nil, m, nil, nil, nil}

	// Calls the actual module function
	resp, err := s.Get(query)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, users, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", users, resp)
	}
}

func TestGetOneSuccess(t *testing.T) {
	var query = "111"
	var user = &User{ID: "111", Name: "Shourie"}

	// Defines input and return type
	m.On("GetOne", query).Return(repoUser, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Users{nil, m, nil, nil, nil}

	// Calls the actual module function
	resp, err := s.GetOne(query)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, user, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", user, resp)
	}
}

func TestGetAllSuccess(t *testing.T) {
	var users = []*User{{ID: "111", Name: "Shourie"}}

	// Defines input and return type
	m.On("GetAll").Return(repoUsers, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Users{nil, m, nil, nil, nil}

	// Calls the actual module function
	resp, err := s.GetAll()

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, users, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", users, resp)
	}
}

func TestInsertSuccess(t *testing.T) {
	var user = User{ID: "111", Name: "Shourie"}

	// Defines input and return type
	m.On("Insert", *repoUser).Return(nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Users{nil, m, nil, nil, nil}

	// Calls the actual module function
	err := s.Insert(user)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}

func TestGetWithInfoSuccess(t *testing.T) {
	// declarations
	var id = "111"
	var ratResponse = &rating.GetResponse{ID: id, Stars: "true"}
	var favResponse = &favourite.GetResponse{ID: id, Beers: []string{"Moon Shine", "Bira", "Simba"}}
	fav := Favourite{Beers: favResponse.Beers}
	// This is gonna be returned ultimately
	var user = &User{ID: id, Name: "Shourie", Stars: ratResponse.Stars, Favourite: fav}

	// Create specific mocks objects only for this test
	m2 := new(MockStoreRating)
	m3 := new(MockStoreFavourite)

	// Defines input and return type
	m.On("GetOne", id).Return(repoUser, nil)
	m2.On("Get", rating.GetRequest{ID: id}).Return(ratResponse, nil)
	m3.On("Get", favourite.GetRequest{ID: id}).Return(favResponse, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Users{nil, m, m2, m3, nil}

	// Calls the actual module function
	resp, err := s.GetWithInfo(id)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)
	m2.AssertExpectations(t)
	m3.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, user, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", user, resp)
	}
}

///////
