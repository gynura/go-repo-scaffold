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

func TestUserServiceTime_FromTime(t *testing.T) {
	inputTime := time.Date(2024, 3, 29, 12, 30, 45, 123456789, time.FixedZone("Custom", -18000))
	expectedTime := inputTime.UTC()

	libTime, err := serviceTime.FromTime(inputTime)
	require.NoError(t, err)
	require.Equal(t, expectedTime, libTime.Time)
}

func TestUserServiceTime_Now(t *testing.T) {
	// Call the Now function to get the current ServiceTime
	currentTime := serviceTime.Now()

	// Check if the current time is close to the expected time (within 1 second tolerance)
	expectedTime := time.Now().UTC().Format(serviceTime.CustomTimeLayout)
	actualTime := currentTime.String()

	// Allow a small tolerance for differences in time format and system clock
	require.True(t, actualTime == expectedTime || actualTime == time.Now().UTC().Format(serviceTime.CustomTimeLayout))
}
