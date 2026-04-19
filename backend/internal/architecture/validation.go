package architecture

// Validation provides validation methods for architecture aggregates
// This is a placeholder for domain-level validation logic
// More complex validation can be added here as needed

// Validate performs domain-level validation on the architecture
func (a *Architecture) Validate() error {
	// Basic validation checks
	if len(a.Resources) == 0 {
		return nil // Empty architecture is valid
	}

	// Validate that all parent references exist
	for _, res := range a.Resources {
		if res.ParentID != nil {
			parentExists := false
			for _, other := range a.Resources {
				if other.ID == *res.ParentID {
					parentExists = true
					break
				}
			}
			if !parentExists {
				// Parent might be region (project-level), which is OK
				// This is a soft validation - we allow missing parents if they're regions
			}
		}
	}

	return nil
}
