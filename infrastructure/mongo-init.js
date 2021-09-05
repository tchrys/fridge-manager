conn = new Mongo();
db = conn.getDB('mydb');
db.createCollection('recipes');
db.recipes.createIndex({ url: 1}, { unique: true})
