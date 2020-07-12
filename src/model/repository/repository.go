package repository

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

const (
	// StatusInitiated describe status "initiated"
	StatusInitiated = "initiated"
	// StatusPendingPull describe status "pending_pull"
	StatusPendingPull = "pending_pull"
	// StatusPulled describe status "pulled"
	StatusPulled = "pulled"
	// StatusPullFailed describe status "pull_failed"
	StatusPullFailed = "pull_failed"
	// StatusPendingPush describe status "pending_push"
	StatusPendingPush = "pending_push"
	// StatusPushed describe status "pushed"
	StatusPushed = "pushed"
	// StatusPushFailed describe status "push_failed"
	StatusPushFailed = "push_failed"
	// StatusPendingClean describe status "pending_clear"
	StatusPendingClean = "pending_clear"
	// StatusCleaned describe status "cleared"
	StatusCleaned = "cleared"
	// StatusCleanFailed describe status "clear_failed"
	StatusCleanFailed = "clear_failed"
	// StatusSynchronized describe status "synced"
	StatusSynchronized = "synced"
	// StatusFailed describe status "failed"
	StatusFailed = "failed"

	// StateActive describe state "active"
	StateActive = "active"
	// StateBlocked describe state "blocked"
	StateBlocked = "blocked"
)

var (
	timeFormat = "02-01-2006 15:04"
)

// RepositoryConfig struct describes repository config.
type Repository struct {
	UUID                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
	Status                  string `json:"status"`
	State                   string `json:"state"`
	UpdatedAt               string `json:"updated_at"`
}

// AddRepositoryRequest struct describes add repository request model.
type AddRepositoryRequest struct {
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

// EditRepositoryRequest struct describes edit repository request model.
type EditRepositoryRequest struct {
	UUID                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

// RemoveRepositoryRequest struct describes remove repository request model.
type RemoveRepositoryRequest struct {
	UUID string `json:"uuid"`
}

func (r *Repository) GetUUID() string {
	return r.UUID
}
func (r *Repository) SetUUID(UUID string) {
	r.UUID = UUID
}
func (r *Repository) GetName() string {
	return r.Name
}
func (r *Repository) SetName(name string) {
	r.Name = name
}

func (r *Repository) GetStatus() string {
	return r.Status
}
func (r *Repository) SetStatus(status string) {
	r.Status = status
}

// Get returns repository config by repository UUID.
func Get(UUID string) *Repository {

	repositoriesList := GetAll()
	for _, repository := range repositoriesList {
		if repository.UUID == UUID {
			return &repository
		}
	}

	return nil
}

// GetAll returns repositories config.
func GetAll() []Repository {

	repositoriesConfig := make([]Repository, 0)
	err := conf.Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		fmt.Printf("Error while loading repositories config file! %s", err.Error())
		err = conf.Save("Repositories.json", []map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error while creating new repositories config file! %s", err.Error())
		}
		return []Repository{}
	}

	return repositoriesConfig
}

// Update describes update repository status action.
func Create(r *Repository) error {

	UUID, _ := uuid.NewV4()
	r.UUID = UUID.String()
	r.Status = StatusInitiated
	r.State = "active"
	r.UpdatedAt = ""
	repositories := GetAll()
	repositories = append(repositories, *r)
	err := conf.SaveRepositories(repositories)
	if err != nil {
		return err
	}

	return nil
}

// Update describes update repository status action.
func (r *Repository) Update() error {

	oldRepositoriesList := GetAll()
	for i, repository := range oldRepositoriesList {
		if repository.UUID == r.UUID {
			oldRepositoriesList[i].Name = r.Name
			oldRepositoriesList[i].SourcePlatformUUID = r.SourcePlatformUUID
			oldRepositoriesList[i].SourcePlatformPath = r.SourcePlatformPath
			oldRepositoriesList[i].DestinationPlatformUUID = r.DestinationPlatformUUID
			oldRepositoriesList[i].DestinationPlatformPath = r.DestinationPlatformPath
			oldRepositoriesList[i].Status = r.Status
			oldRepositoriesList[i].State = r.State
			if r.Status == StatusPulled || r.Status == StatusPushed || r.Status == StatusSynchronized {
				t := time.Now()
				oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
			}
		}
	}

	err := conf.SaveRepositories(oldRepositoriesList)
	if err != nil {
		return err
	}

	return nil
}

// GetSourceRepositoryName parses source repository name.
func (r *Repository) GetSourceRepositoryName() string {

	spp := strings.Trim(strings.TrimRight(r.SourcePlatformPath, "git"), ".")
	return strings.Split(spp, "/")[len(strings.Split(spp, "/"))-1]
}

// GetDestinationRepositoryName parses destination repository name.
func (r *Repository) GetDestinationRepositoryName() string {

	dpp := strings.Trim(strings.TrimRight(r.DestinationPlatformPath, "git"), ".")
	return strings.Split(dpp, "/")[len(strings.Split(dpp, "/"))-1]
}
