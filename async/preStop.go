package async

// PreStop async的任务统一在这里close chan，否则go线程永远不会停止
func PreStop() {
	close(MailChan)
	close(ChatRecordChan)
}
