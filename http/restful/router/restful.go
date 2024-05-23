package router

import (
	"strings"

	"github.com/agclqq/prow-framework/http/restful"

	"github.com/gin-gonic/gin"
)

type RestfulRouteName int
type RestfulEntry struct {
	Middlewares []gin.HandlerFunc
	Only        []RestfulRouteName
	Exclude     []RestfulEntry
}

const (
	INDEX RestfulRouteName = iota
	CREATE
	STORE
	SHOW
	EDIT
	UPDATE
	DESTROY
)

func ApiResource(fatherGroup *gin.RouterGroup, relativePath string, resource restful.ApiResource, only ...RestfulRouteName) {
	relativePath = strings.TrimRight(relativePath, "/")
	ps := strings.Split(relativePath, "/")
	lastPath := ps[len(ps)-1]
	group := fatherGroup.Group(relativePath)
	{
		if len(only) > 0 {
			for _, v := range only {
				switch v {
				case INDEX:
					group.GET("", resource.Index) //列表页
				case STORE:
					group.POST("", resource.Store) //上传，保存
				case SHOW:
					group.GET(":"+lastPath, resource.Show) //单资源查询
				case UPDATE:
					group.PUT(":"+lastPath, resource.Update) //更新
				case DESTROY:
					group.DELETE(":"+lastPath, resource.Destroy) //删除
				}

			}
		} else {
			group.GET("", resource.Index)                //列表页
			group.POST("", resource.Store)               //上传，保存
			group.GET(":"+lastPath, resource.Show)       //单资源查询
			group.PUT(":"+lastPath, resource.Update)     //更新
			group.DELETE(":"+lastPath, resource.Destroy) //删除
		}
	}
}
func Resource(fatherGroup *gin.RouterGroup, relativePath string, resource restful.Resource, only ...RestfulRouteName) {
	ApiResource(fatherGroup, relativePath, resource, only...)
	relativePath = strings.TrimRight(relativePath, "/")
	ps := strings.Split(relativePath, "/")
	lastPath := ps[len(ps)-1]
	group := fatherGroup.Group(relativePath)
	{
		if len(only) > 0 {
			for _, v := range only {
				switch v {
				case CREATE:
					group.GET("/create", resource.Create) //创建页
				case EDIT:
					group.PUT("/:"+lastPath+"/edit", resource.Edit) //更新
				}

			}
		} else {
			group.GET("/create", resource.Create)           //创建页
			group.PUT("/:"+lastPath+"/edit", resource.Edit) //更新
		}
	}
}
