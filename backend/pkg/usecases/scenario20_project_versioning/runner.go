// scenario20_project_versioning demonstrates the immutable versioning system:
//  1. Create a project with an initial architecture (v1: VPC + EC2 + S3)
//  2. Verify version 1 is recorded in project_versions
//  3. Save an updated architecture (v2: adds ALB + RDS, upgrades EC2)
//     → SaveArchitecture returns a NEW project_id (immutable snapshot)
//  4. Verify the version chain (v1 → v2) using GetVersions
//  5. Verify the old project_id still loads v1's architecture unchanged
//
// Run from repo root:
//
//	go run backend/pkg/usecases/scenario20_project_versioning/runner.go
//
// Or from backend dir:
//
//	go run pkg/usecases/scenario20_project_versioning/runner.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"

	// Register AWS cloud types
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"

	infrastructurerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/infrastructure"
	projectrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/project"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
	"gorm.io/gorm"
)

// ─── local repo adapters (same pattern as scenario15) ──────────────────────

type projectRepoAdapter struct{ *projectrepo.ProjectRepository }

func (a *projectRepoAdapter) BeginTransaction(ctx context.Context) (interface{}, context.Context) {
	return a.ProjectRepository.BeginTransaction(ctx)
}
func (a *projectRepoAdapter) CommitTransaction(tx interface{}) error {
	if g, ok := tx.(*gorm.DB); ok {
		return a.ProjectRepository.CommitTransaction(g)
	}
	return fmt.Errorf("not a *gorm.DB")
}
func (a *projectRepoAdapter) RollbackTransaction(tx interface{}) error {
	if g, ok := tx.(*gorm.DB); ok {
		return a.ProjectRepository.RollbackTransaction(g)
	}
	return fmt.Errorf("not a *gorm.DB")
}

type depTypeRepoAdapter struct {
	*resourcerepo.DependencyTypeRepository
}

func (a *depTypeRepoAdapter) Create(_ context.Context, _ *models.DependencyType) error {
	return fmt.Errorf("Create not implemented in adapter")
}

// ─── helpers ───────────────────────────────────────────────────────────────

// loadArchFile reads a JSON file and unmarshals it into UpdateArchitectureRequest.
// Tries multiple paths so the runner can be invoked from repo-root or backend dir.
func loadArchFile(candidates []string) (*dto.UpdateArchitectureRequest, error) {
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			var req dto.UpdateArchitectureRequest
			if err := json.Unmarshal(data, &req); err != nil {
				return nil, fmt.Errorf("parse %s: %w", p, err)
			}
			fmt.Printf("  loaded %s (%d bytes)\n", p, len(data))
			return &req, nil
		}
	}
	return nil, fmt.Errorf("could not find architecture file (tried: %v)", candidates)
}

func printSep(title string) {
	fmt.Printf("\n══════════════════════════════════════════\n  %s\n══════════════════════════════════════════\n", title)
}

func check(label string, err error) {
	if err != nil {
		log.Fatalf("❌  %s: %v", label, err)
	}
}

// ─── main ──────────────────────────────────────────────────────────────────

func main() {
	if err := run(); err != nil {
		log.Fatalf("FAILED: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	printSep("Scenario 20 – Project Versioning Demo")

	// ── 1. Database ─────────────────────────────────────────────────────────
	db, err := database.Connect()
	check("database connect", err)
	fmt.Println("✓ Database connected")

	// Clean up scenario data from previous runs (keep resource types etc.)
	for _, q := range []string{
		"DELETE FROM resource_dependencies",
		"DELETE FROM resource_containments",
		"DELETE FROM resources",
		"DELETE FROM project_versions WHERE project_id IN (SELECT id FROM projects WHERE name LIKE 'Versioning Demo%')",
		"DELETE FROM projects WHERE name LIKE 'Versioning Demo%'",
	} {
		if err := db.Exec(q).Error; err != nil {
			fmt.Printf("  ⚠ cleanup: %v\n", err)
		}
	}
	fmt.Println("✓ Cleaned up previous scenario20 data")

	// Ensure required resource types exist
	for _, rt := range []string{"vpc", "subnet", "ec2", "s3", "security-group", "region", "alb", "rds", "route-table"} {
		var count int64
		db.Model(&models.ResourceType{}).Where("name = ? AND cloud_provider = 'aws'", rt).Count(&count)
		if count == 0 {
			db.Create(&models.ResourceType{Name: rt, CloudProvider: "aws", IsRegional: true})
			fmt.Printf("  ✓ seeded ResourceType: %s\n", rt)
		}
	}
	// Ensure depends_on dependency type
	var depTypeCount int64
	db.Model(&models.DependencyType{}).Where("name = ?", "depends_on").Count(&depTypeCount)
	if depTypeCount == 0 {
		db.Create(&models.DependencyType{Name: "depends_on"})
		fmt.Println("  ✓ seeded DependencyType: depends_on")
	}

	// ── 2. Repositories ──────────────────────────────────────────────────────
	logger := slog.Default()

	projectRepoRaw, err := projectrepo.NewProjectRepository(logger)
	check("project repo", err)
	projectRepo := &projectRepoAdapter{projectRepoRaw}

	versionRepoRaw, err := projectrepo.NewProjectVersionRepository()
	check("version repo", err)
	versionRepo := &services.ProjectVersionRepositoryAdapter{Repo: versionRepoRaw}

	resourceRepo, err := resourcerepo.NewResourceRepository(logger)
	check("resource repo", err)
	resourceTypeRepo, err := resourcerepo.NewResourceTypeRepository()
	check("resource type repo", err)
	containmentRepo, err := resourcerepo.NewResourceContainmentRepository()
	check("containment repo", err)
	depRepo, err := resourcerepo.NewResourceDependencyRepository()
	check("dependency repo", err)
	depTypeRepoRaw, err := resourcerepo.NewDependencyTypeRepository()
	check("dep type repo", err)
	depTypeRepo := &depTypeRepoAdapter{depTypeRepoRaw}
	userRepo, err := userrepo.NewUserRepository()
	check("user repo", err)
	iacTargetRepo, err := infrastructurerepo.NewIACTargetRepository()
	check("iac target repo", err)
	variableRepo, err := projectrepo.NewProjectVariableRepository()
	check("variable repo", err)
	outputRepo, err := projectrepo.NewProjectOutputRepository()
	check("output repo", err)

	// ── 3. Service ──────────────────────────────────────────────────────────
	projectService := services.NewProjectService(
		projectRepo,
		versionRepo,
		resourceRepo,
		resourceTypeRepo,
		containmentRepo,
		depRepo,
		depTypeRepo,
		userRepo,
		iacTargetRepo,
		variableRepo,
		outputRepo,
	)
	fmt.Println("✓ Services initialized")

	// ── 4. Test user ─────────────────────────────────────────────────────────
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if _, err := userRepo.FindByID(ctx, userID); err != nil {
		newUser := &models.User{ID: userID, Name: "Demo User", Email: "demo@example.com", IsVerified: true}
		if e := userRepo.Create(ctx, newUser); e != nil {
			fmt.Printf("  ⚠ create user: %v\n", e)
		} else {
			fmt.Println("  ✓ Created demo user")
		}
	}

	// ── 5. Create initial project (no architecture yet) ─────────────────────
	printSep("Phase 1 – Create Project (v1 Snapshot)")

	proj, err := projectService.Create(ctx, &serverinterfaces.CreateProjectRequest{
		UserID:        userID,
		Name:          "Versioning Demo Project",
		Description:   "Tests immutable versioning",
		CloudProvider: "aws",
		Region:        "us-east-1",
		IACTargetID:   1,
	})
	check("create project", err)
	rootProjectID := proj.ID
	fmt.Printf("✓ Project created:  ID=%s  Name=%s\n", rootProjectID, proj.Name)

	// ── 6. Save v1 architecture (triggers first cloneProjectSnapshot) ────────
	v1Req, err := loadArchFile([]string{
		"pkg/usecases/scenario20_project_versioning/architecture_v1.json",
		"backend/pkg/usecases/scenario20_project_versioning/architecture_v1.json",
	})
	check("load v1 json", err)

	v1Detail, err := projectService.CreateVersion(ctx, rootProjectID, &serverinterfaces.CreateVersionRequest{
		Nodes:     v1Req.Nodes,
		Edges:     v1Req.Edges,
		Variables: v1Req.Variables,
		Outputs:   v1Req.Outputs,
		Message:   "Initial architecture (v1 – VPC + EC2 + S3)",
	})
	check("CreateVersion v1", err)

	v1VersionID := v1Detail.ID
	v1ProjectID := v1Detail.ProjectID
	fmt.Printf("✓ v1 saved:  project_id=%s  version_id=%s  version_number=%d\n",
		v1ProjectID, v1VersionID, v1Detail.VersionNumber)

	if v1ProjectID == rootProjectID {
		return fmt.Errorf("FAIL: v1 project_id should differ from root (CreateVersion must clone)")
	}
	fmt.Printf("✓ Root project ID %s ≠ v1 snapshot ID %s (correct)\n", rootProjectID, v1ProjectID)

	// Verify v1 architecture has expected node count
	v1Arch, err := projectService.GetArchitecture(ctx, v1ProjectID)
	check("GetArchitecture v1", err)
	fmt.Printf("✓ v1 architecture loaded: %d nodes, %d edges\n", len(v1Arch.Nodes), len(v1Arch.Edges))
	if len(v1Arch.Nodes) == 0 {
		return fmt.Errorf("FAIL: v1 architecture should have nodes")
	}

	// ── 7. Save v2 architecture (adds ALB + RDS) ────────────────────────────
	printSep("Phase 2 – Update Architecture (v2 Snapshot)")

	v2Req, err := loadArchFile([]string{
		"pkg/usecases/scenario20_project_versioning/architecture_v2.json",
		"backend/pkg/usecases/scenario20_project_versioning/architecture_v2.json",
	})
	check("load v2 json", err)

	v2Detail, err := projectService.CreateVersion(ctx, v1ProjectID, &serverinterfaces.CreateVersionRequest{
		Nodes:     v2Req.Nodes,
		Edges:     v2Req.Edges,
		Variables: v2Req.Variables,
		Outputs:   v2Req.Outputs,
		Message:   "Added ALB + RDS, upgraded EC2 (v2)",
	})
	check("CreateVersion v2", err)

	v2ProjectID := v2Detail.ProjectID
	fmt.Printf("✓ v2 saved:  project_id=%s  version_id=%s  version_number=%d\n",
		v2ProjectID, v2Detail.ID, v2Detail.VersionNumber)

	if v2ProjectID == v1ProjectID {
		return fmt.Errorf("FAIL: v2 project_id should differ from v1 (CreateVersion must clone)")
	}
	fmt.Printf("✓ v1 ID %s ≠ v2 ID %s (correct)\n", v1ProjectID, v2ProjectID)

	// Verify v2 has more nodes than v1
	v2Arch, err := projectService.GetArchitecture(ctx, v2ProjectID)
	check("GetArchitecture v2", err)
	fmt.Printf("✓ v2 architecture loaded: %d nodes, %d edges\n", len(v2Arch.Nodes), len(v2Arch.Edges))
	if len(v2Arch.Nodes) <= len(v1Arch.Nodes) {
		fmt.Printf("  ⚠ v2 should have more nodes than v1 (%d vs %d)\n", len(v2Arch.Nodes), len(v1Arch.Nodes))
	} else {
		fmt.Printf("  ✓ v2 has more nodes than v1 (%d > %d)\n", len(v2Arch.Nodes), len(v1Arch.Nodes))
	}

	// ── 8. Verify version chain ──────────────────────────────────────────────
	printSep("Phase 3 – Verify Version Chain")

	// GetVersions can take any version in the chain — let's pass the v2 snapshot
	versions, err := projectService.GetVersions(ctx, v2ProjectID)
	check("GetVersions", err)

	fmt.Printf("✓ Version chain has %d entries:\n", len(versions))
	for _, v := range versions {
		parentStr := "<root>"
		if v.ParentVersionID != nil {
			parentStr = v.ParentVersionID.String()
		}
		fmt.Printf("    [v%d] version_id=%-36s  project_id=%-36s  parent=%s\n",
			v.VersionNumber, v.ID, v.ProjectID, parentStr)
	}

	if len(versions) < 2 {
		return fmt.Errorf("FAIL: expected at least 2 versions in chain, got %d", len(versions))
	}

	// Verify chain ordering (ascending by version_number)
	for i := 1; i < len(versions); i++ {
		if versions[i].VersionNumber <= versions[i-1].VersionNumber {
			return fmt.Errorf("FAIL: versions not ordered ascending (got %d after %d)",
				versions[i].VersionNumber, versions[i-1].VersionNumber)
		}
	}
	fmt.Println("✓ Versions ordered ascending by version_number")

	// Verify parent chain integrity: each version (except root) must point to previous
	for i := 1; i < len(versions); i++ {
		if versions[i].ParentVersionID == nil {
			return fmt.Errorf("FAIL: version %d has nil parent_version_id; expected the previous version", i)
		}
		if *versions[i].ParentVersionID != versions[i-1].ID {
			return fmt.Errorf("FAIL: version %d parent_version_id=%s does not match previous version id=%s",
				i, *versions[i].ParentVersionID, versions[i-1].ID)
		}
	}
	fmt.Println("✓ Parent chain integrity validated")

	// ── 9. Verify v1 snapshot is immutable (old project_id still works) ──────
	printSep("Phase 4 – Immutability Verification")

	v1ArchRecheck, err := projectService.GetArchitecture(ctx, v1ProjectID)
	check("GetArchitecture v1 (recheck after v2 save)", err)
	if len(v1ArchRecheck.Nodes) != len(v1Arch.Nodes) {
		return fmt.Errorf("FAIL: v1 snapshot was mutated by v2 save! (nodes: %d → %d)",
			len(v1Arch.Nodes), len(v1ArchRecheck.Nodes))
	}
	fmt.Printf("✓ v1 snapshot unchanged after v2 save (%d nodes, immutable ✓)\n", len(v1ArchRecheck.Nodes))

	// ── 10. Restore v1 as a new version ─────────────────────────────────────
	printSep("Phase 5 – Restore v1 (creates v3)")

	// Restore v1 by creating a new version with v1's state
	v1ArchForRestore, err := projectService.GetArchitecture(ctx, v1ProjectID)
	check("GetArchitecture v1 for restore", err)

	v3Detail, err := projectService.CreateVersion(ctx, v2ProjectID, &serverinterfaces.CreateVersionRequest{
		Nodes:     v1ArchForRestore.Nodes,
		Edges:     v1ArchForRestore.Edges,
		Variables: v1ArchForRestore.Variables,
		Outputs:   v1ArchForRestore.Outputs,
		Message:   "Restored to v1 state",
	})
	check("CreateVersion (restore to v1)", err)
	v3ProjectID := v3Detail.ProjectID
	fmt.Printf("✓ Restore created:  project_id=%s  version_id=%s  version_number=%d\n",
		v3ProjectID, v3Detail.ID, v3Detail.VersionNumber)

	if v3ProjectID == v1ProjectID || v3ProjectID == v2ProjectID {
		return fmt.Errorf("FAIL: restored project_id should be brand new")
	}

	// Final version count
	versionsAfterRestore, err := projectService.GetVersions(ctx, v2ProjectID)
	check("GetVersions after restore", err)
	fmt.Printf("✓ Version chain after restore: %d entries\n", len(versionsAfterRestore))

	// ── Summary ──────────────────────────────────────────────────────────────
	printSep("RESULTS")
	fmt.Println("Root project ID :", rootProjectID)
	fmt.Println("v1 snapshot ID  :", v1ProjectID)
	fmt.Println("v2 snapshot ID  :", v2ProjectID)
	fmt.Println("v3 (restore) ID :", v3ProjectID)
	fmt.Printf("Total versions  : %d\n", len(versionsAfterRestore))
	fmt.Println()
	fmt.Println("✅  Scenario 20 PASSED – immutable versioning works correctly")
	return nil
}
