package sha256

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Sha256HashService struct {
	Key string
}

func NewSha256HashService(key string) authlib.AuthService {
	return Sha256HashService{Key: key}
}

func (h Sha256HashService) GetHash(pwd string) string {
	data := hmac.New(sha256.New, []byte(h.Key))
	data.Write([]byte(pwd))
	return fmt.Sprintf("%x", data.Sum(nil))
}

func (h Sha256HashService) Equals(hash, pwd string) bool {
	mHash := []byte(h.GetHash(pwd))
	return bytes.Equal(mHash, []byte(hash))
}

func (h Sha256HashService) Validate(hash, pwd string) error {
	if !h.Equals(hash, pwd) {
		return errors.New("invalid hash")
	}
	return nil
}

func (h Sha256HashService) DecodeToken(token string) (string, string, error) {
	data, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		log.WithFields(log.Fields{"func": "DecodeToken"}).Errorln(err)
		return "", "", err
	}

	loginPwdArr := strings.Split(string(data), ":")
	login, pwdHash := loginPwdArr[0], loginPwdArr[1]

	return login, pwdHash, nil
}

func (h Sha256HashService) EncodeToken(login, pwdHash string) string {
	data := fmt.Sprintf("%s:%s", login, pwdHash)
	return base64.StdEncoding.EncodeToString([]byte(data))
}
