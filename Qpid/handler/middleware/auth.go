package middleware

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

type loginUserRetrieverImpl struct{}

func (r *loginUserRetrieverImpl) IsUserLoggedIn() bool {
	return true
}

func (r *loginUserRetrieverImpl) GetLoginUser() (string, error) {
	return "test-user", nil
}

func GetLoginUserRetriever() domain.LoginUserRetriever {
	return &loginUserRetrieverImpl{}
}
