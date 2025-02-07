package models

import (
	"bytes"
	"context"
	"embed"
	"math/rand"
	"sync"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/slam"
	"go.viam.com/rdk/spatialmath"
)

const (
	chunkSizeBytes = 1 * 1024 * 1024
)

var (
	Fake = resource.NewModel("cjnj193", "fake-slam", "fake")

	//go:embed pointcloud_2.pcd
	pcFile2 embed.FS
	//go:embed pointcloud_3.pcd
	pcFile3 embed.FS
	//go:embed pointcloud_4.pcd
	pcFile4 embed.FS
	//go:embed internalState4.pbstream
	internalStateFile4 embed.FS
)

func init() {
	resource.RegisterService(slam.API, Fake,
		resource.Registration[slam.Service, *Config]{
			Constructor: newFakeSlamFake,
		},
	)
}

type Config struct {
	NextMap      int  `json:"calls_till_next_map"`
	IsLocalizing bool `json:"is_localizing"`
	/*
		Put config attributes here. There should be public/exported fields
		with a `json` parameter at the end of each attribute.

		Example config struct:
			type Config struct {
				Pin   string `json:"pin"`
				Board string `json:"board"`
				MinDeg *float64 `json:"min_angle_deg,omitempty"`
			}

		If your model does not need a config, replace *Config in the init
		function with resource.NoNativeConfig
	*/

	/* Uncomment this if your model does not need to be validated
	   and has no implicit dependecies. */
	// resource.TriviallyValidateConfig
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, error) {
	// Add config validation code here
	return nil, nil
}

type fakeSlamFake struct {
	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
	nextMap    int
	mode       slam.MappingMode

	mu            sync.Mutex
	pcd           [][]byte
	internalState []byte
	cnt           int

	/* Uncomment this if your model does not need to reconfigure. */
	resource.TriviallyReconfigurable

	// Uncomment this if the model does not have any goroutines that
	// need to be shut down while closing.
	resource.TriviallyCloseable
}

func newFakeSlamFake(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (slam.Service, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &fakeSlamFake{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		nextMap:    conf.NextMap,
	}
	if s.nextMap == 0 {
		s.nextMap = 5
	}
	if conf.IsLocalizing {
		s.mode = slam.MappingModeLocalizationOnly
	} else {
		s.mode = slam.MappingModeNewMap
	}

	// using this as a placeholder image. need to determine the right way to have the module use it
	s.pcd = make([][]byte, 3)
	s.pcd[0], err = pcFile2.ReadFile("pointcloud_2.pcd")
	if err != nil {
		return nil, err
	}
	s.pcd[1], err = pcFile3.ReadFile("pointcloud_3.pcd")
	if err != nil {
		return nil, err
	}
	s.pcd[2], err = pcFile4.ReadFile("pointcloud_4.pcd")
	if err != nil {
		return nil, err
	}
	s.internalState, err = internalStateFile4.ReadFile("internalState4.pbstream")
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *fakeSlamFake) Name() resource.Name {
	return s.name
}

func (s *fakeSlamFake) Position(ctx context.Context) (spatialmath.Pose, error) {

	return spatialmath.NewPose(r3.Vector{X: -rand.Float64()*12000 + 2000, Y: -rand.Float64()*12000 + 2000}, &spatialmath.EulerAngles{Yaw: rand.Float64() * 360}), nil
}

func (s *fakeSlamFake) PointCloudMap(ctx context.Context, returnEditedMap bool) (func() ([]byte, error), error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cnt++
	whichPCD := s.cnt / s.nextMap % 3
	return toChunkedFunc(s.pcd[whichPCD]), nil
}

func (s *fakeSlamFake) InternalState(ctx context.Context) (func() ([]byte, error), error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return toChunkedFunc(s.internalState), nil
}

// toChunkedFunc takes binary data and wraps it in a helper function that converts it into chunks for streaming APIs.
func toChunkedFunc(b []byte) func() ([]byte, error) {
	chunk := make([]byte, chunkSizeBytes)

	reader := bytes.NewReader(b)

	f := func() ([]byte, error) {
		bytesRead, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		return chunk[:bytesRead], err
	}
	return f
}

func (s *fakeSlamFake) Properties(ctx context.Context) (slam.Properties, error) {
	return slam.Properties{CloudSlam: false, MappingMode: s.mode, InternalStateFileType: ".pbstream"}, nil
}

func (s *fakeSlamFake) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
