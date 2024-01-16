package favourite

import (
	user "go-boilerplate-api/apis/grpc/generated/user"
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
func (m *MockStore) GetFav(grpcCon grpcConnectioner) (*user.Users, error) {
	// This allows us to pass in mocked results, so that the mock store will return whatever we define
	returnVals := m.Called(grpcCon)
	// return the values which we define
	return returnVals.Get(0).(*user.Users), returnVals.Error(1)
}

//////

// TESTS //////

//declarations
var m *MockStore

// global inits
func setup() {
	m = new(MockStore)
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

func TestModuleSuccess(t *testing.T) {
	// Declaractions
	req := GetRequest{ID: "1111"}
	res := GetResponse{}
	res.Beers = []string{"Moon Shine", "Bira", "Simba"}

	// Defines input and return type
	m.On("GetFav", nil).Return(&user.Users{}, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Favourite{nil, nil, m}

	// Calls the actual module function
	resp, err := s.Get(req)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.Equal(t, res, *resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should be same : res = %v  resp = %v", res, resp)
	}
}

func TestModuleFail(t *testing.T) {
	// Declaractions
	req := GetRequest{ID: "1111"}
	res := GetResponse{}
	res.Beers = []string{"Moon brew", "Bira", "Simba"}

	// Defines input and return type
	m.On("GetFav", nil).Return(&user.Users{}, nil)

	// Next, we create a new instance of our module with the mock store as its "Favourite" dependency
	s := Favourite{nil, nil, m}

	// Calls the actual module function
	resp, err := s.Get(req)

	// The expectations that we defined for our mock store earlier are asserted here
	m.AssertExpectations(t)

	// Assert response object
	success := assert.NotEqual(t, res, *resp)

	// Finally, we assert that we should'nt get any error
	if err != nil {
		t.Errorf("error should be nil, got: %v", err)
	}

	if !success {
		t.Errorf("assert failed, result should not be same : res = %v  resp = %v", res, resp)
	}
}

////////////
