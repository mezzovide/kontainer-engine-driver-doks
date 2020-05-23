package digitalocean

import (
	"context"
	"errors"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/service"
	"github.com/sirupsen/logrus"

	"github.com/rancher/kontainer-engine/types"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/options"
	"github.com/ribeiro-rodrigo/kontainer-engine-driver-digitalocean/digitalocean/state"
)

type Driver struct {
	stateBuilder       state.Builder
	optionsBuilder     options.Builder
	digitalOceanFactory service.DigitalOceanFactory
	driverCapabilities types.Capabilities
}

func NewDriver() Driver {
	driver := Driver{
		stateBuilder:   state.NewBuilder(),
		optionsBuilder: options.NewBuilder(),
		digitalOceanFactory: service.NewDigitalOceanFactory(),
	}

	return driver
}

func (driver *Driver) GetDriverCreateOptions(ctx context.Context) (*types.DriverFlags, error) {
	logrus.Debug("DigitalOcean.Driver.GetDriverCreateOptions(...) called")
	return driver.optionsBuilder.BuildCreateOptions(), nil
}

func (*Driver) GetDriverUpdateOptions(ctx context.Context) (*types.DriverFlags, error) {
	return nil, nil
}

func (driver *Driver) Create(ctx context.Context, opts *types.DriverOptions, info *types.ClusterInfo) (*types.ClusterInfo, error) {
	logrus.Debug("DigitalOcean.Driver.Create(...) called")
	state, err := driver.stateBuilder.BuildStateFromOpts(opts)

	if err != nil{
		logrus.Debugf("Error building state: %v",err)
		return nil, err
	}

	if state.Token == ""{
		logrus.Debugf("Error token not found: %v",err)
		err = errors.New("token was not reported")
		return nil, err
	}

	digitalOceanService := driver.digitalOceanFactory(state.Token)

	clusterID, err := digitalOceanService.CreateCluster(ctx,state)

	if err != nil {
		logrus.Debugf("Error crate cluster: %v",err)
		return nil, err
	}

	state.ClusterID = clusterID

	err = state.Save(info)

	if err != nil {
		logrus.Debugf("Error save state: %v",err)
		return nil, err
	}

	return info, nil
}

func (*Driver) PostCheck(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.ClusterInfo, error) {

	/*kubeConfig, err := digitalOceanService.GetKubeConfig(clusterID)

	if err != nil {
		logrus.Debugf("Error get kubeConfig %v",err)
		return nil, err
	} */


	return nil, nil
}

func (*Driver) Update(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions) (*types.ClusterInfo, error) {

	return nil, nil
}

func (*Driver) Remove(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*Driver) GetVersion(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.KubernetesVersion, error) {
	return nil, nil
}

func (*Driver) SetVersion(ctx context.Context, clusterInfo *types.ClusterInfo, version *types.KubernetesVersion) error {
	return nil
}

func (*Driver) GetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo) (*types.NodeCount, error) {
	return nil, nil
}

func (*Driver) SetClusterSize(ctx context.Context, clusterInfo *types.ClusterInfo, count *types.NodeCount) error {

	return nil
}

func (*Driver) GetCapabilities(ctx context.Context) (*types.Capabilities, error) {
	return nil, nil
}

func (*Driver) RemoveLegacyServiceAccount(ctx context.Context, clusterInfo *types.ClusterInfo) error {
	return nil
}

func (*Driver) ETCDSave(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}

func (*Driver) ETCDRestore(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) (*types.ClusterInfo, error) {
	return nil, nil
}

func (*Driver) GetK8SCapabilities(ctx context.Context, opts *types.DriverOptions) (*types.K8SCapabilities, error) {
	return nil, nil
}

func (*Driver) ETCDRemoveSnapshot(ctx context.Context, clusterInfo *types.ClusterInfo, opts *types.DriverOptions, snapshotName string) error {
	return nil
}

