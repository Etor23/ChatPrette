db = db.getSiblingDB('chat_app');

// ==================== COLECCIONES ====================
db.createCollection('users');
db.createCollection('conversations');
db.createCollection('messages');

// ==================== ÍNDICES: users ====================
db.users.createIndex({ "firebase_uid": 1 }, { unique: true });
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });

// ==================== ÍNDICES: conversations ====================
db.conversations.createIndex({ "members": 1 });
db.conversations.createIndex({ "lastMessageAt": -1 });

// ==================== ÍNDICES: messages ====================
db.messages.createIndex({ "conversationId": 1, "createdAt": -1 });

// ==================== DATOS DE PRUEBA (opcional) ====================
db.users.insertOne({
    _id: ObjectId(),
    firebase_uid: "test_firebase_uid_123",
    email: "test@mail.com",
    username: "testuser",
    displayName: "Test User",
    avatarUrl: null,
    createdAt: new Date(),
    updatedAt: new Date()
});

print('');
print('Base de datos chat_app inicializada');
print('Colecciones: users, conversations, messages');
print('Índices creados');
print('');