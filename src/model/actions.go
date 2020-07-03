package model

// PullActionRequest struct describes pull action request.
type PullActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// PushActionRequest struct describes push action request.
type PushActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// CleanActionRequest struct describes clean action request.
type CleanActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// InfoActionRequest struct describes info action request.
type InfoActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// BlockActionRequest struct describes block action request.
type BlockActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// ActivateActionRequest struct describes activate action request.
type ActivateActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}

// OpenDirActionRequest struct describes open fs dir action request.
type OpenDirActionRequest struct {
	Path string `json:"path"`
}

// SyncTagsActionRequest struct describes sync tags action request.
type SyncTagsActionRequest struct {
	RepositoryUUID string `json:"uuid"`
}
