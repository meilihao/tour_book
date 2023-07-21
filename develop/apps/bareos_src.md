# bareos_src
version: 22.1.0

## filed
ref:
- [Full connection overview](https://docs.bareos.org/TasksAndConcepts/NetworkSetup.html)

1. filed.cc

	```c++
	int main(int argc, char* argv[])
	{
		...
		LoadFdPlugins(me->plugin_directory, me->plugin_names); // 加载plugin
		...
	  	// if configured, start threads and connect to Director.
	  	StartConnectToDirectorThreads(); // 连接到dir

	  	// start socket server to listen for new connections.
	  	StartSocketServer(me->FDaddrs);  // 等待连接, 核心逻辑
	  	...
	}
	```

	StartConnectToDirectorThreads(主动连接到dir): handle_connection_to_director->StartProcessDirectorCommands->process_director_commands(void* p_jcr)->process_director_commands(JobControlRecord* jcr, BareosSocket* dir)->`if (!cmds[i].func(jcr)) { /* do command */`

	StartSocketServer: HandleConnectionRequest->handle_director_connection(处理来自dir的连接)->process_director_commands(JobControlRecord* jcr, BareosSocket* dir)/handle_stored_connection(处理来自sd的连接)->