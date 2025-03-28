package time_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	serviceTime "userservice/internal/pkg/commons/time"
)

func TestUserServiceTime_String(t *testing.T) {
	currentTime := serviceTime.Now()
	expectedTime := time.Now().UTC().Format(serviceTime.CustomTimeLayout)

	require.Equal(t, expectedTime, currentTime.String())
}

func TestUserServiceTime_FromString(t *testing.T) {
	validTimeStr := "2024-07-15T07:25:55.32Z"
	invalidTimeStr := "2024-07-15 07:25:55"

	// Test valid case
	st, err := serviceTime.FromString(validTimeStr)
	require.NoError(t, err)
	require.Equal(t, st.String(), validTimeStr)

	_, err = serviceTime.FromString(invalidTimeStr)
	require.Error(t, err)
}
