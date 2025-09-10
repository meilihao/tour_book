# airflow
ref:
- [Airflow从零到神](https://www.bilibili.com/video/BV19f4y1V7UG)
- [**Airflow 实践笔记-从入门到精通**](https://zhuanlan.zhihu.com/p/517364346)
- [airflow 2.10.5 翻译文档](https://www.aidoczh.com/airflow/index.html)

ver: v3.0.4

一个可编程, 调度和监控的工作流平台, 核心是DAG工作流

> 竞品: [**DolphinScheduler**](https://dolphinscheduler.apache.org/zh-cn/), XXL-Job

## [安装](https://airflow.apache.org/docs/apache-airflow/stable/start.html)

> `uv pip install "apache-airflow==${AIRFLOW_VERSION}" --constraint "${CONSTRAINT_URL}"`可能卡住, 换`uv -v pip install ...`排查即可

`airflow standalone`初始密码在`AIRFLOW_HOME/simple_auth_manager_passwords.json.generated`里

> 在 Airflow 3.x 中，初始账户的创建方式和旧版本不同，关键取决于使用的 认证方式（AuthManager）

> airflow内部是统一使用UTC来传递参数，本地时间的转换只是方便展示

## 概念
- Data Pipeline：数据管道或者数据流水线，可以理解为贯穿数据处理分析过程中不同工作环节的流程，例如加载不同的数据源，数据加工以及可视化
- DAG：有向非循环图（directed acyclic graphs），可以理解为有先后顺序任务的多个Tasks的组合. 图的概念是由节点组成的，有向的意思就是说节点之间是有方向的，即有依赖关系；非循环的意思就是说节点直接的依赖关系只能是单向的, 不会出现循环依赖. 每个 Dag 都有唯一的 DagId，当一个 DAG 启动的时候，Airflow 都将在数据库中创建一个DagRun记录，相当于一个日志

    DAG是多个脚本处理任务组成的工作流pipeline
- Task：是包含一个具体Operator的对象，operator实例化的时候称为task. DAG图中的每个节点都是一个任务，可以是一条命令行（BashOperator），也可以是一段 Python 脚本（PythonOperator）等，然后这些节点根据依赖关系构成了一个图，称为一个 DAG。当一个任务执行的时候，实际上是创建了一个 Task实例运行，它运行在 DagRun 的上下文中
- Connections：是管理外部系统的连接对象，如外部MySQL、HTTP服务等，连接信息包括conn_id／hostname／login／password／schema等，可以通过界面查看和管理，编排workflow时，使用conn_id进行使用
- Pools: 用来控制tasks执行的并行数。将一个task赋给一个指定的pool，并且指明priority_weight权重，从而干涉tasks的执行顺序
- XComs：在airflow中，operator一般是原子的，也就是它们一般是独立执行，不需要和其他operator共享信息。但是如果两个operators需要共享信息，例如filename之类的，则推荐将这两个operators组合成一个operator；如果一定要在不同的operator实现，则使用XComs (cross-communication)来实现在不同tasks之间交换信息. 在airflow 2.0以后，因为task的函数跟python常规函数的写法一样，operator之间可以传递参数，但本质上还是使用XComs，只是不需要在语法上具体写XCom的相关代码
- [Trigger Rules](https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/dags.html#trigger-rules)：指task的触发条件。默认情况下是task的直接上游执行成功后开始执行，airflow允许更复杂的依赖设置，包括all_success(所有的父节点执行成功)，all_failed(所有父节点处于failed或upstream_failed状态)，all_done(所有父节点执行完成)，one_failed(一旦有一个父节点执行失败就触发，不必等所有父节点执行完成)，one_success(一旦有一个父节点执行成功就触发，不必等所有父节点执行完成)，dummy(依赖关系只是用来查看的，可以任意触发)。另外，airflow提供了depends_on_past，设置为True时，只有上一次调度成功了，才可以触发。

  一般来讲，只有当上游任务“执行成功”时，才会开始执行下游任务. 但是除了“执行成功all_success”这个条件以外，还有其他的trigger rule，例如one_success, one_failed（至少一个上游失败），none_failed ，none_skipped
- Backfill: 可以支持重跑历史任务，例如当ETL代码修改后，把上周或者上个月的数据处理任务重新跑一遍

## 配置
### env
AIRFLOW_HOME 是 Airflow 寻找 DAG 和插件的基准目录.

airlfow的配置是通过`AIRFLOW_HOME`下的airflow.cfg配置文件进行读取的.

### airflow.cfg
> 每个配置项都有注释和相关的官方文档链接

- dags_folder = /data/tmpfs/airflow/dags : dags存放目录
- load_examples : web ui默认加载airflow自带的dag案例
- min_file_process_interval : 同步dags_folder的间隔
- [executor](https://airflow.apache.org/docs/apache-airflow/3.0.4/core-concepts/executor/index.html) : 使用哪种executor

    - 集群executor: CeleryExecutor(推荐rabbitmq)/KubernetesExecutor

## 命令
```bash
airflow dags list # 查看当前所有 dag 任务
airflow tasks list -v example_bash_operator # 查看某个 dag 任务
airflow tasks test example_bash_operator run_after_loop 2015-01-01 # 调试任务, 仅终端有日志, web上没有相关日志
airflow dags list-import-errors # 查看导入错误
```

在手动触发任务的时候，是无法执行还未到执行时间的定时任务的。例如 每天下午3点的任务，如果是今天下午2点30手动触发，实际上跑的是昨天下午3点的脚本。手动点`Trigger`触发，有一个选项`Advanced Options`，点击进去可以选择`Logical date`，例如前天或者大前天的任务；但是如果选择明天的时间，会发现这个任务会处于排队queue的状态

## dag
ref:
- [Home/Core Concepts/Dags](https://airflow.apache.org/docs/apache-airflow/stable/core-concepts/dags.html)
- [配置DAG的参数](https://zhuanlan.zhihu.com/p/517364346)
- [Templates reference](https://airflow.apache.org/docs/apache-airflow/3.0.4/templates-ref.html)

参数(也可通过web ui中Dag的`Details`标签页查看):
- dag_id: DAG 的 dag_id 都是唯一的，重复会报错或覆盖旧的 DAG
- schedule: cron调度rule
- start_date: 开始时间
- end_date: 结束时间, 没有则表示一直运行
- catchup: 是否自动执行历史任务
- tags: 标签

    在dag_tag表里
- params: 设置默认参数
- `>>` : 任务依赖关系

    > 等同set_upstream和set_downstream方法

    `start >> [fetch_weather, fetch_sales] >> join_datasets`: start执行完以后，同时执行fetch_weather和fetch_sales，然后执行join_datasets

- BashOperator : 具体执行的bash任务, 结果为 true 才会走入下一个依赖任务，如果为 false 则忽略
- task_id : 在一个 DAG（即 dag_id）中必须唯一，但在多个 DAG 中可以重复
- bash_command : 具体任务执行命令

DAG在配置的时候，可以配置同时运行的任务数concurrency，默认是16个. 这个16，就是task slot，可以理解为资源，如果资源满了，具备运行条件的task就需要等待

定义DAG的方式有三种:
1. 使用with语法
1. 标准构造函数
1. 使用装饰器@dag.

当Airflow从Python文件加载DAG时，它只会拉取顶层的任何DAG实例对象, 忽略其他嵌套的DAG实例对象.

DAG Run是DAG运行一次的对象（记录），记录所包含任务的状态信息。如果所有的任务状态是success或者skipped，就是success；如果任务有failed或者upstream_failed，就是falied. 其中的run_id的前缀会有如下几个:
- scheduled__ 表明是不是定时的
- backfill__ 表明是不是回填的
- manual__ 表明是不是手动或者trigger的

运行dags方法:
1. 手动
2. 通过api
3. 通过dag定义中的schedule

模板变量:
- task_instance_key_str

    task_instance_key_str 是一个由 dag_id、task_id、execution_date 组成的字符串，用于唯一标识一个任务实例, 常用于在调度器、日志系统、数据库等地方精确标记一个特定的任务执行

    ```py
    also_run_this = BashOperator(
            task_id="also_run_this",
            bash_command='echo "ti_key={{ task_instance_key_str }}"',
        )
    ```

## Operator
Operator的类型有以下几种:
1. DummyOperator

    作为一个虚拟的任务节点，使得DAG有一个起点，但实际不执行任务；或者是在上游几个分支任务的合并节点，为了清楚的现实数据逻辑

1. BashOperator

    当一个任务是执行一个shell命令，就可以用BashOperator。 可以是一个命令，也可以指向一个具体的脚本文件

1. 条件分支判断

    1. BranchDateTimeOperator

        在一个时间段内执行一种任务，否则执行另一个任务。Target_lower可以设置为None

    1. BranchDayOfWeekOperator

        根据是哪一天来选择跑哪个任务

    1. BranchPythonOperator

        根据业务逻辑条件，选择下游的一个task运行

1. PythonOperator

    用的最广泛的Operator

    dags间通信:
    ```py
    # down_dag.py: python_callable = calc
    # down_body.py for down_dag.py
    def handle_dag_param(**kwargs):
    param = kwargs.get("params")
    print(param)

    def calc(**kwargs):
        time_range, handle_mode, owner = handle_dag_param(**kwargs)

    # up_body.py
    kwargs["ti"].xcom_push(
        key="stat_payload",
        value={
            "start": payload["start"],
            "end": payload["end"],
        },
    )
    return time_range, handle_mode, owner


    def deal(**kwargs):
        time_range, handle_mode, owner = handle_dag_param(**kwargs)

    # up_dag.py: python_callable = deal
    up_trigger = TriggerDagRunOperator(
        task_id="up_trigger",
        trigger_dag_id="down_dag",
        conf={"payload": "{{ task_instance.xcom_pull('up_dag', key='stat_payload') }}"},
    )

    up_dag.set_downstream(up_trigger)
    ```
 
1. 图之间依赖关系的operator

    如果两个任务流之间，存在一些依赖关系

    使用ExternalTaskSensor，根据另一个DAG中的某一个任务的执行情况来决定是否执行当前任务

1. LatestOnlyOperator

    LatestOnlyOperator，是为了标识该DAG是不是最新的执行时间，只有在最新的时候才有必要执行下游任务，例如部署模型的任务，只需要在最近一次的时间进行部署即可

1. Sensor

    Sensor 是用来判断外部条件是否成熟的感应器，例如判断输入文件是否到位（可以设置一个时间窗口内，例如到某个时间点之前检查文件是否到位），但是sensor很耗费计算资源（设置mode为reschedule可以减少开销，默认是poke），DAG会设置concurrency约定同时最多有多少个任务可以运行，称为task slot. 所以一种办法是使用Deferrable Operators

    Deferrable Operators是sensor operator的替代品，其本质是应用了python的异步机制，有一个trigger process，当其符合特定条件后，就会将任务调度起来. Deferrable Operators相比Sensor Operator很节省资源

1. 自定义Operator

    Hook是一种自定义的operator

### PythonOperator
```py3
    def my_sleeping_function(random_base):
        """This is a function that will run within the DAG execution"""
        time.sleep(random_base)

    for i in range(5):
        sleeping_task = PythonOperator(
            task_id=f"sleep_for_{i}", python_callable=my_sleeping_function, op_kwargs={"random_base": i / 10} # python_callable是回调函数, op_kwargs是python_callable的入参
        )

        run_this >> log_the_sql >> sleeping_task

    
    def log_sql(**kwargs):
        log.info("Python task decorator query: %s", str(kwargs["templates_dict"]["query"]))

    log_the_sql = PythonOperator(
        task_id="log_sql_query",
        python_callable=log_sql,
        templates_dict={"query": "sql/sample.sql"}, # 非PythonOperator定义的入参将作为python_callable的入参
        templates_exts=[".sql"],
    )
```

op_kwargs: 在 PythonOperator 实例中定义, 用于向 Python 函数传递静态的关键字参数, 直接作为函数参数 `my_func(arg1=val1, ...)`
params: 在 DAG 实例中定义, 在 DAG 运行时传递动态参数, 通过函数上下文中的 `**kwargs 或 dag_run.conf` 字典访问

## 组件
- scheduler: 一种使用DAG定义结合元数据中的任务状态来决定哪些任务需要被执行以及任务执行优先级的过程

    定期扫描dags_folder
- api: web ui, 监控DAG运行状态, 或对DAG操作
- db: 数据库, 默认sqlite, 存储所有的DAG, 任务定义, 运行的历史, 用户, 权限等等

    默认db为airflow.db
- worker: 用来执行Executor接收的任务, 是实际执行任务逻辑的进程

## FAQ
### [`/airflow/settings.py raises TypeError: 'module' object is not callable`](https://github.com/apache/airflow/issues/44933)
Airflow before 2.8.1 doesn't support pendulum 3

解决:
1. 安装airflow 2.8.1及以上
1. 安装时使用官方指定的[CONSTRAINT_URL](https://airflow.apache.org/docs/apache-airflow/stable/start.html)
1. `pip install "pendulum<3.0.0"`

### `ValueError: The value api_auth/jwt_secret must be set!`
挪了AIRFLOW_HOME目录, `export AIRFLOW_HOME`前还需要更新里面的airflow.cfg中的相关路径