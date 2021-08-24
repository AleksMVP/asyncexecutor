package asyncexecutor

type IAsyncExecutor interface {
	Start() 		     error 
	Stop()               error
	Wait()               error
	AddTask(task func()) error
}