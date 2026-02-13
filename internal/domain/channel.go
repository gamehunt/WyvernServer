package domain

import "wyvern/server/internal/types"

type ChannelType int 

const (
	DM        ChannelType = iota
	GroupDM
	GuildText  
	GuildVoice
	GuildCategory
)

type Channel struct {
	Id         types.ID    `bson:"id"`
	Type       ChannelType `bson:"type"`
	GuildId    *types.ID   `bson:"guild_id,omitempty"`
	Position   *int        `bson:"position,omitempty"`
	OwnerId    *types.ID   `bson:"owner_id,omitempty"`
	Name       *string     `bson:"name,omitempty"`
	ParentId   *types.ID   `bson:"parent_id,omitempty"`
	Recipients []types.ID  `bson:"recipients,omitempty"`
}
