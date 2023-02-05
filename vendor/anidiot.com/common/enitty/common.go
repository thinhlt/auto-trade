package enitty

type Signal struct {
  Price       float64 `json:"price,omitempty"`
  Order       string  `json:"order,omitempty"`
  CurrentTime int64   `json:"currentTime,omitempty"`
  Name        string  `json:"name,omitempty"`
  Msg         string  `json:"msg,omitempty"`
  J           float64 `json:"j,omitempty"`
  InList           bool `json:"in_list,omitempty"`
}
