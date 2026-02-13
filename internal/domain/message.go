package domain

import "wyvern/server/internal/types"

type Message struct {
	Id        types.ID `bson:"id"`
	ChannelId types.ID `bson:"channel_id"`
	AuthorId  types.ID `bson:"author_id"`
	Contents  []byte   `bson:"contents"`
}
