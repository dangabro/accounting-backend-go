package data

type IdError interface {
	IsSystem() bool
	Error() string
}

func CreateIdError(system bool, message string) IdError {
	return &idError{
		system:  system,
		message: message,
	}
}

type idError struct {
	system  bool
	message string
}

func (d *idError) IsSystem() bool {
	return d.system
}

func (d *idError) Error() string {
	return d.message
}

func GetIdError(err error) IdError {
	iderr, ok := err.(IdError)
	if ok {
		return iderr
	} else {
		res := CreateIdError(true, err.Error())
		return res
	}
}

type PayloadError struct {
	Message string `json:"message"`
}
