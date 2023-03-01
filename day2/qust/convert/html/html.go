package html

type Data struct {
	data []byte
}

func New(body []byte) *Data {
	return &Data{body}
}
