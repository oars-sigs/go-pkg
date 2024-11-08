package idaas

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID string
}

type UserObject struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	DeptId   string `json:"deptId"`
	DeptPath string `json:"deptPath"`
	AppId    string `json:"appId"`
	RealId   string `json:"realId"`
	Kind     int    `json:"kind"`
	Verified bool   `json:"verified"`
	PostId   string `json:"postId"`
	RoleId   string `json:"roleId"`
	GroupId  string `json:"groupId"`
}

func CreateToken(u *UserInfo, secret string, expiration int64) string {
	uo := &UserObject{
		UserId:   u.UserId,
		Username: u.Username,
		NickName: u.NickName,
		DeptId:   u.DeptId,
		AppId:    u.AppID,
		RealId:   u.RealId,
		Kind:     u.Kind,
		Verified: u.Verified,
		RoleId:   u.RoleId,
		PostId:   u.PostId,
		GroupId:  u.GroupId,
	}
	if u.Dept != nil {
		uo.DeptPath = u.Dept.Path
	}
	return CreateTokenWithObj(uo, secret, expiration)
}

func CreateTokenWithObj(uo *UserObject, secret string, expiration int64) string {
	expireToken := time.Now().Add(time.Second * time.Duration(expiration)).Unix()
	data, _ := json.Marshal(uo)
	claims := TokenClaims{
		UserID: uo.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "oars-idaas",
			Subject:   string(data),
		},
	}
	s, _ := base64.URLEncoding.DecodeString(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(s)
	return signedToken
}

func Valid(req *http.Request, secret string) (string, error) {
	tokenStr := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if tokenStr == "" {
		tokenStr = req.FormValue("token")
	}
	u, err := Parse(tokenStr, secret)
	if err != nil {
		return "", err
	}
	return u.UserId, nil
}

func ParseReq(req *http.Request, secret string) (*UserObject, error) {
	tokenStr := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if tokenStr == "" {
		tokenStr = req.FormValue("token")
	}
	u, err := Parse(tokenStr, secret)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func Parse(tokenStr, secret string) (*UserObject, error) {
	s, _ := base64.URLEncoding.DecodeString(secret)
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return s, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		var u UserObject
		err := json.Unmarshal([]byte(claims.Subject), &u)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}
	return nil, errors.New("must auth")
}
