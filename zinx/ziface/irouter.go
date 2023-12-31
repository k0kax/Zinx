package ziface

/*
	路由的抽象接口，
	路由里的数据都是request了
*/

type IRouter interface {
	//在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)

	//在处理conn业务的主方法Hook
	Handle(requset IRequest)

	//在处理conn业务的钩子方法Hook
	PostHandle(request IRequest)
}
