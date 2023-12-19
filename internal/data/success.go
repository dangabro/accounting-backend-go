package data

type SuccessData struct {
	Success bool `json:"success"`
}

func GetSuccessData() SuccessData {
	return SuccessData{
		Success: true,
	}
}
