package restful

import "github.com/gin-gonic/gin"

type ApiResource interface {
	Index(*gin.Context)
	Store(*gin.Context)
	Show(*gin.Context)
	Update(*gin.Context)
	Destroy(*gin.Context)
}
type Resource interface {
	ApiResource
	Create(*gin.Context)
	Edit(*gin.Context)
}
