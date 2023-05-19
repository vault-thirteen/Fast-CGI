package h

type Data struct {
	StatusCode uint
	StatusText string
	Headers    []*Header
	Body       []byte
}

type Header struct {
	Name  string
	Value string
}
