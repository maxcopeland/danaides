import aws_cdk as core
import aws_cdk.assertions as assertions

from danaides.danaides_stack import DanaidesStack

# example tests. To run these tests, uncomment this file along with the example
# resource in danaides_tmp/danaides_stack.py
def test_sqs_queue_created():
    app = core.App()
    stack = DanaidesStack(app, "danaides_tmp")
    template = assertions.Template.from_stack(stack)

#     template.has_resource_properties("AWS::SQS::Queue", {
#         "VisibilityTimeout": 300
#     })
