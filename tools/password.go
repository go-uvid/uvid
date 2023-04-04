package tools

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	// Generate a salt with a cost of 10
	salt, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func ComparePassword(hashedPassword, password string) error {
	// Compare the password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
