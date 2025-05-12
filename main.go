package main

import (
    "fmt"
    "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {

    // Create a VPC
    vpc, error := aws.ec2.NewVpc(ctx, "my-vpc", &aws.ec2.NewVpcArgs{
        CidrBlock: pulumi.String("10.0.0.0/16"),
    })
    if error != nil {
        return fmt.Errorf("failed to create VPC: %v", error)}
    }

    // Create a public subnet
    subnet, error := aws.ec2.NewSubnet(ctx, "my-subnet", &aws.ec2.NewSubnetArgs{
        VpcId:               vpc.ID(),
        CidrBlock:           pulumi.String("10.0.1.0/24"),
        AvailabilityZone:    pulumi.String("us-east-1a"),
        MapPublicIpOnLaunch: pulumi.Bool(true),
    })
    if error != nil {
        return fmt.Errorf("failed to create subnet: %v", error)
    }

    // create an Internet Gateway
    internet_gateway, error := aws.ec2.NewInternetGateway(ctx, "my-internet-gateway", &aws.ec2.NewInternetGatewayArgs{
        VpcId: vpc.ID(),
    })
    if error != nil {
        return fmt.Errorf("failed to create internet gateway: %v", error)
    }

    // create security group allowing SSH (port 22) and HTTP (port 80) access
    sec_group, error := aws.ec2.NewSecurityGroup(ctx, "my-sec-group", &aws.ec2.NewSecurityGroupArgs{
        VpcId: vpc.ID(),
        Ingress: []aws.ec2.SecurityGroupIngressArgs{
            {
                FromPort:   pulumi.Int(22),
                ToPort:     pulumi.Int(22),
                Protocol:   pulumi.String("tcp"),
                CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
            },
            {
                FromPort:   pulumi.Int(80),
                ToPort:     pulumi.Int(80),
                Protocol:   pulumi.String("tcp"),
                CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
            },
        },
    })
    if error != nil {
        return fmt.Errorf("failed to create security group: %v", error)
    }

    // Export VPC ID and the subnet ID for future reference
    ctx.Export("vpc_id", vpc.ID())
    ctx.Export("subnet_id", subnet.ID())
    ctx.Export("internet_gateway_id", internet_gateway.ID())
    ctx.Export("sec_group_id", sec_group.ID())

    return nil
}
