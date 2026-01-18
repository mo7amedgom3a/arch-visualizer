package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// mockAWSStorageService is a mock implementation of AWSStorageService for testing
type mockAWSStorageService struct {
	volume      *awsebs.Volume
	createError error
	getError    error
}

// Ensure mockAWSStorageService implements AWSStorageService
var _ awsservice.AWSStorageService = (*mockAWSStorageService)(nil)

// Helper function to convert Volume input to output
func volumeToOutput(volume *awsebs.Volume) *awsebsoutputs.VolumeOutput {
	if volume == nil {
		return nil
	}
	return &awsebsoutputs.VolumeOutput{
		ID:               "vol-mock-1234567890abcdef0",
		ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-mock-1234567890abcdef0",
		Name:             volume.Name,
		AvailabilityZone: volume.AvailabilityZone,
		Size:             volume.Size,
		VolumeType:       volume.VolumeType,
		IOPS:             volume.IOPS,
		Throughput:       volume.Throughput,
		Encrypted:        volume.Encrypted,
		KMSKeyID:         volume.KMSKeyID,
		SnapshotID:       volume.SnapshotID,
		State:            "available",
		AttachedTo:       nil,
		CreateTime:       time.Now(),
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: volume.Name},
		},
	}
}

// EBS Volume Operations

func (m *mockAWSStorageService) CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.volume = volume
	return volumeToOutput(volume), nil
}

func (m *mockAWSStorageService) GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return volumeToOutput(m.volume), nil
}

func (m *mockAWSStorageService) UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	m.volume = volume
	return volumeToOutput(volume), nil
}

func (m *mockAWSStorageService) DeleteEBSVolume(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSStorageService) ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	if m.volume != nil {
		return []*awsebsoutputs.VolumeOutput{volumeToOutput(m.volume)}, nil
	}
	return []*awsebsoutputs.VolumeOutput{}, nil
}

// Volume Attachment Operations

func (m *mockAWSStorageService) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	return nil
}

func (m *mockAWSStorageService) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	return nil
}

// Tests

func TestAWSStorageAdapter_CreateEBSVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	domainVolume := &domainstorage.EBSVolume{
		Name:             "test-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		Type:             "gp3",
		Encrypted:        false,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	if createdVolume.Name != domainVolume.Name {
		t.Errorf("Expected name %s, got %s", domainVolume.Name, createdVolume.Name)
	}

	if createdVolume.ID == "" {
		t.Error("Expected volume ID to be populated")
	}
}

func TestAWSStorageAdapter_CreateEBSVolume_ValidationError(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	invalidVolume := &domainstorage.EBSVolume{
		Name:   "", // Invalid: empty name
		Region: "us-east-1",
		Size:   40,
		Type:   "gp3",
	}

	ctx := context.Background()
	_, err := adapter.CreateEBSVolume(ctx, invalidVolume)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSStorageAdapter_GetEBSVolume(t *testing.T) {
	mockService := &mockAWSStorageService{
		volume: &awsebs.Volume{
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			Encrypted:        false,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	volume, err := adapter.GetEBSVolume(ctx, "vol-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if volume == nil {
		t.Fatal("Expected volume, got nil")
	}

	if volume.Name != "test-volume" {
		t.Errorf("Expected name test-volume, got %s", volume.Name)
	}
}

func TestAWSStorageAdapter_ListEBSVolumes(t *testing.T) {
	mockService := &mockAWSStorageService{
		volume: &awsebs.Volume{
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			Encrypted:        false,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	volumes, err := adapter.ListEBSVolumes(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(volumes))
	}

	if volumes[0].Name != "test-volume" {
		t.Errorf("Expected name test-volume, got %s", volumes[0].Name)
	}
}

func TestAWSStorageAdapter_AttachVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	err := adapter.AttachVolume(ctx, "vol-123", "i-123", "/dev/xvdf")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSStorageAdapter_DetachVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	err := adapter.DetachVolume(ctx, "vol-123", "i-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSStorageAdapter_ErrorHandling(t *testing.T) {
	mockService := &mockAWSStorageService{
		getError: errors.New("aws service error"),
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetEBSVolume(ctx, "vol-123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}
