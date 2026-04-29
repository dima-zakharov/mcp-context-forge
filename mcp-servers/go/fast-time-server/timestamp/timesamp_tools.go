package timestamp

import (
	"context"

	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/maypok86/otter"
)

type TimestampManager struct {
	cache otter.Cache[string, string]
}

func NewTimestampManager(s *server.MCPServer) *TimestampManager {
	cache, _ := otter.MustBuilder[string, string](1_000_000).
		WithTTL(time.Duration(42 * time.Second)).
		Build()
	manager := &TimestampManager{cache: cache}
	manager.AddTimestampTools(s)
	return manager
}

func (m *TimestampManager) setTimestamp( //
	ctx context.Context, request mcp.CallToolRequest, //
) (*mcp.CallToolResult, error) {
	note := request.GetString("note", "")
	now := time.Now().Format(time.RFC3339)
	session := server.ClientSessionFromContext(ctx)
	txt := now + " - " + note
	m.cache.Set(session.SessionID(), txt)

	return mcp.NewToolResultText(txt), nil
}

func (m *TimestampManager) getTimestamp( //
	ctx context.Context, request mcp.CallToolRequest, //
) (*mcp.CallToolResult, error) {

	session := server.ClientSessionFromContext(ctx)
	if session == nil {
		return mcp.NewToolResultError("session not found"), nil
	}
	now, ok := m.cache.Get(session.SessionID())
	if !ok {
		now = "No timestamp set"
	}

	return mcp.NewToolResultText(now), nil
}

func (m *TimestampManager) clearTimestamp( //
	ctx context.Context, request mcp.CallToolRequest, //
) (*mcp.CallToolResult, error) {

	session := server.ClientSessionFromContext(ctx).SessionID()
	m.cache.Delete(session)

	return mcp.NewToolResultText("Timestamp cleared"), nil
}

func (m *TimestampManager) AddTimestampTools(s *server.MCPServer) {
	set_timestamp := mcp.NewTool( //
		"set_timestamp", //
		mcp.WithDescription("Sets session timestamp with a note"),
		mcp.WithString("note",
			mcp.Required(),
			mcp.Description("Note to attach to timestamp"),
		),
	)
	get_timestamp := mcp.NewTool( //
		"get_timestamp", mcp.WithDescription("Gets session timestamp"),
	)
	clear_timestamp := mcp.NewTool( //
		"clear_timestamp", mcp.WithDescription("Clears session timestamp"),
	)
	s.AddTool(set_timestamp, m.setTimestamp)
	s.AddTool(get_timestamp, m.getTimestamp)
	s.AddTool(clear_timestamp, m.clearTimestamp)
}
