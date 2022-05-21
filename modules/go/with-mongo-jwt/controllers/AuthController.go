package controllers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	mongoconnect "github.com/ruthv1k/flock/modules/go/with-mongo-jwt/database"
	"github.com/ruthv1k/flock/modules/go/with-mongo-jwt/models"
	"github.com/ruthv1k/flock/modules/go/with-mongo-jwt/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register allows user to register into the application with credentials and returns a success message if valid data is provided
//
// controllers.Register
func Register(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return err
	}

	if u.Email == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Email and password cannot be empty"})
	}

	usersCollection, ctx, cancel := mongoconnect.GetCollection("users")
	defer cancel()

	// check if the given user is an existing user
	userDetails := usersCollection.FindOne(ctx, bson.M{"email": u.Email})

	var user models.User
	userDetails.Decode(&user)

	if user.Email == u.Email {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "User already exists"})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 13)
	if err != nil {
		return err
	}

	user = models.User{
		UserId:      uuid.NewString(),
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Password:    string(password),
		// assign a default role
		Role: "writer",
	}

	if _, err := usersCollection.InsertOne(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Account created"})
}

// Login allows user to login into the application with credentials and returns a JWT token string if valid credentails are passed.
//
// controllers.Login
func Login(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return err
	}

	if u.Email == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Email and password cannot be empty"})
	}

	usersCollection, ctx, cancel := mongoconnect.GetCollection("users")
	defer cancel()

	// check if the given user is an existing user
	userDetails := usersCollection.FindOne(ctx, bson.M{"email": u.Email})

	if userDetails.Err() != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
	}

	var user models.User
	userDetails.Decode(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid email / password"})
	}

	claims := &models.UserClaims{
		UserId:      user.UserId,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SECRET, _ := utils.GetEnv("SECRET_KEY")

	signedToken, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})

}
