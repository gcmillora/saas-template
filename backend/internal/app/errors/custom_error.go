package errors

import "encoding/json"

type CustomErr struct {
	Message string
	Data    interface{}
}

func (c *CustomErr) Error() string {
	return c.Message
}

func (c *CustomErr) Render(expected map[string]interface{}) error {
	errJsonStr, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(errJsonStr, &expected)

	if err != nil {
		return err
	}

	return nil
}

type BadRequestError struct {
	CustomErr
}

type UnauthorizedError struct {
	CustomErr
}

type ForbiddenError struct {
	CustomErr
}

type ValidationError struct {
	CustomErr
}

type NotFoundError struct {
	CustomErr
}
