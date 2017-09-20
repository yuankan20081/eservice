package writer

type Writer interface {
	WriteResponse(proto uint16, body interface{})
}
