package utils_test

import (
	"backend-assignment-avito/internal/utils/jwt_utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyJWTtoken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MTc2NTE2OTksImlhdCI6MTc3NDQ1MTY5OSwiVXNlcklkIjoiMjQ5YmU3Y2YtZDQxOS00YzU0LTk3ZjItZDA0MTA3ODA2ZTM2IiwiUm9sZSI6ImFkbWluIn0.XZWnw20_E6hyxpCwgB6bdMkcu-mJmDumiOwMukwyyj4"

	claims, err := jwt_utils.VerifyJWTtoken(token)

	require.NoError(t, err)
	assert.Equal(t, "admin", claims.Role)
}

func TestVerifyJWTtokenFail(t *testing.T) {
	token := "1yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQ0MDI2NzAsImlhdCI6MTc3NDM1OTQ3MCwiVXNlcklkIjoiMjQ5YmU3Y2YtZDQxOS00YzU0LTk3ZjItZDA0MTA3ODA2ZTM2IiwiUm9sZSI6ImFkbWluIn0.jcPenIgCiji4xyP8t5zm-Tyc3Q8BMFPmo_GlTQZ6vBE"

	_, err := jwt_utils.VerifyJWTtoken(token)

	require.Error(t, err)
}

func TestVerifyJWTtokenInvalid(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQzNTk3MzYsImlhdCI6MTc3NDM1OTczNiwiVXNlcklkIjoiMjQ5YmU3Y2YtZDQxOS00YzU0LTk3ZjItZDA0MTA3ODA2ZTM2IiwiUm9sZSI6ImFkbWluIn0.NOeOjklw4TrVitMCmp9rkd9L-EwRij8P__jxvi6kRX4"

	_, err := jwt_utils.VerifyJWTtoken(token)

	require.Error(t, err)
}
