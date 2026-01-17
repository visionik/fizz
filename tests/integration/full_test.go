//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visionik/fizz/internal/client"
	"github.com/visionik/fizz/internal/config"
	"github.com/visionik/libfizz-go/fizzy"
)

// TestFullIntegration runs a complete integration test of all fizz features
// It creates a dedicated test board and cleans it up afterward
func TestFullIntegration(t *testing.T) {
	// Setup: Create client
	cfg, err := config.LoadFromEnv()
	require.NoError(t, err, "FIZZY_TOKEN and FIZZY_ACCOUNT must be set")

	c, err := client.New(cfg, true) // debug mode
	require.NoError(t, err, "Failed to create client")

	ctx := context.Background()

	// Create a unique test board name with timestamp
	testBoardName := fmt.Sprintf("FIZZ_INTEGRATION_TEST_%d", time.Now().Unix())
	
	t.Run("Identity", func(t *testing.T) {
		testIdentity(t, ctx, c)
	})

	var testBoardID string
	t.Run("Boards", func(t *testing.T) {
		testBoardID = testBoards(t, ctx, c, testBoardName)
	})

	// Skip subsequent tests if board creation failed
	if testBoardID == "" {
		t.Fatal("Test board creation failed, cannot continue")
	}

	// Ensure cleanup happens even if tests fail
	defer func() {
		t.Logf("Cleaning up test board: %s", testBoardID)
		err := c.Boards.Delete(ctx, testBoardID)
		if err != nil {
			t.Logf("Warning: Failed to delete test board %s: %v", testBoardID, err)
		}
	}()

	var testCardID string
	t.Run("Cards", func(t *testing.T) {
		testCardID = testCards(t, ctx, c, testBoardID)
	})

	if testCardID != "" {
		t.Run("Comments", func(t *testing.T) {
			testComments(t, ctx, c, testCardID)
		})

		t.Run("Steps", func(t *testing.T) {
			testSteps(t, ctx, c, testCardID)
		})
	}

	t.Run("Columns", func(t *testing.T) {
		testColumns(t, ctx, c, testBoardID)
	})

	t.Run("Tags", func(t *testing.T) {
		testTags(t, ctx, c)
	})

	t.Run("Users", func(t *testing.T) {
		testUsers(t, ctx, c)
	})

	t.Run("Notifications", func(t *testing.T) {
		testNotifications(t, ctx, c)
	})
}

func testIdentity(t *testing.T, ctx context.Context, c *client.Client) {
	t.Log("Testing identity get...")
	
	identity, err := c.Identity.Get(ctx)
	require.NoError(t, err, "Failed to get identity")
	assert.NotNil(t, identity, "Identity should not be nil")
	assert.NotEmpty(t, identity.Accounts, "Should have at least one account")
	
	t.Logf("✓ Identity verified: %d account(s)", len(identity.Accounts))
}

func testBoards(t *testing.T, ctx context.Context, c *client.Client, boardName string) string {
	t.Log("Testing boards...")
	
	// List boards
	t.Log("  - Listing boards")
	boards, err := c.Boards.List(ctx)
	require.NoError(t, err, "Failed to list boards")
	initialCount := len(boards)
	t.Logf("  ✓ Found %d existing boards", initialCount)

	// Create test board
	t.Logf("  - Creating test board: %s", boardName)
	description := "Automated integration test board - safe to delete"
	board, err := c.Boards.Create(ctx, &fizzy.BoardCreateOptions{
		Name:        boardName,
		Description: &description,
	})
	require.NoError(t, err, "Failed to create board")
	require.NotEmpty(t, board.ID, "Board ID should not be empty")
	assert.Equal(t, boardName, board.Name, "Board name mismatch")
	testBoardID := board.ID
	t.Logf("  ✓ Created board: %s", testBoardID)

	// Get board
	t.Log("  - Getting board by ID")
	fetchedBoard, err := c.Boards.Get(ctx, testBoardID)
	require.NoError(t, err, "Failed to get board")
	assert.Equal(t, testBoardID, fetchedBoard.ID, "Board ID mismatch")
	assert.Equal(t, boardName, fetchedBoard.Name, "Board name mismatch")
	t.Log("  ✓ Board retrieved successfully")

	// Update board
	t.Log("  - Updating board")
	newName := boardName + "_UPDATED"
	newDesc := "Updated description"
	updatedBoard, err := c.Boards.Update(ctx, testBoardID, &fizzy.BoardUpdateOptions{
		Name:        &newName,
		Description: &newDesc,
	})
	require.NoError(t, err, "Failed to update board")
	assert.Equal(t, newName, updatedBoard.Name, "Board name not updated")
	if updatedBoard.Description != nil {
		assert.Equal(t, newDesc, *updatedBoard.Description, "Board description not updated")
	} else {
		t.Log("  ⚠ Warning: Description not returned by API")
	}
	t.Log("  ✓ Board updated successfully")

	// Verify board in list
	t.Log("  - Verifying board appears in list")
	boards, err = c.Boards.List(ctx)
	require.NoError(t, err, "Failed to list boards")
	assert.Equal(t, initialCount+1, len(boards), "Board count should increase by 1")
	found := false
	for _, b := range boards {
		if b.ID == testBoardID {
			found = true
			break
		}
	}
	assert.True(t, found, "Test board not found in list")
	t.Log("  ✓ Board appears in list")

	t.Log("✓ Boards tests passed")
	return testBoardID
}

func testCards(t *testing.T, ctx context.Context, c *client.Client, boardID string) string {
	t.Log("Testing cards...")

	// Create card
	t.Log("  - Creating card")
	body := "This is a test card created by integration tests"
	card, err := c.Cards.Create(ctx, &fizzy.CardCreateOptions{
		BoardID: boardID,
		Title:   "Integration Test Card",
		Body:    &body,
	})
	require.NoError(t, err, "Failed to create card")
	require.NotEmpty(t, card.ID, "Card ID should not be empty")
	assert.Equal(t, "Integration Test Card", card.Title, "Card title mismatch")
	testCardID := card.ID
	t.Logf("  ✓ Created card: %s (number: %d)", testCardID, card.Number)

	// List cards (all)
	t.Log("  - Listing all cards")
	allCards, err := c.Cards.ListAll(ctx, nil)
	require.NoError(t, err, "Failed to list all cards")
	t.Logf("  ✓ Listed %d total cards", len(allCards))

	// List cards by board
	t.Log("  - Listing cards by board")
	boardCards, err := c.Cards.ListAll(ctx, &fizzy.CardListOptions{
		BoardID: boardID,
	})
	require.NoError(t, err, "Failed to list board cards")
	assert.GreaterOrEqual(t, len(boardCards), 1, "Should have at least our test card")
	t.Logf("  ✓ Found %d cards on test board", len(boardCards))

	// Get card
	t.Log("  - Getting card by ID")
	fetchedCard, err := c.Cards.Get(ctx, testCardID)
	require.NoError(t, err, "Failed to get card")
	assert.Equal(t, testCardID, fetchedCard.ID, "Card ID mismatch")
	t.Log("  ✓ Card retrieved successfully")

	// Update card
	t.Log("  - Updating card")
	newTitle := "Updated Test Card"
	newBody := "Updated body content"
	updatedCard, err := c.Cards.Update(ctx, testCardID, &fizzy.CardUpdateOptions{
		Title: &newTitle,
		Body:  &newBody,
	})
	require.NoError(t, err, "Failed to update card")
	assert.Equal(t, newTitle, updatedCard.Title, "Card title not updated")
	assert.Equal(t, newBody, *updatedCard.Body, "Card body not updated")
	t.Log("  ✓ Card updated successfully")

	// Test card actions
	t.Log("  - Testing card actions")

	// Close card
	t.Log("    - Closing card")
	err = c.Cards.Close(ctx, testCardID)
	require.NoError(t, err, "Failed to close card")
	closedCard, _ := c.Cards.Get(ctx, testCardID)
	assert.NotNil(t, closedCard.ClosedAt, "Card should have closed_at timestamp")
	t.Log("    ✓ Card closed")

	// Reopen card
	t.Log("    - Reopening card")
	err = c.Cards.Reopen(ctx, testCardID)
	require.NoError(t, err, "Failed to reopen card")
	reopenedCard, _ := c.Cards.Get(ctx, testCardID)
	assert.Nil(t, reopenedCard.ClosedAt, "Card should not have closed_at after reopen")
	t.Log("    ✓ Card reopened")

	// Watch card
	t.Log("    - Watching card")
	err = c.Cards.Watch(ctx, testCardID)
	require.NoError(t, err, "Failed to watch card")
	t.Log("    ✓ Card watched")

	// Unwatch card
	t.Log("    - Unwatching card")
	err = c.Cards.Unwatch(ctx, testCardID)
	require.NoError(t, err, "Failed to unwatch card")
	t.Log("    ✓ Card unwatched")

	// Mark golden
	t.Log("    - Marking card as golden")
	err = c.Cards.MarkGolden(ctx, testCardID)
	require.NoError(t, err, "Failed to mark card as golden")
	goldenCard, _ := c.Cards.Get(ctx, testCardID)
	assert.True(t, goldenCard.Golden, "Card should be golden")
	t.Log("    ✓ Card marked golden")

	// Unmark golden
	t.Log("    - Unmarking golden")
	err = c.Cards.UnmarkGolden(ctx, testCardID)
	require.NoError(t, err, "Failed to unmark golden")
	normalCard, _ := c.Cards.Get(ctx, testCardID)
	assert.False(t, normalCard.Golden, "Card should not be golden")
	t.Log("    ✓ Card unmarked golden")

	t.Log("✓ Cards tests passed")
	return testCardID
}

func testComments(t *testing.T, ctx context.Context, c *client.Client, cardID string) {
	t.Log("Testing comments...")

	// Create comment
	t.Log("  - Creating comment")
	comment, err := c.Comments.Create(ctx, cardID, &fizzy.CommentCreateOptions{
		Body: "This is a test comment",
	})
	require.NoError(t, err, "Failed to create comment")
	require.NotEmpty(t, comment.ID, "Comment ID should not be empty")
	commentID := comment.ID
	t.Logf("  ✓ Created comment: %s", commentID)

	// List comments
	t.Log("  - Listing comments")
	comments, err := c.Comments.List(ctx, cardID)
	require.NoError(t, err, "Failed to list comments")
	assert.GreaterOrEqual(t, len(comments), 1, "Should have at least one comment")
	t.Logf("  ✓ Listed %d comment(s)", len(comments))

	// Verify our comment is in the list
	found := false
	for _, cmt := range comments {
		if cmt.ID == commentID {
			found = true
			break
		}
	}
	assert.True(t, found, "Created comment should be in list")
	t.Log("  ✓ Comment found in list")

	// Update comment
	t.Log("  - Updating comment")
	updatedComment, err := c.Comments.Update(ctx, cardID, commentID, &fizzy.CommentUpdateOptions{
		Body: "Updated comment text",
	})
	require.NoError(t, err, "Failed to update comment")
	assert.Contains(t, updatedComment.Body, "Updated", "Comment not updated")
	t.Log("  ✓ Comment updated")

	// Test reactions on comment
	t.Log("  - Testing reactions")
	t.Log("    - Creating reaction")
	reaction, err := c.Reactions.Create(ctx, cardID, commentID, &fizzy.ReactionCreateOptions{
		Content: "👍",
	})
	require.NoError(t, err, "Failed to create reaction")
	reactionID := reaction.ID
	t.Logf("    ✓ Created reaction: %s", reactionID)

	t.Log("    - Listing reactions")
	reactions, err := c.Reactions.List(ctx, cardID, commentID)
	require.NoError(t, err, "Failed to list reactions")
	assert.GreaterOrEqual(t, len(reactions), 1, "Should have at least one reaction")
	t.Log("    ✓ Reactions listed")

	t.Log("    - Deleting reaction")
	err = c.Reactions.Delete(ctx, cardID, commentID, reactionID)
	require.NoError(t, err, "Failed to delete reaction")
	t.Log("    ✓ Reaction deleted")

	// Delete comment
	t.Log("  - Deleting comment")
	err = c.Comments.Delete(ctx, cardID, commentID)
	require.NoError(t, err, "Failed to delete comment")
	t.Log("  ✓ Comment deleted")

	t.Log("✓ Comments tests passed")
}

func testSteps(t *testing.T, ctx context.Context, c *client.Client, cardID string) {
	t.Log("Testing steps...")

	// Create step
	t.Log("  - Creating step")
	completed := false
	step, err := c.Steps.Create(ctx, cardID, &fizzy.StepCreateOptions{
		Content:   "Test checklist item",
		Completed: &completed,
	})
	require.NoError(t, err, "Failed to create step")
	require.NotEmpty(t, step.ID, "Step ID should not be empty")
	stepID := step.ID
	t.Logf("  ✓ Created step: %s", stepID)

	// List steps
	t.Log("  - Listing steps")
	steps, err := c.Steps.List(ctx, cardID)
	require.NoError(t, err, "Failed to list steps")
	assert.GreaterOrEqual(t, len(steps), 1, "Should have at least one step")
	t.Logf("  ✓ Listed %d step(s)", len(steps))

	// Get step
	t.Log("  - Getting step")
	fetchedStep, err := c.Steps.Get(ctx, cardID, stepID)
	require.NoError(t, err, "Failed to get step")
	assert.Equal(t, stepID, fetchedStep.ID, "Step ID mismatch")
	t.Log("  ✓ Step retrieved")

	// Update step - mark as completed
	t.Log("  - Updating step (marking complete)")
	completedTrue := true
	updatedStep, err := c.Steps.Update(ctx, cardID, stepID, &fizzy.StepUpdateOptions{
		Completed: &completedTrue,
	})
	require.NoError(t, err, "Failed to update step")
	assert.True(t, updatedStep.Completed, "Step should be completed")
	t.Log("  ✓ Step marked complete")

	// Update step content
	t.Log("  - Updating step content")
	newContent := "Updated checklist item"
	updatedStep, err = c.Steps.Update(ctx, cardID, stepID, &fizzy.StepUpdateOptions{
		Content: &newContent,
	})
	require.NoError(t, err, "Failed to update step content")
	assert.Equal(t, newContent, updatedStep.Content, "Step content not updated")
	t.Log("  ✓ Step content updated")

	// Delete step
	t.Log("  - Deleting step")
	err = c.Steps.Delete(ctx, cardID, stepID)
	require.NoError(t, err, "Failed to delete step")
	t.Log("  ✓ Step deleted")

	t.Log("✓ Steps tests passed")
}

func testColumns(t *testing.T, ctx context.Context, c *client.Client, boardID string) {
	t.Log("Testing columns...")

	// List columns
	t.Log("  - Listing columns")
	initialColumns, err := c.Columns.List(ctx, boardID)
	require.NoError(t, err, "Failed to list columns")
	initialCount := len(initialColumns)
	t.Logf("  ✓ Found %d existing column(s)", initialCount)

	// Create column
	t.Log("  - Creating column")
	column, err := c.Columns.Create(ctx, boardID, &fizzy.ColumnCreateOptions{
		Name: "Test Column",
	})
	require.NoError(t, err, "Failed to create column")
	require.NotEmpty(t, column.ID, "Column ID should not be empty")
	columnID := column.ID
	t.Logf("  ✓ Created column: %s", columnID)

	// Get column
	t.Log("  - Getting column")
	fetchedColumn, err := c.Columns.Get(ctx, boardID, columnID)
	require.NoError(t, err, "Failed to get column")
	assert.Equal(t, columnID, fetchedColumn.ID, "Column ID mismatch")
	t.Log("  ✓ Column retrieved")

	// Update column
	t.Log("  - Updating column")
	newName := "Updated Column Name"
	updatedColumn, err := c.Columns.Update(ctx, boardID, columnID, &fizzy.ColumnUpdateOptions{
		Name: &newName,
	})
	require.NoError(t, err, "Failed to update column")
	assert.Equal(t, newName, updatedColumn.Name, "Column name not updated")
	t.Log("  ✓ Column updated")

	// Delete column
	t.Log("  - Deleting column")
	err = c.Columns.Delete(ctx, boardID, columnID)
	require.NoError(t, err, "Failed to delete column")
	t.Log("  ✓ Column deleted")

	// Verify column deleted
	finalColumns, err := c.Columns.List(ctx, boardID)
	require.NoError(t, err, "Failed to list columns after delete")
	assert.Equal(t, initialCount, len(finalColumns), "Column count should be back to initial")
	t.Log("  ✓ Column deletion verified")

	t.Log("✓ Columns tests passed")
}

func testTags(t *testing.T, ctx context.Context, c *client.Client) {
	t.Log("Testing tags...")

	// List tags
	t.Log("  - Listing tags")
	initialTags, err := c.Tags.List(ctx)
	require.NoError(t, err, "Failed to list tags")
	initialCount := len(initialTags)
	t.Logf("  ✓ Found %d existing tag(s)", initialCount)

	// Create tag
	t.Log("  - Creating tag")
	color := "#FF5733"
	tag, err := c.Tags.Create(ctx, &fizzy.TagCreateOptions{
		Name:  "test-tag",
		Color: &color,
	})
	require.NoError(t, err, "Failed to create tag")
	require.NotEmpty(t, tag.ID, "Tag ID should not be empty")
	assert.Equal(t, "test-tag", tag.Name, "Tag name mismatch")
	t.Logf("  ✓ Created tag: %s", tag.ID)

	// Verify tag in list
	tags, err := c.Tags.List(ctx)
	require.NoError(t, err, "Failed to list tags after creation")
	assert.Equal(t, initialCount+1, len(tags), "Tag count should increase by 1")
	t.Log("  ✓ Tag appears in list")

	t.Log("✓ Tags tests passed")
}

func testUsers(t *testing.T, ctx context.Context, c *client.Client) {
	t.Log("Testing users...")

	// List users
	t.Log("  - Listing users")
	users, err := c.Users.List(ctx)
	require.NoError(t, err, "Failed to list users")
	assert.NotEmpty(t, users, "Should have at least one user")
	t.Logf("  ✓ Listed %d user(s)", len(users))

	// Verify user data
	for _, user := range users {
		assert.NotEmpty(t, user.ID, "User ID should not be empty")
		assert.NotEmpty(t, user.Name, "User name should not be empty")
	}
	t.Log("  ✓ User data validated")

	t.Log("✓ Users tests passed")
}

func testNotifications(t *testing.T, ctx context.Context, c *client.Client) {
	t.Log("Testing notifications...")

	// List notifications
	t.Log("  - Listing notifications")
	notifications, err := c.Notifications.List(ctx)
	require.NoError(t, err, "Failed to list notifications")
	t.Logf("  ✓ Listed %d notification(s)", len(notifications))

	// If there are notifications, test read operations
	if len(notifications) > 0 {
		notif := notifications[0]
		t.Logf("  - Testing notification operations on ID: %s", notif.ID)

		// Mark as read
		t.Log("    - Marking notification as read")
		err = c.Notifications.Read(ctx, notif.ID)
		require.NoError(t, err, "Failed to mark notification as read")
		t.Log("    ✓ Notification marked read")

		// Mark as unread
		t.Log("    - Marking notification as unread")
		err = c.Notifications.Unread(ctx, notif.ID)
		require.NoError(t, err, "Failed to mark notification as unread")
		t.Log("    ✓ Notification marked unread")
	} else {
		t.Log("  ℹ No notifications to test read/unread operations")
	}

	t.Log("✓ Notifications tests passed")
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	// Verify environment variables are set
	if os.Getenv("FIZZY_TOKEN") == "" || os.Getenv("FIZZY_ACCOUNT") == "" {
		fmt.Println("ERROR: FIZZY_TOKEN and FIZZY_ACCOUNT environment variables must be set")
		fmt.Println("Set them before running integration tests:")
		fmt.Println("  export FIZZY_TOKEN=your_token")
		fmt.Println("  export FIZZY_ACCOUNT=your_account")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}
