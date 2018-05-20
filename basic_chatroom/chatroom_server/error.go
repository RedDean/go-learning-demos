package main

// 错误处理
func HandleError(err error)  {
	if err!=nil {
		panic(err)
	}
	return
}