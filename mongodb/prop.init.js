db = db.getSiblingDB('property-db')

db.createUser(
    {
        user : "admin",
        pwd : "pass",
        roles : [
            {
                role : "readWrite",
                db : "property-db"
            }
        ]
    }
)

db = db.getSiblingDB('chat-db')

db.createUser(
    {
        user : "admin",
        pwd : "pass",
        roles : [
            {
                role : "readWrite",
                db : "chat-db"
            }
        ]
    }
)


