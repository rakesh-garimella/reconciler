package status

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	e "github.com/kyma-incubator/reconciler/pkg/error"
	log "github.com/kyma-incubator/reconciler/pkg/logger"
	"github.com/kyma-incubator/reconciler/pkg/reconciler"
	"github.com/kyma-incubator/reconciler/pkg/test"
	"github.com/stretchr/testify/require"
)

//testCallbackHandler is tracking fired status-updates in an env-var (allows a stateless callback implementation)
//This implementation CAN NOT RUN IN PARALLEL!
type testCallbackHandler struct {
}

func newTestCallbackHandler(t *testing.T) *testCallbackHandler {
	require.NoError(t, os.Unsetenv("_testCallbackHandlerStatuses"))
	return &testCallbackHandler{}
}

func (cb *testCallbackHandler) Callback(status reconciler.Status) error {
	statusList := os.Getenv("_testCallbackHandlerStatuses")
	if statusList == "" {
		statusList = string(status)
	} else {
		statusList = fmt.Sprintf("%s,%s", statusList, status)
	}
	return os.Setenv("_testCallbackHandlerStatuses", statusList)
}

func (cb *testCallbackHandler) Statuses() []reconciler.Status {
	statuses := strings.Split(os.Getenv("_testCallbackHandlerStatuses"), ",")
	var result []reconciler.Status
	for _, status := range statuses {
		result = append(result, reconciler.Status(status))
	}
	return result
}

func (cb *testCallbackHandler) LatestStatus() reconciler.Status {
	statuses := strings.Split(os.Getenv("_testCallbackHandlerStatuses"), ",")
	return reconciler.Status(statuses[len(statuses)-1])
}

func TestStatusUpdater(t *testing.T) { //DO NOT RUN THIS TEST CASES IN PARALLEL!
	if !test.RunExpensiveTests() {
		return
	}

	t.Parallel()

	logger := log.NewOptionalLogger(true)
	t.Run("Test status updater without timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		callbackHdlr := newTestCallbackHandler(t)

		statusUpdater, err := NewStatusUpdater(ctx, callbackHdlr, logger, Config{
			Interval: 1 * time.Second,
			Timeout:  10 * time.Second,
		})
		require.NoError(t, err)
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.NotStarted)

		require.NoError(t, statusUpdater.Running())
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.Running)
		time.Sleep(2 * time.Second)

		require.NoError(t, statusUpdater.Failed())
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.Failed)
		time.Sleep(2 * time.Second)

		require.NoError(t, statusUpdater.Success())
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.Success)
		time.Sleep(2 * time.Second)

		//check fired status updates
		require.GreaterOrEqual(t, len(callbackHdlr.Statuses()), 4) //anything >= 4 is sufficient to ensure the statusUpdaters works
		require.Equal(t, callbackHdlr.LatestStatus(), reconciler.Success)
	})

	t.Run("Test status updater with context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		callbackHdlr := newTestCallbackHandler(t)

		statusUpdater, err := NewStatusUpdater(ctx, callbackHdlr, logger, Config{
			Interval: 1 * time.Second,
			Timeout:  10 * time.Second,
		})
		require.NoError(t, err)
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.NotStarted)

		require.NoError(t, statusUpdater.Running())
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.Running)

		time.Sleep(3 * time.Second) //wait longer than timeout to simulate expired context

		require.True(t, statusUpdater.isContextClosed()) //verify that status-updater received timeout

		//check fired status updates
		require.GreaterOrEqual(t, len(callbackHdlr.Statuses()), 2) //anything > 1 is sufficient to ensure the statusUpdaters worked

		err = statusUpdater.Failed()
		require.Error(t, err)
		require.IsType(t, &e.ContextClosedError{}, err) //status changes have to fail after status-updater was interrupted
	})

	t.Run("Test status updater with status updater timeout", func(t *testing.T) {
		callbackHdlr := newTestCallbackHandler(t)

		statusUpdater, err := NewStatusUpdater(context.Background(), callbackHdlr, logger, Config{
			Interval: 500 * time.Millisecond,
			Timeout:  1 * time.Second,
		})
		require.NoError(t, err)
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.NotStarted)

		require.NoError(t, statusUpdater.Running())
		require.Equal(t, statusUpdater.CurrentStatus(), reconciler.Running)

		time.Sleep(2 * time.Second) //wait longer than status update timeout to timeout

		//check fired status updates
		require.LessOrEqual(t, len(callbackHdlr.Statuses()), 2) //anything <= 2 is sufficient to ensure the statusUpdaters worked

		err = statusUpdater.Failed()
		require.Error(t, err)
		require.IsType(t, &e.ContextClosedError{}, err) //status changes have to fail after status-updater was interrupted
	})

}
