package timestamp

import (
	"context"

	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/maypok86/otter"
)

type TimestampManager struct {
	cache otter.Cache[string, time.Time]
}

func NewTimestampManager() *TimestampManager {
	cache, _ := otter.MustBuilder[string, time.Time](1_000_000).
		WithTTL(time.Duration(42 * time.Second)).
		Build()
	return &TimestampManager{cache: cache}
}

func (m *TimestampManager) handleRecord(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("12:00"), nil
}

func (m *TimestampManager) AddTimestampTools(s *server.MCPServer) {
	set_timestamp := mcp.NewTool("set_timestamp")
	s.AddTool(set_timestamp, m.handleRecord)

}
