package external

import (
	"cloud-render/internal/dto"
	"fmt"
)

type UserInfoResponse struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

func UserInfo(url string, id int64) (*dto.GetUserDTO, error) {
	requestURL := fmt.Sprintf(url, id)

	var res UserInfoResponse

	err := get(requestURL, &res)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserDTO{Login: res.Login, Email: res.Email}, nil
}
