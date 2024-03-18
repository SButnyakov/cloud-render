package suite

import (
	"cloud-render/internal/http/api"
	"cloud-render/internal/http/buffer"
	"cloud-render/internal/lib/config"
	"cloud-render/internal/lib/response"
	"cloud-render/test"
	apiTest "cloud-render/test/api"
	authTest "cloud-render/test/auth"
	bufferTest "cloud-render/test/buffer"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type InegrationTestsSuite struct {
	authClient       *test.AuthClient
	authTestClient   *authTest.AuthTestClient
	apiTestClient    *apiTest.APITestClient
	bufferTestClient *bufferTest.BufferTestClient
}

func SetupSuite(authCfg, apiCfg, bufferCfg *config.Config) *InegrationTestsSuite {
	authClient := &test.AuthClient{Config: authCfg}

	return &InegrationTestsSuite{
		authClient:       authClient,
		authTestClient:   authTest.NewAuthTestClient(authCfg, authClient),
		apiTestClient:    apiTest.NewAPITestClient(apiCfg, authClient),
		bufferTestClient: bufferTest.NewBufferTestClient(bufferCfg),
	}
}

func (s *InegrationTestsSuite) RegisterAndSignIn(login, email, password string) (bool, error) {
	resParams, err := s.authClient.SignUp(login, email, password)
	if err != nil {
		return false, err
	}

	if resParams.Code != http.StatusCreated {
		errString, err := unmarshalError(resParams.Body)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal error: %s", err.Error())
		}
		return false, fmt.Errorf("failed to sign up. code: %d, error: %s", resParams.Code, errString)
	}

	resParams, err = s.authClient.SignIn(login, password)
	if err != nil {
		return false, err
	}

	if resParams.Code != http.StatusOK {
		errString, err := unmarshalError(resParams.Body)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal error: %s", err.Error())
		}
		return false, fmt.Errorf("failed to sign in. code: %d, error: %s", resParams.Code, errString)
	}

	return true, nil
}

func (s *InegrationTestsSuite) TestEditSignIn() (bool, string, error) {
	defer s.logout()

	login := "editSignInLogin"
	email := "editSignIn@email.com"
	password := "editSignInPassword"

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	newPassword := "newEditSignInPassword"

	resParams, err := s.authTestClient.Edit(login, email, newPassword)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to edit user info: %s", string(resParams.Body))
	}

	resParams, err = s.authClient.SignIn(login, password)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusBadRequest {
		return false, "", fmt.Errorf("failed to edit user info: %s", string(resParams.Body))
	}

	msg, err := unmarshalError(resParams.Body)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal response: %s", err.Error())
	}

	expectedMsg := "invalid credentials"
	if msg != expectedMsg {
		return false, fmt.Sprintf("wrong error msg. expected: %s, actual: %s", expectedMsg, msg), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) TestOrdersSend() (bool, string, error) {
	defer s.logout()

	login := "orderSendLogin"
	email := "orderSend@email.com"
	password := "orderSendPassword"

	// Register and sign in

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	// Orders

	resParams, err := s.apiTestClient.Orders()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var ordersResponse api.GetManyOrdersResponse

	err = json.Unmarshal(resParams.Body, &ordersResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if len(ordersResponse.Orders) != 0 {
		return false, fmt.Sprintf("wrong number of orders. expected: 0, actual: %d", len(ordersResponse.Orders)), nil
	}

	// Send

	filename := "temp.blend"
	format := "png"
	resolution := "1920x1080"

	resParams, err = s.apiTestClient.Send(filename, format, resolution)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusCreated {
		return false, "", fmt.Errorf("failed to create new order. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	// Orders

	resParams, err = s.apiTestClient.Orders()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var newOrdersResponse api.GetManyOrdersResponse

	err = json.Unmarshal(resParams.Body, &newOrdersResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if len(newOrdersResponse.Orders) != 1 {
		return false, fmt.Sprintf("wrong number of orders. expected: 1, actual: %d", len(newOrdersResponse.Orders)), nil
	}

	order := newOrdersResponse.Orders[0]
	status := "in queue"

	if order.FileName != filename {
		return false, fmt.Sprintf("wrong filename. expected: %s, actual: %s", filename, order.FileName), nil
	}
	if order.Status != status {
		return false, fmt.Sprintf("wrong status. expected: %s, actual: %s", status, order.Status), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) TestOrdersUpdate(uid string) (bool, string, error) {
	defer s.logout()

	login := "orUpdateLogin"
	email := "orderUpdate@email.com"
	password := "orderUpdatePassword"

	// Register and sign in

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	// Send

	filename := "temp.blend"
	format := "png"
	resolution := "1920x1080"

	resParams, err := s.apiTestClient.Send(filename, format, resolution)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusCreated {
		return false, "", fmt.Errorf("failed to create new order. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	// Orders

	resParams, err = s.apiTestClient.Orders()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var ordersResponse api.GetManyOrdersResponse

	err = json.Unmarshal(resParams.Body, &ordersResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if len(ordersResponse.Orders) != 1 {
		return false, fmt.Sprintf("wrong number of orders. expected: 1, actual: %d", len(ordersResponse.Orders)), nil
	}

	order := ordersResponse.Orders[0]
	status := "in queue"

	if order.FileName != filename {
		return false, fmt.Sprintf("wrong filename. expected: %s, actual: %s", filename, order.FileName), nil
	}
	if order.Status != status {
		return false, fmt.Sprintf("wrong status. expected: %s, actual: %s", status, order.Status), nil
	}

	// Update

	resParams, err = s.bufferTestClient.Update(uid, strconv.FormatInt(order.Id, 10), "in-progress")
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to update order status. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	// Orders

	resParams, err = s.apiTestClient.Orders()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var newOrdersResponse api.GetManyOrdersResponse

	err = json.Unmarshal(resParams.Body, &newOrdersResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if len(newOrdersResponse.Orders) != 1 {
		return false, fmt.Sprintf("wrong number of orders. expected: 1, actual: %d", len(newOrdersResponse.Orders)), nil
	}

	order = newOrdersResponse.Orders[0]
	status = "in progress"

	if order.FileName != filename {
		return false, fmt.Sprintf("wrong filename. expected: %s, actual: %s", filename, order.FileName), nil
	}
	if order.Status != status {
		return false, fmt.Sprintf("wrong status. expected: %s, actual: %s", status, order.Status), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) TestSendRequest() (bool, string, error) {
	defer s.logout()

	login := "sendReqLogin"
	email := "sendReq@email.com"
	password := "sendReqPassword"

	// Register and sign in

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	// Request

	resParams, err := s.bufferTestClient.Request()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get queued orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var requestResponse buffer.RequestResponse

	err = json.Unmarshal(resParams.Body, &requestResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if requestResponse.Status != response.StatusEmpty {
		return false, fmt.Sprintf("wrong response status. expected: %s, actual: %s", response.StatusEmpty, requestResponse.Status), nil
	}

	// Send

	filename := "temp.blend"
	format := "png"
	resolution := "1920x1080"

	resParams, err = s.apiTestClient.Send(filename, format, resolution)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusCreated {
		return false, "", fmt.Errorf("failed to create new order. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	// Request

	resParams, err = s.bufferTestClient.Request()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get queued orders. code: %d, body: %s", resParams.Code, string(resParams.Body))
	}

	var newRequestResponse buffer.RequestResponse

	err = json.Unmarshal(resParams.Body, &newRequestResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal orders response: %s", err.Error())
	}

	if newRequestResponse.Status != response.StatusOK {
		return false, fmt.Sprintf("wrong response status. expected: %s, actual: %s", response.StatusOK, newRequestResponse.Status), nil
	}
	if newRequestResponse.Format != format {
		return false, fmt.Sprintf("wrong format. expected: %s, actual: %s", format, newRequestResponse.Format), nil
	}
	if newRequestResponse.Resolution != resolution {
		return false, fmt.Sprintf("wrong resolution. expected: %s, actual: %s", resolution, newRequestResponse.Resolution), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) TestSubscribeUser() (bool, string, error) {
	defer s.logout()

	login := "subUserLogin"
	email := "subUser@email.com"
	password := "subUserPassword"

	// Register and sign in

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	// User

	resParams, err := s.apiTestClient.User()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user info: %s", string(resParams.Body))
	}

	var userResponse api.UserResposne

	err = json.Unmarshal(resParams.Body, &userResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal user response: %s", err.Error())
	}

	if login != userResponse.Login {
		return false, fmt.Sprintf("invalid login. expected: %s, actual: %s", login, userResponse.Login), nil
	}
	if email != userResponse.Email {
		return false, fmt.Sprintf("invalid email. expected: %s, actual: %s", email, userResponse.Email), nil
	}
	if userResponse.ExpireDate != nil {
		return false, fmt.Sprintf("invalid exp date. expected: %s, actual: %s", nil, *userResponse.ExpireDate), nil
	}

	// Subscribe

	resParams, err = s.apiTestClient.Subscribe()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to subscribe user: %s", string(resParams.Body))
	}

	// User

	resParams, err = s.apiTestClient.User()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user info: %s", string(resParams.Body))
	}

	var newUserResponse api.UserResposne

	err = json.Unmarshal(resParams.Body, &newUserResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal user response: %s", err.Error())
	}

	if login != newUserResponse.Login {
		return false, fmt.Sprintf("invalid login. expected: %s, actual: %s", login, newUserResponse.Login), nil
	}
	if email != newUserResponse.Email {
		return false, fmt.Sprintf("invalid email. expected: %s, actual: %s", email, newUserResponse.Email), nil
	}
	if newUserResponse.ExpireDate == nil {
		return false, fmt.Sprintf("invalid exp date: %s", *newUserResponse.ExpireDate), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) TestUserEdit() (bool, string, error) {
	defer s.logout()

	login := "userEditLogin"
	email := "userEdit@email.com"
	password := "userEditPassword"

	_, err := s.RegisterAndSignIn(login, email, password)
	if err != nil {
		return false, "", fmt.Errorf("failed to log in: %s", err.Error())
	}

	resParams, err := s.apiTestClient.User()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user info: %s", string(resParams.Body))
	}

	var userResponse api.UserResposne

	err = json.Unmarshal(resParams.Body, &userResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal user response: %s", err.Error())
	}

	if login != userResponse.Login {
		return false, fmt.Sprintf("invalid login. expected: %s, actual: %s", login, userResponse.Login), nil
	}
	if email != userResponse.Email {
		return false, fmt.Sprintf("invalid email. expected: %s, actual: %s", email, userResponse.Email), nil
	}

	login = "newUserLogin"
	email = "newUserEdit@email.com"
	password = "newUserPassword"

	resParams, err = s.authTestClient.Edit(login, email, password)
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to edit user info: %s", string(resParams.Body))
	}

	var newUserResponse api.UserResposne

	resParams, err = s.apiTestClient.User()
	if err != nil {
		return false, "", err
	}

	if resParams.Code != http.StatusOK {
		return false, "", fmt.Errorf("failed to get user info: %s", string(resParams.Body))
	}

	err = json.Unmarshal(resParams.Body, &newUserResponse)
	if err != nil {
		return false, "", fmt.Errorf("failed to unmarshal user response: %s", err.Error())
	}

	if login != newUserResponse.Login {
		return false, fmt.Sprintf("invalid login. expected: %s, actual: %s", login, newUserResponse.Login), nil
	}
	if email != newUserResponse.Email {
		return false, fmt.Sprintf("invalid email. expected: %s, actual: %s", email, newUserResponse.Email), nil
	}

	return true, "", nil
}

func (s *InegrationTestsSuite) logout() {
	s.authClient.AccessToken = ""
	s.authClient.RefreshToken = ""
}

func unmarshalError(data []byte) (string, error) {
	var resp response.Response

	err := json.Unmarshal(data, &resp)
	if err != nil {
		return "", err
	}

	return resp.Error, nil
}
