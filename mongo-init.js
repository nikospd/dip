db = db.getSiblingDB('resources')
db.users.ensureIndex({"email": 1}, {"unique": true});
db.users.ensureIndex({"username": 1}, {"unique": true});
db.filters.ensureIndex({"storage_id": 1}, {"unique": true});
