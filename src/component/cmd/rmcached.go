package cmd

// RmCached flushes cache.
func RmCached(path string) bool {

	if path == "" {
		return false
	}

	cmd := command("git", "rm --cached -r .")
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		return false
	}

	return true
}
