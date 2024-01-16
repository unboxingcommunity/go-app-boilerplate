package repo

import (
	"go-boilerplate-api/pkg/clients/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create a MockStore struct with an embedded mock instance
type MockStore struct {
	mock.Mock
}

// MOCKS -----------------
func (m *MockStore) Get(query string) []db.MimicUser {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(query)
	// return the values which we define
	return returnVals.Get(0).([]db.MimicUser)
}

func (m *MockStore) GetOne(id string) db.MimicUser {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(id)
	// return the values which we define
	return returnVals.Get(0).(db.MimicUser)
}

func (m *MockStore) GetAll() []db.MimicUser {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called()
	// return the values which we define
	return returnVals.Get(0).([]db.MimicUser)
}

func (m *MockStore) Insert(obj interface{}) error {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(obj)
	// return the values which we define
	return returnVals.Error(0)
}

//////

// TESTS ---------------------

// declarations (common)
var mUsers []db.MimicUser
var mUser db.MimicUser
var m *MockStore

// global inits
func setup() {
	m = new(MockStore)
	mUsers = []db.MimicUser{{ID: "111", Name: "Shourie"}}
	mUser = db.MimicUser{ID: "111", Name: "Shourie"}
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
	var repoUsers = []*User{{ID: "111", Name: "Shourie"}}

	// Defines input and return type
	m.On("Get", query).Return(mUsers)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	repo := UserRepo{nil, m}

	// Calls the actual module function
	resp, err := repo.Get(query)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, repoUsers, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", repoUsers, resp)
	}
}

func TestGetOneSuccess(t *testing.T) {
	var query = "111"
	var repoUser = &User{ID: "111", Name: "Shourie"}

	// Defines input and return type
	m.On("GetOne", query).Return(mUser)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	repo := UserRepo{nil, m}

	// Calls the actual module function
	resp, err := repo.GetOne(query)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, repoUser, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", repoUser, resp)
	}
}

func TestGetAllSuccess(t *testing.T) {
	var repoUsers = []*User{{ID: "111", Name: "Shourie"}}

	// Defines input and return type
	m.On("GetAll").Return(mUsers)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	repo := UserRepo{nil, m}

	// Calls the actual module function
	resp, err := repo.GetAll()

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, repoUsers, resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", repoUsers, resp)
	}
}

func TestInsertSuccess(t *testing.T) {
	var repoUser = User{ID: "111", Name: "Shourie"}

	// Defines input and return type
	m.On("Insert", repoUser).Return(nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	repo := UserRepo{nil, m}

	// Calls the actual module function
	err := repo.Insert(repoUser)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}
}
