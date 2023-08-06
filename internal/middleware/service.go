package middleware

import "github.com/7Maliko7/forms-api/internal/service"

// Middleware describes a service middleware.
type Middleware func(service service.Service) service.Service
