package safe


const RealGitPath = "/opt/local/bin/git"

var CommandFilter = map[string]bool{
	//These are common Git commands used in various situations:

   // start a working area (see also: git help tutorial)
   "clone":   true,  // Clone a repository into a new directory
   "init":    true,  // Create an empty Git repository or reinitialize an existing one

   // work on the current change (see also: git help everyday)
   "add":     true,  // Add file contents to the index
   "mv":      true,  // Move or rename a file, a directory, or a symlink
   "restore": true,  // Restore working tree files
   "rm":      true,  // Remove files from the working tree and from the index

   // examine the history and state (see also: git help revisions)
   "bisect":  false,  // Use binary search to find the commit that introduced a bug
   "diff":    false,  // Show changes between commits, commit and working tree, etc
   "grep":    false,  // Print lines matching a pattern
   "log":     false,  // Show commit logs
   "show":    false,  // Show various types of objects
   "status":  false, // Show the working tree status

   // grow, mark and tweak your common history
   "branch":  true, // List, create, or delete branches
   "commit":  true, // Record changes to the repository
   "merge":   true,  // Join two or more development histories together
   "rebase":  true, // Reapply commits on top of another base tip
   "reset":   true, // Reset current HEAD to the specified state
   "switch":  true, // Switch branches
   "tag":     true, // Create, list, delete or verify a tag object signed with GPG

   // collaborate (see also: git help workflows)
   "fetch":   true, // Download objects and refs from another repository
   "pull":    true, // Fetch from and integrate with another repository or a local branch
   "push":    true, // Update remote refs along with associated objects

   // possible shortcuts
   "co":      true,  // Shortcut for checkout
   "br":      true,  // Shortcut for branch
   "ci":      true,  // Shortcut for commit
   "st":      true,  // Shortcut for status
   "lg":      true,  // Shortcut for log
}

func GitMain() {
	SafeRun(SafeCfg{
		RealPath:  RealGitPath,
		Name:      "git",
		CmdFilter: CommandFilter,
	})
}
