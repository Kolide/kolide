package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/errors"
	"github.com/kolide/kolide-ose/kolide"
)

// swagger:parameters GetUser
type GetUserRequestBody struct {
	ID uint `json:"id"`
}

// swagger:response GetUserResponseBody
type GetUserResponseBody struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Admin              bool   `json:"admin"`
	Enabled            bool   `json:"enabled"`
	NeedsPasswordReset bool   `json:"needs_password_reset"`
}

// swagger:route POST /api/v1/kolide/user GetUser
//
// Get information on a user
//
// Using this API will allow the requester to inspect and get info on
// other users in the application.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func GetUser(c *gin.Context) {
	var body GetUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.IsLoggedIn() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformReadActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters CreateUser
type CreateUserRequestBody struct {
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Admin              bool   `json:"admin"`
	NeedsPasswordReset bool   `json:"needs_password_reset"`
}

// swagger:route PUT /api/v1/kolide/user CreateUser
//
// Create a new user
//
// Using this API will allow the requester to create a new user with the ability
// to control various user settings
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func CreateUser(c *gin.Context) {
	var body CreateUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.IsAdmin() {
		UnauthorizedError(c)
		return
	}

	svc := GetSVC(c)
	user := &kolide.User{
		Username:           body.Username,
		Password:           []byte(body.Password),
		Email:              body.Email,
		Admin:              body.Admin,
		Enabled:            true,
		NeedsPasswordReset: body.NeedsPasswordReset,
	}

	user, err = svc.NewUser(user)
	if err != nil {
		logrus.Errorf("Error creating new user: %s", err.Error())
		errors.ReturnError(c, err)
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters ModifyUser
type ModifyUserRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// swagger:route PATCH /api/v1/kolide/user ModifyUser
//
// Update a user's basic information and settings
//
// Using this API will allow the requester to update a user's basic settings.
// Note that updating administrative settings are not exposed via this endpoint
// as this is primarily intended to be used by users to update their own
// settings.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func ModifyUser(c *gin.Context) {
	var body ModifyUserRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters ChangeUserPassword
type ChangePasswordRequestBody struct {
	ID                 uint   `json:"id"`
	CurrentPassword    string `json:"current_password"`
	PasswordResetToken string `json:"password_reset_token"`
	NewPassword        string `json:"new_password" validate:"required"`
	NewPasswordConfim  string `json:"new_password_confirm" validate:"required"`
}

// swagger:route PATCH /api/v1/kolide/user/password ChangeUserPassword
//
// Change a user's password
//
// Using this API will allow the requester to change their password. Users
// should include their own user id as the "id" paramater and/or their own
// username as the "username" parameter. Admins can change the passords for
// other users by defining their ID or username.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func ChangeUserPassword(c *gin.Context) {
	var body ChangePasswordRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	if body.NewPassword != body.NewPasswordConfim {
		errors.ReturnError(c, errors.NewWithStatus(http.StatusBadRequest, "Submitted passwords do not match", ""))
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	svc := GetSVC(c)
	user, err := svc.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	deleteResetRequest := func(reset *kolide.PasswordResetRequest) {
		if err != nil { //TODO: this shadows error returns
			err = db.DeletePasswordResetRequest(reset)
			if err != nil {
				errors.ReturnError(c, errors.DatabaseError(err))
				return
			}
		}
	}

	if body.PasswordResetToken != "" {
		reset, err := db.FindPassswordResetByToken(body.PasswordResetToken)
		if err != nil {
			UnauthorizedError(c)
			return
		}

		if time.Now().After(reset.ExpiresAt) {
			deleteResetRequest(reset)
			UnauthorizedError(c)
			return
		}
		defer deleteResetRequest(reset)
	} else if !vc.IsAdmin() {
		if body.CurrentPassword != "" {
			if err := user.ValidatePassword(body.CurrentPassword); err != nil {
				UnauthorizedError(c)
				return
			}
		} else {
			UnauthorizedError(c)
			return
		}
	}

	err = svc.SetPassword(user.ID, body.NewPassword)
	if err != nil {
		logrus.Errorf("Error setting user password: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err)) // probably not this
		return
	}

	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters SetUserAdminState
type SetUserAdminStateRequestBody struct {
	ID    uint `json:"id"`
	Admin bool `json:"admin"`
}

// swagger:route PATCH /api/v1/kolide/user/admin SetUserAdminState
//
// Modify a user's admin settings
//
// This endpoint allows an existing admin to promote a non-admin to admin or
// demote a current admin to non-admin.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func SetUserAdminState(c *gin.Context) {
	var body SetUserAdminStateRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.IsAdmin() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	user.Admin = body.Admin
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

// swagger:parameters SetUserEnabledState
type SetUserEnabledStateRequestBody struct {
	ID              uint   `json:"id"`
	Enabled         bool   `json:"enabled"`
	CurrentPassword string `json:"current_password"`
}

// swagger:route PATCH /api/v1/kolide/user/enabled SetUserEnabledState
//
// Enable or disable a user.
//
// This endpoint allows an existing admin to enable a disabled user or
// disable an enabled user. If a user calls this endpoint, to disable,
// their own account, they must also submit their current password, to
// verify their request.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: GetUserResponseBody
func SetUserEnabledState(c *gin.Context) {
	var body SetUserEnabledStateRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		errors.ReturnError(c, err)
		return
	}

	vc := VC(c)
	if !vc.CanPerformActions() {
		UnauthorizedError(c)
		return
	}

	db := GetDB(c)
	user, err := db.UserByID(body.ID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	if !vc.IsAdmin() {
		if user.ValidatePassword(body.CurrentPassword) != nil {
			UnauthorizedError(c)
			return
		}
	}

	user.Enabled = body.Enabled
	err = db.SaveUser(user)
	if err != nil {
		logrus.Errorf("Error updating user in database: %s", err.Error())
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}
	c.JSON(http.StatusOK, GetUserResponseBody{
		ID:                 user.ID,
		Username:           user.Username,
		Name:               user.Name,
		Email:              user.Email,
		Admin:              user.Admin,
		Enabled:            user.Enabled,
		NeedsPasswordReset: user.NeedsPasswordReset,
	})
}

////////////////////////////////////////////////////////////////////////////////
// Password Reset HTTP endpoints
////////////////////////////////////////////////////////////////////////////////

// swagger:parameters ResetUserPassword
type ResetPasswordRequestBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ResetPasswordResponseBody ...
// swagger:response ResetPasswordResponseBody
type ResetPasswordResponseBody struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// swagger:route POST /api/v1/kolide/user/password/reset ResetUserPassword
//
// Reset a user's password
//
// Using this API will allow the requester to reset their password. Users
// should include their own user id as the "id" paramater and/or their own
// username as the "username" parameter. Admins can change the passwords for
// other users by defining their ID or username. Logged out users can reset
// their own password by including their email in addition to either their
// user id or their username.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: ResetPasswordResponseBody
func ResetUserPassword(c *gin.Context) {
	var body ResetPasswordRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing ResetPassword post body: %s", err.Error())
		return
	}

	db := GetDB(c)
	vc := VC(c)

	if !vc.IsLoggedIn() {
		if body.Email == "" || body.Username == "" {
			UnauthorizedError(c)
			errors.NewWithStatus(http.StatusBadRequest, "Must submit email and username", "Required parameters were not submitted")
			return
		}
	}

	user, err := db.User(body.Username)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			errors.NewWithStatus(http.StatusBadRequest, "Invalid username", "Username not found in the database")
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}
	if user.Email != body.Email {
		errors.NewWithStatus(http.StatusUnauthorized, "Invalid email", "Supplied email does not match user record")
		return
	}

	// If an admin is requesting a password reset for a user, require that the
	// user change their password to log in
	if vc.IsAdmin() && !vc.IsUserID(user.ID) {
		user.NeedsPasswordReset = true

		err = db.SaveUser(user)
		if err != nil {
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	if vc.IsAdmin() || vc.IsUserID(user.ID) || !vc.IsLoggedIn() {
		// logged-in admin user resetting a password or logged-in user
		// resetting their own password or logged-out user presumably resetting
		// their own password

		request, err := kolide.NewPasswordResetRequest(GetDB(c), user.ID, time.Now().Add(time.Hour*24))
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Database error"))
			return
		}

		html, text, err := kolide.GetEmailBody(
			kolide.PasswordResetEmail,
			kolide.PasswordResetRequestEmailParameters{
				Name:  user.Name,
				Token: request.Token,
			},
		)
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}

		subject, err := kolide.GetEmailSubject(kolide.PasswordResetEmail)
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}

		err = kolide.SendEmail(GetSMTPConnectionPool(c), user.Email, subject, html, text)
		if err != nil {
			errors.ReturnError(c, errors.NewFromError(err, http.StatusInternalServerError, "Email error"))
			return
		}
	} else {
		// Logged-in user trying to reset another user's password
		UnauthorizedError(c)
		return
	}

	c.JSON(http.StatusOK, ResetPasswordResponseBody{
		ID:       user.ID,
		Username: user.Username,
	})
}

// swagger:parameters VerifyPasswordResetRequest
type VerifyPasswordResetRequestRequestBody struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

// swagger:parameters VerifyPasswordResetRequestResponseBody
type VerifyPasswordResetRequestResponseBody struct {
	Valid bool `json:"valid"`
}

// swagger:route POST /api/v1/kolide/user/password/verify VerifyPasswordResetRequest
//
// Verify an email campaign before it is used.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: VerifyPasswordResetRequestResponseBody
func VerifyPasswordResetRequest(c *gin.Context) {
	var body VerifyPasswordResetRequestRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing request body: %s", err.Error())
		return
	}

	db := GetDB(c)
	_, err = db.UserByID(body.UserID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusBadRequest, &VerifyPasswordResetRequestResponseBody{
				Valid: false,
			})
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	reset, err := db.FindPassswordResetByTokenAndUserID(body.Token, body.UserID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusBadRequest, VerifyPasswordResetRequestResponseBody{
				Valid: false,
			})
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, VerifyPasswordResetRequestResponseBody{
			Valid: false,
		})
		return
	}

	c.JSON(http.StatusOK, VerifyPasswordResetRequestResponseBody{
		Valid: true,
	})
}

// swagger:parameters DeletePasswordResetRequest
type DeletePasswordResetRequestRequestBody struct {
	ID uint `json:"id" validate:"required"`
}

// swagger:route DELETE /api/v1/kolide/user/password/reset DeletePasswordResetRequest
//
// Delete an email campaign after it has been used.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Security:
//       authenticated: yes
//
//     Responses:
//       200: nil
func DeletePasswordResetRequest(c *gin.Context) {
	var body DeletePasswordResetRequestRequestBody
	err := ParseAndValidateJSON(c, &body)
	if err != nil {
		logrus.Errorf("Error parsing request body: %s", err.Error())
		return
	}

	db := GetDB(c)
	campaign, err := db.FindPassswordResetByID(body.ID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			NotFoundRequestError(c)
			return
		default:
			errors.ReturnError(c, errors.DatabaseError(err))
			return
		}
	}

	user, err := db.UserByID(campaign.UserID)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	vc := VC(c)
	if !vc.CanPerformWriteActionOnUser(user) {
		UnauthorizedError(c)
		return
	}

	err = db.DeletePasswordResetRequest(campaign)
	if err != nil {
		errors.ReturnError(c, errors.DatabaseError(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
