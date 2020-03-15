package Models

type PullActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type PushActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type BlockActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}

type ActiveActionRequest struct {
	RepositoryUuid string `json:"uuid"`
}