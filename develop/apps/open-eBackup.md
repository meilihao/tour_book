# [open-eBackup](https://gitcode.com/eBackup/open-eBackup)
ref:
- [open-eBackup 1.0.6 Release](https://mp.weixin.qq.com/s/UHyeJ-jdnHXxXRiZ1seWcQ)

date: 2024.10.9 [文档中的路径](https://gitcode.com/eBackup/open-eBackup/blob/master/doc/quick_guide/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8.md)和开源代码的布局以及各组件的路径还是存在差异

## src
- src/AppPlugins : 生态接入层
- src/ProtectAgent : 备份客户端插件化框架, 常驻进程

## env
- jdk 1.8 by `src/Infrastructure_OM/infrastructure/script/commParam.sh:jdk_version="jre1.8.0_422"`

## demo
- [华为HCIE实验-eBackup安装与备份](https://blog.13x.cc/post/ebackup145.html/)

    FusionSphere
- [下载eBackup镜像模板](https://support.huaweicloud.com/intl/zh-cn/hcbkp-fg-cbr/cbr_03_0059.html)

    可注册账户下载eBackup vm, 是v8.0.0, 但[`OceanStor BCManager`](https://support.huawei.com/enterprise/en/flash-storage/oceanstor-bcmanager-pid-21597093/soft)最新是v8.6.x

## ProtectAgent
main在[Agent/src/src/agent/Agent.cpp](https://gitcode.com/eBackup/open-eBackup/blob/master/src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/Agent.cpp)

信息来源:
- [`ReceiveThreadFunc` by fastcgi](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/Communication.cpp)

    pInstance.PushReqMsgQueue(stReqMsg)
- [`FTExceptionHandle::PushUnFreezeReq(MONITOR_OBJ& pMonitorObj)`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/FTExceptionHandle.cpp)

    Communication::GetInstance().PushReqMsgQueue(stReqMsg)
- [`mp_int32 FTExceptionHandle::PushQueryStatusReq(MONITOR_OBJ& pMonitorObj)`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/FTExceptionHandle.cpp)
    Communication::GetInstance().PushReqMsgQueue(stReqMsg)
- [`CMpThread::Create(&recvTid, RecvThread, this);`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/message/tcp/TCPClientHandler.cpp)-> [`mp_int32 TCPClientHandler::HandleRecvDppMessage(CDppMessage &message)`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/message/tcp/TCPClientHandler.cpp)

    msgList.push_back(&message);
    MessageHandler::GetInstance().PushReqMsg(msgPair)
- [BusinessClient::HeartBeat()](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/message/tcp/BusinessClient.cpp)

    MessageHandler::GetInstance().PushRspMsg(busiPair)

消息处理:
- [`mp_int32 iRet = (this->*iter->second)(*message);`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/message/tcp/TCPClientHandler.cpp)

    `TCPClientHandler::HandleMessage()`处理msgList
- [`TaskDispatchWorker::DispacthProc(void* pThis)`->`pDispathWorker->PushMsgToWorker(msg)`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/TaskDispatchWorker.cpp)

    都是先Front(), 处理, 再Pop():
    1. MessageHandler::GetInstance().GetFrontReqMsg()->pDispathWorker->PushMsgToWorker(msg)->PopReqMsg()
    1. Communication::GetInstance().GetFrontReqMsgQueue(msg)->pDispathWorker->PushMsgToWorker(msg)->PopReqMsgQueue()

    pDispathWorker->PushMsgToWorker(msg)->pCurrWorker->PushRequest(msg)->`TaskWorker::ReqProc()`:
    - DispatchRestMsg() by REQMESSAGE_TYPE

        pPlugin = GetPlugin(strService) + pPlugin->Invoke(requestMsg, responseMsg)

        > pPlugin由`TaskPool::Init()`的`CreatePlgConfParse()`+`CreatePlugMgr`加载
    - DispatchTcpMsg() by DPPMESSAGE_TYPE

### plugin
**注意**:
1. 调用plugin的源头(ProtectManager的实现部分, 其仅开源了部分相关接口和定义)没有开源

构建的插件列表见: `src/ProtectAgent/component/protectagent/protectagent/Agent/CMakeLists.txt`的`lib<xxx>-${AGENT_BUILD_NUM}.so`, 但`libdevice-${AGENT_BUILD_NUM}.so`没有对应代码

其他查找插件的方法:
1. 找CMakeLists.txt是否包含`add_library(plugins`
1. 配置文件`src/ProtectAgent/component/protectagent/protectagent/Agent/conf/pluginmgr.xml`(已涵盖`src/ProtectAgent/component/protectagent/protectagent/Agent/conf/backup/pluginmgr.xml`)

根据`class CServicePlugin : public IPlugin`的`DoAction()`, 所有插件均实现了该接口.

`src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/plugins/appprotect/AppProtectPlugin.cpp`的`AppProtectPlugin::InitializeExternalPlugMgr()`又扩展了ExternalPluginManager, 支持`src/AppPlugins_NAS`+`src/AppPlugins_virtualization`

#### vmwarenative
plugin是`src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/plugins/vmwarenative`, 它会调用其实现`src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/apps/vmwarenative`.
其实就是`plugins/vmwarenative`-> `src/apps/vmwarenative`, 以VMwareNativeInitVddkLibTask举例:
1. `VMwareNativeBackupPlugin::InitVMwareNativeVddkLib(CDppMessage &reqMsg, CDppMessage &rspMsg)`
1. `RunSyncTask<VMwareNativeInitVddkLibTask>(msgBody, taskId, respMsg)`

    1. `T* task = new (std::nothrow) T(taskId)`

        1. CreateTaskStep(): 添加step
    1. RunSyncSubTask(taskId, taskName, taskInfo, msgBody, task)
    1. task->InitTaskStep(params[MANAGECMD_KEY_BODY]) = `VMwareNativeInitVddkLibTask::InitTaskStep(const Json::Value &param)`

        1. step->Init(stepParam)
    1. task->RunTask(respMsg)

        1. Task::RunTask(Json::Value &respMsg): 执行step
            1. Task::DoTaskStep(TaskStep* taskStep)
                1. taskStep->Run() = `TaskStepInitVMwareDiskLib::Run()`

VMwareNativeBackupPlugin调用顺序:
1. MANAGE_CMD_NO_VMWARENATIVE_INIT_VDDKLIB
1. 业务

    备份:
    - MANAGE_CMD_NO_VMWARENATIVE_PREPARE_BACKUP
    - MANAGE_CMD_NO_VMWARENATIVE_OPENDISK_BACKUP
    - MANAGE_CMD_NO_VMWARENATIVE_RUN_BACKUP
    - MANAGE_CMD_NO_VMWARENATIVE_QUERY_BACKUP_PROGRESS
    - MANAGE_CMD_NO_VMWARENATIVE_FINISH_DISK_BACKUP
    - MANAGE_CMD_NO_VMWARENATIVE_FINISH_BACKUP_TASK
    - MANAGE_CMD_NO_VMWARENATIVE_CANCEL_BACKUP_TASK
    - MANAGE_CMD_NO_VMWARENATIVE_CLEANUP_RESOURCES

    还原:
    - MANAGE_CMD_NO_VMWARENATIVE_PREPARE_RECOVERY
    - MANAGE_CMD_NO_VMWARENATIVE_RUN_RECOVERY
    - MANAGE_CMD_NO_VMWARENATIVE_QUERY_RECOVERY_PROGRESS
    - MANAGE_CMD_NO_VMWARENATIVE_FINISH_DISK_RECOVERY
    - MANAGE_CMD_NO_VMWARENATIVE_FINISH_RECOVERY_TASK
    - MANAGE_CMD_NO_VMWARENATIVE_CANCEL_RECOVERY_TASK

1. MANAGE_CMD_NO_VMWARENATIVE_CLEANUP_VDDKLIB
1. 辅助

    vmfs:
    1. MANAGE_CMD_NO_VMWARENATIVE_CHECK_VMFS_TOOL
    1. MANAGE_CMD_NO_VMWARENATIVE_VMFS_MOUNT
    1. MANAGE_CMD_NO_VMWARENATIVE_VMFS_UMOUNT

    nas:
    1. MANAGE_CMD_NO_VMWARENATIVE_SLNAS_MOUNT
    1. MANAGE_CMD_NO_VMWARENATIVE_SLNAS_UMOUNT

    iscsi:
    1. MANAGE_CMD_NO_HOST_LINK_ISCSI

    dataprocess:
    1. MANAGE_CMD_NO_VMWARENATIVE_ALLDISK_AFS_BITMAP

        1. PrepareAllDisksAfsBitmap
            1. TaskStepPrepareAfsBitmap::Run()
                1. TaskStepVMwareNative::DataProcessLogic
                    1. ExchangeMsgWithDataProcessService(param, respParam, reqCmd, rspCmd, m_internalTimeout)

                        1. SendDPMessage()
                        1. VMwareNativeDataPathProcess::HandleDisksAfsBitmap # by CMD_VMWARENATIVE_BACKUP_AFS

dataprocess处理过程:
1. VMwareNativeDataPathProcess::ExtCmdProcess(CDppMessage &message)
    1. VMwareNativeDpCmdParse(message, currentTaskID, parentTaskID)
        1. handerThread->InsertTask()
        1. VMwareNativeDataPathProcess::HandleDisksAfsBitmap()
            处理入口: HanderThread::StartThread()->`std::make_unique<std::thread>(&HanderThread::FuncHander, this)` # 处理CDppMessage

            处理函数由VMwareNativeDataPathProcess::GenerateCmdMsgHandlerFunMapBackup()注册
        1. VMwareNativeDataPathImpl::VMwareNativeBackupAfsBitmap()


## FAQ
ProtectAgent缺失的`securec.h`, 大概是[platform/huaweisecurec/include/securec.h](https://github.com/huaweicloud/huaweicloud-sdk-c-obs/blob/master/platform/huaweisecurec/include/securec.h)

[src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/securecom/SDPFunc.cpp]()部分缺失的kmc header在[gitee.com/openeuler/cantian/tree/master/library/kmc](https://gitee.com/openeuler/cantian/tree/master/library/kmc), 但未找到其kmc的源码.