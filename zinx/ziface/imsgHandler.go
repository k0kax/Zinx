package ziface

/*
	消息管理的抽象层
*/

type IMsgHandle interface {
	//调度执行对应的Router消息处理方法
	DoMsgHandle(request IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	//启动一个worker工作池 开启工作池的动作只能发生一次，一个zinx框架只能存在一个worker工作池
	StartWorkerPool()

	//将消息发送给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
