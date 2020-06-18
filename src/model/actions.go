package model

type PullActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

type PushActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

type CleanActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

type InfoActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

type BlockActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

type ActivateActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}
