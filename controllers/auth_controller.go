package controllers

import (
	"fiber-api/config"
	"fiber-api/models"
	"fiber-api/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Register(c fiber.Ctx) error {
	var user models.User
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var dbUser models.User
	config.DB.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID != 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	if err := utils.Validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters, contain at least one uppercase letter, one lowercase letter, and one number",
		})
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	config.DB.Create(&user)
	return c.JSON(user)
}

func Login(c fiber.Ctx) error {
	var user models.LoginRequest
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var dbUser models.User
	config.DB.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = dbUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": signedToken,
	})

}
