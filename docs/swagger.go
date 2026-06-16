// Package docs provides API documentation annotations for Swagger generation.
package docs

// API Documentation Annotations
// These comments are used by swaggo to generate OpenAPI documentation.

// @title Admin System API
// @version 1.0
// @description A production-ready admin panel built with Gai Framework.
// @description Provides user management, role-based access control, and system administration.

// @contact.name API Support
// @contact.url https://github.com/user/admin-system

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Authentication
// @tag.description Login, logout, and token management

// @tag.name Users
// @tag.description User management operations

// @tag.name Roles
// @tag.description Role management operations

// @tag.name Permissions
// @tag.description Permission management operations

// @tag.name Menus
// @tag.description Menu management operations

// @tag.name Logs
// @tag.description Operation log viewing
type SwaggerDocs struct{}
