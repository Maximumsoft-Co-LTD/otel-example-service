function oid() { return ObjectId() }

const now = new Date();

const adminRoleId = oid();
const userRoleId  = oid();

db.roles.insertMany([
    {
        _id: adminRoleId,
        slug: "admin",
        name: "Administrator",
        permissions: [
            "auth.login",
            "auth.logout",
            "auth.refresh",
            "user.create",
            "user.update",
            "user.delete",
            "user.list"
        ],
        created_at: now,
        updated_at: now
    },
    {
        _id: userRoleId,
        slug: "user",
        name: "User",
        permissions: [
            "auth.login",
            "auth.logout",
            "auth.refresh",
            "profile.view"
        ],
        created_at: now,
        updated_at: now
    }
]);

print("\n=== Inserted Roles ===");
print("AdminRoleID:", adminRoleId);
print("UserRoleID :", userRoleId);

const hashedPassword = "$2a$10$qVzqVxz9kC9DqUCtkjZKpOMnZQub2fXLdR6JyH5ewwAstZg5tr1nK";

db.users.insertMany([
    {
        _id: oid(),
        name: "Administrator",
        username: "admin",
        password: hashedPassword,
        email: "admin@example.com",
        roles: [adminRoleId],
        email_verified_at: now,
        active: true,
        created_at: now,
        updated_at: now
    },
    {
        _id: oid(),
        name: "Normal User",
        username: "user",
        password: hashedPassword,
        email: "user@example.com",
        roles: [userRoleId],
        email_verified_at: now,
        active: true,
        created_at: now,
        updated_at: now
    }
]);

print("\n=== Inserted Users Success ===\n");