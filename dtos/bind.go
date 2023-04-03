package dtos

import "github.com/labstack/echo/v4"

func BindAndValidateDTO(c echo.Context, dto interface{}) error {
	if err := c.Bind(dto); err != nil {
		c.Logger().Error(err)
		return err
	}
	if err := c.Validate(dto); err != nil {
		c.Logger().Error(err)
		return err
	}
	return nil
}
