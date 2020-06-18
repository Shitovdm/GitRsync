package model

type PullActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type PushActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type CleanActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type InfoActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type BlockActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type ActivateActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}
