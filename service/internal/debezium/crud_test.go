package debezium_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium"
)

func Test(t *testing.T) {
	payload, err := os.ReadFile("./fixtures/sample_snapshot_auth_user_payload.json")
	require.NoError(t, err)

	var e debezium.DebeziumEvent
	require.NoError(
		t,
		json.Unmarshal(payload, &e),
	)

	assert.Equal(t, debezium.OPERATION_SNAPSHOT, e.Payload.Op)
}
