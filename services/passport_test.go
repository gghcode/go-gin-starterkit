package services_test

import (
	"testing"

	"github.com/gghcode/go-gin-starterkit/services"
	"github.com/stretchr/testify/suite"
)

type passportUnit struct {
	suite.Suite

	passport services.Passport
}

func (suite *passportUnit) SetupTest() {
	suite.passport = services.NewPassport()
}

func TestPassportUnit(t *testing.T) {
	suite.Run(t, new(passportUnit))
}

func (suite *passportUnit) TestPasswordVerfication() {
	testCases := []struct {
		description    string
		password       string
		verifyPassword string
		expected       bool
	}{
		{
			description:    "ShouldBeValid",
			password:       "12345678",
			verifyPassword: "12345678",
			expected:       true,
		},
		{
			description:    "ShouldBeInvalid",
			password:       "12345678910",
			verifyPassword: "12345",
			expected:       false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			passwordHash, err := suite.passport.HashPassword(tc.password)
			suite.NoError(err)

			actual := suite.passport.IsValidPassword(tc.verifyPassword, passwordHash)
			suite.Equal(tc.expected, actual)
		})
	}
}
