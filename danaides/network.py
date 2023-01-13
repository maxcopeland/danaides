from constructs import Construct
from aws_cdk import (
    aws_ec2 as ec2,
    CfnOutput
)


class Network(Construct):

    @property
    def s3_endpoint(self):
        return self._s3_endpoint

    @property
    def sqs_endpoint(self):
        return self._sqs_endpoint

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        vpc = ec2.Vpc(self, "DanaidesVPC",
                      ip_addresses=ec2.IpAddresses.cidr("10.0.0.0/16"),
                      subnet_configuration=[
                          {
                              "cidrMask": 24,
                              "name": "Public",
                              "subnetType": ec2.SubnetType.PUBLIC
                          },
                          {
                              "cidrMask": 24,
                              "name": "Private",
                              "subnetType": ec2.SubnetType.PRIVATE_WITH_EGRESS
                          },
                        ],
                      )
        self._s3_endpoint = vpc.add_gateway_endpoint(id="S3GatewayEndpoint",
                                 service=ec2.GatewayVpcEndpointAwsService.S3)

        self._sqs_endpoint = vpc.add_interface_endpoint(id="SQSInterfaceEndpoint",
                                   service=ec2.InterfaceVpcEndpointAwsService.SQS)

        vpc.add_interface_endpoint(id="SSMInterfaceEndpoint",
                                   service=ec2.InterfaceVpcEndpointAwsService.SSM)

        vpc.add_interface_endpoint(id="SSMMessagesInterfaceEndpoint",
                                   service=ec2.InterfaceVpcEndpointAwsService.SSM_MESSAGES)

        vpc.add_interface_endpoint(id="EC2MessagesInterfaceEndpoint",
                                   service=ec2.InterfaceVpcEndpointAwsService.EC2_MESSAGES)

        bastion_host = ec2.BastionHostLinux(self, "BastionHost",
                                            vpc=vpc,
                                            instance_type=ec2.InstanceType.of(
                                                ec2.InstanceClass.T2,
                                                ec2.InstanceSize.MICRO,
                                                ),
                                            )

        private_server = ec2.Instance(self, "PrivateServer",
                                      vpc=vpc,
                                      instance_name="PrivateServer",
                                      instance_type=ec2.InstanceType.of(
                                          ec2.InstanceClass.T2,
                                          ec2.InstanceSize.MICRO,
                                          ),
                                      machine_image=ec2.MachineImage.latest_amazon_linux(),
                                      )

        CfnOutput(self, id="PrivateServerId", value=private_server.instance_id)
