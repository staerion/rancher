package nodescaling

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/gke"
	"github.com/rancher/rancher/tests/framework/extensions/scalinginput"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GKENodeScalingTestSuite struct {
	suite.Suite
	client        *rancher.Client
	session       *session.Session
	scalingConfig *scalinginput.Config
}

func (s *GKENodeScalingTestSuite) TearDownSuite() {
	s.session.Cleanup()
}

func (s *GKENodeScalingTestSuite) SetupSuite() {
	testSession := session.NewSession()
	s.session = testSession

	s.scalingConfig = new(scalinginput.Config)
	config.LoadConfig(scalinginput.ConfigurationFileKey, s.scalingConfig)

	client, err := rancher.NewClient("", testSession)
	require.NoError(s.T(), err)

	s.client = client
}

func (s *GKENodeScalingTestSuite) TestScalingGKENodePools() {
	scaleOneNode := gke.NodePool{
		InitialNodeCount: &oneNode,
	}

	scaleTwoNodes := gke.NodePool{
		InitialNodeCount: &twoNodes,
	}

	tests := []struct {
		name     string
		gkeNodes gke.NodePool
		client   *rancher.Client
	}{
		{"Scaling node group by 1", scaleOneNode, s.client},
		{"Scaling node group by 2", scaleTwoNodes, s.client},
	}

	for _, tt := range tests {
		clusterID, err := clusters.GetClusterIDByName(s.client, s.client.RancherConfig.ClusterName)
		require.NoError(s.T(), err)

		s.Run(tt.name, func() {
			ScalingGKENodePools(s.T(), s.client, clusterID, &tt.gkeNodes)
		})
	}
}

func (s *GKENodeScalingTestSuite) TestScalingGKENodePoolsDynamicInput() {
	if s.scalingConfig.GKENodePool == nil {
		s.T().Skip()
	}

	clusterID, err := clusters.GetClusterIDByName(s.client, s.client.RancherConfig.ClusterName)
	require.NoError(s.T(), err)

	ScalingGKENodePools(s.T(), s.client, clusterID, s.scalingConfig.GKENodePool)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGKENodeScalingTestSuite(t *testing.T) {
	suite.Run(t, new(GKENodeScalingTestSuite))
}
