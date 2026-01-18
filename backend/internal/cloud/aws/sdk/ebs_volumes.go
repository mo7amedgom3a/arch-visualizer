package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
)

// CreateEBSVolume creates a new EBS volume using AWS SDK
func CreateEBSVolume(ctx context.Context, client *AWSClient, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("volume validation failed: %w", err)
	}

	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(volume.AvailabilityZone),
		Size:             aws.Int32(int32(volume.Size)),
		VolumeType:       types.VolumeType(volume.VolumeType),
		Encrypted:        aws.Bool(volume.Encrypted),
	}

	// IOPS
	if volume.IOPS != nil {
		input.Iops = aws.Int32(int32(*volume.IOPS))
	}

	// Throughput (gp3 only)
	if volume.Throughput != nil {
		input.Throughput = aws.Int32(int32(*volume.Throughput))
	}

	// KMS Key ID
	if volume.KMSKeyID != nil {
		input.KmsKeyId = volume.KMSKeyID
	}

	// Snapshot ID
	if volume.SnapshotID != nil {
		input.SnapshotId = volume.SnapshotID
	}

	// Tags
	if len(volume.Tags) > 0 {
		var tagSpecs []types.TagSpecification
		var tags []types.Tag
		for _, tag := range volume.Tags {
			tags = append(tags, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		tagSpecs = append(tagSpecs, types.TagSpecification{
			ResourceType: types.ResourceTypeVolume,
			Tags:         tags,
		})
		input.TagSpecifications = tagSpecs
	}

	// Create the volume
	result, err := client.EC2.CreateVolume(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create EBS volume: %w", err)
	}

	if result == nil || result.VolumeId == nil {
		return nil, fmt.Errorf("volume creation returned nil")
	}

	// Get the created volume to return full details
	return GetEBSVolume(ctx, client, *result.VolumeId)
}

// GetEBSVolume retrieves an EBS volume by ID
func GetEBSVolume(ctx context.Context, client *AWSClient, id string) (*awsebsoutputs.VolumeOutput, error) {
	input := &ec2.DescribeVolumesInput{
		VolumeIds: []string{id},
	}

	output, err := client.EC2.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe EBS volume: %w", err)
	}

	if len(output.Volumes) == 0 {
		return nil, fmt.Errorf("EBS volume not found: %s", id)
	}

	return convertVolumeToOutput(&output.Volumes[0]), nil
}

// UpdateEBSVolume modifies an existing EBS volume (size, type, IOPS, throughput)
func UpdateEBSVolume(ctx context.Context, client *AWSClient, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("volume validation failed: %w", err)
	}

	// Modify volume attributes
	input := &ec2.ModifyVolumeInput{
		VolumeId: aws.String(id),
	}

	// Size (can only increase)
	if volume.Size > 0 {
		input.Size = aws.Int32(int32(volume.Size))
	}

	// Volume type
	if volume.VolumeType != "" {
		input.VolumeType = types.VolumeType(volume.VolumeType)
	}

	// IOPS
	if volume.IOPS != nil {
		input.Iops = aws.Int32(int32(*volume.IOPS))
	}

	// Throughput (gp3 only)
	if volume.Throughput != nil {
		input.Throughput = aws.Int32(int32(*volume.Throughput))
	}

	_, err := client.EC2.ModifyVolume(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to modify EBS volume: %w", err)
	}

	// Get the updated volume
	return GetEBSVolume(ctx, client, id)
}

// DeleteEBSVolume deletes an EBS volume
func DeleteEBSVolume(ctx context.Context, client *AWSClient, id string) error {
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	}

	_, err := client.EC2.DeleteVolume(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete EBS volume: %w", err)
	}

	return nil
}

// ListEBSVolumes lists EBS volumes with optional filters
func ListEBSVolumes(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	input := &ec2.DescribeVolumesInput{}

	// Convert filters to AWS filter format
	if len(filters) > 0 {
		var awsFilters []types.Filter
		for key, values := range filters {
			awsFilters = append(awsFilters, types.Filter{
				Name:   aws.String(key),
				Values: values,
			})
		}
		input.Filters = awsFilters
	}

	output, err := client.EC2.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe EBS volumes: %w", err)
	}

	volumes := make([]*awsebsoutputs.VolumeOutput, len(output.Volumes))
	for i, vol := range output.Volumes {
		volumes[i] = convertVolumeToOutput(&vol)
	}

	return volumes, nil
}

// AttachVolume attaches an EBS volume to an EC2 instance
func AttachVolume(ctx context.Context, client *AWSClient, volumeID, instanceID, deviceName string) error {
	input := &ec2.AttachVolumeInput{
		VolumeId:   aws.String(volumeID),
		InstanceId: aws.String(instanceID),
		Device:     aws.String(deviceName),
	}

	_, err := client.EC2.AttachVolume(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to attach volume: %w", err)
	}

	return nil
}

// DetachVolume detaches an EBS volume from an EC2 instance
func DetachVolume(ctx context.Context, client *AWSClient, volumeID, instanceID string) error {
	input := &ec2.DetachVolumeInput{
		VolumeId:   aws.String(volumeID),
		InstanceId: aws.String(instanceID),
	}

	_, err := client.EC2.DetachVolume(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach volume: %w", err)
	}

	return nil
}

// convertVolumeToOutput converts AWS SDK Volume to output model
func convertVolumeToOutput(vol *types.Volume) *awsebsoutputs.VolumeOutput {
	// Construct ARN manually since Volume doesn't have VolumeArn field
	// Format: arn:aws:ec2:region:account-id:volume/vol-id
	arn := ""
	if vol.VolumeId != nil {
		arn = fmt.Sprintf("arn:aws:ec2:region:account:volume/%s", *vol.VolumeId)
	}

	output := &awsebsoutputs.VolumeOutput{
		ID:               aws.ToString(vol.VolumeId),
		ARN:              arn,
		Name:             "", // Name from tags
		AvailabilityZone: aws.ToString(vol.AvailabilityZone),
		Size:             int(aws.ToInt32(vol.Size)),
		VolumeType:       string(vol.VolumeType),
		State:            string(vol.State),
		Encrypted:        aws.ToBool(vol.Encrypted),
		CreateTime:       aws.ToTime(vol.CreateTime),
	}

	// IOPS
	if vol.Iops != nil {
		iops := int(*vol.Iops)
		output.IOPS = &iops
	}

	// Throughput
	if vol.Throughput != nil {
		throughput := int(*vol.Throughput)
		output.Throughput = &throughput
	}

	// KMS Key ID
	if vol.KmsKeyId != nil {
		output.KMSKeyID = vol.KmsKeyId
	}

	// Snapshot ID
	if vol.SnapshotId != nil {
		output.SnapshotID = vol.SnapshotId
	}

	// Attached instance
	if len(vol.Attachments) > 0 {
		attachedTo := vol.Attachments[0].InstanceId
		output.AttachedTo = attachedTo
	}

	// Tags
	if vol.Tags != nil {
		output.Tags = make([]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}, len(vol.Tags))
		for i, tag := range vol.Tags {
			output.Tags[i] = struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			}
			// Extract Name from tags
			if aws.ToString(tag.Key) == "Name" {
				output.Name = aws.ToString(tag.Value)
			}
		}
	}

	return output
}
