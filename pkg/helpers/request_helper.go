package helpers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUintParam(c *gin.Context, paramName string) (uint, error) {
	param := c.Param(paramName)
	if param == "" {
		return 0, fmt.Errorf("missing %s", paramName)
	}

	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", paramName)
	}

	return uint(id), nil
}
