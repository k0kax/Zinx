package znet

import "Zinx/zinx/ziface"

/*
定义router时先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写就好了,暂时不用写具体的实现方法，别的方法继承重写它即可
*/
type BaseRouter struct{}

// 模板设计模式
// 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(requset ziface.IRequest) {}

// 在处理conn业务的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
