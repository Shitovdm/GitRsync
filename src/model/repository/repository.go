package repository

import (
	"errors"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

const (
	//	ConfigFileName describe repositories config file name.
	ConfigFileName = "Repositories.json"

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

// Repository struct describes repository config.
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

// GetUUID returns repository UUID.
func (r *Repository) GetUUID() string {
	return r.UUID
}

// SetUUID sets repository UUID.
func (r *Repository) SetUUID(UUID string) {
	r.UUID = UUID
}

// GetName returns repository Name.
func (r *Repository) GetName() string {
	return r.Name
}

// SetName sets repository Name.
func (r *Repository) SetName(name string) {
	r.Name = name
}

// GetSourcePlatformUUID returns repository SourcePlatformUUID.
func (r *Repository) GetSourcePlatformUUID() string {
	return r.SourcePlatformUUID
}

// SetSourcePlatformUUID sets repository SourcePlatformUUID.
func (r *Repository) SetSourcePlatformUUID(sourcePlatformUUID string) {
	r.SourcePlatformUUID = sourcePlatformUUID
}

// GetSourcePlatformPath returns repository SourcePlatformPath.
func (r *Repository) GetSourcePlatformPath() string {
	return r.SourcePlatformPath
}

// SetSourcePlatformPath sets repository SourcePlatformPath.
func (r *Repository) SetSourcePlatformPath(sourcePlatformPath string) {
	r.SourcePlatformPath = sourcePlatformPath
}

// GetDestinationPlatformUUID returns repository DestinationPlatformUUID.
func (r *Repository) GetDestinationPlatformUUID() string {
	return r.DestinationPlatformUUID
}

// SetDestinationPlatformUUID sets repository DestinationPlatformUUID.
func (r *Repository) SetDestinationPlatformUUID(destinationPlatformUUID string) {
	r.DestinationPlatformUUID = destinationPlatformUUID
}

// GetDestinationPlatformPath returns repository DestinationPlatformPath.
func (r *Repository) GetDestinationPlatformPath() string {
	return r.DestinationPlatformPath
}

// SetDestinationPlatformPath sets repository DestinationPlatformPath.
func (r *Repository) SetDestinationPlatformPath(destinationPlatformPath string) {
	r.DestinationPlatformPath = destinationPlatformPath
}

// GetStatus returns repository Status.
func (r *Repository) GetStatus() string {
	return r.Status
}

// SetStatus sets repository Status.
func (r *Repository) SetStatus(status string) {
	r.Status = status
}

// GetState returns repository State.
func (r *Repository) GetState() string {
	return r.State
}

// SetState sets repository State.
func (r *Repository) SetState(state string) {
	r.State = state
}

// GetUpdatedAt returns repository UpdatedAt.
func (r *Repository) GetUpdatedAt() string {
	return r.UpdatedAt
}

// SetUpdatedAt sets repository UpdatedAt.
func (r *Repository) SetUpdatedAt(updatedAt string) {
	r.UpdatedAt = updatedAt
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
	err := conf.Load(ConfigFileName, &repositoriesConfig)
	if err != nil {
		fmt.Printf("Error while loading repositories config file! %s", err.Error())
		err = conf.Save(ConfigFileName, []map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error while creating new repositories config file! %s", err.Error())
		}
		return []Repository{}
	}

	return repositoriesConfig
}

// GetAllInInterface returns repositories config data.
func GetAllInInterface() ([]map[string]interface{}, error) {

	var repositoriesConfig []map[string]interface{}
	err := conf.Load(ConfigFileName, &repositoriesConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load repositories configuration")
	}
	return repositoriesConfig, nil
}

// GetActive returns active repositories config.
func GetActive() ([]map[string]interface{}, error) {

	repositoriesConfig, _ := GetAllInInterface()
	var activeRepositories []map[string]interface{}
	for _, repo := range repositoriesConfig {
		if repo["state"] == StateActive {
			activeRepositories = append(activeRepositories, repo)
		}
	}

	return activeRepositories, nil
}

// GetBlocked returns blocked repositories config.
func GetBlocked() ([]map[string]interface{}, error) {

	repositoriesConfig, _ := GetAllInInterface()
	var blockedRepositories []map[string]interface{}
	for _, repo := range repositoriesConfig {
		if repo["state"] == StateBlocked {
			blockedRepositories = append(blockedRepositories, repo)
		}
	}

	return blockedRepositories, nil
}

// Create creates new repository..
func Create(r *Repository) error {

	UUID, _ := uuid.NewV4()
	r.UUID = UUID.String()
	r.Status = StatusInitiated
	r.State = "active"
	r.UpdatedAt = ""
	repositories := GetAll()
	repositories = append(repositories, *r)
	err := saveRepositories(repositories)
	if err != nil {
		return err
	}

	return nil
}

// Update describes update repository status action.
func (r *Repository) Update() error {

	t := time.Now()
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
			oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
		}
	}

	err := saveRepositories(oldRepositoriesList)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes repository config.
func (r *Repository) Delete() error {

	oldRepositoriesList := GetAll()
	newRepositoriesList := make([]Repository, 0)
	for _, repository := range oldRepositoriesList {
		if repository.UUID != r.UUID {
			newRepositoriesList = append(newRepositoriesList, repository) //nolint:staticcheck
		}
	}

	err := saveRepositories(newRepositoriesList)
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

// saveRepositories stores repositories config data.
func saveRepositories(repositories []Repository) error {
	err := conf.Save(ConfigFileName, &repositories)
	if err != nil {
		return fmt.Errorf("Error while saving repositories config file! %s ", err.Error())
	}
	return nil
}
