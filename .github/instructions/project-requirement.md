# Project Title: GoCMS API
Goal: To build a robust and scalable backend API for a Content Management System, providing endpoints for managing content (articles, pages, etc.), users, roles, and media files.
Core Features (APIs to be built):

## Content Management:
- Create New Content (e.g., Article, Page)
- Get Content by ID
- Get List of Content (with filtering, pagination, sorting)
- Update Existing Content
- Delete Content
- Publish/Unpublish Content
- Manage Content Revisions (optional, but a good addition)
- Categorization and Tagging of Content
- Content Versioning (more advanced)



## User Management:
- User Registration (with email verification, optional)
- User Login (using various authentication methods)
- Get User Profile
- Update User Profile
- Delete User
- Password Reset Functionality



## Role and Permissions Management:
- Create New Roles (e.g., Admin, Editor, Author, Guest)
- Assign Roles to Users
- Define Permissions for Roles (e.g., can create content, can edit all content, can manage users)
- Check User Permissions for specific actions



## Media Management:
- Upload Media Files (images, documents, videos - with appropriate validation and storage)
- Get Media File by ID
- Get List of Media Files (with filtering)
- Delete Media File
- Serve Media Files (securely, perhaps with resizing/optimization - more advanced)


# Technical Stack:

Technical Stack:
Language: Go (Golang)
Web Framework/Router: Gin
Database: sqlite
Database Driver/ORM: Gorm
Authentication: JWT
Authorization: Libraries like Casbin
File Storage: Cloud Storage (AWS S3)
