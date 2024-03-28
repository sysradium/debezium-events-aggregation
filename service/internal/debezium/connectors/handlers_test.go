package connectors_test

import (
	"os"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors"
)

func TestAuthHandlerBasicTest(t *testing.T) {
	payload, err := os.ReadFile("../fixtures/sample_snapshot_auth_user_payload.json")
	require.NoError(t, err)

	ch := make(chan *message.Message)
	conn := connectors.NewAuthUserHandler(ch)

	msg := message.NewMessage(
		"e35e3746-43e2-4c38-97eb-90d6dfcab239",
		message.Payload(payload),
	)
	require.NoError(
		t,
		conn.Handle(msg),
	)
}
