from aws_cdk import (
    Stack,
    aws_s3
)
from constructs import Construct

from .network import Network


class DanaidesStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:

        super().__init__(scope, construct_id, **kwargs)

        my_network = Network(self, "DanaidesNetwork")

        my_bucket = aws_s3.Bucket(self, "TestBucket")
