package account

type Donator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewDonator(n string) Donator {
	result := Donator{}
	result.Name = n
	result.ID = GenerateID()
	return result
}

func (d Donator) GetName() string {
	return d.Name
}
