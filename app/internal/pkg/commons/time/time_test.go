package time_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	
	serviceTime "userservice/internal/pkg/commons/time"
)

func TestUserServiceTime_String(t *testing.T) {
	currentTime := serviceTime.Now()
	expectedTime := time.Now().UTC().Format(serviceTime.CustomTimeLayout)

	require.Equal(t, expectedTime, currentTime.String())
}
