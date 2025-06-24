package uservalidator_test

import (
	"errors"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByPhoneNumber(phone string) (*entity.User, error) {
	args := m.Called(phone)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockRepository) IsPhoneNumberUniq(phone string) (bool, error) {
	args := m.Called(phone)

	return args.Bool(0), args.Error(1)
}

func TestValidator_ValidateLoginRequest_Success(t *testing.T) {

	mockRepo := new(MockRepository)
	mockRepo.On("GetUserByPhoneNumber", "09123456789").
		Return(
			&entity.User{
				Id:             1,
				Name:           "",
				PhoneNumber:    "",
				HashedPassword: "",
				Role:           0,
			}, nil,
		)

	validator := uservalidator.NewValidator(mockRepo)

	req := &userparam.LoginRequest{
		PhoneNumber: "09123456789",
		Password:    "Valid@123",
	}

	fieldErrors, err := validator.ValidateLoginRequest(req)

	assert.Nil(t, err)
	assert.Nil(t, fieldErrors)
	mockRepo.AssertExpectations(t)
}

func TestValidator_ValidateLoginRequest_InvalidPhoneNumber(t *testing.T) {

	t.Run("Invalid_PhoneNumber_NotExistce", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetUserByPhoneNumber", "09123456789").
			Return(nil, errors.New("not found"))

		validator := uservalidator.NewValidator(mockRepo)

		req := &userparam.LoginRequest{
			PhoneNumber: "09123456789",
			Password:    "Valid@123",
		}

		fieldErrors, err := validator.ValidateLoginRequest(req)

		assert.NotNil(t, err)
		assert.NotNil(t, fieldErrors)

		assert.Contains(t, err.Error(), errormessage.ErrorMsgInvalidRequest)
		assert.Contains(t, fieldErrors["phone_number"], errormessage.ErrorMsgNotExistPhoneNumber)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid_PhoneNumber_Regex", func(t *testing.T) {

		type testCase struct {
			UserId      uint
			Phone       string
			Password    string
			FieldErrors map[string]string
			Err         string
		}

		testCases := []testCase{
			{
				UserId:      1,
				Phone:       "1234567891",
				Password:    "Valid@123",
				FieldErrors: map[string]string{"phone_number": errormessage.ErrorMsgInvalidPhoneNumberRegex1},
				Err:         errormessage.ErrorMsgInvalidRequest,
			},
			{
				UserId:      1,
				Phone:       "09191234ffdgd",
				Password:    "Valid@123",
				FieldErrors: map[string]string{"phone_number": errormessage.ErrorMsgInvalidPhoneNumberRegex2},
				Err:         errormessage.ErrorMsgInvalidRequest,
			},
		}

		for _, tc := range testCases {

			mockRepo := new(MockRepository)
			mockRepo.On("GetUserByPhoneNumber", tc.Phone).
				Return(
					&entity.User{
						Id: tc.UserId,
					},
					nil,
				)

			validator := uservalidator.NewValidator(mockRepo)

			req := &userparam.LoginRequest{
				PhoneNumber: tc.Phone,
				Password:    tc.Password,
			}

			fieldErrors, err := validator.ValidateLoginRequest(req)

			assert.NotNil(t, err)
			assert.NotNil(t, fieldErrors)

			assert.Contains(t, err.Error(), tc.Err)
			assert.Contains(t, fieldErrors["phone_number"], tc.FieldErrors["phone_number"])

			mockRepo.AssertExpectations(t)
		}
	})
}

func TestValidator_ValidateRegisterRequest_WeakPassword(t *testing.T) {

	type infoTestCase struct {
		Name        string
		Phone       string
		Password    string
		FieldErrors map[string]string
		Err         string
	}
	testCases := []infoTestCase{
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "p",
			FieldErrors: map[string]string{"password": "the length must be between 8 and 50"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "",
			FieldErrors: map[string]string{"password": "cannot be blank"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "valid@123",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "valid1234",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "valid@valid",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "VALID@123",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "VALIDVALID",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "validvalid",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "12345678",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "user user",
			Phone:       "09191234567",
			Password:    "@##&@#@#",
			FieldErrors: map[string]string{"password": errormessage.ErrorMsgInvalidPasswordRegex},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
	}

	for _, tc := range testCases {

		mockRepo := new(MockRepository)
		mockRepo.On("IsPhoneNumberUniq", tc.Phone).
			Return(true, nil)

		validator := uservalidator.NewValidator(mockRepo)

		req := &userparam.RegisterRequest{
			Name:        tc.Name,
			PhoneNumber: tc.Phone,
			Password:    tc.Password,
		}

		fieldErrors, err := validator.ValidateRegisterRequest(req)

		assert.NotNil(t, err)
		assert.NotNil(t, fieldErrors)

		assert.Contains(t, err.Error(), tc.Err)
		assert.Contains(t, fieldErrors["password"], tc.FieldErrors["password"])

		mockRepo.AssertExpectations(t)
	}
}

func TestValidator_ValidateRegisterRequest_InvalidName(t *testing.T) {
	type infoTestCase struct {
		Name        string
		Phone       string
		Password    string
		FieldErrors map[string]string
		Err         string
	}
	testCases := []infoTestCase{
		{
			Name:        "user user1",
			Phone:       "09191234567",
			Password:    "Valid@123",
			FieldErrors: map[string]string{"name": "must be in a valid format"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "1234",
			Phone:       "09191234567",
			Password:    "Valid@123",
			FieldErrors: map[string]string{"name": "must be in a valid format"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "1234 4321",
			Phone:       "09191234567",
			Password:    "Valid@123",
			FieldErrors: map[string]string{"name": "must be in a valid format"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        "useruser ",
			Phone:       "09191234567",
			Password:    "Valid@123",
			FieldErrors: map[string]string{"name": "must be in a valid format"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
		{
			Name:        " useruser",
			Phone:       "09191234567",
			Password:    "Valid@123",
			FieldErrors: map[string]string{"name": "must be in a valid format"},
			Err:         errormessage.ErrorMsgInvalidRequest,
		},
	}

	for _, tc := range testCases {

		mockRepo := new(MockRepository)
		mockRepo.On("IsPhoneNumberUniq", tc.Phone).
			Return(true, nil)

		validator := uservalidator.NewValidator(mockRepo)

		req := &userparam.RegisterRequest{
			Name:        tc.Name,
			PhoneNumber: tc.Phone,
			Password:    tc.Password,
		}

		fieldErrors, err := validator.ValidateRegisterRequest(req)

		assert.NotNil(t, err)
		assert.NotNil(t, fieldErrors)

		assert.Contains(t, err.Error(), tc.Err)
		assert.Contains(t, fieldErrors["name"], tc.FieldErrors["name"])

		mockRepo.AssertExpectations(t)
	}
}
