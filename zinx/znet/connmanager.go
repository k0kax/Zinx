package znet

import (
	"Zinx/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
链接管理模块
*/
type ConnManager struct {
	//管理的链接信息
	connections map[uint32]ziface.IConnection
	//读写链接集合保护的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//从conn加入到ConnManager
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID:", conn.GetConnID(), "add to ConnManager successfully:conn num = ", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID:", conn.GetConnID(), "remove to ConnManager successfully:conn num = ", connMgr.Len())

}

// 根据ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		//找到链接
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前连接的总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清楚并终止所有链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear all connections succ!! conn num=", connMgr.Len())
}
