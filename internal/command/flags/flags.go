package flags

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weaveworks/flintlock/internal/config"
	"github.com/weaveworks/flintlock/pkg/defaults"
)

const (
	grpcEndpointFlag      = "grpc-endpoint"
	httpEndpointFlag      = "http-endpoint"
	parentIfaceFlag       = "parent-iface"
	disableReconcileFlag  = "disable-reconcile"
	disableAPIFlag        = "disable-api"
	firecrackerBinFlag    = "firecracker-bin"
	firecrackerDetachFlag = "firecracker-detach"
	containerdSocketFlag  = "containerd-socket"
	kernelSnapshotterFlag = "containerd-kernel-ss"
	containerdNamespace   = "containerd-ns"
	maximumRetryFlag      = "maximum-retry"
	authTokenFlag         = "auth-token"
	insecureFlag          = "insecure"
	tlsCertFlag           = "tls-cert"
	tlsKeyFlag            = "tls-key"
	tlsClientValidateFlag = "tls-client-validate"
	tlsClientCAFlag       = "tls-client-ca"
)

// AddGRPCServerFlagsToCommand will add gRPC server flags to the supplied command.
func AddGRPCServerFlagsToCommand(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().StringVar(&cfg.GRPCAPIEndpoint,
		grpcEndpointFlag,
		defaults.GRPCAPIEndpoint,
		"The endpoint for the gRPC server to listen on.")

	cmd.Flags().StringVar(&cfg.AuthToken,
		authTokenFlag,
		"",
		"The token to use for very basic token based authentication.")
}

// AddGWServerFlagsToCommand will add gRPC HTTP gateway flags to the supplied command.
func AddGWServerFlagsToCommand(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().StringVar(&cfg.GRPCAPIEndpoint,
		grpcEndpointFlag,
		defaults.GRPCAPIEndpoint,
		"The address of the gRPC server to act as a gateway for.")

	cmd.Flags().StringVar(&cfg.HTTPAPIEndpoint,
		httpEndpointFlag,
		defaults.HTTPAPIEndpoint,
		"The endpoint for the HTTP proxy to the gRPC service to listen on.")
}

func AddNetworkFlagsToCommand(cmd *cobra.Command, cfg *config.Config) error {
	cmd.Flags().StringVar(&cfg.ParentIface,
		parentIfaceFlag,
		"",
		"The parent iface for the network interfaces. Note it could also be a bond")

	if err := cmd.MarkFlagRequired(parentIfaceFlag); err != nil {
		return fmt.Errorf("setting %s as required: %w", parentIfaceFlag, err)
	}

	return nil
}

func AddHiddenFlagsToCommand(cmd *cobra.Command, cfg *config.Config) error {
	cmd.Flags().BoolVar(&cfg.DisableReconcile,
		disableReconcileFlag,
		false,
		"Set to true to stop the reconciler running")

	cmd.Flags().IntVar(&cfg.MaximumRetry,
		maximumRetryFlag,
		defaults.MaximumRetry,
		"Number of times to retry failed reconciliation")

	cmd.Flags().BoolVar(&cfg.DisableAPI,
		disableAPIFlag,
		false,
		"Set to true to stop the api server running")

	if err := cmd.Flags().MarkHidden(disableReconcileFlag); err != nil {
		return fmt.Errorf("setting %s as hidden: %w", disableReconcileFlag, err)
	}

	if err := cmd.Flags().MarkHidden(maximumRetryFlag); err != nil {
		return fmt.Errorf("setting %s as hidden: %w", maximumRetryFlag, err)
	}

	if err := cmd.Flags().MarkHidden(disableAPIFlag); err != nil {
		return fmt.Errorf("setting %s as hidden: %w", disableAPIFlag, err)
	}

	return nil
}

// AddFirecrackerFlagsToCommand will add the firecracker provider specific flags to the supplied cobra command.
func AddFirecrackerFlagsToCommand(cmd *cobra.Command, cfg *config.Config) error {
	cmd.Flags().StringVar(&cfg.FirecrackerBin,
		firecrackerBinFlag,
		defaults.FirecrackerBin,
		"The path to the firecracker binary to use.")
	cmd.Flags().BoolVar(&cfg.FirecrackerDetatch,
		firecrackerDetachFlag,
		defaults.FirecrackerDetach,
		"If true the child firecracker processes will be detached from the parent flintlock process.")

	return nil
}

// AddContainerDFlagsToCommand will add the containerd specific flags to the supplied cobra command.
func AddContainerDFlagsToCommand(cmd *cobra.Command, cfg *config.Config) error {
	cmd.Flags().StringVar(&cfg.CtrSocketPath,
		containerdSocketFlag,
		defaults.ContainerdSocket,
		"The path to the containerd socket.")

	cmd.Flags().StringVar(&cfg.CtrSnapshotterKernel,
		kernelSnapshotterFlag,
		defaults.ContainerdKernelSnapshotter,
		"The name of the snapshotter to use with containerd for kernel/initrd images.")

	cmd.Flags().StringVar(&cfg.CtrNamespace,
		containerdNamespace,
		defaults.ContainerdNamespace,
		"The name of the containerd namespace to use.")

	return nil
}

func AddTLSFlagsToCommand(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().BoolVar(&cfg.TLS.Insecure,
		insecureFlag,
		false,
		"Run the gRPC server insecurely (i.e. without TLS). Not recommended.")

	cmd.Flags().StringVar(&cfg.TLS.CertFile,
		tlsCertFlag,
		"",
		"Path to the certificate to use for TLS.")

	cmd.Flags().StringVar(&cfg.TLS.KeyFile,
		tlsKeyFlag,
		"",
		"Path to the key to use for TLS.")

	cmd.Flags().BoolVar(&cfg.TLS.ValidateClient,
		tlsClientValidateFlag,
		false,
		"Validate the certificates of clients calling the gRPC server.")

	cmd.Flags().StringVar(&cfg.TLS.ClientCAFile,
		tlsClientCAFlag,
		"",
		"Path to the certificate to use when validating client certificates.")
}
