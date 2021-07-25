// When work with services layer we need to use interface instace of concrete class
// which allow us to do mocking later
package services

// step 4:expose interface type of variable, which is instance of your struct
var (
	ItemServices itemsServiceInterface = &itemServices{}
)

// step 1: define your interface first
type itemsServiceInterface interface {
	GetItem()
	SetItem()
}

// step 2: create service struct
type itemServices struct{}

// step 3: implemented all of methord defined in your service struct
func (s *itemServices) GetItem() {

}

func (s *itemServices) SetItem() {

}
