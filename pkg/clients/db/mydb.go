package db

// MyDBInterface ..
type MyDBInterface interface {
	GetOne(id string) MimicUser
	Get(query string) []MimicUser
	GetAll() []MimicUser
	Insert(obj interface{}) error
}

// NewMyDB ..
func NewMyDB() MyDBInterface {
	return &MyDB{connection: "connection established"}
}

// MyDB ..
type MyDB struct {
	connection string
}

// MimicUser Just to minic user collection
type MimicUser struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetOne ..
func (m *MyDB) GetOne(id string) MimicUser {
	return MimicUser{ID: id, Name: "Shepard"}
}

// Get ..
func (m *MyDB) Get(query string) []MimicUser {

	return []MimicUser{{ID: "1", Name: "Shepard"}, {ID: "2", Name: "Miranda"}}
}

// GetAll ..
func (m *MyDB) GetAll() []MimicUser {
	return []MimicUser{{ID: "1", Name: "Shepard"},
		{ID: "2", Name: "Miranda"}, {ID: "3", Name: "Tali"}}
}

// Insert ..
func (m *MyDB) Insert(obj interface{}) error {
	return nil
}
