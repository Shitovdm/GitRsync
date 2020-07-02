package cmd

import (
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"os"
)

// CopyRepository copies repository.
func CopyRepository(repositoryFullPath string, destinationRepositoryName string, sourceRepositoryName string) bool {

	//	Moving .git to temporary folder.
	err := TemporaryMoveGitFolder(repositoryFullPath, sourceRepositoryName)
	if err != nil {
		return false
	}

	//	Removing .git folder from source repository.
	err = os.RemoveAll(repositoryFullPath + "/source/" + sourceRepositoryName + "/.git")
	if err != nil {
		return false
	}

	//	Copy all repository files from source repo to destination repo.
	err = helper.CopyDirContent(repositoryFullPath+"/source/"+sourceRepositoryName, repositoryFullPath+"/destination/"+destinationRepositoryName)
	if err != nil {
		return false
	}

	//	Rewrite .git folder.
	err = RewriteGitFiles(repositoryFullPath, destinationRepositoryName)
	if err != nil {
		return false
	}

	//	Restore .git folder for source repository.
	err = RestoreGitFolder(repositoryFullPath, sourceRepositoryName)
	if err != nil { //nolint:gosimple
		return false
	}

	//	Remove temporary folder.
	err = RemoveTemporaryGitFolder(repositoryFullPath)

	return err == nil
}

// TemporaryMoveGitFolder moves git folder.
func TemporaryMoveGitFolder(repositoryFullPath string, sourceRepositoryName string) error {
	if !helper.IsDirExists(repositoryFullPath + "/tmp/.git") {
		err := helper.CreateNewDir(repositoryFullPath + "/tmp/.git")
		if err != nil {
			return err
		}
	}

	err := helper.CopyDirContent(repositoryFullPath+"/source/"+sourceRepositoryName+"/.git", repositoryFullPath+"/tmp/.git")
	if err != nil {
		return err
	}
	return nil
}

// RestoreGitFolder restores git folder.
func RestoreGitFolder(repositoryFullPath string, sourceRepositoryName string) error {

	if !helper.IsDirExists(repositoryFullPath + "/source/" + sourceRepositoryName + "/.git") {
		err := helper.CreateNewDir(repositoryFullPath + "/source/" + sourceRepositoryName + "/.git")
		if err != nil {
			return err
		}
	}

	err := helper.CopyDirContent(repositoryFullPath+"/tmp/.git", repositoryFullPath+"/source/"+sourceRepositoryName+"/.git")
	if err != nil {
		return err
	}
	return nil
}

// RemoveTemporaryGitFolder removes temp git folder.
func RemoveTemporaryGitFolder(repositoryFullPath string) error {
	err := os.RemoveAll(repositoryFullPath + "/tmp")
	if err != nil {
		return err
	}
	return nil
}

// RewriteGitFiles rewrites git files.
func RewriteGitFiles(repositoryFullPath string, destinationRepositoryName string) error {
	tmpGitFolder := repositoryFullPath + "/tmp/.git"
	destinationGitFolder := repositoryFullPath + "/destination/" + destinationRepositoryName + "/.git"

	//	Rewrite folders.
	_ = helper.CopyDirContent(tmpGitFolder+"/logs", destinationGitFolder+"/logs")
	_ = helper.CopyDirContent(tmpGitFolder+"/objects", destinationGitFolder+"/objects")
	_ = helper.CopyDirContent(tmpGitFolder+"/smartgit", destinationGitFolder+"/smartgit")
	_ = helper.CopyDirContent(tmpGitFolder+"/refs/heads", destinationGitFolder+"/refs/heads")
	_ = helper.CopyDirContent(tmpGitFolder+"/refs/tags", destinationGitFolder+"/refs/tags")

	//	Rewrite files.
	_ = helper.CopyFile(tmpGitFolder+"/index", destinationGitFolder+"/index")
	_ = helper.CopyFile(tmpGitFolder+"/HEAD", destinationGitFolder+"/HEAD")
	//_ = helper.CopyFile(tmpGitFolder+"/packed-refs", destinationGitFolder+"/packed-refs")

	return nil
}
