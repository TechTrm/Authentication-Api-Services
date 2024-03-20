package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	mail "github.com/TechTrm/Authentication-Api-Services/mail"
	"github.com/TechTrm/Authentication-Api-Services/util"
	"github.com/gin-gonic/gin"
)


type sendEmailRequest struct {
	UserEmail string `json:"user_email" binding:"required,email"`
}

func(server *Server) sendEmail(ctx *gin.Context){
	var req sendEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	requestedUserNameOrEmail := db.GetUserByNameOrEmailParams{
		Username: req.UserEmail,
		Email:   req.UserEmail,
 }
	user, err := server.store.GetUserByNameOrEmail(ctx, requestedUserNameOrEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}


	verifyEmail, err := server.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserID:      user.ID,
		SecretCode: util.RandomString(64),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sender := mail.NewGmailSender(server.config.EmailSenderName , server.config.EmailSenderAddress , server.config.EmailSenderPassword)


	subject := "Comfirm your email address to reset your password"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/api/reset_password/verify_email?verify_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Welcome to Authencation API Service!<br/>
	Please <a href="%s">click here</a> to verify your email address to reset your password.<br/>API-URL=%s
	`, user.FullName, verifyUrl, verifyUrl)
	to := []string{user.Email}

    if err := sender.SendEmail(subject, content, to, nil, nil, nil); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	msg := map[string]string{
		"message": "Successfully send and email to reset your password. Please check your email to verify your email address.",
	}
	
	ctx.JSON(http.StatusOK, msg)

}

type verifyEmailRequest struct {
	VerifyID    int64    `json:"verify_id" binding:"required"`
	SecretCode string `json:"secret_code" binding:"required"`
}


func(server *Server) verifyEmail(ctx *gin.Context){
	var req verifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	verifyEmail, err := server.store.GetVerifyEmail(ctx, req.VerifyID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if verifyEmail.SecretCode != req.SecretCode {
		err := fmt.Errorf("invalid secret code")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if time.Now().After(verifyEmail.ExpiredAt) {
		err := fmt.Errorf("email verification token has expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if verifyEmail.IsUsed {
	  	err := fmt.Errorf("token has been used")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err = server.store.UpdateVerifyEmail(ctx, db.UpdateVerifyEmailParams{
		ID: verifyEmail.ID,
		SecretCode: verifyEmail.SecretCode,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: render A new page to reset password

	msg := map[string]interface{}{
		"message": "Email Successfully Verified Now You Can reset your password. Using Blow Link within 30 minutes, Otherwise this link will be expired.",
		"reset_password_url": "http://localhost:8080/v1/api/reset_password",
		"field_name": map[string]any{
			"password": "example@123",
			"confirm_password": "example@123",
			"verify_id": verifyEmail.ID,
		},
	}
	
	ctx.JSON(http.StatusOK, msg)


}


type resetPasswordRequest struct {
	VerifyID    int64    `json:"verify_id" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}


func(server *Server) resetPassword(ctx *gin.Context){
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	verifyEmail, err := server.store.GetVerifyEmail(ctx, req.VerifyID)
	if err != nil {
		if err == sql.ErrNoRows {
			err := fmt.Errorf("verify id not found")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !verifyEmail.IsUsed {
		err := fmt.Errorf("email not verified")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if time.Now().After(verifyEmail.ExpiredAt) {
		err := fmt.Errorf("reset password time has expired, Max time is 30 minutes")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if req.Password != req.ConfirmPassword {
		err := fmt.Errorf("password and confirm password does not match")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		HashedPassword: hashedPassword,
		PasswordChangedAt: time.Now(),
		ID: verifyEmail.UserID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	msg := map[string]string{
		"message": "Password Successfully Reset",
	}

	ctx.JSON(http.StatusOK, msg)

}
