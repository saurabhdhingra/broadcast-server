package cmd

import "errors"


var validTokens = map[string]string{
	"user1": "token123",
	"user2": "token456",
}

func AuthenticateUser(username, token string) error {
	if validTokens[username] != token {
		return errors.New("invalid authentication token")
	}
	return nil
}
