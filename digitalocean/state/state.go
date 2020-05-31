package state

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"

	"github.com/rancher/kontainer-engine/drivers/options"
	"github.com/rancher/kontainer-engine/types"
)

type State struct {
	ClusterID string `json:"cluster_id,omitempty"`
	Token string `json:"token,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Name        string `json:"name,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	AutoUpgrade *bool `json:"auto_upgrade,omitempty"`
	RegionSlug  string `json:"region_slug,omitempty"`
	VPCID       string `json:"vpc_id,omitempty"`
	VersionSlug string `json:"version_slug,omitempty"`
	NodePool    NodePool `json:"node_pool,omitempty"`
}

type NodePool struct {
	Name      string            `json:"name,omitempty"`
	Size      string            `json:"size,omitempty"`
	Count     int               `json:"count,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	AutoScale *bool              `json:"auto_scale,omitempty"`
	MinNodes  int               `json:"min_nodes,omitempty"`
	MaxNodes  int               `json:"max_nodes,omitempty"`
}

func (state *State) Save(clusterInfo *types.ClusterInfo) error{
	bytes, err := json.Marshal(state)

	if err != nil {
		return errors.Wrap(err, "could not marshal state")
	}

	if clusterInfo.Metadata == nil {
		clusterInfo.Metadata = make(map[string]string)
	}

	clusterInfo.Metadata["state"] = string(bytes)

	return nil
}

type Builder interface {
	BuildStateFromOpts(driverOptions *types.DriverOptions) (State, error)
	BuildStateFromClusterInfo(clusterInfo *types.ClusterInfo)(State,error)
}

type builderImpl struct{}

func NewBuilder() Builder {
	return builderImpl{}
}

func (builderImpl) BuildStateFromOpts(driverOptions *types.DriverOptions) (State, error) {

	state := State{
		Tags:     []string{},
		NodePool: NodePool{},
	}

	getValue := func(typ string, keys ...string) interface{} {
		return options.GetValueFromDriverOptions(driverOptions, typ, keys...)
	}

	state.Token = getValue(types.StringType, "token").(string)
	state.DisplayName = getValue(types.StringType, "display-name", "displayName").(string)
	state.Name = getValue(types.StringType, "name").(string)
	state.Tags = getTagsFromStringSlice(getValue(types.StringSliceType, "tags").(*types.StringSlice))
	state.AutoUpgrade = getBoolPointer(getValue(types.BoolPointerType, "auto-upgraded", "autoUpgraded"))
	state.RegionSlug = getValue(types.StringType, "region-slug", "regionSlug").(string)
	state.VPCID = getValue(types.StringType, "vpc-id", "vpcID").(string)
	state.VersionSlug = getValue(types.StringType, "version-slug", "versionSlug").(string)
	state.NodePool.Name = getValue(types.StringType, "node-pool-name", "nodePoolName").(string)
	state.NodePool.AutoScale = getBoolPointer(
		getValue(types.BoolPointerType, "node-pool-autoscale", "nodePoolAutoscale"),
	)

	if state.NodePool.AutoScale != nil && *state.NodePool.AutoScale {
		state.NodePool.MaxNodes = int(getValue(types.IntType, "node-pool-max", "nodePoolMax").(int64))
		state.NodePool.MinNodes = int(getValue(types.IntType, "node-pool-min", "nodePoolMin").(int64))
	}

	state.NodePool.Count = int(getValue(types.IntType, "node-pool-count", "nodePoolCount").(int64))

	nodePoolLabels := getLabelsFromStringSlice(
		getValue(types.StringSliceType, "node-pool-labels", "nodePoolLabels").(*types.StringSlice),
	)

	state.NodePool.Labels = nodePoolLabels
	state.NodePool.Size = getValue(types.StringType, "node-pool-size", "nodePoolSize").(string)

	return state, nil
}

func (builderImpl) BuildStateFromClusterInfo(clusterInfo *types.ClusterInfo)(State, error){
	stateJson, ok := clusterInfo.Metadata["state"]
	state := State{}

	if !ok{
		return state, errors.New("there is no state in the clusterInfo")
	}

	err := json.Unmarshal([]byte(stateJson),&state)

	return state, err

}

func getTagsFromStringSlice(tagsString *types.StringSlice)[]string{
	if tagsString.Value == nil {
		return []string{}
	}

	return tagsString.Value
}

func getBoolPointer(boolPointer interface{})*bool{

	if boolPointer == nil {
		return nil
	}
	return boolPointer.(*bool)
}

func getLabelsFromStringSlice(labelsString *types.StringSlice) map[string]string {

	labels := map[string]string{}

	if labelsString == nil {
		return labels
	}

	for _, part := range labelsString.Value {
		kv := strings.Split(part, "=")

		if len(kv) == 2 {
			labels[kv[0]] = kv[1]
		}
	}

	return labels
}
