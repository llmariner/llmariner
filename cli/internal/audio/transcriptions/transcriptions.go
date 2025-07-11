package transcriptions

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/audio/transcriptions"
)

// Cmd is the root command for transcriptions.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "transcriptions",
		Short:              "Transcriptions commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	return cmd
}

func createCmd() *cobra.Command {
	req := &iv1.CreateAudioTranscriptionRequest{}
	var (
		filename string
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				return fmt.Errorf("file %q does not exist", filename)
			}

			return create(cmd.Context(), req, filename)
		},
	}
	cmd.Flags().StringVar(&filename, "file", "", "Audio file")
	cmd.Flags().StringVar(&req.Model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&req.Prompt, "prompt", "", "Optional text to guide the model's style or continue a previous audio segment")
	cmd.Flags().Float64Var(&req.Temperature, "temperature", 0.0, "Temperature")

	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func create(
	ctx context.Context,
	req *iv1.CreateAudioTranscriptionRequest,
	filename string,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	reqB, contentType, err := buildRequest(req, filename)
	if err != nil {
		return err
	}

	var resp iv1.Transcription
	if err := ihttp.NewClient(env).SendMultipart(path, reqB, contentType, &resp); err != nil {
		return err
	}

	fmt.Println(resp.Text)
	return nil

}

func buildRequest(req *iv1.CreateAudioTranscriptionRequest, filename string) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer func() {
		_ = w.Close()
	}()

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", err
	}

	fb, err := os.ReadFile(filename)
	if err != nil {
		return nil, "", fmt.Errorf("read file %q: %s", filename, err)
	}

	if _, err := fw.Write(fb); err != nil {
		return nil, "", err
	}

	fw, err = w.CreateFormField("model")
	if err != nil {
		return nil, "", err
	}
	if _, err := fw.Write([]byte(req.Model)); err != nil {
		return nil, "", err
	}

	if req.Prompt != "" {
		fw, err = w.CreateFormField("prompt")
		if err != nil {
			return nil, "", err
		}
		if _, err := fw.Write([]byte(req.Prompt)); err != nil {
			return nil, "", err
		}
	}

	fw, err = w.CreateFormField("temperature")
	if err != nil {
		return nil, "", err
	}
	if _, err := fmt.Fprintf(fw, "%f", req.Temperature); err != nil {
		return nil, "", err
	}

	return &b, w.FormDataContentType(), nil

}
