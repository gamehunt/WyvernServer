db.createCollection('users')
db.users.createIndex({"id": 1}, {unique: true})
db.users.createIndex({"identity": 1}, {unique: true})

db.createCollection('sessions')
db.sessions.createIndex({"id": 1}, {unique: true})

db.createCollection('guilds')
db.guilds.createIndex({"id": 1}, {unique: true})

db.createCollection('channels')
db.channels.createIndex({"id": 1}, {unique: true})
db.channels.createIndex({"guild_id": 1})
db.channels.createIndex({"parent_id": 1})

db.createCollection('members')
db.members.createIndex({ "user_id":  1, "guild_id": 1 }, { unique: true })
db.members.createIndex({ "guild_id": 1 })
db.members.createIndex({ "user_id":  1 })
