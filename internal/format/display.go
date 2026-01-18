package format

import (
	"strings"

	"github.com/visionik/libfizz-go/fizzy"
)

// CardDisplay represents a card with fields suitable for table display
// Designed to fit within 120 character terminal width
type CardDisplay struct {
	Num    int    `json:"num"`
	Title  string `json:"title"`
	Desc   string `json:"desc,omitempty"`
	Status string `json:"status"`
	Board  string `json:"board,omitempty"`
	Tags   string `json:"tags,omitempty"`
}

// CardDetailDisplay represents a single card with more detail
type CardDetailDisplay struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Board       string `json:"board,omitempty"`
	Golden      bool   `json:"golden"`
	Closed      bool   `json:"closed"`
	Assignees   string `json:"assignees,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Created     string `json:"created"`
	URL         string `json:"url"`
}

// BoardDisplay represents a board with fields suitable for table display
type BoardDisplay struct {
	Name      string `json:"name"`
	Desc      string `json:"desc,omitempty"`
	Access    string `json:"access"`
	Creator   string `json:"creator,omitempty"`
	Created   string `json:"created"`
}

// BoardDetailDisplay represents a single board with more detail
type BoardDetailDisplay struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	AllAccess   bool   `json:"all_access"`
	Creator     string `json:"creator,omitempty"`
	Created     string `json:"created"`
	URL         string `json:"url,omitempty"`
}

// truncate truncates a string to maxLen characters, adding "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// formatNames joins names with commas and truncates
func formatNames(names []string, maxLen int) string {
	if len(names) == 0 {
		return ""
	}
	joined := strings.Join(names, ", ")
	return truncate(joined, maxLen)
}

// ToCardDisplay converts a Card to CardDisplay for compact table output
func ToCardDisplay(card fizzy.Card) CardDisplay {
	display := CardDisplay{
		Num:    card.Number,
		Title:  truncate(card.Title, 40),
		Status: card.Status,
	}

	// Add description if present (truncated)
	if card.Description != nil && *card.Description != "" {
		display.Desc = truncate(*card.Description, 30)
	}

	// Add board name if available
	if card.Board != nil {
		display.Board = truncate(card.Board.Name, 20)
	}

	// Format tags (compact)
	if len(card.Tags) > 0 {
		names := make([]string, 0, len(card.Tags))
		for _, t := range card.Tags {
			names = append(names, t.Name)
		}
		display.Tags = formatNames(names, 15)
	}

	return display
}

// ToCardDetailDisplay converts a Card to CardDetailDisplay for single card view
func ToCardDetailDisplay(card fizzy.Card) CardDetailDisplay {
	display := CardDetailDisplay{
		Number:  card.Number,
		Title:   card.Title,
		Status:  card.Status,
		Golden:  card.Golden,
		Closed:  card.Closed,
		Created: card.CreatedAt.Format("2006-01-02 15:04"),
		URL:     card.URL,
	}

	// Add description if present
	if card.Description != nil && *card.Description != "" {
		display.Description = *card.Description
	}

	// Add board name if available
	if card.Board != nil {
		display.Board = card.Board.Name
	}

	// Format assignees
	if len(card.Assignees) > 0 {
		names := make([]string, 0, len(card.Assignees))
		for _, a := range card.Assignees {
			names = append(names, a.Name)
		}
		display.Assignees = strings.Join(names, ", ")
	}

	// Format tags
	if len(card.Tags) > 0 {
		names := make([]string, 0, len(card.Tags))
		for _, t := range card.Tags {
			names = append(names, t.Name)
		}
		display.Tags = strings.Join(names, ", ")
	}

	return display
}

// ToCardDisplaySlice converts a slice of Cards to CardDisplay
func ToCardDisplaySlice(cards []fizzy.Card) []CardDisplay {
	displays := make([]CardDisplay, len(cards))
	for i, card := range cards {
		displays[i] = ToCardDisplay(card)
	}
	return displays
}

// ToBoardDisplay converts a Board to BoardDisplay for compact table output
func ToBoardDisplay(board fizzy.Board) BoardDisplay {
	display := BoardDisplay{
		Name:    truncate(board.Name, 30),
		Created: board.CreatedAt.Format("2006-01-02"),
	}

	if board.Description != nil && *board.Description != "" {
		display.Desc = truncate(*board.Description, 40)
	}

	if board.AllAccess {
		display.Access = "all"
	} else {
		display.Access = "restricted"
	}

	if board.Creator != nil {
		display.Creator = truncate(board.Creator.Name, 20)
	}

	return display
}

// ToBoardDetailDisplay converts a Board to BoardDetailDisplay for single board view
func ToBoardDetailDisplay(board fizzy.Board) BoardDetailDisplay {
	display := BoardDetailDisplay{
		ID:        board.ID,
		Name:      board.Name,
		AllAccess: board.AllAccess,
		Created:   board.CreatedAt.Format("2006-01-02 15:04"),
	}

	if board.Description != nil && *board.Description != "" {
		display.Description = *board.Description
	}

	if board.URL != "" {
		display.URL = board.URL
	}

	if board.Creator != nil {
		display.Creator = board.Creator.Name
	}

	return display
}

// ToBoardDisplaySlice converts a slice of Boards to BoardDisplay
func ToBoardDisplaySlice(boards []fizzy.Board) []BoardDisplay {
	displays := make([]BoardDisplay, len(boards))
	for i, board := range boards {
		displays[i] = ToBoardDisplay(board)
	}
	return displays
}
