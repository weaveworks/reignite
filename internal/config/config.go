package config

import (
	"errors"
	"os"
	"time"

	"github.com/weaveworks/flintlock/pkg/log"
)

// Config represents the flintlockd configuration.
type Config struct {
	// ConfigFilePath is the path to the shared configuration file.
	ConfigFilePath string
	// Logging contains the logging related config.
	Logging log.Config
	// GRPCEndpoint is the endpoint for the gRPC server.
	GRPCAPIEndpoint string
	// HTTPAPIEndpoint is the endpoint for the HTTP proxy for the gRPC service..
	HTTPAPIEndpoint string
	// FirecrackerBin is the firecracker binary to use.
	FirecrackerBin string
	// FirecrackerDetatch indicates if the child firecracker processes should be detached from their parent.
	FirecrackerDetatch bool
	// StateRootDir is the directory to act as the root for the runtime state of flintlock.
	StateRootDir string
	// ParentIface is the name of the network interface to use for the parent in macvtap interfaces.
	ParentIface string
	// CtrSnapshotterKernel is the name of the containerd snapshotter to use for kernel images.
	CtrSnapshotterKernel string
	// CtrSocketPath is the path to the containerd socket.
	CtrSocketPath string
	// CtrNamespace is the default containerd namespace to use
	CtrNamespace string
	// DisableReconcile is used to stop the reconcile part from running.
	DisableReconcile bool
	// DisableAPI is used to disable the api server.
	DisableAPI bool
	// ResyncPeriod defines the period when we should do a reconcile of the microvms (even if there are no events).
	ResyncPeriod time.Duration
	// MaximumRetry defined how many times we retry if reconciliation failed.
	MaximumRetry int
	// DeleteVMTimeout defines the timeout for the delete vm operation.
	DeleteVMTimeout time.Duration
	// AuthToken is the static token to use for very basic authentication.
	AuthToken string
	// TLS holds the TLS related configuration.
	TLS TLSConfig
}

// TLSConfig holds the configuration for TLS.
type TLSConfig struct {
	// Insecure indicates if we should start the server insecurely (i.e. without TLS).
	Insecure bool
	// CertFile is the path to the certificate file to use.
	CertFile string
	// KeyFile is the path to the certificate key file to use.
	KeyFile string
	// ValidateClient indicates if the client certificates should be validated.
	ValidateClient bool
	// ClientCAFile is the path to a CA certificate file to use when validating client certificates.
	ClientCAFile string
}

// Validate will validate the TLS config.
func (t TLSConfig) Validate() error {
	if t.Insecure {
		if t.KeyFile != "" || t.CertFile != "" {
			return errNoCertWhenInsecure
		}

		return nil
	}

	if t.CertFile == "" {
		return errCertRequired
	}

	if t.KeyFile == "" {
		return errKeyRequired
	}

	if _, err := os.Stat(t.CertFile); errors.Is(err, os.ErrNotExist) {
		return newCertMissingError("certificate file", t.CertFile)
	}

	if _, err := os.Stat(t.KeyFile); errors.Is(err, os.ErrNotExist) {
		return newCertMissingError("key file", t.KeyFile)
	}

	if t.ClientCAFile != "" {
		if _, err := os.Stat(t.ClientCAFile); errors.Is(err, os.ErrNotExist) {
			return newCertMissingError("client CA file", t.ClientCAFile)
		}
	}

	return nil
}
