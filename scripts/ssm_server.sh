export AWS_REGION=us-east-1
export BASTION_INSTANCE_ID=$(aws ec2 describe-instances \
                             --region $AWS_REGION \
                             --filter "Name=tag:Name,Values=PrivateServer" \
                             --query "Reservations[].Instances[?State.Name == 'running'].InstanceId[]" \
                             --output text \
                             --profile dev)
aws ssm start-session --target $BASTION_INSTANCE_ID --region=$AWS_REGION --profile dev