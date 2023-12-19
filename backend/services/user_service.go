package services

import (
	"net/mail"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"umigame-api/models"
	"umigame-api/myerrors"
	"umigame-api/repositories"
	"umigame-api/utils"
)

func (s *Service) RegisterUserService(user models.User) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		err = myerrors.EmailInvalid.Wrap(err, "email invalid")
		return err
	}

	hash, err := utils.PasswordEncrypt(user.Password)
	if err != nil {
		err = myerrors.EncryptFailed.Wrap(err, "internal server error")
		return err
	}
	user.PasswordHash = hash

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		err = myerrors.GenUUIDFailed.Wrap(err, "internal server error")
		return err
	}
	mailAuthUuid := uuidObj.String()
	user.UUID = mailAuthUuid

	if err = repositories.RegisterUser(s.db, &user); err != nil {
		return err
	}

	mailMessage := utils.MailMessage(s.port, mailAuthUuid)
	myMail := utils.Mail{
		Host:     "smtp.gmail.com",
		Port:     "587",
		From:     os.Getenv("EMAIL_ADDRESS"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		To:       []string{user.Email},
		Subject:  "メールアドレスの認証",
		Message:  mailMessage,
	}
	if err := myMail.SendMail(); err != nil {
		err = myerrors.SendMailFailed.Wrap(err, "internal server error")
		return err
	}

	return nil
}

func (s *Service) MailCheckService(uuid string) error {
	err := repositories.UpdateActivate(s.db, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) LoginService(email string, password string) (models.User, error) {
	user, err := repositories.GetAuthInfo(s.db, email)
	if err != nil {
		return models.User{}, err
	}

	if !user.IsValid {
		err := myerrors.NotActivate.Wrap(NotActivate, "email is not valified")
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		err = myerrors.InvalidPassword.Wrap(err, "password is not correnct")
		return models.User{}, err
	}

	response := models.User{
		UUID: user.UUID,
	}

	return response, nil
}