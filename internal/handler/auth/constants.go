package auth

import (
	"fmt"
)

const (
	// user input validation constants
	passwordMinLength  = 8
	onlyAlphabetsRegex = `^[a-zA-Z]+$`
	emailRegex         = `[a-zA-Z0-9.*%±]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}`

	// email related constants
	verificationEmailSubject  string = "Robin Hood Army | Email Verification"
	resetPasswordEmailSubject string = "Robin Hood Army | Reset Password"
)

func getVerificationEmailBody(otp uint64) string {
	return fmt.Sprintf("Your verification code is :%d", otp)
}

func getResetPasswordEmailBody(link string) string {
	return fmt.Sprintf("Your reset password link in :%s", link)
}
