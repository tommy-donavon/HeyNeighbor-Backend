package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func (s *testSuite) TestValidateAddress() {
	validate := NewValidator()

	prop := &Address{
		StreetAddress: "123 south",
		City:          "Salt lake City",
		State:         "UT",
		ZipCode:       "84101-1111",
	}
	prop2 := &Address{
		StreetAddress: "123 south",
		City:          "Salt lake City",
		State:         "UT",
		ZipCode:       "84101",
	}
	if !assert.NoError(s.T(), validate.Validate(prop)) {
		s.T().Fail()
	}
	if !assert.NoError(s.T(), validate.Validate(prop2)) {
		s.T().Fail()
	}
}

func (s *testSuite) TestPropertyValidation() {
	validate := NewValidator()

	// should pass
	prop := &Property{
		PropertyName: "pic",
		PropertyAddress: &Address{
			StreetAddress: "24 South 650 East",
			City:          "Salt Lake City",
			State:         "UT",
			ZipCode:       "84444",
		},
		PropertyManager: "shell",
		NumOfUnits:      400,
	}
	//should fail
	prop2 := &Property{
		PropertyName: "pic",
		PropertyAddress: &Address{
			StreetAddress: "24 South 650 East",
			City:          "Salt Lake City",
			State:         "UT",
			ZipCode:       "84444-12",
		},
		PropertyManager: "shell",
	}
	if !assert.NoError(s.T(), validate.Validate(prop)) {
		s.T().Fail()
	}
	assert.Error(s.T(), validate.Validate(prop2))

}

func (s *testSuite) TestTenantValidation() {
	validate := NewValidator()
	prop := &Property{
		PropertyName: "test",
		PropertyAddress: &Address{
			StreetAddress: "24 South 650 East",
			City:          "Salt Lake City",
			State:         "UT",
			ZipCode:       "84444-1211",
		},
		PropertyManager: "shell",
		NumOfUnits:      400,
		Tenants: []*Tenant{
			{
				Username:   "timmy",
				Nickname:   "timtim",
				UnitNumber: 102,
			},
		},
	}
	prop2 := &Property{
		PropertyName: "test",
		PropertyAddress: &Address{
			StreetAddress: "24 South 650 East",
			City:          "Salt Lake City",
			State:         "UT",
			ZipCode:       "84444-1211",
		},
		PropertyManager: "shell",
		NumOfUnits:      400,
	}

	assert.NoError(s.T(), validate.Validate(prop))
	assert.NoError(s.T(), validate.Validate(prop2))
}

func TestInit(t *testing.T) {
	suite.Run(t, new(testSuite))
}
