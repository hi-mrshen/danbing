## 重构
1. scheduler模块抽出来, 进行统计和调度，要不然job太重了
2. record添加数据字段的功能, coloum 已经传入，并且支持并发
3. 函数太长, 基于流程进行切分
4. task , taskgroup, 分界线不清晰优点混乱
5. 添加统计消耗时间和最后更新时间功能
6. 调度中存在历史记录，分析历史记录功能，保留接口不实现


job  
    - collect config 


scheduler 
    - communication (top)
    - persistence
    - group task    -> 基于持久化数据重新分配grouptask
    - tasks
    - taskgroup -> scheduler  
    - taskgroup run 
    - imlement static 
      - time 

taskgroup 
    - go reader -> channel -> go writer  : task 
    - communication 
   

channel 
    - implement static 
      - msg: update_time 
      - byte 
      - record 

