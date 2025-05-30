package bootstrap

import (
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/service"
)

func Init(port int) error {
	api := service.NewRestApiService()
	return api.ServeContent(port)
}
