package domain

type Permission uint64

const(
	ViewChannel  Permission = 1 << iota
	SendMessages 
	ManageMessages
	ManageChannels
	ManageRoles
	ManageGuild
	KickMembers
	BanMembers
	Administrator
)

func (s Permission) Has(p Permission) bool {
	return (s & p) != 0
}

func (s Permission) Add(p Permission) Permission {
	return s | p;
}

func (s Permission) Remove(p Permission) Permission {
	return s & ^p
}
