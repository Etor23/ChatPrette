package dto

type CreateDmRequest struct {
	OtherUserID string `json:"other_user_id" binding:"required"`
}

type CreateGroupRequest struct {
	Name	string   `json:"name" binding:"required"`
	Members []string `json:"members" binding:"required,min=1"`
}

type UpdateGroupNameRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

type AddGroupMembersRequest struct {
	Members []string `json:"members" binding:"required,min=1"`
}

type RemoveGroupMembersRequest struct {
	Members []string `json:"members" binding:"required,min=1"`
}

type ConversationResponse struct {
	ID					string   	`json:"id"`
	Type				string   	`json:"type"`
	Members 			[]string 	`json:"members"`
	Name				string   	`json:"name,omitempty"`
	LastMessageAt 		*string 	`json:"last_message_at,omitempty"`
	LastMessagePreview 	string 		`json:"last_message_preview,omitempty"`
}
