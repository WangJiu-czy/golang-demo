package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"testing"
	"time"
)

func TestToken(t *testing.T) {
	claims := make(jwt.MapClaims)
	// 有效期
	claims[TokenClaimEXP] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims[TokenClaimUID] = uuid.New()

	token, err := GenJwtToken(claims)
	if err != nil {
		t.Error(err)
	}

	t.Log("token:", token)

	isToken := CheckJwtToken(token)
	t.Log("isToken:", isToken)

	if uid, found := GetUIDFromToken(token); found {
		t.Log("用户id：", uid)
	}
}
