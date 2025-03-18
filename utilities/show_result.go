package utilities

import "github.com/gofiber/fiber/v2"

// ShowSuccess responds with a success message and optional data
func ShowSuccess(c *fiber.Ctx, msg interface{}, code int, data interface{}) error {
    response := fiber.Map{
        "success":      true,
        "response_code": code,
        "message":      msg,
    }
    if data != nil {
        response["data"] = data
    }
    return c.JSON(response)
}

// ShowMessage responds with a simple success message
func ShowMessage(c *fiber.Ctx, msg interface{}, code int) error {
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success":      true,
        "response_code": code,
        "message":      msg,
    })
}

// ShowError responds with an error message and nested field-specific errors
func ShowError(c *fiber.Ctx, message string, code int, errors map[string][]string) error {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "success":      false,
        "response_code": code,
        "message":      message,
        "errors":       errors,
    })
}
