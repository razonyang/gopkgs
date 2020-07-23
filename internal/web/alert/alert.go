package alert

import "encoding/gob"

func init() {
	gob.Register(Alert{})
}

const (
	StyleSuccess = "success"
	StyletDanger = "danger"
	StyleWarning = "warning"
	StyleInfo    = "info"
)

type Alert struct {
	Style   string
	Message string
}

func New(style, message string) Alert {
	return Alert{
		Style:   style,
		Message: message,
	}
}

func NewSuccess(message string) Alert {
	return New(StyleSuccess, message)
}

func NewDanger(message string) Alert {
	return New(StyletDanger, message)
}
