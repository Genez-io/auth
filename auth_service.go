package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type User struct {
	Email             *string                 `json:"email"`
	UserId            string                  `json:"userId"`
	AuthProvider      string                  `json:"authProvider"`
	CreatedAt         time.Time               `json:"createdAt"`
	Verified          *bool                   `json:"verified"`
	Name              *string                 `json:"name"`
	Address           *string                 `json:"address"`
	ProfilePictureUrl *string                 `json:"profilePictureUrl"`
	CustomInfo        *map[string]interface{} `json:"customInfo"`
}

func GetUserByToken(token string) (*User, error) {
	remote := NewRemote(os.Getenv("GNZ_AUTH_FUNCTION_URL"))
	userInterface, err := remote.Call("AuthService.userInfo", token)
	if err != nil {
		return nil, err
	}
	// convert interface to User
	jsonUser, err := json.Marshal(userInterface)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal(jsonUser, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type ReqBody struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  interface{}   `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type ErrorStruct struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Info    *map[string]interface{} `json:"info,omitempty"`
}

func (e ErrorStruct) Error() string {
	return e.Message
}

type ResBody struct {
	Jsonrpc string       `json:"jsonrpc"`
	Error   *ErrorStruct `json:"error"`
	Result  interface{}  `json:"result"`
	Id      int          `json:"id"`
}

type Remote struct {
	URL string
}

func NewRemote(url string) Remote {
	return Remote{URL: url}
}

func (r Remote) Call(method interface{}, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		args = make([]interface{}, 0)
	}
	reqBody := ReqBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  args,
		Id:      3}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resBody ResBody
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Error != nil {
		return nil, resBody.Error
	}

	return resBody.Result, nil
}
