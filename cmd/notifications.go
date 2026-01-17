package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/visionik/fizz/internal/format"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Notification management",
	Long:  "List and manage notifications",
}

var notificationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		notifications, err := client.Notifications.List(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to list notifications: %w", err)
		}

		formatter, err := format.NewFormatter(GetFormat(), cmd.OutOrStdout())
		if err != nil {
			return err
		}

		return formatter.Format(notifications)
	},
}

var notificationsReadCmd = &cobra.Command{
	Use:   "read <notification-id>",
	Short: "Mark a notification as read",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		notificationID := args[0]

		err := client.Notifications.Read(cmd.Context(), notificationID)
		if err != nil {
			return fmt.Errorf("failed to mark notification as read: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Notification marked as read\n")
		return nil
	},
}

var notificationsUnreadCmd = &cobra.Command{
	Use:   "unread <notification-id>",
	Short: "Mark a notification as unread",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()
		notificationID := args[0]

		err := client.Notifications.Unread(cmd.Context(), notificationID)
		if err != nil {
			return fmt.Errorf("failed to mark notification as unread: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Notification marked as unread\n")
		return nil
	},
}

var notificationsReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all notifications as read",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := GetClient()

		err := client.Notifications.ReadAll(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to mark all notifications as read: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "All notifications marked as read\n")
		return nil
	},
}

func init() {
	notificationsCmd.AddCommand(notificationsListCmd)
	notificationsCmd.AddCommand(notificationsReadCmd)
	notificationsCmd.AddCommand(notificationsUnreadCmd)
	notificationsCmd.AddCommand(notificationsReadAllCmd)
	rootCmd.AddCommand(notificationsCmd)
}
