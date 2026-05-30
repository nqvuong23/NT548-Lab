package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
//
//	Helpers
//
// ─────────────────────────────────────────────
func newEC2Client(t *testing.T, region string) *ec2.EC2 {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	require.NoError(t, err)
	return ec2.New(sess)
}

// ─────────────────────────────────────────────
//
//	Shared Terraform options (used by every test)
//
// ─────────────────────────────────────────────
func terraformOptions(t *testing.T) *terraform.Options {
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../", // thư mục gốc Terraform (terraform/)
		Vars: map[string]interface{}{
			"project_name":                   "nt548-test",
			"aws_region":                     "ap-southeast-1",
			"environment":                    "test",
			"vpc_cidr":                       "10.0.0.0/16",
			"public_subnet_cidrs":            []string{"10.0.1.0/24", "10.0.2.0/24"},
			"private_subnet_cidrs":           []string{"10.0.3.0/24", "10.0.4.0/24"},
			"availability_zones":             []string{"ap-southeast-1a", "ap-southeast-1b"},
			"ec2_ami_id":                     "ami-0df7a207adb9748c7", // Amazon Linux 2 ap-southeast-1
			"ec2_instance_type":              "t3.micro",
			"ec2_volume_type":                "gp3",
			"ec2_volume_size":                8,
			"ec2_enable_detailed_monitoring": false,
		},

		EnvVars: map[string]string{
            "TF_CLI_ARGS_init": "-backend=false",
        },
	})
}

// ═══════════════════════════════════════════════════════════
// TestAll — apply một lần, chạy tất cả checks
// ═══════════════════════════════════════════════════════════
func TestAll(t *testing.T) {
	t.Parallel()

	opts := terraformOptions(t)
	defer terraform.Destroy(t, opts)
	terraform.RunTerraformCommand(t, opts, "init", "-reconfigure", "-upgrade")
	terraform.Apply(t, opts)

	region := "ap-southeast-1"
	client := newEC2Client(t, region)

	// ── Lấy outputs ──────────────────────────────────────────────────────────
	vpcID             := terraform.Output(t, opts, "vpc_id")
	publicSubnetIDs   := terraform.OutputList(t, opts, "public_subnet_ids")
	privateSubnetIDs  := terraform.OutputList(t, opts, "private_subnet_ids")
	natGatewayID      := terraform.Output(t, opts, "nat_gateway_id")
	publicInstanceID  := terraform.Output(t, opts, "public_instance_id")
	publicInstanceIP  := terraform.Output(t, opts, "public_instance_public_ip")
	privateInstanceID := terraform.Output(t, opts, "private_instance_id")

	// ═══════════════════════════════════════════════════════════
	// 1. VPC
	// ═══════════════════════════════════════════════════════════

	t.Run("TC01_VPC_CIDR", func(t *testing.T) {
		resp, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{
			VpcIds: aws.StringSlice([]string{vpcID}),
		})
		require.NoError(t, err)
		require.Len(t, resp.Vpcs, 1)
		assert.Equal(t, "10.0.0.0/16", aws.StringValue(resp.Vpcs[0].CidrBlock),
			"VPC CIDR phải đúng với giá trị khai báo")
	})

	t.Run("TC02_VPC_DNS", func(t *testing.T) {
		respHostnames, err := client.DescribeVpcAttribute(&ec2.DescribeVpcAttributeInput{
			VpcId:     aws.String(vpcID),
			Attribute: aws.String("enableDnsHostnames"),
		})
		require.NoError(t, err)
		assert.True(t, aws.BoolValue(respHostnames.EnableDnsHostnames.Value),
			"DNS hostnames phải được bật")

		respSupport, err := client.DescribeVpcAttribute(&ec2.DescribeVpcAttributeInput{
			VpcId:     aws.String(vpcID),
			Attribute: aws.String("enableDnsSupport"),
		})
		require.NoError(t, err)
		assert.True(t, aws.BoolValue(respSupport.EnableDnsSupport.Value),
			"DNS support phải được bật")
	})

	t.Run("TC03_PublicSubnet_Count", func(t *testing.T) {
		assert.Len(t, publicSubnetIDs, 2,
			"Số lượng Public Subnet phải bằng số CIDR đầu vào (2)")
	})

	t.Run("TC04_PrivateSubnet_Count", func(t *testing.T) {
		assert.Len(t, privateSubnetIDs, 2,
			"Số lượng Private Subnet phải bằng số CIDR đầu vào (2)")
	})

	t.Run("TC05_PublicSubnet_MapPublicIP", func(t *testing.T) {
		resp, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(publicSubnetIDs),
		})
		require.NoError(t, err)
		for _, subnet := range resp.Subnets {
			assert.True(t, aws.BoolValue(subnet.MapPublicIpOnLaunch),
				"Public Subnet %s phải có map_public_ip_on_launch = true",
				aws.StringValue(subnet.SubnetId))
		}
	})

	t.Run("TC06_PrivateSubnet_NoMapPublicIP", func(t *testing.T) {
		resp, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice(privateSubnetIDs),
		})
		require.NoError(t, err)
		for _, subnet := range resp.Subnets {
			assert.False(t, aws.BoolValue(subnet.MapPublicIpOnLaunch),
				"Private Subnet %s phải có map_public_ip_on_launch = false",
				aws.StringValue(subnet.SubnetId))
		}
	})

	t.Run("TC07_InternetGateway_AttachedToVPC", func(t *testing.T) {
		resp, err := client.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("attachment.vpc-id"), Values: aws.StringSlice([]string{vpcID})},
			},
		})
		require.NoError(t, err)
		assert.Len(t, resp.InternetGateways, 1,
			"Phải có đúng 1 Internet Gateway được gắn vào VPC")
		igwState := aws.StringValue(resp.InternetGateways[0].Attachments[0].State)
		assert.Equal(t, "available", igwState,
			"Internet Gateway phải ở trạng thái attached (available)")
	})

	t.Run("TC08_AZ_Count_GTE_Subnet_Count", func(t *testing.T) {
		azCount := 2 // khớp với availability_zones đầu vào
		assert.GreaterOrEqual(t, azCount, len(publicSubnetIDs),
			"Số AZ phải >= số Public Subnet để tránh index out of range")
	})

	// ═══════════════════════════════════════════════════════════
	// 2. Route Tables
	// ═══════════════════════════════════════════════════════════

	t.Run("TC09_PublicRouteTable_HasIGWRoute", func(t *testing.T) {
		resp, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("association.subnet-id"), Values: aws.StringSlice(publicSubnetIDs)},
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.RouteTables, "Phải tìm thấy Public Route Table")

		found := false
		for _, route := range resp.RouteTables[0].Routes {
			if aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" &&
				aws.StringValue(route.GatewayId) != "" {
				found = true
				break
			}
		}
		assert.True(t, found, "Public Route Table phải có route 0.0.0.0/0 → Internet Gateway")
	})

	t.Run("TC10_PrivateRouteTable_HasNATRoute", func(t *testing.T) {
		resp, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("association.subnet-id"), Values: aws.StringSlice(privateSubnetIDs)},
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.RouteTables, "Phải tìm thấy Private Route Table")

		found := false
		for _, route := range resp.RouteTables[0].Routes {
			if aws.StringValue(route.DestinationCidrBlock) == "0.0.0.0/0" &&
				aws.StringValue(route.NatGatewayId) == natGatewayID {
				found = true
				break
			}
		}
		assert.True(t, found, "Private Route Table phải có route 0.0.0.0/0 → NAT Gateway")
	})

	t.Run("TC11_PublicSubnets_AssociatedWith_PublicRT", func(t *testing.T) {
		resp, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("association.subnet-id"), Values: aws.StringSlice(publicSubnetIDs)},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.RouteTables, 1, "Tất cả Public Subnet phải dùng chung 1 Route Table")

		associatedSubnets := []string{}
		for _, assoc := range resp.RouteTables[0].Associations {
			if assoc.SubnetId != nil {
				associatedSubnets = append(associatedSubnets, aws.StringValue(assoc.SubnetId))
			}
		}
		for _, subnetID := range publicSubnetIDs {
			assert.Contains(t, associatedSubnets, subnetID,
				"Public Subnet %s phải được associate với Public Route Table", subnetID)
		}
	})

	t.Run("TC12_PrivateSubnets_AssociatedWith_PrivateRT", func(t *testing.T) {
		resp, err := client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("association.subnet-id"), Values: aws.StringSlice(privateSubnetIDs)},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.RouteTables, 1, "Tất cả Private Subnet phải dùng chung 1 Route Table")

		associatedSubnets := []string{}
		for _, assoc := range resp.RouteTables[0].Associations {
			if assoc.SubnetId != nil {
				associatedSubnets = append(associatedSubnets, aws.StringValue(assoc.SubnetId))
			}
		}
		for _, subnetID := range privateSubnetIDs {
			assert.Contains(t, associatedSubnets, subnetID,
				"Private Subnet %s phải được associate với Private Route Table", subnetID)
		}
	})

	// ═══════════════════════════════════════════════════════════
	// 3. NAT Gateway
	// ═══════════════════════════════════════════════════════════

	t.Run("TC13_EIP_Domain_VPC", func(t *testing.T) {
		resp, err := client.DescribeAddresses(&ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("domain"), Values: aws.StringSlice([]string{"vpc"})},
				{Name: aws.String("tag:Project"), Values: aws.StringSlice([]string{"nt548-test"})},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Addresses, "Phải có EIP được tạo trong domain vpc")
	})

	t.Run("TC14_NATGateway_InPublicSubnet", func(t *testing.T) {
		resp, err := client.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			NatGatewayIds: aws.StringSlice([]string{natGatewayID}),
		})
		require.NoError(t, err)
		require.Len(t, resp.NatGateways, 1)

		natSubnetID := aws.StringValue(resp.NatGateways[0].SubnetId)
		assert.Contains(t, publicSubnetIDs, natSubnetID,
			"NAT Gateway phải được đặt trong Public Subnet, không phải Private Subnet")
	})

	t.Run("TC15_NATGateway_HasCorrectEIP", func(t *testing.T) {
		resp, err := client.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
			NatGatewayIds: aws.StringSlice([]string{natGatewayID}),
		})
		require.NoError(t, err)
		require.Len(t, resp.NatGateways, 1)
		require.NotEmpty(t, resp.NatGateways[0].NatGatewayAddresses)

		allocationID := aws.StringValue(resp.NatGateways[0].NatGatewayAddresses[0].AllocationId)
		assert.NotEmpty(t, allocationID,
			"NAT Gateway phải được gắn đúng Elastic IP (AllocationId không được rỗng)")
	})

	// ═══════════════════════════════════════════════════════════
	// 4. EC2
	// ═══════════════════════════════════════════════════════════

	t.Run("TC16_PublicInstance_InPublicSubnet", func(t *testing.T) {
		resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{publicInstanceID}),
		})
		require.NoError(t, err)
		require.Len(t, resp.Reservations, 1)
		instance := resp.Reservations[0].Instances[0]

		assert.Contains(t, publicSubnetIDs, aws.StringValue(instance.SubnetId),
			"Public Instance phải nằm trong Public Subnet")
	})

	t.Run("TC17_PublicSubnet_MapPublicIP_ForPublicInstance", func(t *testing.T) {
		resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{publicInstanceID}),
		})
		require.NoError(t, err)
		instance := resp.Reservations[0].Instances[0]
		subnetID := aws.StringValue(instance.SubnetId)

		subnetResp, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			SubnetIds: aws.StringSlice([]string{subnetID}),
		})
		require.NoError(t, err)
		assert.True(t, aws.BoolValue(subnetResp.Subnets[0].MapPublicIpOnLaunch),
			"Subnet của Public Instance phải có map_public_ip_on_launch = true")

		assert.NotEmpty(t, publicInstanceIP,
			"Public Instance phải có public IP")
	})

	t.Run("TC18_PrivateInstance_InPrivateSubnet", func(t *testing.T) {
		resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{privateInstanceID}),
		})
		require.NoError(t, err)
		require.Len(t, resp.Reservations, 1)
		instance := resp.Reservations[0].Instances[0]

		assert.Contains(t, privateSubnetIDs, aws.StringValue(instance.SubnetId),
			"Private Instance phải nằm trong Private Subnet")
	})

	t.Run("TC19_PrivateInstance_NoPublicIP", func(t *testing.T) {
		resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{privateInstanceID}),
		})
		require.NoError(t, err)
		instance := resp.Reservations[0].Instances[0]

		assert.Empty(t, aws.StringValue(instance.PublicIpAddress),
			"Private Instance không được có public IP")
	})

	t.Run("TC20_BothInstances_RootVolume_Encrypted", func(t *testing.T) {
		for _, instanceID := range []string{publicInstanceID, privateInstanceID} {
			resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice([]string{instanceID}),
			})
			require.NoError(t, err)
			instance := resp.Reservations[0].Instances[0]

			// Lấy root volume từ block device mappings
			require.NotEmpty(t, instance.BlockDeviceMappings,
				"Instance %s phải có block device mapping", instanceID)

			rootVolumeID := aws.StringValue(instance.BlockDeviceMappings[0].Ebs.VolumeId)
			volResp, err := client.DescribeVolumes(&ec2.DescribeVolumesInput{
				VolumeIds: aws.StringSlice([]string{rootVolumeID}),
			})
			require.NoError(t, err)
			require.Len(t, volResp.Volumes, 1)

			assert.True(t, aws.BoolValue(volResp.Volumes[0].Encrypted),
				"Root volume của instance %s phải được mã hóa (encrypted = true)", instanceID)
		}
	})

	t.Run("TC21_BothInstances_IMDSv2_Required", func(t *testing.T) {
		for _, instanceID := range []string{publicInstanceID, privateInstanceID} {
			resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice([]string{instanceID}),
			})
			require.NoError(t, err)
			instance := resp.Reservations[0].Instances[0]

			require.NotNil(t, instance.MetadataOptions,
				"Instance %s phải có MetadataOptions", instanceID)
			assert.Equal(t, "required",
				aws.StringValue(instance.MetadataOptions.HttpTokens),
				"Instance %s phải dùng IMDSv2 (http_tokens = required)", instanceID)
		}
	})

	// ═══════════════════════════════════════════════════════════
	// 5. Security Groups
	// ═══════════════════════════════════════════════════════════

	// Lấy SG IDs từ instance thực tế
	getInstanceSGIDs := func(t *testing.T, instanceID string) []string {
		resp, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instanceID}),
		})
		require.NoError(t, err)
		sgIDs := []string{}
		for _, sg := range resp.Reservations[0].Instances[0].SecurityGroups {
			sgIDs = append(sgIDs, aws.StringValue(sg.GroupId))
		}
		return sgIDs
	}

	publicSGIDs  := getInstanceSGIDs(t, publicInstanceID)
	privateSGIDs := getInstanceSGIDs(t, privateInstanceID)
	require.NotEmpty(t, publicSGIDs,  "Public Instance phải có Security Group")
	require.NotEmpty(t, privateSGIDs, "Private Instance phải có Security Group")
	publicSGID  := publicSGIDs[0]
	privateSGID := privateSGIDs[0]

	t.Run("TC22_PublicSG_SSH_From_Anywhere", func(t *testing.T) {
		resp, err := client.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("group-id"), Values: aws.StringSlice([]string{publicSGID})},
			},
		})
		require.NoError(t, err)

		found := false
		for _, rule := range resp.SecurityGroupRules {
			if aws.BoolValue(rule.IsEgress) {
				continue
			}
			if aws.Int64Value(rule.FromPort) == 22 &&
				aws.Int64Value(rule.ToPort) == 22 &&
				aws.StringValue(rule.IpProtocol) == "tcp" &&
				aws.StringValue(rule.CidrIpv4) == "0.0.0.0/0" {
				found = true
				break
			}
		}
		assert.True(t, found,
			"Public SG phải có ingress rule TCP port 22 từ 0.0.0.0/0")
	})

	t.Run("TC23_PublicSG_AllowAllEgress", func(t *testing.T) {
		resp, err := client.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("group-id"), Values: aws.StringSlice([]string{publicSGID})},
			},
		})
		require.NoError(t, err)

		found := false
		for _, rule := range resp.SecurityGroupRules {
			if aws.BoolValue(rule.IsEgress) &&
				aws.StringValue(rule.IpProtocol) == "-1" &&
				aws.StringValue(rule.CidrIpv4) == "0.0.0.0/0" {
				found = true
				break
			}
		}
		assert.True(t, found,
			"Public SG phải có egress rule cho phép toàn bộ traffic ra ngoài")
	})

	t.Run("TC24_PrivateSG_SSH_OnlyFrom_PublicSG", func(t *testing.T) {
		resp, err := client.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("group-id"), Values: aws.StringSlice([]string{privateSGID})},
			},
		})
		require.NoError(t, err)

		found := false
		for _, rule := range resp.SecurityGroupRules {
			if aws.BoolValue(rule.IsEgress) {
				continue
			}
			if aws.Int64Value(rule.FromPort) == 22 &&
				aws.Int64Value(rule.ToPort) == 22 &&
				aws.StringValue(rule.IpProtocol) == "tcp" &&
				rule.ReferencedGroupInfo != nil &&
				aws.StringValue(rule.ReferencedGroupInfo.GroupId) == publicSGID {
				found = true
				break
			}
		}
		assert.True(t, found,
			"Private SG phải có ingress SSH chỉ từ Public SG (referenced_security_group_id)")
	})

	t.Run("TC25_PrivateSG_AllowAllEgress", func(t *testing.T) {
		resp, err := client.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("group-id"), Values: aws.StringSlice([]string{privateSGID})},
			},
		})
		require.NoError(t, err)

		found := false
		for _, rule := range resp.SecurityGroupRules {
			if aws.BoolValue(rule.IsEgress) &&
				aws.StringValue(rule.IpProtocol) == "-1" &&
				aws.StringValue(rule.CidrIpv4) == "0.0.0.0/0" {
				found = true
				break
			}
		}
		assert.True(t, found,
			"Private SG phải có egress rule cho phép toàn bộ traffic ra ngoài")
	})

	t.Run("TC26_PrivateSG_NoIngress_FromAnywhere", func(t *testing.T) {
		resp, err := client.DescribeSecurityGroupRules(&ec2.DescribeSecurityGroupRulesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("group-id"), Values: aws.StringSlice([]string{privateSGID})},
			},
		})
		require.NoError(t, err)

		for _, rule := range resp.SecurityGroupRules {
			if !aws.BoolValue(rule.IsEgress) {
				assert.NotEqual(t, "0.0.0.0/0", aws.StringValue(rule.CidrIpv4),
					"Private SG không được có ingress rule nào từ 0.0.0.0/0")
			}
		}
	})
}