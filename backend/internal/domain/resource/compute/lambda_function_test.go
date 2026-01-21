package compute

import (
	"testing"
)

func TestLambdaFunction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		function *LambdaFunction
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-s3-function",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
			},
			wantErr: false,
		},
		{
			name: "valid-container-image-function",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				PackageType:  stringPtr("Image"),
				ImageURI:     stringPtr("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest"),
			},
			wantErr: false,
		},
		{
			name: "missing-function-name",
			function: &LambdaFunction{
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
				Region:  "us-east-1",
				S3Bucket: stringPtr("my-bucket"),
				S3Key:    stringPtr("code.zip"),
				Runtime:  stringPtr("python3.9"),
				Handler:  stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "function_name is required",
		},
		{
			name: "missing-role-arn",
			function: &LambdaFunction{
				FunctionName: "test-function",
				Region:      "us-east-1",
				S3Bucket:    stringPtr("my-bucket"),
				S3Key:       stringPtr("code.zip"),
				Runtime:     stringPtr("python3.9"),
				Handler:     stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "role_arn is required",
		},
		{
			name: "missing-region",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:     "arn:aws:iam::123456789012:role/test-role",
				S3Bucket:    stringPtr("my-bucket"),
				S3Key:       stringPtr("code.zip"),
				Runtime:     stringPtr("python3.9"),
				Handler:     stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "region is required",
		},
		{
			name: "no-code-source",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
			},
			wantErr: true,
			errMsg:  "either S3 code",
		},
		{
			name: "both-code-sources",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				PackageType:  stringPtr("Image"),
				ImageURI:     stringPtr("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest"),
			},
			wantErr: true,
			errMsg:  "cannot specify both",
		},
		{
			name: "s3-without-runtime",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Handler:      stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "runtime is required",
		},
		{
			name: "container-with-runtime",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				PackageType:  stringPtr("Image"),
				ImageURI:     stringPtr("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest"),
				Runtime:      stringPtr("python3.9"),
			},
			wantErr: true,
			errMsg:  "runtime should not be set",
		},
		{
			name: "invalid-memory-size-too-small",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				MemorySize:   intPtr(64),
			},
			wantErr: true,
			errMsg:  "memory_size must be between",
		},
		{
			name: "invalid-memory-size-not-multiple-of-64",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				MemorySize:   intPtr(200),
			},
			wantErr: true,
			errMsg:  "memory_size must be a multiple of 64",
		},
		{
			name: "invalid-timeout-too-small",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				Timeout:      intPtr(0),
			},
			wantErr: true,
			errMsg:  "timeout must be between",
		},
		{
			name: "invalid-function-name-too-long",
			function: &LambdaFunction{
				FunctionName: "this-is-a-very-long-function-name-that-exceeds-the-maximum-length-of-sixty-four-characters",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "function name must be at most 64 characters",
		},
		{
			name: "invalid-role-arn-format",
			function: &LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "invalid-arn",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
			},
			wantErr: true,
			errMsg:  "role ARN must be in format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.function.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errMsg)
				} else if err.Error() != "" && tt.errMsg != "" {
					// Check if error message contains expected substring
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLambdaFunction_HasS3Code(t *testing.T) {
	tests := []struct {
		name     string
		function *LambdaFunction
		want     bool
	}{
		{
			name: "has-s3-code",
			function: &LambdaFunction{
				S3Bucket: stringPtr("my-bucket"),
				S3Key:    stringPtr("code.zip"),
			},
			want: true,
		},
		{
			name: "no-s3-code",
			function: &LambdaFunction{
				PackageType: stringPtr("Image"),
				ImageURI:    stringPtr("image-uri"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function.HasS3Code()
			if got != tt.want {
				t.Errorf("HasS3Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambdaFunction_HasContainerImage(t *testing.T) {
	tests := []struct {
		name     string
		function *LambdaFunction
		want     bool
	}{
		{
			name: "has-container-image",
			function: &LambdaFunction{
				PackageType: stringPtr("Image"),
				ImageURI:    stringPtr("image-uri"),
			},
			want: true,
		},
		{
			name: "no-container-image",
			function: &LambdaFunction{
				S3Bucket: stringPtr("my-bucket"),
				S3Key:    stringPtr("code.zip"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function.HasContainerImage()
			if got != tt.want {
				t.Errorf("HasContainerImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
