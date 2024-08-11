package auth

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/RobinHoodArmyHQ/robin-api/internal/util"
)

const (
	passwordMinLength  = 8
	onlyAlphabetsRegex = `^[a-zA-Z]+$`
	emailRegex         = `^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`
)

func validateUserInputs(req *RegisterUserRequest) error {
	var err error
	// firstName and lastName validation
	if err = validateUserName(req.FirstName, req.LastName); err != nil {
		return err
	}

	// email validation
	if err = validateUserEmailID(req.EmailId); err != nil {
		return err
	}

	// password validation
	if err = validateUserPassword(req.Password); err != nil {
		return err
	}

	return nil
}

func validateUserPassword(password string) error {
	var hasNumbers, hasLetters, hasSpecialCharacter bool

	if len(password) < passwordMinLength {
		return fmt.Errorf("password should be of at least %d characters", passwordMinLength)
	}

	// special character match
	for _, c := range password {
		if unicode.IsNumber(c) {
			hasNumbers = true
			continue
		}

		if unicode.IsLetter(c) {
			hasLetters = true
			continue
		}

		if unicode.IsSymbol(c) || unicode.IsPunct(c) {
			hasSpecialCharacter = true
			continue
		}
	}

	if !hasNumbers {
		return fmt.Errorf("password should contain a number")
	}

	if !hasLetters {
		return fmt.Errorf("password should contain alphabets")
	}

	if !hasSpecialCharacter {
		return fmt.Errorf("password should contain a special character")
	}

	return nil
}

func validateUserEmailID(emailID string) error {
	r, err := regexp.Compile(emailRegex)
	if err != nil {
		return fmt.Errorf("error compile email regex")
	}

	if matched := r.MatchString(emailID); !matched {
		return fmt.Errorf("emailID should contain only alphabets")
	}

	return nil
}

func validateUserName(firstName, lastName string) error {
	r, err := regexp.Compile(onlyAlphabetsRegex)

	if err != nil {
		return fmt.Errorf("error compiling name regex")
	}

	if matched := r.MatchString(firstName); !matched {
		return fmt.Errorf("first name should contain only alphabets")
	}

	if matched := r.MatchString(lastName); !matched {
		return fmt.Errorf("last name should contain only alphabets")
	}

	return nil
}

func validateResetLink(userID, token, code string, timestamp int64) error {
	s := util.GetUserInfoStr(userID, token, timestamp)
	hash := util.GenerateHashCode(s)

	if hash != code {
		return fmt.Errorf("invalid link")
	}

	return nil
}
