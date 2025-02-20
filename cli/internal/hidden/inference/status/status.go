package status

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/inference/status"
)

// Cmd is the root command for inference status.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "status",
		Short:              "Status commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(getCmd())
	return cmd
}

func getCmd() *cobra.Command {
	var ()
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get status from inference-manager-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getStatus(cmd.Context())
		},
	}
	return cmd
}

func getStatus(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var req iv1.GetInferenceStatusRequest
	var resp iv1.InferenceStatus
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
