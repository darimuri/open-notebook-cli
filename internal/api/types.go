package api

// Config response
type ConfigResponse struct {
	Version string `json:"version"`
	Health   bool   `json:"health"`
}

// Notebook types
type NotebookCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type NotebookUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Archived    *bool   `json:"archived,omitempty"`
}

type NotebookResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Archived    bool   `json:"archived"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	SourceCount int    `json:"source_count"`
	NoteCount   int    `json:"note_count"`
}

type NotebookDeletePreview struct {
	NotebookID         string `json:"notebook_id"`
	NotebookName       string `json:"notebook_name"`
	NoteCount          int    `json:"note_count"`
	ExclusiveSourceCount int   `json:"exclusive_source_count"`
	SharedSourceCount  int    `json:"shared_source_count"`
}

type NotebookDeleteResponse struct {
	Message       string `json:"message"`
	DeletedNotes  int    `json:"deleted_notes"`
	DeletedSources int   `json:"deleted_sources"`
	UnlinkedSources int  `json:"unlinked_sources"`
}

// Note types
type NoteCreate struct {
	NotebookID string `json:"notebook_id"`
	Content    string `json:"content"`
	Metadata   string `json:"metadata,omitempty"`
}

type NoteUpdate struct {
	Content  *string `json:"content,omitempty"`
	Metadata *string `json:"metadata,omitempty"`
}

type NoteResponse struct {
	ID         string `json:"id"`
	NotebookID string `json:"notebook_id"`
	Content    string `json:"content"`
	Metadata   string `json:"metadata"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

// Source types
type SourceResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Size     int    `json:"size"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

type SourceCreate struct {
	Notebooks []string `json:"notebooks,omitempty"`
	Type     string   `json:"type"`
	URL      string   `json:"url,omitempty"`
	Content  string   `json:"content,omitempty"`
	Embed    *bool    `json:"embed,omitempty"`
}

type SourceStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Embed types
type EmbedRequest struct {
	ItemID          string `json:"item_id"`
	ItemType        string `json:"item_type"`
	AsyncProcessing bool   `json:"async_processing"`
}

type EmbedResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ItemID  string `json:"item_id"`
	ItemType string `json:"item_type"`
}

// Search types
type SearchRequest struct {
	Query     string   `json:"query"`
	NotebookIDs []string `json:"notebook_ids,omitempty"`
	Limit     int      `json:"limit,omitempty"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

type SearchResult struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	NotebookID string `json:"notebook_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Score      float64 `json:"score"`
}

// Ask types
type AskRequest struct {
	Question          string   `json:"question"`
	StrategyModel     string   `json:"strategy_model"`
	AnswerModel       string   `json:"answer_model"`
	FinalAnswerModel  string   `json:"final_answer_model"`
	NotebookIDs       []string `json:"notebook_ids,omitempty"`
}

type AskResponse struct {
	Answer    string   `json:"answer"`
	Question  string   `json:"question"`
	Sources   []string `json:"sources,omitempty"`
	ModelUsed string   `json:"model_used,omitempty"`
}

// Job types
type JobResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	Result    any    `json:"result,omitempty"`
}

// Model types
type ModelResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Type       string `json:"type"`
	Credential string `json:"credential"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

type DefaultModelsResponse struct {
	DefaultChatModel            string `json:"default_chat_model"`
	DefaultTransformationModel  string `json:"default_transformation_model"`
	LargeContextModel          string `json:"large_context_model"`
	DefaultTextToSpeechModel    string `json:"default_text_to_speech_model"`
	DefaultSpeechToTextModel    string `json:"default_speech_to_text_model"`
	DefaultEmbeddingModel      string `json:"default_embedding_model"`
	DefaultToolsModel          string `json:"default_tools_model"`
}

type CredentialResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`
	Modalities      []string `json:"modalities"`
	BaseURL         string   `json:"base_url"`
	Endpoint        string   `json:"endpoint"`
	APIVersion      string   `json:"api_version"`
	EndpointLLM    string   `json:"endpoint_llm"`
	EndpointEmbed  string   `json:"endpoint_embedding"`
	EndpointSTT     string   `json:"endpoint_stt"`
	EndpointTTS     string   `json:"endpoint_tts"`
	Project         string   `json:"project"`
	Location        string   `json:"location"`
	CredentialsPath string   `json:"credentials_path"`
	HasAPIKey       bool     `json:"has_api_key"`
	Created         string   `json:"created"`
	Updated         string   `json:"updated"`
	ModelCount      int      `json:"model_count"`
	DecryptionError string   `json:"decryption_error"`
}