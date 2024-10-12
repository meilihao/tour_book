# [open-eBackup](https://gitcode.com/eBackup/open-eBackup)
date: 2024.10.9 [文档中的路径](https://gitcode.com/eBackup/open-eBackup/blob/master/doc/quick_guide/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8.md)和开源代码的布局以及各组件的路径还是存在差异

## ProtectAgent
main在[Agent/src/src/agent/Agent.cpp](https://gitcode.com/eBackup/open-eBackup/blob/master/src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/Agent.cpp)

消息处理:
- [`mp_int32 iRet = (this->*iter->second)(*message);`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/message/tcp/TCPClientHandler.cpp)
- [`TaskDispatchWorker::DispacthProc(void* pThis)`->`pDispathWorker->PushMsgToWorker(msg)`](src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/agent/TaskDispatchWorker.cpp)

## FAQ
ProtectAgent缺失的`securec.h`, 大概是[platform/huaweisecurec/include/securec.h](https://github.com/huaweicloud/huaweicloud-sdk-c-obs/blob/master/platform/huaweisecurec/include/securec.h)

[src/ProtectAgent/component/protectagent/protectagent/Agent/src/src/securecom/SDPFunc.cpp]()部分缺失的kmc header在[gitee.com/openeuler/cantian/tree/master/library/kmc](https://gitee.com/openeuler/cantian/tree/master/library/kmc), 但未找到其kmc的源码.