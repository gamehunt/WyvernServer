package domain

import "wyvern/server/internal/types"

type OverrideType int

const(
	RoleOverride   OverrideType = iota
	MemberOverride
)

type PermissionOverride struct {
	Id        types.ID      `bson:"id"`
	TargetId  types.ID      `bson:"target_id"`
	ChannelId types.ID      `bson:"channel_id"`
	Type      OverrideType  `bson:"type"`
	Allow     Permission    `bson:"allow"`
	Deny      Permission    `bson:"deny"`
}
