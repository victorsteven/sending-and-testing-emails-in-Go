package welcome

import (
	"github.com/gin-gonic/gin"
	"mail-sending/emails/welcome_mail"
	"mail-sending/helpers"
	"net/http"
)

type service struct {
	email welcome_mail.Service
}

func NewWelcome(email welcome_mail.Service) *service {
	return &service{email}
}

func (h *service) WelcomeMail(c *gin.Context) {

	welcome := &helpers.WelcomeModel{}

	if err := c.ShouldBindJSON(&welcome); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"body":   "Please provide valid inputs",
		})
		return
	}

	if err := helpers.ValidateInputs(*welcome); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"body":   err.Error(),
		})
		return
	}

	cred := &helpers.WelcomeMail{
		Name:  welcome.Name,
		Email: welcome.Email,
	}

	_, err := h.email.SendWelcomeMail(cred)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"body":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body":   "Please check your mail",
	})
	return

}
