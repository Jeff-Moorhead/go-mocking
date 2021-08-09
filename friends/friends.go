package friends

type Friend struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

func (self *Friend) Equals(other Friend) bool {
	return self.Name == other.Name && self.Age == other.Age && self.Occupation == other.Occupation
}
