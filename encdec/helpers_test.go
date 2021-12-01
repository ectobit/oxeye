package encdec_test

type msg struct {
	Name string `json:"name"`
}

func newMsg() *msg {
	return &msg{Name: "John Doe"}
}
