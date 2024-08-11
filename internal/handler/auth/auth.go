package auth

import (
	"fmt"

	"net/http"
	"strconv"
	"time"

	"github.com/RobinHoodArmyHQ/robin-api/internal/env"
	userrepo "github.com/RobinHoodArmyHQ/robin-api/internal/repositories/user"
	userverification "github.com/RobinHoodArmyHQ/robin-api/internal/repositories/userVerification"
	"github.com/RobinHoodArmyHQ/robin-api/internal/util"
	"github.com/RobinHoodArmyHQ/robin-api/models"
	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func AuthHandler(c *gin.Context) {

	// validate country isd code
	countryCode, err := strconv.ParseUint(c.PostForm("country_code"), 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Status: models.StatusFailed(fmt.Sprintf("invalid country code %d", countryCode)),
		})
		return
	}

	// validate mobile number
	mobileNumber, err := strconv.ParseUint(c.PostForm("mobile_number"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Status: models.StatusFailed(fmt.Sprintf("invalid mobile number %d", mobileNumber)),
		})
		return
	}

	// generate request id and send response
	requestId := uuid.Must(uuid.NewRandom())
	c.JSON(http.StatusOK, models.AuthResponse{
		Status:    models.StatusSuccess(),
		RequestID: requestId,
	})
}

func RegisterUser(c *gin.Context) {
	var request RegisterUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, RegisterUserResponse{
			Status: models.StatusFailed(err.Error()),
		})
		return
	}

	// validate user inputs
	err := validateUserInputs(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, RegisterUserResponse{
			Status: models.StatusFailed(err.Error()),
		})
		return
	}

	registerUser := &userrepo.GetUserByEmailRequest{
		EmailId: request.EmailId,
	}

	userRepo := env.FromContext(c).UserRepository
	user, err := userRepo.GetUserByEmail(registerUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if user != nil {
		c.JSON(http.StatusOK, RegisterUserResponse{
			Status:    models.StatusSuccess(),
			IsNewUser: 0,
		})
		return
	}

	// creating hashed password
	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// generate 6 digit OTP
	otp := util.GenerateOtp(6)

	// convert otp string to uint64
	uiOtp, err := strconv.ParseUint(otp, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// create a new user in user_verificatons table
	userInfo := models.UserInfo{
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		PasswordHash: hashedPassword,
	}

	newUserData := &models.UserVerification{
		EmailId:        request.EmailId,
		Otp:            uiOtp,
		OtpGeneratedAt: time.Now(),
		OtpExpiresAt:   time.Now().Add(10 * time.Minute),
		UserInfo:       userInfo,
	}

	userVerificationRepo := env.FromContext(c).UserVerificationRepository

	newUser, err := userVerificationRepo.CreateUser(&userverification.CreateUserRequest{
		User: newUserData,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, RegisterUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// TO-DO: send verification otp via aws-ses

	c.JSON(http.StatusCreated, RegisterUserResponse{
		UserID: newUser.UserID.String(),
		Status: models.StatusSuccess(),
	})
}

func LoginUser(c *gin.Context) {
	var request LoginUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, LoginUserResponse{
			Status: models.StatusFailed(fmt.Sprintln("Invalid credentials")),
		})
		return
	}

	userRepo := env.FromContext(c).UserRepository

	loginUser := &userrepo.GetUserByEmailRequest{
		EmailId: request.EmailId,
	}
	user, err := userRepo.GetUserByEmail(loginUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, LoginUserResponse{
			Status: models.StatusFailed(fmt.Sprintln("Incorrect email or password")),
		})
		return
	}

	// verify user password
	ok := util.CheckPasswordHash(request.Password, user.User.PasswordHash)

	if !ok {
		c.JSON(http.StatusBadRequest, LoginUserResponse{
			Status: models.StatusFailed(fmt.Sprintln("Incorrect email or password")),
		})
		return
	}

	// create a JWT token
	jwtToken, err := util.GenerateJwt(user.User.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginUserResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginUserResponse{
		Status: models.StatusSuccess(),
		Token:  jwtToken,
	})
}

func VerifyOtp(c *gin.Context) {
	var request VerifyOtpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, VerifyOtpResponse{
			Status: models.StatusFailed(err.Error()),
		})
		return
	}

	userRepo := env.FromContext(c).UserRepository
	userVerificationRepo := env.FromContext(c).UserVerificationRepository

	// get user by user_id
	user, err := userVerificationRepo.GetUserByUserID(&userverification.GetUserByUserIdRequest{UserID: nanoid.NanoID(request.UserID)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, VerifyOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, VerifyOtpResponse{
			Status: models.StatusFailed(fmt.Sprintf("No user found with given user_id: %s", request.UserID)),
		})
		return
	}

	// check if we have already created a user with this users email_id
	existingUser, err := userRepo.GetUserByEmail(&userrepo.GetUserByEmailRequest{EmailId: user.User.EmailId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, VerifyOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, VerifyOtpResponse{
			Status: models.StatusFailed("User already verified, please login to continue"),
		})
		return
	}

	// check if otp has expired
	currTime := time.Now()
	if currTime.After(user.User.OtpExpiresAt) {
		c.JSON(http.StatusOK, VerifyOtpResponse{
			Status: models.StatusFailed("Verification code expired, re-send code to continue"),
		})
		return
	}

	// match the otp
	if request.Otp != user.User.Otp {
		c.JSON(http.StatusBadRequest, VerifyOtpResponse{
			Status: models.StatusFailed("Wrong Otp"),
		})
		return
	}

	newUser := &models.User{
		FirstName:    user.User.UserInfo.FirstName,
		LastName:     user.User.UserInfo.LastName,
		EmailId:      user.User.EmailId,
		PasswordHash: user.User.UserInfo.PasswordHash,
		UserID:       user.User.UserID,
	}

	// now create a new entry in users table
	createdUser, err := userRepo.CreateUser(&userrepo.CreateUserRequest{
		User: newUser,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, VerifyOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// set user verified in user_verifications table
	updateUser := &userverification.UpdateUserRequest{
		UserID: user.User.UserID,
		Values: map[string]interface{}{
			"is_verified": 1,
		},
	}

	if _, err := userVerificationRepo.UpdateUser(updateUser); err != nil {
		c.JSON(http.StatusInternalServerError, VerifyOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// create a new jwt-token
	token, err := util.GenerateJwt(createdUser.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, VerifyOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	c.JSON(http.StatusOK, VerifyOtpResponse{
		Status: models.StatusSuccess(),
		Token:  token,
	})
}

func ResendOtp(c *gin.Context) {
	var request ResendOtpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ResendOtpResponse{
			Status: models.StatusFailed("Missing Params"),
		})
		return
	}

	userVerificationRepo := env.FromContext(c).UserVerificationRepository

	user, err := userVerificationRepo.GetUserByUserID(&userverification.GetUserByUserIdRequest{
		UserID: nanoid.NanoID(request.UserID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ResendOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// update otp_retry_count
	updateUser := &userverification.UpdateUserRequest{
		UserID: nanoid.NanoID(request.UserID),
		Values: map[string]interface{}{
			"otp_expires_at":  time.Now().Add(10 * time.Minute),
			"otp_retry_count": user.User.OtpRetryCount + 1,
		},
	}

	if _, err := userVerificationRepo.UpdateUser(updateUser); err != nil {
		c.JSON(http.StatusInternalServerError, ResendOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// TO-DO resend verification code

	c.JSON(http.StatusOK, ResendOtpResponse{
		Status: models.StatusSuccess(),
	})
}

func SendPasswordResetLink(c *gin.Context) {
	var request SendResetPasswordLinkRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, SendResetPasswordLinkResponse{
			Status: models.StatusFailed("Invalid inputs"),
		})
		return
	}

	// get user by email
	userRepo := env.FromContext(c).UserRepository
	user, err := userRepo.GetUserByEmail(&userrepo.GetUserByEmailRequest{EmailId: request.EmailId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, SendResetPasswordLinkResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, SendResetPasswordLinkResponse{
			Status: models.StatusFailed("No user exist with such email"),
		})
		return
	}

	// generate nanoid
	randomID, err := nanoid.GetID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, SendResetPasswordLinkResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	linkExpireTimeStamp := time.Now().Add(2 * time.Minute).UTC().Unix()
	userInfoStr := util.GetUserInfoStr(user.User.UserID.String(), randomID.String(), linkExpireTimeStamp)
	userHashValue := util.GenerateHashCode(userInfoStr)

	// add user_id in password reset link
	resetPasswordLink := fmt.Sprintf("%s?code=%s&user_id=%s&timestamp=%d&token=%s", viper.GetString("auth.password_reset_link"), userHashValue, user.User.UserID, linkExpireTimeStamp, randomID)

	// TO-DO send link on the registered/verified email
	fmt.Println(resetPasswordLink)

	c.JSON(http.StatusOK, SendResetPasswordLinkResponse{
		Status: models.StatusSuccess(),
	})
}

func ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ResetPasswordResponse{
			Status: models.StatusFailed("Invalid inputs"),
		})
		return
	}

	// password validations
	err := validateUserPassword(request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResetPasswordResponse{
			Status: models.StatusFailed(err.Error()),
		})
		return
	}

	userRepo := env.FromContext(c).UserRepository
	user, err := userRepo.GetUser(&userrepo.GetUserRequest{UserID: nanoid.NanoID(request.UserID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResetPasswordResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	if user == nil || user.User == nil {
		c.JSON(http.StatusBadRequest, ResetPasswordResponse{
			Status: models.StatusFailed("No user found with this user_id"),
		})
		return
	}

	timestamp, err := strconv.ParseInt(request.Timestamp, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResetPasswordResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	// check if link is expired.
	currTimestamp := time.Now().UTC().Unix()
	if currTimestamp > timestamp {
		c.JSON(http.StatusBadRequest, ResendOtpResponse{
			Status: models.StatusFailed("Reset password link expired"),
		})
		return
	}

	// validate link
	err = validateResetLink(request.UserID, request.Token, request.Code, timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResetPasswordResponse{
			Status: models.StatusFailed(err.Error()),
		})
		return
	}

	passwordHash, err := util.HashPassword(request.NewPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ResendOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	updateUser := &userrepo.UpdateUserRequest{
		UserID: nanoid.NanoID(request.UserID),
		Values: map[string]interface{}{
			"password_hash": passwordHash,
		},
	}

	if _, err := userRepo.UpdateUser(updateUser); err != nil {
		c.JSON(http.StatusInternalServerError, ResendOtpResponse{
			Status: models.StatusSomethingWentWrong(),
		})
		return
	}

	c.JSON(http.StatusOK, ResetPasswordResponse{
		Status: models.StatusSuccess(),
	})
}
