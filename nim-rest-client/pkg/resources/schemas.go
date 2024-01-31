package resources

import "time"

// ErrorDetail represents a detailed error message returned by the server.
type ErrorDetail struct {
	Description string `json:"description"`
}

// ErrorModel represents an error, including a numeric error code, details, and a human-readable message.
type ErrorModel struct {
	Code    int           `json:"code"`
	Details []ErrorDetail `json:"details"`
	Message string        `json:"message"`
}

// PaginationResponse is used in list responses that support pagination.
type PaginationResponse struct {
	Count        int `json:"count"`
	ItemsPerPage int `json:"itemsPerPage"`
	StartIndex   int `json:"startIndex"`
}

// ExternalId represents the commit ID of the config.
type ExternalId string

// ExternalIdType represents the type of commit for config update.
type ExternalIdType string

// ConfigVersion represents the configuration version for an NGINX instance or instance group.
type ConfigVersion struct {
	CreateTime     time.Time    `json:"createTime"`
	ExternalId     ExternalId   `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	UID            string       `json:"uid"`
	VersionHash    string       `json:"versionHash"`
}

// GroupConfigVersion represents the configuration version for an instance group.
type GroupConfigVersion struct {
	Instances []ConfigVersion `json:"instances"`
	Versions  []ConfigVersion `json:"versions"`
}

// InstanceStatus represents a success or failure deployment status for an NGINX instance.
type InstanceStatus struct {
	FailMessage string `json:"failMessage"`
	FailType    string `json:"failType"`
	Name        string `json:"name"`
}

// InstanceActivityStatus represents a pending deployment status for an NGINX instance.
type InstanceActivityStatus struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// DeploymentDetail represents additional deployment details for NGINX instances in the group.
type DeploymentDetail struct {
	Error   string                  `json:"error"`
	Failure []InstanceStatus        `json:"failure"`
	Pending []InstanceActivityStatus `json:"pending"`
	Success []InstanceStatus        `json:"success"`
}

// DeploymentDetails represents the outcome of a Publish request for an NGINX instance or Instance Group.
type DeploymentDetails struct {
	CreateTime time.Time        `json:"createTime"`
	Details    DeploymentDetail `json:"details"`
	ID         string           `json:"id"`
	Message    string           `json:"message"`
	Status     string           `json:"status"`
	UpdateTime time.Time        `json:"updateTime"`
}

// SelfLinks contains a link from the resource to itself.
type SelfLinks struct {
	Rel string `json:"rel"`
}

// NamedLinks contains information about the object being referred to.
type NamedLinks struct {
	SelfLinks
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

// InstanceGroup defines an instance group and its member instances.
type InstanceGroup struct {
	BaseConfigName        string              `json:"baseConfigName"`
	ConfigVersion         GroupConfigVersion  `json:"configVersion"`
	Description           string              `json:"description"`
	DisplayName           string              `json:"displayName"`
	Instances             []string            `json:"instances"`
	LastDeploymentDetails DeploymentDetails   `json:"lastDeploymentDetails"`
	Links                 []NamedLinks        `json:"links"`
	Name                  string              `json:"name"`
	UID                   string              `json:"uid"`
}

// InstanceGroupListResponse represents a list of NGINX instance group descriptions.
type InstanceGroupListResponse struct {
	PaginationResponse
	Items []InstanceGroup `json:"items"`
}

// InstanceGroupSimple provides a simplified representation of an NGINX instance group.
type InstanceGroupSimple struct {
	Description   string      `json:"description"`
	DisplayName   string      `json:"displayName"`
	ExternalId     ExternalId `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	Name          string      `json:"name"`
	UID           string      `json:"uid"`
}

// LoadInstanceGroups represents a summary list of NGINX instance groups.
type LoadInstanceGroups struct {
	Count int                   `json:"count"`
	Items []InstanceGroupSimple `json:"items"`
}

// FileData represents file contents.
type FileData struct {
	Contents string `json:"contents"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
}

// AuxData represents auxiliary configuration files.
type AuxData struct {
	Files   []FileData `json:"files"`
	RootDir string     `json:"rootDir"`
}

// ConfigData represents NGINX configuration files.
type ConfigData struct {
	Files   []FileData `json:"files"`
	RootDir string     `json:"rootDir"`
}

// Directory represents information about a directory.
type Directory struct {
	Files       []FileData `json:"files"`
	Name        string     `json:"name"`
	Permissions string     `json:"permissions"`
	Size        int        `json:"size"`
	UpdateTime  time.Time  `json:"updateTime"`
}

// DirectoryMap represents a map of directories.
type DirectoryMap map[string]Directory

// InstanceGroupConfigResponse represents the configuration details for an NGINX instance group.
type InstanceGroupConfigResponse struct {
	AuxFiles    AuxData     `json:"auxFiles"`
	ConfigFiles ConfigData  `json:"configFiles"`
	CreateTime  time.Time   `json:"createTime"`
	DirectoryMap DirectoryMap `json:"directoryMap"`
	ExternalId  ExternalId `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	Instances   []string    `json:"instances"`
	UpdateTime  time.Time   `json:"updateTime"`
}

// PublishConfigRequest represents a configuration publish request.
type PublishConfigRequest struct {
	AuxFiles       AuxData      `json:"auxFiles"`
	ConfigFiles    ConfigData   `json:"configFiles"`
	ConfigUID      string       `json:"configUID"`
	ExternalId     ExternalId   `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	IgnoreConflict bool         `json:"ignoreConflict"`
	UpdateTime     time.Time    `json:"updateTime"`
	ValidateConfig bool         `json:"validateConfig"`
}

// PublishConfigResponse represents the outcome of a publish config.
type PublishConfigResponse struct {
	DeploymentUID string   `json:"deploymentUID"`
	Links         SelfLinks `json:"links"`
	Result        string   `json:"result"`
}

// SavedNginxConfigSummary defines configuration summary data.
type SavedNginxConfigSummary struct {
	AuxFiles      AuxData      `json:"auxFiles"`
	ConfigFiles   ConfigData   `json:"configFiles"`
	CreateTime    time.Time    `json:"createTime"`
	ExternalId    ExternalId   `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	Instances     []string    `json:"instances"`
	Name          string      `json:"name"`
	UpdateTime    time.Time    `json:"updateTime"`
}

// ListConfigsResponse represents a list of saved NGINX configuration summaries.
type ListConfigsResponse struct {
	PaginationResponse
	Items []SavedNginxConfigSummary `json:"items"`
}

// AuxUploadResult represents details of aux file actions and results for an upload.
type AuxUploadResult struct {
	FileName string `json:"fileName"`
	Result   string `json:"result"`
	Size     int    `json:"size"`
}

// ConfigUploadResult represents details of config file actions and results for an upload.
type ConfigUploadResult struct {
	FileName string `json:"fileName"`
	Result   string `json:"result"`
	Size     int    `json:"size"`
}

// UploadResult represents details of file actions and results for different types of uploads.
type UploadResult struct {
	AuxFiles  []AuxUploadResult    `json:"auxFiles"`
	ConfigFiles []ConfigUploadResult `json:"configFiles"`
}

// UploadConfigRequest represents a configuration upload request.
type UploadConfigRequest struct {
	AuxFiles       AuxData      `json:"auxFiles"`
	ConfigFiles    ConfigData   `json:"configFiles"`
	ConfigName     string       `json:"configName"`
	ExternalId     ExternalId   `json:"externalId"`
	ExternalIdType ExternalIdType `json:"externalIdType"`
	InstGroupUID   string       `json:"instGroupUid"`
	NginxUID       string       `json:"nginxUid"`
}

// UploadConfigResponse defines staged configuration upload results.
type UploadConfigResponse struct {
	Results UploadResult `json:"results"`
	Summary SavedNginxConfigSummary `json:"summary"`
}

