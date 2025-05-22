package hasher

import (
	"github.com/Chaitu35/claimeasy/backend/utils"
	"golang.org/x/crypto/bcrypt"
)
type PasswordHasher interface{
	Hash(Password string) (string, error)
	Compare(hashedPassword, password string) error

}
 
type BcryptHasher struct {}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}


func (b *BcryptHasher) Hash(password string) (string, error) {
	// Implement the hashing logic here
	return utils.HashPassword(password)
}

func (b *BcryptHasher) Compare(hashedPassword, password string) error {
	err:=bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {	return err}
	return nil 
}