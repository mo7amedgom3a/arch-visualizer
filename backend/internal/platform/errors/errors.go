package errors

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
	"gorm.io/gorm"
)

// Platform-specific error codes
const (
	// Database errors
	CodeDatabaseConnectionFailed  = "DATABASE_CONNECTION_FAILED"
	CodeDatabaseQueryFailed       = "DATABASE_QUERY_FAILED"
	CodeDatabaseTransactionFailed = "DATABASE_TRANSACTION_FAILED"
	CodeDatabaseMigrationFailed   = "DATABASE_MIGRATION_FAILED"
	CodeDatabaseConfigError       = "DATABASE_CONFIG_ERROR"

	// Repository errors
	CodeRepositoryNotFound       = "REPOSITORY_NOT_FOUND"
	CodeRepositoryCreateFailed   = "REPOSITORY_CREATE_FAILED"
	CodeRepositoryUpdateFailed   = "REPOSITORY_UPDATE_FAILED"
	CodeRepositoryDeleteFailed   = "REPOSITORY_DELETE_FAILED"
	CodeRepositoryQueryFailed    = "REPOSITORY_QUERY_FAILED"
	CodeRepositoryDuplicateEntry = "REPOSITORY_DUPLICATE_ENTRY"

	// Resource repository errors
	CodeResourceNotFound     = "RESOURCE_NOT_FOUND"
	CodeResourceCreateFailed = "RESOURCE_CREATE_FAILED"
	CodeResourceUpdateFailed = "RESOURCE_UPDATE_FAILED"
	CodeResourceDeleteFailed = "RESOURCE_DELETE_FAILED"
	CodeResourceInvalidID    = "RESOURCE_INVALID_ID"

	// Project repository errors
	CodeProjectNotFound     = "PROJECT_NOT_FOUND"
	CodeProjectCreateFailed = "PROJECT_CREATE_FAILED"
	CodeProjectUpdateFailed = "PROJECT_UPDATE_FAILED"
	CodeProjectDeleteFailed = "PROJECT_DELETE_FAILED"

	// User repository errors
	CodeUserNotFound       = "USER_NOT_FOUND"
	CodeUserCreateFailed   = "USER_CREATE_FAILED"
	CodeUserUpdateFailed   = "USER_UPDATE_FAILED"
	CodeUserDeleteFailed   = "USER_DELETE_FAILED"
	CodeUserDuplicateEmail = "USER_DUPLICATE_EMAIL"

	// Pricing repository errors
	CodePricingNotFound     = "PRICING_NOT_FOUND"
	CodePricingCreateFailed = "PRICING_CREATE_FAILED"
	CodePricingUpdateFailed = "PRICING_UPDATE_FAILED"

	// Auth errors
	CodeAuthUnauthorized       = "AUTH_UNAUTHORIZED"
	CodeAuthForbidden          = "AUTH_FORBIDDEN"
	CodeAuthTokenInvalid       = "AUTH_TOKEN_INVALID"
	CodeAuthTokenExpired       = "AUTH_TOKEN_EXPIRED"
	CodeAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
)

// NewDatabaseConnectionFailed creates an error for database connection failures
func NewDatabaseConnectionFailed(cause error) *errors.AppError {
	return errors.Wrap(cause, CodeDatabaseConnectionFailed, errors.KindInternal, "Failed to connect to database")
}

// NewDatabaseQueryFailed creates an error for database query failures
func NewDatabaseQueryFailed(operation string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeDatabaseQueryFailed, errors.KindInternal, "Database query failed").
		WithOp(operation)
}

// NewDatabaseTransactionFailed creates an error for transaction failures
func NewDatabaseTransactionFailed(operation string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeDatabaseTransactionFailed, errors.KindInternal, "Database transaction failed").
		WithOp(operation)
}

// NewDatabaseConfigError creates an error for database configuration issues
func NewDatabaseConfigError(reason string) *errors.AppError {
	return errors.New(CodeDatabaseConfigError, errors.KindInternal, "Database configuration error").
		WithMeta("reason", reason)
}

// NewRepositoryNotFound creates an error for when a repository record is not found
func NewRepositoryNotFound(resourceType string, id interface{}) *errors.AppError {
	return errors.New(CodeRepositoryNotFound, errors.KindNotFound, "Repository record not found").
		WithMeta("resource_type", resourceType).
		WithMeta("id", id)
}

// NewRepositoryCreateFailed creates an error for repository create failures
func NewRepositoryCreateFailed(resourceType string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeRepositoryCreateFailed, errors.KindInternal, "Failed to create repository record").
		WithMeta("resource_type", resourceType)
}

// NewRepositoryUpdateFailed creates an error for repository update failures
func NewRepositoryUpdateFailed(resourceType string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeRepositoryUpdateFailed, errors.KindInternal, "Failed to update repository record").
		WithMeta("resource_type", resourceType)
}

// NewRepositoryDeleteFailed creates an error for repository delete failures
func NewRepositoryDeleteFailed(resourceType string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeRepositoryDeleteFailed, errors.KindInternal, "Failed to delete repository record").
		WithMeta("resource_type", resourceType)
}

// NewRepositoryDuplicateEntry creates an error for duplicate entry violations
func NewRepositoryDuplicateEntry(resourceType string, field string, value interface{}) *errors.AppError {
	return errors.New(CodeRepositoryDuplicateEntry, errors.KindConflict, "Duplicate entry in repository").
		WithMeta("resource_type", resourceType).
		WithMeta("field", field).
		WithMeta("value", value)
}

// HandleGormError converts GORM errors to AppError
func HandleGormError(err error, resourceType string, operation string) *errors.AppError {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	if appErr := errors.AsAppError(err); appErr != nil {
		return appErr.WithOp(operation)
	}

	// Handle GORM-specific errors
	if err == gorm.ErrRecordNotFound {
		return NewRepositoryNotFound(resourceType, "unknown").
			WithOp(operation)
	}

	// For other GORM errors, wrap them
	return errors.Wrap(err, CodeDatabaseQueryFailed, errors.KindInternal, "Database operation failed").
		WithOp(operation).
		WithMeta("resource_type", resourceType)
}

// NewResourceNotFound creates an error for when a resource is not found
func NewResourceNotFound(resourceID interface{}) *errors.AppError {
	return errors.New(CodeResourceNotFound, errors.KindNotFound, "Resource not found").
		WithMeta("resource_id", resourceID)
}

// NewResourceCreateFailed creates an error for resource creation failures
func NewResourceCreateFailed(cause error) *errors.AppError {
	return errors.Wrap(cause, CodeResourceCreateFailed, errors.KindInternal, "Failed to create resource")
}

// NewResourceUpdateFailed creates an error for resource update failures
func NewResourceUpdateFailed(cause error) *errors.AppError {
	return errors.Wrap(cause, CodeResourceUpdateFailed, errors.KindInternal, "Failed to update resource")
}

// NewResourceDeleteFailed creates an error for resource deletion failures
func NewResourceDeleteFailed(cause error) *errors.AppError {
	return errors.Wrap(cause, CodeResourceDeleteFailed, errors.KindInternal, "Failed to delete resource")
}

// NewResourceInvalidID creates an error for invalid resource ID
func NewResourceInvalidID(id interface{}) *errors.AppError {
	return errors.New(CodeResourceInvalidID, errors.KindValidation, "Invalid resource ID").
		WithMeta("id", id)
}

// NewProjectNotFound creates an error for when a project is not found
func NewProjectNotFound(projectID interface{}) *errors.AppError {
	return errors.New(CodeProjectNotFound, errors.KindNotFound, "Project not found").
		WithMeta("project_id", projectID)
}

// NewProjectCreateFailed creates an error for project creation failures
func NewProjectCreateFailed(cause error) *errors.AppError {
	return errors.Wrap(cause, CodeProjectCreateFailed, errors.KindInternal, "Failed to create project")
}

// NewUserNotFound creates an error for when a user is not found
func NewUserNotFound(userID interface{}) *errors.AppError {
	return errors.New(CodeUserNotFound, errors.KindNotFound, "User not found").
		WithMeta("user_id", userID)
}

// NewUserDuplicateEmail creates an error for duplicate email
func NewUserDuplicateEmail(email string) *errors.AppError {
	return errors.New(CodeUserDuplicateEmail, errors.KindConflict, "User with this email already exists").
		WithMeta("email", email)
}

// NewAuthUnauthorized creates an error for unauthorized access
func NewAuthUnauthorized(reason string) *errors.AppError {
	return errors.New(CodeAuthUnauthorized, errors.KindUnauthorized, "Unauthorized access").
		WithMeta("reason", reason)
}

// NewAuthForbidden creates an error for forbidden access
func NewAuthForbidden(reason string) *errors.AppError {
	return errors.New(CodeAuthForbidden, errors.KindForbidden, "Forbidden access").
		WithMeta("reason", reason)
}

// NewAuthTokenInvalid creates an error for invalid token
func NewAuthTokenInvalid(reason string) *errors.AppError {
	return errors.New(CodeAuthTokenInvalid, errors.KindUnauthorized, "Invalid authentication token").
		WithMeta("reason", reason)
}

// NewAuthInvalidCredentials creates an error for invalid credentials
func NewAuthInvalidCredentials() *errors.AppError {
	return errors.New(CodeAuthInvalidCredentials, errors.KindUnauthorized, "Invalid credentials")
}
