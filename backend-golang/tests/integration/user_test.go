package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sample-stack-golang/internal/di"
	"sample-stack-golang/internal/modules/user/domain"
	"sample-stack-golang/tests/fixtures"
	"sample-stack-golang/tests/utils"
)

type UserTestSuite struct {
	suite.Suite
	ctx         context.Context
	container   *di.Container
	userService domain.UserUseCase
	db          *sql.DB
}

func (s *UserTestSuite) SetupSuite() {
	s.ctx = utils.TestContext(s.T())
	
	// Setup container
	s.container = di.NewContainer()
	
	// Get database connection
	var err error
	s.db, err = utils.GetTestDB()
	require.NoError(s.T(), err)
	
	// Run migrations
	err = utils.MigrateTestDB(s.db)
	require.NoError(s.T(), err)
	
	// Register services
	di.RegisterServices(s.container, s.db, nil)
	
	// Get user service
	serviceContainer := di.NewServiceContainer(s.container)
	s.userService = serviceContainer.GetUserService()
	require.NotNil(s.T(), s.userService, "User service should not be nil")
}

func (s *UserTestSuite) TearDownSuite() {
	// Cleanup test data
	err := utils.CleanupTestDB(s.db)
	require.NoError(s.T(), err)
	
	// Close database connection
	if s.db != nil {
		s.db.Close()
	}
}

func (s *UserTestSuite) TestGetUser() {
	// Arrange
	testUsers := fixtures.GetTestUsers()
	expectedUser := testUsers[0]
	
	// Create test user first
	user := domain.User{
		Name:     expectedUser.Username,
		Email:    expectedUser.Email,
		Password: "testpassword123",
	}
	createdUser, err := s.userService.CreateUser(user)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), createdUser.ID)

	// Act
	foundUser, err := s.userService.GetUser(createdUser.ID)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), createdUser.ID, foundUser.ID)
	require.Equal(s.T(), expectedUser.Username, foundUser.Name)
	require.Equal(s.T(), expectedUser.Email, foundUser.Email)
}

func (s *UserTestSuite) TestCreateUser() {
	// Arrange
	newUser := domain.User{
		Name:     "newuser",
		Email:    "new@example.com",
		Password: "testpassword123",
	}

	// Act
	createdUser, err := s.userService.CreateUser(newUser)

	// Assert
	require.NoError(s.T(), err)
	require.NotZero(s.T(), createdUser.ID)
	require.Equal(s.T(), newUser.Name, createdUser.Name)
	require.Equal(s.T(), newUser.Email, createdUser.Email)
	require.Empty(s.T(), createdUser.Password) // Password should not be returned
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
} 