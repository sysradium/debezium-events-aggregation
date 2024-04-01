package connectors_test

import (
	"os"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"
	"github.com/sysradium/debezium-events-aggregation/service/internal/adapters/ephemeral"
	"github.com/sysradium/debezium-events-aggregation/service/internal/app"
	"github.com/sysradium/debezium-events-aggregation/service/internal/debezium/connectors"
)

func TestAuthHandlerBasicTest(t *testing.T) {
	payload, err := os.ReadFile("../fixtures/sample_snapshot_auth_user_payload.json")
	require.NoError(t, err)

	db, err := ephemeral.New()
	require.NoError(t, err)

	conn := connectors.NewAuthUserHandler(
		make(chan *message.Message),
		app.NewApplication(db),
	)

	msg := message.NewMessage(
		"e35e3746-43e2-4c38-97eb-90d6dfcab239",
		message.Payload(payload),
	)
	require.NoError(t, conn.Handle(msg))
}
