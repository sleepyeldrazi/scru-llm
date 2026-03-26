package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sleepyeldrazi/scru-llm/internal/types"
)

// EventLog defines the interface for event logging
type EventLog interface {
	// Append writes an event to the log
	Append(event *types.Event) error

	// Read returns events, optionally filtered
	Read(opts ReadOptions) ([]*types.Event, error)

	// Close closes the event log
	Close() error
}

// ReadOptions defines options for reading events
type ReadOptions struct {
	TaskID    string
	SprintID  string
	EventType string
	Since     time.Time
	Limit     int
}

// JSONLFileEventLog implements EventLog using a JSON Lines file
type JSONLFileEventLog struct {
	mu     sync.Mutex
	file   *os.File
	writer *bufio.Writer
	path   string
}

// NewJSONLFileEventLog creates a new JSONL file-based event log
func NewJSONLFileEventLog(path string) (*JSONLFileEventLog, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create event log directory: %w", err)
	}

	// Open file for append (create if doesn't exist)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open event log file: %w", err)
	}

	return &JSONLFileEventLog{
		file:   file,
		writer: bufio.NewWriter(file),
		path:   path,
	}, nil
}

// Append writes an event to the log
func (l *JSONLFileEventLog) Append(event *types.Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Ensure timestamp is set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Marshal to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Write line
	if _, err := l.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	if err := l.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// Flush to ensure durability
	if err := l.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush event: %w", err)
	}

	return nil
}

// Read returns events matching the options
func (l *JSONLFileEventLog) Read(opts ReadOptions) ([]*types.Event, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Flush any pending writes
	if err := l.writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush before read: %w", err)
	}

	// Open file for reading
	file, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.Event{}, nil
		}
		return nil, fmt.Errorf("failed to open event log for reading: %w", err)
	}
	defer file.Close()

	var events []*types.Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var event types.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue // Skip malformed lines
		}

		// Apply filters
		if opts.TaskID != "" && event.TaskID != opts.TaskID {
			continue
		}
		if opts.SprintID != "" && event.SprintID != opts.SprintID {
			continue
		}
		if opts.EventType != "" && event.Type != opts.EventType {
			continue
		}
		if !opts.Since.IsZero() && event.Timestamp.Before(opts.Since) {
			continue
		}

		events = append(events, &event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading event log: %w", err)
	}

	// Apply limit (get most recent)
	if opts.Limit > 0 && len(events) > opts.Limit {
		start := len(events) - opts.Limit
		events = events[start:]
	}

	return events, nil
}

// Close closes the event log
func (l *JSONLFileEventLog) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush on close: %w", err)
	}

	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close event log: %w", err)
	}

	return nil
}
