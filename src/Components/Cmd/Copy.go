package Cmd

import (
	"github.com/Shitovdm/GitRsync/src/Components/Helpers"
	"os"
)

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
	err = Helpers.CopyDirContent(repositoryFullPath+"/source/"+sourceRepositoryName, repositoryFullPath+"/destination/"+destinationRepositoryName)
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
	if err != nil {
		return false
	}

	//	Remove temporary folder.
	err = RemoveTemporaryGitFolder(repositoryFullPath)
	if err != nil {
		return false
	}

	return true
}

func TemporaryMoveGitFolder(repositoryFullPath string, sourceRepositoryName string) error {
	if !Helpers.IsDirExists(repositoryFullPath + "/tmp/.git") {
		err := Helpers.CreateNewDir(repositoryFullPath + "/tmp/.git")
		if err != nil {
			return err
		}
	}

	err := Helpers.CopyDirContent(repositoryFullPath+"/source/"+sourceRepositoryName+"/.git", repositoryFullPath+"/tmp/.git")
	if err != nil {
		return err
	}
	return nil
}

func RestoreGitFolder(repositoryFullPath string, sourceRepositoryName string) error {

	if !Helpers.IsDirExists(repositoryFullPath + "/source/" + sourceRepositoryName + "/.git") {
		err := Helpers.CreateNewDir(repositoryFullPath + "/source/" + sourceRepositoryName + "/.git")
		if err != nil {
			return err
		}
	}

	err := Helpers.CopyDirContent(repositoryFullPath+"/tmp/.git", repositoryFullPath+"/source/"+sourceRepositoryName+"/.git")
	if err != nil {
		return err
	}
	return nil
}

func RemoveTemporaryGitFolder(repositoryFullPath string) error {
	err := os.RemoveAll(repositoryFullPath + "/tmp")
	if err != nil {
		return err
	}
	return nil
}

func RewriteGitFiles(repositoryFullPath string, destinationRepositoryName string) error {
	tmpGitFolder := repositoryFullPath + "/tmp/.git"
	destinationGitFolder := repositoryFullPath + "/destination/" + destinationRepositoryName + "/.git"

	//	Rewrite folders.
	_ = Helpers.CopyDirContent(tmpGitFolder+"/logs", destinationGitFolder+"/logs")
	_ = Helpers.CopyDirContent(tmpGitFolder+"/objects", destinationGitFolder+"/objects")
	_ = Helpers.CopyDirContent(tmpGitFolder+"/smartgit", destinationGitFolder+"/smartgit")
	_ = Helpers.CopyDirContent(tmpGitFolder+"/refs/heads", destinationGitFolder+"/refs/heads")
	_ = Helpers.CopyDirContent(tmpGitFolder+"/refs/tags", destinationGitFolder+"/refs/tags")

	//	Rewrite files.
	_ = Helpers.CopyFile(tmpGitFolder+"/index", destinationGitFolder+"/index")
	_ = Helpers.CopyFile(tmpGitFolder+"/HEAD", destinationGitFolder+"/HEAD")

	return nil
}
