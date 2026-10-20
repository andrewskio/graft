package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// LogConsumer, used by the Compose SDK, processes container log events
type LogConsumer interface {
	Log(containerName, message string)
	Err(containerName, message string)
	Status(containerName, message string)
}

const (
	maxLogEntries = 1000
	// Bound the queue and UI batches so bursts cannot grow memory without limit.
	maxLogBatch      = 256
	logBatchInterval = 50 * time.Millisecond
)

type logEntriesMsg []string

type logStreamFinishedMsg struct {
	err error
}

type logConsumer struct {
	// The stream writes formatted entries here; the pump reads them separately.
	entries chan<- string
}

func (c *logConsumer) Log(containerName, message string) {
	c.entries <- fmt.Sprintf("[INFO] [%s] %s", containerName, strings.TrimSpace(message))
}

func (c *logConsumer) Err(containerName, message string) {
	c.entries <- fmt.Sprintf("[WARN] [%s] %s", containerName, strings.TrimSpace(message))
}

func (c *logConsumer) Status(containerName, message string) {
	c.entries <- fmt.Sprintf("[STATUS] [%s] is now: %s", containerName, message)
}

type logModel struct {
	title          string
	entries        []string
	viewport       viewport.Model
	streamErr      error
	quitRequested  bool
	darkBackground bool
}

func newLogModel(title string) logModel {
	view := viewport.New()
	view.SoftWrap = true
	return logModel{title: title, viewport: view}
}

func (m logModel) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logEntriesMsg:
		// Apply a whole batch at once, retaining only recent history.
		wasAtBottom := m.viewport.AtBottom()
		m.entries = append(m.entries, msg...)
		if len(m.entries) > maxLogEntries {
			m.entries = append([]string(nil), m.entries[len(m.entries)-maxLogEntries:]...)
		}
		m.viewport.SetContent(strings.Join(m.entries, "\n"))
		if wasAtBottom {
			m.viewport.GotoBottom()
		}
	case logStreamFinishedMsg:
		if msg.err != nil && !errors.Is(msg.err, context.Canceled) {
			m.streamErr = msg.err
		}
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(max(1, msg.Height-2))
		m.viewport.GotoBottom()
	case tea.BackgroundColorMsg:
		m.darkBackground = msg.IsDark()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitRequested = true
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m logModel) View() tea.View {
	footer := "q/ctrl+c quit  up/down scroll  pgup/pgdn page"
	if m.streamErr != nil {
		footer = "log stream ended: " + m.streamErr.Error() + "  (q to quit)"
	}
	title := lipgloss.NewStyle().Bold(true).Render(m.title)
	footer = lipgloss.NewStyle().Render(footer)
	view := tea.NewView(fmt.Sprintf("%s\n%s\n%s", title, m.viewport.View(), footer))
	view.AltScreen = true
	return view
}

// RunLogViewer displays a scrollable log view while stream sends entries to consumer.
func RunLogViewer(ctx context.Context, title string, stream func(context.Context, LogConsumer) error) error {
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	program := tea.NewProgram(newLogModel(title), tea.WithContext(streamCtx))
	streamDone := make(chan error, 1)
	// The buffered channel lets the stream briefly get ahead of the UI pump.
	entries := make(chan string, maxLogBatch)
	pumpDone := make(chan struct{})
	go func() {
		defer close(pumpDone)
		ticker := time.NewTicker(logBatchInterval)
		defer ticker.Stop()

		batch := make([]string, 0, maxLogBatch)
		flush := func() {
			if len(batch) == 0 {
				return
			}
			// Send a copy because this goroutine reuses batch for the next group.
			program.Send(logEntriesMsg(append([]string(nil), batch...)))
			batch = batch[:0]
		}
		for {
			select {
			case entry, ok := <-entries:
				if !ok {
					// Drain the final partial batch before telling the UI the stream ended.
					flush()
					err := <-streamDone
					program.Send(logStreamFinishedMsg{err: err})
					return
				}
				batch = append(batch, entry)
				if len(batch) == maxLogBatch {
					flush()
				}
			case <-ticker.C:
				// Also flush on a timer so low-volume logs appear promptly.
				flush()
			}
		}
	}()
	go func() {
		// The stream runs independently; closing entries signals the pump to drain.
		err := stream(streamCtx, &logConsumer{entries: entries})
		close(entries)
		streamDone <- err
	}()

	finalModel, runErr := program.Run()
	cancel()
	// Wait until the pump has finished sending its final UI messages.
	<-pumpDone
	if runErr != nil {
		return runErr
	}
	if model, ok := finalModel.(logModel); ok {
		if model.quitRequested {
			return nil
		}
		return model.streamErr
	}
	return nil
}
