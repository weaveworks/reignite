// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"github.com/liquidmetal-dev/flintlock/core/application"
	"github.com/liquidmetal-dev/flintlock/core/ports"
	"github.com/liquidmetal-dev/flintlock/infrastructure/containerd"
	"github.com/liquidmetal-dev/flintlock/infrastructure/controllers"
	"github.com/liquidmetal-dev/flintlock/infrastructure/godisk"
	"github.com/liquidmetal-dev/flintlock/infrastructure/grpc"
	"github.com/liquidmetal-dev/flintlock/infrastructure/microvm"
	"github.com/liquidmetal-dev/flintlock/infrastructure/network"
	"github.com/liquidmetal-dev/flintlock/infrastructure/ulid"
	"github.com/liquidmetal-dev/flintlock/internal/config"
	"github.com/liquidmetal-dev/flintlock/pkg/defaults"
	"github.com/spf13/afero"
	"time"
)

// Injectors from wire.go:

func InitializePorts(cfg *config.Config) (*ports.Collection, error) {
	config2 := containerdConfig(cfg)
	microVMRepository, err := containerd.NewMicroVMRepo(config2)
	if err != nil {
		return nil, err
	}
	config3 := networkConfig(cfg)
	networkService := network.New(config3)
	fs := afero.NewOsFs()
	diskService := godisk.New(fs)
	v, err := microvm.NewFromConfig(cfg, networkService, diskService, fs)
	if err != nil {
		return nil, err
	}
	eventService, err := containerd.NewEventService(config2)
	if err != nil {
		return nil, err
	}
	idService := ulid.New()
	imageService, err := containerd.NewImageService(config2)
	if err != nil {
		return nil, err
	}
	collection := appPorts(microVMRepository, v, eventService, idService, networkService, imageService, fs, diskService)
	return collection, nil
}

func InitializeApp(cfg *config.Config, ports2 *ports.Collection) application.App {
	applicationConfig := appConfig(cfg)
	app := application.New(applicationConfig, ports2)
	return app
}

func InializeController(app application.App, ports2 *ports.Collection) *controllers.MicroVMController {
	eventService := eventSvcFromScope(ports2)
	reconcileMicroVMsUseCase := reconcileUCFromApp(app)
	microVMQueryUseCases := queryUCFromApp(app)
	microVMController := controllers.New(eventService, reconcileMicroVMsUseCase, microVMQueryUseCases)
	return microVMController
}

func InitializeGRPCServer(app application.App) ports.MicroVMGRPCService {
	microVMCommandUseCases := commandUCFromApp(app)
	microVMQueryUseCases := queryUCFromApp(app)
	microVMGRPCService := grpc.NewServer(microVMCommandUseCases, microVMQueryUseCases)
	return microVMGRPCService
}

// wire.go:

func containerdConfig(cfg *config.Config) *containerd.Config {
	return &containerd.Config{
		SnapshotterKernel: cfg.CtrSnapshotterKernel,
		SnapshotterVolume: defaults.ContainerdVolumeSnapshotter,
		SocketPath:        cfg.CtrSocketPath,
		Namespace:         cfg.CtrNamespace,
	}
}

func networkConfig(cfg *config.Config) *network.Config {
	return &network.Config{
		ParentDeviceName: cfg.ParentIface,
		BridgeName:       cfg.BridgeName,
	}
}

func appConfig(cfg *config.Config) *application.Config {
	return &application.Config{
		RootStateDir:    cfg.StateRootDir,
		MaximumRetry:    cfg.MaximumRetry,
		DefaultProvider: cfg.DefaultVMProvider,
	}
}

func appPorts(repo ports.MicroVMRepository, providers map[string]ports.MicroVMService, es ports.EventService, is ports.IDService, ns ports.NetworkService, ims ports.ImageService, fs afero.Fs, ds ports.DiskService) *ports.Collection {
	return &ports.Collection{
		Repo:              repo,
		MicrovmProviders:  providers,
		EventService:      es,
		IdentifierService: is,
		NetworkService:    ns,
		ImageService:      ims,
		FileSystem:        fs,
		Clock:             time.Now,
		DiskService:       ds,
	}
}

func eventSvcFromScope(ports2 *ports.Collection) ports.EventService {
	return ports2.EventService
}

func reconcileUCFromApp(app application.App) ports.ReconcileMicroVMsUseCase {
	return app
}

func queryUCFromApp(app application.App) ports.MicroVMQueryUseCases {
	return app
}

func commandUCFromApp(app application.App) ports.MicroVMCommandUseCases {
	return app
}
