package gin

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func queryBool(key string, c *gin.Context) (bool, bool, error) {
	v := c.Query(key)
	if v == "" {
		return false, false, nil
	}

	b, err := strconv.ParseBool(v)
	return b, true, err
}

func queryInt(key string, c *gin.Context) (int, bool, error) {
	v := c.Query(key)
	if v == "" {
		return 0, false, nil
	}

	i, err := strconv.Atoi(v)
	return i, true, err
}

func queryUInt64(key string, c *gin.Context) (uint64, bool, error) {
	v := c.Query(key)
	if v == "" {
		return 0, false, nil
	}

	i, err := strconv.ParseUint(v, 10, 64)
	return i, true, err
}
