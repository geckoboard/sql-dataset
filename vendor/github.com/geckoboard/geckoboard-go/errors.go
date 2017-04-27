package geckoboard

type Error struct {
	InnerError InnerError `json:"error"`
}

type InnerError struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.InnerError.Message
}
