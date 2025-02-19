package eks

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/provisioning"
	"github.com/rancher/rancher/tests/framework/extensions/provisioninginput"
	"github.com/rancher/rancher/tests/framework/extensions/users"
	password "github.com/rancher/rancher/tests/framework/extensions/users/passwordgenerator"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HostedEKSClusterProvisioningTestSuite struct {
	suite.Suite
	client             *rancher.Client
	session            *session.Session
	standardUserClient *rancher.Client
}

func (h *HostedEKSClusterProvisioningTestSuite) TearDownSuite() {
	h.session.Cleanup()
}

func (h *HostedEKSClusterProvisioningTestSuite) SetupSuite() {
	testSession := session.NewSession()
	h.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(h.T(), err)

	h.client = client

	enabled := true
	var testuser = namegen.AppendRandomString("testuser-")
	var testpassword = password.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(h.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(h.T(), err)

	h.standardUserClient = standardUserClient
}

func (h *HostedEKSClusterProvisioningTestSuite) TestProvisioningHostedEKS() {
	tests := []struct {
		name   string
		client *rancher.Client
	}{
		{provisioninginput.AdminClientName.String(), h.client},
		{provisioninginput.StandardClientName.String(), h.standardUserClient},
	}

	for _, tt := range tests {
		clusterObject, err := provisioning.CreateProvisioningEKSHostedCluster(tt.client)
		require.NoError(h.T(), err)

		provisioning.VerifyHostedCluster(h.T(), tt.client, clusterObject)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHostedEKSClusterProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(HostedEKSClusterProvisioningTestSuite))
}
