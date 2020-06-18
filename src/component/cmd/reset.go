package cmd

// ResetHard resets hard repository.
func ResetHard(path string) bool {

	if path == "" {
		return false
	}

	cmd := command("git", "reset --hard")
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		return false
	}

	return true
}
