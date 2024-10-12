# [open-eBackup](https://gitcode.com/eBackup/open-eBackup)
date: 2024.10.9 [文档中的路径](https://gitcode.com/eBackup/open-eBackup/blob/master/doc/quick_guide/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8.md)和开源代码的布局以及各组件的路径还是存在差异

## env
- jdk 1.8 by `src/Infrastructure_OM/infrastructure/script/commParam.sh:jdk_version="jre1.8.0_422"`

## demo
- [华为HCIE实验-eBackup安装与备份](https://blog.13x.cc/post/ebackup145.html/)

    FusionSphere

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
构建的插件列表见: `src/ProtectAgent/component/protectagent/protectagent/Agent/CMakeLists.txt`的`lib<xxx>-${AGENT_BUILD_NUM}.so`, 但`libdevice-${AGENT_BUILD_NUM}.so`没有对应代码

其他查找插件的方法:
1. 找CMakeLists.txt是否包含`add_library(plugins`
1. 配置文件`src/ProtectAgent/component/protectagent/protectagent/Agent/conf/pluginmgr.xml`(已涵盖`src/ProtectAgent/component/protectagent/protectagent/Agent/conf/backup/pluginmgr.xml`)

根据`class CServicePlugin : public IPlugin`的`DoAction()`, 所有插件均实现了该接口.

`src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/plugins/appprotect/AppProtectPlugin.cpp`的`AppProtectPlugin::InitializeExternalPlugMgr()`又扩展了ExternalPluginManager, 支持`src/AppPlugins_NAS`+`src/AppPlugins_virtualization`

## FAQ
ProtectAgent缺失的`securec.h`, 大概是[platform/huaweisecurec/include/securec.h](https://github.com/huaweicloud/huaweicloud-sdk-c-obs/blob/master/platform/huaweisecurec/include/securec.h)

[src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/securecom/SDPFunc.cpp]()部分缺失的kmc header在[gitee.com/openeuler/cantian/tree/master/library/kmc](https://gitee.com/openeuler/cantian/tree/master/library/kmc), 但未找到其kmc的源码.