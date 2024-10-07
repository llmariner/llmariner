package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Cmd is the root command for usage.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "usage",
		Short:              "Usage commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(summaryCmd())
	return cmd
}

func summaryCmd() *cobra.Command {
	var (
		server string

		start string
		ende  string

		last7d  bool
		last30d bool
	)
	cmd := &cobra.Command{
		Use:   "summary <TENANT_ID>",
		Short: "Get aggregated API usage summary",
		Example: `# [Prerequisite] Start the API usage server.
  kubectl port-forward service/api-usage-server-admin-grpc 8084:8084 -n llmariner &

  # Get the summary in the specified time range.
  llma hidden usage summary default-tenant-id --start="2006-01-02T15:04:05Z" --end="2006-01-30T15:04:05Z"

  # Get the summary in the last 7 days.
  llma hidden usage summary default-tenant-id

  # Get the summary in the last 30 days.
  llma hidden usage summary default-tenant-id --last-30d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("tenant ID is required")
			}
			if server == "" {
				return fmt.Errorf("server address is required")
			}

			var startTime, endTime int64
			if last7d || last30d {
				now := time.Now()
				var t time.Time
				if last7d {
					t = now.Add(-7 * 24 * time.Hour)
				} else {
					t = now.Add(-30 * 24 * time.Hour)
				}
				startTime = t.UnixNano()
				endTime = now.UnixNano()
			} else {
				if start != "" {
					s, err := time.Parse(time.RFC3339, start)
					if err != nil {
						return fmt.Errorf("failed to parse start time: %s", err)
					}
					startTime = s.UnixNano()
				}
				if ende != "" {
					e, err := time.Parse(time.RFC3339, ende)
					if err != nil {
						return fmt.Errorf("failed to parse end time: %s", err)
					}
					endTime = e.UnixNano()
				}
			}

			return getSummary(cmd.Context(), server, args[0], startTime, endTime)
		},
	}
	cmd.Flags().StringVar(&server, "server", "localhost:8084", "Server address")
	cmd.Flags().StringVar(&start, "start", "", "Start date for the summary")
	cmd.Flags().StringVar(&ende, "end", "", "End date for the summary")
	cmd.Flags().BoolVar(&last7d, "last-7d", false, "Get summary for the last 7 days")
	cmd.Flags().BoolVar(&last30d, "last-30d", false, "Get summary for the last 30 days")
	return cmd
}

func getSummary(ctx context.Context, server, tenantID string, start, end int64) error {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.NewClient(server, opt)
	if err != nil {
		return fmt.Errorf("failed to create client: %s", err)
	}
	client := v1.NewAPIUsageServiceClient(cc)

	resp, err := client.GetAggregatedSummary(ctx, &v1.GetAggregatedSummaryRequest{
		TenantId:  tenantID,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		return fmt.Errorf("failed to get summary: %s", err)
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
