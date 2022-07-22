package ports

import (
	"context"
	"time"

	mvmv1 "github.com/weaveworks-liquidmetal/flintlock/api/services/microvm/v1alpha1"
	"github.com/weaveworks-liquidmetal/flintlock/core/models"
)

// MicroVMService is the port definition for a microvm service.
type MicroVMService interface {
	// Capabilities returns a list of the capabilities the provider supports.
	Capabilities() models.Capabilities

	// Create will create a new microvm.
	Create(ctx context.Context, vm *models.MicroVM) error
	// Delete will delete a VM and its runtime state.
	Delete(ctx context.Context, id string) error
	// Start will start a created microvm.
	Start(ctx context.Context, vm *models.MicroVM) error
	// State returns the state of a microvm.
	State(ctx context.Context, id string) (MicroVMState, error)
	// Metrics returns with the metrics of a microvm.
	Metrics(ctx context.Context, id models.VMID) (MachineMetrics, error)
}

// This state represents the state of the Firecracker MVM process itself
// The state for the entire Flintlock MVM is represented in models.MicroVMState.
type MicroVMState string

// MachineMetrics is a metrics interface for providers.
type MachineMetrics interface {
	ToPrometheus() []byte
}

const (
	MicroVMStateUnknown    MicroVMState = "unknown"
	MicroVMStatePending    MicroVMState = "pending"
	MicroVMStateConfigured MicroVMState = "configured"
	MicroVMStateRunning    MicroVMState = "running"
)

// MicroVMGRPCService is a port for a microvm grpc service.
type MicroVMGRPCService interface {
	mvmv1.MicroVMServer
}

// IDService is a port for a service for working with identifiers.
type IDService interface {
	// GenerateRandom will create a random identifier.
	GenerateRandom() (string, error)
}

// EventService is a port for a service that acts as a event bus.
type EventService interface {
	// Publish will publish an event to a specific topic.
	Publish(ctx context.Context, topic string, eventToPublish interface{}) error
	// SubscribeTopic will subscribe to events on a named topic..
	SubscribeTopic(ctx context.Context, topic string) (ch <-chan *EventEnvelope, errs <-chan error)
	// SubscribeTopics will subscribe to events on a set of named topics.
	SubscribeTopics(ctx context.Context, topics []string) (ch <-chan *EventEnvelope, errs <-chan error)
	// Subscribe will subscribe to events on all topics
	Subscribe(ctx context.Context) (ch <-chan *EventEnvelope, errs <-chan error)
}

type EventEnvelope struct {
	Timestamp time.Time
	Namespace string
	Topic     string
	Event     interface{}
}

// ImageService is a port for a service that interacts with OCI images.
type ImageService interface {
	// Pull will get (i.e. pull) the image for a specific owner.
	Pull(ctx context.Context, input *ImageSpec) error
	// PullAndMount will get (i.e. pull) the image for a specific owner and then
	// make it available via a mount point.
	PullAndMount(ctx context.Context, input *ImageMountSpec) ([]models.Mount, error)
	// Exists checks if the image already exists on the machine.
	Exists(ctx context.Context, input *ImageSpec) (bool, error)
	// IsMounted checks if the image is pulled and mounted.
	IsMounted(ctx context.Context, input *ImageMountSpec) (bool, error)
}

type ImageSpec struct {
	// ImageName is the name of the image to get.
	ImageName string
	// Owner is the name of the owner of the image.
	Owner string
}

// ImageMountSpec is the declaration of an image that needs to be pulled and mounted.
type ImageMountSpec struct {
	// ImageName is the name of the image to get.
	ImageName string
	// Owner is the name of the owner of the image.
	Owner string
	// Use is an indicator of what the image will be used for.
	Use models.ImageUse
	// OwnerUsageID is an identifier from the owner.
	OwnerUsageID string
}

// NetworkService is a port for a service that interacts with the network
// stack on the host machine.
type NetworkService interface {
	// IfaceCreate will create the network interface.
	IfaceCreate(ctx context.Context, input IfaceCreateInput) (*IfaceDetails, error)
	// IfaceDelete is used to delete a network interface
	IfaceDelete(ctx context.Context, input DeleteIfaceInput) error
	// IfaceExists will check if an interface with the given name exists
	IfaceExists(ctx context.Context, name string) (bool, error)
	// IfaceDetails will get the details of the supplied network interface.
	IfaceDetails(ctx context.Context, name string) (*IfaceDetails, error)
}

type IfaceCreateInput struct {
	// DeviceName is the name of the network interface to create on the host.
	DeviceName string
	// Type is the type of network interface to create.
	Type models.IfaceType
	// MAC allows the specifying of a specific MAC address to use for the interface. If
	// not supplied a autogenerated MAC address will be used.
	MAC string
	// Attach indicates if this device should be attached to the parent bridge. Only applicable to TAP devices.
	Attach bool
	// BridgeName is the name of the bridge to attach to. Only if this is a tap device and attach is true.
	BridgeName string
}

type IfaceDetails struct {
	// DeviceName is the name of the network interface created on the host.
	DeviceName string
	// Type is the type of network interface created.
	Type models.IfaceType
	// MAC is the MAC address of the created interface.
	MAC string
	// Index is the network interface index on the host.
	Index int
}

type DeleteIfaceInput struct {
	// DeviceName is the name of the network interface to delete from the host.
	DeviceName string
}

// DiskService is a port for a service that creates disk images.
type DiskService interface {
	// Create will create a new disk.
	Create(ctx context.Context, input DiskCreateInput) error
}

// DiskType represents the type of disk.
type DiskType int

const (
	// DiskTypeFat32 is a FAT32 compatible filesystem.
	DiskTypeFat32 DiskType = iota
	// DiskTypeISO9660 is an iso filesystem.
	DiskTypeISO9660
)

// DiskCreateInput are the input options for creating a disk.
type DiskCreateInput struct {
	//Path is the filesystem path of where to create the disk.
	Path string
	// Size is how big the disk should be. It uses human readable formats
	// suck as 8Mb, 10Kb.
	Size string
	// VolumeName is the name to give to the volume.
	VolumeName string
	// Type is the type of disk to create.
	Type DiskType
	// Files are the files to create in the new disk.
	Files []DiskFile
}

// DiskFile represents a file to create in a disk.
type DiskFile struct {
	// Path is the path in the disk image for the file.
	Path string
	// ContentBase64 is the content of the file encoded as base64.
	ContentBase64 string
}
