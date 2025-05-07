package astorage

import "github.com/kataras/iris/v12"

// RegisterClientStorageApi 需要进客户端端对服务端的认证
func (g *AStorage) RegisterClientStorageApi(p iris.Party) {
	p.Get("/astorage/data/{k:string}/{uid:uint64}", g.getClientData) // 单个数据
	p.Get("/astorage/data", g.getClientDataByUid)                    // 数组
	p.Delete("/astorage/data/{k:string}/{uid:uint64}", g.deleteClientData)
	p.Post("/astorage/data", g.postClientData)
}

// RegisterServiceStorageApi 需要进行服务端对服务端的认证
func (g *AStorage) RegisterServiceStorageApi(p iris.Party) {
	p.Get("/astorage/sdata/{k:string}", g.getConfigOfService) // 单个数据
	p.Get("/astorage/sdata", g.getConfigsOfService)           // 数组
	p.Delete("/astorage/sdata/{k:string}", g.deleteConfigOfService)
	p.Post("/astorage/sdata", g.postConfigOfService)
}
