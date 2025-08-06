package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/AlecAivazis/survey/v2"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	mv1 "github.com/llmariner/model-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/models"
)

// Cmd is the root command for models.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "models",
		Short:              "Models commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(updateCmd())
	cmd.AddCommand(activateCmd())
	cmd.AddCommand(deactivateCmd())
	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "create",
		Short:              "Create commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createBaseCmd())
	cmd.AddCommand(createFineTunedCmd())
	return cmd
}

func createBaseCmd() *cobra.Command {
	var (
		repoStr      string
		globalScoped bool

		config = mv1.ModelConfig{
			RuntimeConfig: &mv1.ModelConfig_RuntimeConfig{
				Resources: &mv1.ModelConfig_RuntimeConfig_Resources{},
			},
			ClusterAllocationPolicy: &mv1.ModelConfig_ClusterAllocationPolicy{},
		}
		disableOnDemandAllocation bool
	)
	cmd := &cobra.Command{
		Use:  "base <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := toSourceRepositoryEnum(repoStr)
			if err != nil {
				return err
			}
			config.ClusterAllocationPolicy.EnableOnDemandAllocation = !disableOnDemandAllocation
			return createBase(cmd.Context(), args[0], repo, globalScoped, &config)
		},
	}

	cmd.Flags().StringVar(&repoStr, "source-repository", "", "Source repository. One of 'object-store', 'hugging-face' or 'ollama'.")
	cmd.Flags().BoolVar(&globalScoped, "global-scoped", false, "If true, the model is global scoped. Otherwise, it has a project scope.")

	cmd.Flags().Int32Var(&config.RuntimeConfig.Resources.Gpu, "gpu", 1, "Number of GPUs to use for the model. Default is 1.")
	cmd.Flags().Int32Var(&config.RuntimeConfig.Replicas, "replicas", 1, "Number of replicas to use for the model. Default is 1.")
	cmd.Flags().BoolVar(&disableOnDemandAllocation, "disable-on-demand-allocation", false, "If true, the model is not activated on demand.")

	_ = cmd.MarkFlagRequired("source-repository")
	return cmd
}

func createFineTunedCmd() *cobra.Command {
	var (
		baseModelID       string
		suffix            string
		repoStr           string
		modelFileLocation string

		config = mv1.ModelConfig{
			RuntimeConfig: &mv1.ModelConfig_RuntimeConfig{
				Resources: &mv1.ModelConfig_RuntimeConfig_Resources{},
			},
			ClusterAllocationPolicy: &mv1.ModelConfig_ClusterAllocationPolicy{},
		}
		disableOnDemandAllocation bool
	)

	cmd := &cobra.Command{
		Use:  "fine-tuned",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := toSourceRepositoryEnum(repoStr)
			if err != nil {
				return err
			}
			config.ClusterAllocationPolicy.EnableOnDemandAllocation = !disableOnDemandAllocation
			return createFineTuned(cmd.Context(), baseModelID, suffix, repo, modelFileLocation, &config)
		},
	}

	cmd.Flags().StringVar(&baseModelID, "base-model-id", "", "Base model ID.")
	cmd.Flags().StringVar(&suffix, "suffix", "", "Suffix for the model ID.")
	cmd.Flags().StringVar(&repoStr, "source-repository", "", "Source repository. One of 'object-store', 'hugging-face' or 'ollama'.")
	cmd.Flags().StringVar(&modelFileLocation, "model-file-location", "", "Model file location.")

	cmd.Flags().Int32Var(&config.RuntimeConfig.Resources.Gpu, "gpu", 1, "Number of GPUs to use for the model. Default is 1.")
	cmd.Flags().Int32Var(&config.RuntimeConfig.Replicas, "replicas", 1, "Number of replicas to use for the model. Default is 1.")
	cmd.Flags().BoolVar(&disableOnDemandAllocation, "disable-on-demand-allocation", false, "If true, the model is not activated on demand.")

	_ = cmd.MarkFlagRequired("base-model-id")
	_ = cmd.MarkFlagRequired("suffix")
	_ = cmd.MarkFlagRequired("source-repository")
	_ = cmd.MarkFlagRequired("model-file-location")
	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd.Context())
		},
	}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <ID>",
		Short: "Retrieve details of a specific model by ID",
		Args:  validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args[0])
		},
	}
	return cmd
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0])
		},
	}
}

func updateCmd() *cobra.Command {
	var (
		gpu                       int32
		replicas                  int32
		disableOnDemandAllocation bool
	)
	cmd := &cobra.Command{
		Use:  "update <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			var config mv1.ModelConfig

			if cmd.Flags().Changed("disable-on-demand-allocation") {
				config.ClusterAllocationPolicy = &mv1.ModelConfig_ClusterAllocationPolicy{
					EnableOnDemandAllocation: !disableOnDemandAllocation,
				}
			}
			if cmd.Flags().Changed("gpu") || cmd.Flags().Changed("replicas") {
				config.RuntimeConfig = &mv1.ModelConfig_RuntimeConfig{}
				if cmd.Flags().Changed("gpu") {
					config.RuntimeConfig.Resources = &mv1.ModelConfig_RuntimeConfig_Resources{}
					config.RuntimeConfig.Resources.Gpu = gpu
				}
				if cmd.Flags().Changed("replicas") {
					config.RuntimeConfig.Replicas = replicas
				}
			}

			return update(cmd.Context(), args[0], &config)
		},
	}

	cmd.Flags().Int32Var(&gpu, "gpu", 1, "Number of GPUs to use for the model. Default is 1.")
	cmd.Flags().Int32Var(&replicas, "replicas", 1, "Number of replicas to use for the model. Default is 1.")
	cmd.Flags().BoolVar(&disableOnDemandAllocation, "disable-on-demand-allocation", false, "If true, the model is not activated on demand.")

	return cmd
}

func activateCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "activate <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return activate(cmd.Context(), args[0])
		},
	}
}

func deactivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "deactivate <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deactivate(cmd.Context(), args[0])
		},
	}
}

func createBase(
	ctx context.Context,
	id string,
	repo mv1.SourceRepository,
	globalScoped bool,
	config *mv1.ModelConfig,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := mv1.CreateModelRequest{
		IsFineTunedModel: false,
		Id:               id,
		SourceRepository: repo,
		IsProjectScoped:  !globalScoped,
		Config:           config,
	}
	var resp mv1.Model
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the model (ID: %q).\n", id)
	fmt.Printf("The model becomes available once it is loaded.\n")
	return nil
}

func createFineTuned(
	ctx context.Context,
	baseModelID string,
	suffix string,
	repo mv1.SourceRepository,
	modelFileLocation string,
	config *mv1.ModelConfig,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := mv1.CreateModelRequest{
		IsFineTunedModel:  true,
		BaseModelId:       baseModelID,
		Suffix:            suffix,
		SourceRepository:  repo,
		ModelFileLocation: modelFileLocation,
		Config:            config,
	}
	var resp mv1.Model
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the model (ID: %q).\n", resp.Id)
	fmt.Printf("The model becomes available once it is loaded.\n")
	return nil
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := mv1.ListModelsRequest{
		IncludeLoadingModels: true,
	}
	var resp mv1.ListModelsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	// Sort models by names.
	var ms []*mv1.Model
	ms = append(ms, resp.Data...)
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Id < ms[j].Id
	})

	tbl := table.New("ID", "Owned By", "Loading Status", "Activation Status", "Source Repository", "Created At")
	ui.FormatTable(tbl)

	for _, m := range ms {
		r := toloadingStatusString(m.LoadingStatus)
		switch m.LoadingStatus {
		case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_LOADING:
			if msg := m.LoadingStatusMessage; msg != "" {
				r = fmt.Sprintf("%s (%s)", r, msg)
			}
		case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_FAILED:
			r = fmt.Sprintf("%s (%s)", r, m.LoadingFailureReason)
		}
		a := toActivationStatusString(m.ActivationStatus)
		tbl.AddRow(
			m.Id,
			m.OwnedBy,
			r,
			a,
			toSourceRepositoryString(m.SourceRepository),
			time.Unix(m.Created, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}

func get(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := mv1.GetModelRequest{
		IncludeLoadingModel: true,
	}
	var m mv1.Model
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &req, &m); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&m, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func delete(ctx context.Context, id string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete model %q?", id),
		Default: false,
	}
	var ok bool
	if err := p.Ask(s, &ok); err != nil {
		return err
	} else if !ok {
		return nil
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := mv1.DeleteModelRequest{
		Id: id,
	}
	var resp mv1.DeleteModelResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the model (ID: %q).\n", id)

	return nil
}

func update(
	ctx context.Context,
	id string,
	config *mv1.ModelConfig,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &mv1.Model{
		Config: config,
	}
	var resp mv1.Model
	if err := ihttp.NewClient(env).Send(http.MethodPatch, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Updated the model (ID: %q).\n", id)
	return nil
}

func activate(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &mv1.ActivateModelRequest{
		Id: id,
	}
	var resp mv1.ActivateModelResponse
	if err := ihttp.NewClient(env).Send(http.MethodPost, fmt.Sprintf("%s/%s:activate", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Sent the activaiton request for the model (ID: %q).\n", id)

	return nil
}

func deactivate(ctx context.Context, id string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Deactivate model %q?", id),
		Default: false,
	}
	var ok bool
	if err := p.Ask(s, &ok); err != nil {
		return err
	} else if !ok {
		return nil
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &mv1.DeactivateModelRequest{
		Id: id,
	}
	var resp mv1.DeactivateModelResponse
	if err := ihttp.NewClient(env).Send(http.MethodPost, fmt.Sprintf("%s/%s:deactivate", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Sent the deeactivation request for the model (ID: %q).\n", id)

	return nil
}

func validateIDArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<ID> is required argument")
	}
	return nil
}

func toSourceRepositoryEnum(repoStr string) (mv1.SourceRepository, error) {
	switch repoStr {
	case "object-store":
		return mv1.SourceRepository_SOURCE_REPOSITORY_OBJECT_STORE, nil
	case "hugging-face":
		return mv1.SourceRepository_SOURCE_REPOSITORY_HUGGING_FACE, nil
	case "ollama":
		return mv1.SourceRepository_SOURCE_REPOSITORY_OLLAMA, nil
	default:
		return mv1.SourceRepository_SOURCE_REPOSITORY_UNSPECIFIED, fmt.Errorf("invalid source repository %q. Must be 'object-store', 'hugging-face' or 'ollama'", repoStr)
	}
}

func toSourceRepositoryString(repo mv1.SourceRepository) string {
	switch repo {
	case mv1.SourceRepository_SOURCE_REPOSITORY_OBJECT_STORE:
		return "object-store"
	case mv1.SourceRepository_SOURCE_REPOSITORY_HUGGING_FACE:
		return "hugging-face"
	case mv1.SourceRepository_SOURCE_REPOSITORY_OLLAMA:
		return "ollama"
	case mv1.SourceRepository_SOURCE_REPOSITORY_FINE_TUNING:
		return "fine-tuning"
	default:
		return "Unknown"
	}
}

func toloadingStatusString(status mv1.ModelLoadingStatus) string {
	switch status {
	case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_UNSPECIFIED:
		// Considered as "succeeded" for backward compatibility.
		return "succeeded"
	case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_REQUESTED:
		return "requested"
	case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_LOADING:
		return "loading"
	case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_SUCCEEDED:
		return "succeeded"
	case mv1.ModelLoadingStatus_MODEL_LOADING_STATUS_FAILED:
		return "failed"
	default:
		return "Unknown"
	}
}

func toActivationStatusString(status mv1.ActivationStatus) string {
	switch status {
	case mv1.ActivationStatus_ACTIVATION_STATUS_UNSPECIFIED:
		return "unspecified"
	case mv1.ActivationStatus_ACTIVATION_STATUS_ACTIVE:
		return "active"
	case mv1.ActivationStatus_ACTIVATION_STATUS_INACTIVE:
		return "inactive"
	default:
		return "unknown"
	}
}
