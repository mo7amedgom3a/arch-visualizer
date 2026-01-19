package ebs

import "testing"

func validVolume() *Volume {
	return &Volume{
		Name:             "test-volume",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		VolumeType:       "gp3",
	}
}

func TestVolumeValidate_Success(t *testing.T) {
	v := validVolume()
	if err := v.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVolumeValidate_RequiredFields(t *testing.T) {
	v := validVolume()
	v.Name = ""
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}

	v = validVolume()
	v.AvailabilityZone = ""
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for missing availability zone")
	}
}

func TestVolumeValidate_AvailabilityZoneFormat(t *testing.T) {
	v := validVolume()
	v.AvailabilityZone = "useast1"
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for invalid availability zone format")
	}
}

func TestVolumeValidate_SizeValidation(t *testing.T) {
	v := validVolume()
	v.Size = 0
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for size <= 0")
	}

	v = validVolume()
	v.Size = 20000
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for size too large")
	}
}

func TestVolumeValidate_VolumeTypeValidation(t *testing.T) {
	v := validVolume()
	v.VolumeType = "bad"
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for invalid volume type")
	}
}

func TestVolumeValidate_IOPSValidation(t *testing.T) {
	v := validVolume()
	iops := 2000
	v.IOPS = &iops
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for gp3 IOPS out of range")
	}

	v = validVolume()
	v.VolumeType = "st1"
	iops = 100
	v.IOPS = &iops
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for IOPS on unsupported volume type")
	}
}

func TestVolumeValidate_ThroughputValidation(t *testing.T) {
	v := validVolume()
	throughput := 100
	v.Throughput = &throughput
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for gp3 throughput out of range")
	}

	v = validVolume()
	v.VolumeType = "io1"
	throughput = 200
	v.Throughput = &throughput
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for throughput on non-gp3")
	}
}

func TestVolumeValidate_KMSKeyIDValidation(t *testing.T) {
	v := validVolume()
	key := "bad-key"
	v.KMSKeyID = &key
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for invalid kms key id")
	}
}

func TestVolumeValidate_SnapshotIDValidation(t *testing.T) {
	v := validVolume()
	snap := "snapshot-123"
	v.SnapshotID = &snap
	if err := v.Validate(); err == nil {
		t.Fatal("expected error for invalid snapshot id")
	}
}
