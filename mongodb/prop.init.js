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


