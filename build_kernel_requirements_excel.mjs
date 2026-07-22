import fs from 'node:fs/promises';
import { Workbook, SpreadsheetFile } from '@oai/artifact-tool';

const outDir = 'outputs/019f56e7-d046-7ca3-99f8-a37a45b25b6e';
const wb = Workbook.create();
const req = wb.worksheets.add('需求映射');
const cfg = wb.worksheets.add('内核配置清单');
const dt = wb.worksheets.add('设备树清单');
const test = wb.worksheets.add('验证方案');
for (const s of [req,cfg,dt,test]) { s.showGridLines = false; }

const requirements = [
 [1,'空载情况下 CPU 使用率不超过 5%','CONFIG_NO_HZ_IDLE=y；CONFIG_NO_HZ=y；CONFIG_HIGH_RES_TIMERS=y；可配合裁剪后台服务','设备树无直接阈值配置；CPU/定时器/中断拓扑来自 SoC 公共 dtsi','部分支持，必须实测','静态配置有利于降低空闲唤醒，但不能证明 ≤5%','top/mpstat/pidstat 连续采样；明确“空载”服务集合，统计 1/5/15 分钟及 24h 均值/峰值'],
 [2,'配置为单内核运行模式操作系统','当前 CONFIG_SMP=y、CONFIG_NR_CPUS=24、CONFIG_HOTPLUG_CPU=y；与“编译期单核”不一致。可用启动参数 maxcpus=1（运行期单核），或重新编译 CONFIG_SMP=n','当前提供的板级 dts 未禁用 CPU 节点；CPU 节点来自未提供的 t1022si-pre.dtsi','不满足（当前静态配置）','需明确选择：运行期单核 maxcpus=1，或严格单核内核 CONFIG_SMP=n','检查 /proc/cmdline、nproc、/sys/devices/system/cpu/online；应仅显示 CPU0'],
 [3,'空载情况下内存占用量不超过 128 MB','CONFIG_MEMCG=y 可做分组统计/限制；CONFIG_SWAP=y；CONFIG_SLUB_DEBUG=y 与 CONFIG_DEBUG_KERNEL=y 可能增加开销，建议生产版关闭调试项','共享 dtsi 声明 1GB DDR，预留 BMan/QMan 52MB；CPUB 覆盖为 2GB，另预留 PCIe 32MB/NTB 16MB（后者位于前者范围内）','部分支持，必须实测','物理内存容量不等于空载占用；128MB 取决于用户空间、服务、缓存口径','free -m、/proc/meminfo、smem；规定口径（建议 MemTotal-MemAvailable）并在启动稳定后及 24h 末记录'],
 [4,'无人为干预可稳定运行 24h','通用稳定性相关能力存在，但无单一 CONFIG 项可证明；建议关闭非必要 DEBUG 并保留 watchdog/日志（需核查完整配置）','设备树启用的外设、时钟、中断、内存预留会影响稳定性；需实际启动验证','无法由静态文件判定','必须进行 24h soak test，并记录崩溃、OOM、hung task、I/O/网络错误','连续 24h 运行目标业务/空载场景；采集 dmesg、journal、uptime、温度、内存、CPU 与错误计数'],
 [5,'具有 CPU 资源限制，防止后台线程持续占用 100% CPU','CONFIG_CGROUPS=y；CONFIG_CGROUP_SCHED=y；CONFIG_FAIR_GROUP_SCHED=y；CONFIG_CFS_BANDWIDTH=y；CONFIG_CPUSETS=y','设备树无 cgroup 配额配置','满足内核能力','用户空间需挂载 cgroup 并设置配额；cgroup v1 可用 cpu.cfs_quota_us/period_us，实际接口取决于系统版本','创建测试 cgroup/systemd slice，限制为 20% 后运行 busy loop，确认实际 CPU 不超过配额'],
 [6,'空载情况下进程间切换不超过 50 微秒','CONFIG_HIGH_RES_TIMERS=y；CONFIG_HZ=250；当前 CONFIG_PREEMPT_NONE=y，不利于低调度延迟，建议评估 CONFIG_PREEMPT 或 PREEMPT_RT','MPIC 中断、时钟和 CPU 节点影响延迟；板级文件声明 interrupt-parent=<&mpic>，但关键 SoC CPU/定时器定义不在所给文件内','当前配置不能保证，必须实测','“进程间切换”需定义测试方法；PREEMPT_NONE 下 50µs 上限风险较高','用 cyclictest/oslat 或自定义双进程 ping-pong；锁核运行，报告最大值和 99.9 分位，至少覆盖目标负载'],
 [7,'具有 glibc 系统库的操作功能','CONFIG_MMU=y；CONFIG_BINFMT_ELF=y；CONFIG_FUTEX=y；CONFIG_PROC_FS=y；CONFIG_SYSFS=y 等可支撑常规 glibc Linux 用户空间','设备树不决定 glibc；需正确描述 CPU/内存/外设以启动用户空间','内核侧基础满足，用户空间待确认','glibc 是用户空间组件，不能仅凭内核配置确认已安装','执行 ldd --version、getconf GNU_LIBC_VERSION；运行动态链接 ELF 程序并检查解释器/依赖'],
 [8,'具有 I/O 多路复用功能','CONFIG_EPOLL=y；同时 CONFIG_EVENTFD=y、CONFIG_SIGNALFD=y、CONFIG_TIMERFD=y','设备树无 epoll 配置；具体设备驱动可作为被监控 fd 来源','满足内核能力','select/poll 通常为基础系统调用；此处直接证据为 epoll 及相关 fd 接口','编译运行 epoll 测试程序，对 pipe/socket/timerfd 注册并验证事件唤醒'],
 [9,'具备 System V 信号量功能','CONFIG_SYSVIPC=y；CONFIG_SYSVIPC_SYSCTL=y','设备树无关','满足内核能力','System V IPC 总开关覆盖信号量、消息队列、共享内存','运行 ipcmk -S、ipcs -s、ipcrm；或调用 semget/semop/semctl'],
 [10,'具备 POSIX 共享内存功能','CONFIG_SHMEM=y；CONFIG_TMPFS=y（通常需挂载 /dev/shm）','设备树 reserved-memory/NTB 共享区不是 POSIX shm；二者概念不同','满足内核能力，挂载状态待确认','POSIX shm 由 shm_open + tmpfs(/dev/shm) 实现；CONFIG_POSIX_MQUEUE 与此需求无关','检查 /dev/shm 为 tmpfs；调用 shm_open/ftruncate/mmap 并跨进程读写'],
 [11,'具备 System V 共享内存功能','CONFIG_SYSVIPC=y；CONFIG_SYSVIPC_SYSCTL=y；CONFIG_SHMEM=y','设备树的 PCIe/NTB 物理共享内存不等同于 System V shm','满足内核能力','System V IPC 提供 shmget/shmat/shmctl','运行 ipcmk -M、ipcs -m、ipcrm；或调用 shmget/shmat/shmdt/shmctl']
];

req.getRange('A1:H1').merge(); req.getRange('A1').values=[['Linux 操作系统 11 条需求—内核配置与设备树映射']];
req.getRange('A2:H2').merge(); req.getRange('A2').values=[['依据：配置过滤后(2).txt、t1022d4rdb.dts/.dtsi、t1022d4rdb-cpua.dts、t1022d4rdb-cpub.dts；结论区分静态能力与运行验证。']];
req.getRange('A4:H4').values=[['序号','需求','相关内核配置','相关设备树','静态判断','判断说明','建议验证','责任/结果']];
req.getRange(`A5:H${4+requirements.length}`).values=requirements.map(r=>[...r,'待填写']);

const configs = [
 ['CONFIG_NO_HZ_COMMON=y',49,'空闲/定时器','支持 tickless 基础','需求1'],['CONFIG_NO_HZ_IDLE=y',50,'空闲/定时器','空闲 CPU 停止周期 tick，有利于降低空载开销','需求1'],['CONFIG_NO_HZ=y',51,'空闲/定时器','启用 NO_HZ','需求1'],['CONFIG_HIGH_RES_TIMERS=y',52,'定时器/延迟','高精度定时器','需求1、6'],['CONFIG_PREEMPT_NONE=y',54,'调度/延迟','无内核抢占，低延迟风险项','需求6'],['CONFIG_CGROUPS=y',68,'资源控制','启用控制组','需求5'],['CONFIG_MEMCG=y',70,'内存控制','内存 cgroup','需求3'],['CONFIG_CGROUP_SCHED=y',73,'资源控制','调度器 cgroup','需求5'],['CONFIG_FAIR_GROUP_SCHED=y',74,'资源控制','公平组调度','需求5'],['CONFIG_CFS_BANDWIDTH=y',75,'资源控制','CFS 带宽/CPU 配额关键项','需求5'],['CONFIG_CPUSETS=y',76,'资源控制','CPU/内存节点绑定','需求2、5、6'],['CONFIG_FUTEX=y',105,'用户空间/同步','glibc/pthread 常用快速用户态互斥','需求7'],['CONFIG_EPOLL=y',107,'I/O 多路复用','epoll 系统调用','需求8'],['CONFIG_SIGNALFD=y',108,'I/O 多路复用','信号 fd','需求8'],['CONFIG_TIMERFD=y',109,'I/O 多路复用','定时器 fd','需求8'],['CONFIG_EVENTFD=y',110,'I/O 多路复用','事件 fd','需求8'],['CONFIG_SHMEM=y',111,'共享内存','内核共享内存/tmpfs 基础','需求10、11'],['CONFIG_SLUB_DEBUG=y',125,'内存/调试','可能增加生产系统内存/性能开销','需求3'],['CONFIG_SLUB=y',127,'内存管理','SLUB 分配器','需求3'],['CONFIG_SMP=y',148,'CPU 拓扑','当前为 SMP，与严格编译期单核不一致','需求2'],['CONFIG_NR_CPUS=24',149,'CPU 拓扑','最大 CPU 数为 24','需求2'],['CONFIG_HZ=250',198,'调度/延迟','调度 tick 频率 250Hz','需求6'],['CONFIG_HOTPLUG_CPU=y',200,'CPU 拓扑','支持 CPU 热插拔/下线','需求2'],['CONFIG_MMU=y',154,'用户空间','glibc 常规环境基础','需求7'],['CONFIG_MODULES=y',293,'模块','模块会影响内存/攻击面，按需裁剪','需求3、4'],['CONFIG_BINFMT_ELF=y',322,'用户空间','运行 ELF 程序','需求7'],['CONFIG_PROC_FS=y',878,'文件系统','/proc 支撑运行监控及 glibc 用户空间','需求7'],['CONFIG_SYSFS=y',883,'文件系统','/sys 支撑设备/CPU 状态检查','需求2、7'],['CONFIG_TMPFS=y',884,'共享内存','/dev/shm 常用承载','需求10'],['CONFIG_DEBUG_KERNEL=y',1112,'调试','生产配置可能增加资源占用/干扰时延','需求3、4、6'],['CONFIG_SYSVIPC=y',24,'IPC','System V 信号量/消息队列/共享内存总开关','需求9、11'],['CONFIG_SYSVIPC_SYSCTL=y',25,'IPC','System V IPC 参数 sysctl','需求9、11'],['CONFIG_POSIX_MQUEUE=y',26,'IPC','POSIX 消息队列（不是 POSIX 共享内存）','补充说明'],['CONFIG_SWAP=y',23,'内存','交换支持；是否启用取决于用户空间','需求3']
];
cfg.getRange('A1:F1').merge(); cfg.getRange('A1').values=[['相关内核配置清单']];
cfg.getRange('A3:F3').values=[['配置项','源文件行号','分类','作用/风险','对应需求','建议']];
cfg.getRange(`A4:F${3+configs.length}`).values=configs.map(x=>[...x, x[0].includes('PREEMPT_NONE')?'评估 CONFIG_PREEMPT/PREEMPT_RT':x[0].includes('DEBUG')?'生产版按需关闭':x[0]==='CONFIG_SMP=y'?'严格单核时改为 CONFIG_SMP=n；运行期单核可用 maxcpus=1':'保留并通过运行测试确认']);

const dtsRows = [
 ['t1022d4rdb.dts',44,'interrupt-parent = <&mpic>;','根中断控制器引用','需求6','仅说明中断父节点，不能证明延迟'],['t1022d4rdb.dts',50,'clock-frequency = <66666666>;','系统时钟 66.67MHz','需求1、6','时钟输入信息'],['t1022d4rdb.dtsi',37,'reserved-memory { ... }','BMan/QMan 动态预留区','需求3','16MB + 4MB + 32MB = 52MB'],['t1022d4rdb.dtsi',183,'memory { reg = <... 0x40000000>; }','共享板级内存声明 1GB','需求3','注释称顶端 52MB 预留、可用约 972MB'],['t1022d4rdb-cpua.dts',25,'interrupt-parent = <&mpic>;','CPUA 中断父节点','需求6','与通用板级一致'],['t1022d4rdb-cpua.dts',30,'clock-frequency = <66666666>;','CPUA 系统时钟','需求1、6','静态硬件描述'],['t1022d4rdb-cpub.dts',28,'reserved-memory { ... }','CPUB PCIe/NTB 共享物理内存预留','需求3、10、11','物理共享区不等同于 POSIX/System V IPC'],['t1022d4rdb-cpub.dts',37,'pcie-ep-shared@7c000000: 32MB, no-map','PCIe EP 共享区','需求3','从普通内存映射排除'],['t1022d4rdb-cpub.dts',47,'pcie-ntb-shared@7d000000: 16MB, no-map','NTB 子区域','需求3','位于上述 32MB 区域内，不应重复相加'],['t1022d4rdb-cpub.dts',57,'memory { reg = <... 0x80000000>; }','CPUB 覆盖为 2GB DDR','需求3','覆盖共享 dtsi 的 1GB 声明'],['未提供：t1022si-pre.dtsi / post.dtsi','—','CPU/SoC 定时器/MPIC 具体节点','本次文件通过 include 引用','需求2、6','无法核查 CPU 节点数量/status；建议补充源文件或反编译 DTB']
];
dt.getRange('A1:F1').merge(); dt.getRange('A1').values=[['相关设备树清单']];
dt.getRange('A3:F3').values=[['源文件','行号','节点/属性','含义','对应需求','审查备注']];
dt.getRange(`A4:F${3+dtsRows.length}`).values=dtsRows;

const tests = requirements.map(r=>[r[0],r[1],r[6], r[0]===1?'CPU 空载均值/峰值 ≤5%':r[0]===2?'online CPU 仅 0':r[0]===3?'约定口径占用 ≤128MB':r[0]===4?'24h 无崩溃/OOM/hung task/关键错误':r[0]===5?'受限任务不超过设定配额':r[0]===6?'最大切换延迟 <50µs（需确认是“<”还是“≤”）':r[0]===7?'glibc 可查询且动态程序可运行':r[0]===8?'epoll 事件收发正确':r[0]===9?'semget/semop/semctl 成功':r[0]===10?'/dev/shm + shm_open/mmap 成功':'shmget/shmat/shmctl 成功','待执行','']);
test.getRange('A1:F1').merge(); test.getRange('A1').values=[['验收验证方案']];
test.getRange('A3:F3').values=[['序号','需求','测试方法','建议通过准则','执行状态','实测结果/证据路径']];
test.getRange(`A4:F${3+tests.length}`).values=tests;
test.getRange(`E4:E${3+tests.length}`).dataValidation={rule:{type:'list',values:['待执行','通过','不通过','阻塞']}};

function styleSheet(sheet, titleRange, headerRange, usedRange, widths) {
 sheet.getRange(titleRange).format={fill:'#17365D',font:{bold:true,color:'#FFFFFF',size:16},verticalAlignment:'center'};
 sheet.getRange(titleRange).format.rowHeight=30;
 sheet.getRange(headerRange).format={fill:'#2F75B5',font:{bold:true,color:'#FFFFFF'},wrapText:true,verticalAlignment:'center',borders:{preset:'inside',style:'thin',color:'#D9E2F3'}};
 sheet.getRange(headerRange).format.rowHeight=32;
 sheet.getRange(usedRange).format.font={name:'Microsoft YaHei',size:10};
 sheet.getRange(usedRange).format.wrapText=true;
 sheet.getRange(usedRange).format.verticalAlignment='top';
 sheet.getRange(usedRange).format.borders={insideHorizontal:{style:'thin',color:'#D9E2F3'}};
 widths.forEach((w,i)=>sheet.getRangeByIndexes(0,i,1,1).format.columnWidth=w);
}
styleSheet(req,'A1:H1','A4:H4',`A1:H${4+requirements.length}`,[8,28,46,40,18,38,46,18]);
req.getRange('A2:H2').format={fill:'#D9EAF7',font:{italic:true,color:'#44546A'},wrapText:true}; req.getRange('A2:H2').format.rowHeight=34;
req.freezePanes.freezeRows(4); req.freezePanes.freezeColumns(2);
req.getRange(`A5:A${4+requirements.length}`).format.horizontalAlignment='center';
req.getRange(`E5:E${4+requirements.length}`).conditionalFormats.add('containsText',{text:'满足',format:{fill:'#E2F0D9',font:{color:'#375623'}}});
req.getRange(`E5:E${4+requirements.length}`).conditionalFormats.add('containsText',{text:'不满足',format:{fill:'#FCE4D6',font:{color:'#9C0006'}}});
req.getRange(`E5:E${4+requirements.length}`).conditionalFormats.add('containsText',{text:'必须实测',format:{fill:'#FFF2CC',font:{color:'#7F6000'}}});
styleSheet(cfg,'A1:F1','A3:F3',`A1:F${3+configs.length}`,[30,12,18,44,18,34]); cfg.freezePanes.freezeRows(3);
styleSheet(dt,'A1:F1','A3:F3',`A1:F${3+dtsRows.length}`,[29,10,44,30,18,44]); dt.freezePanes.freezeRows(3);
styleSheet(test,'A1:F1','A3:F3',`A1:F${3+tests.length}`,[8,30,54,40,16,36]); test.freezePanes.freezeRows(3);
test.getRange(`E4:E${3+tests.length}`).conditionalFormats.add('containsText',{text:'通过',format:{fill:'#E2F0D9'}});
test.getRange(`E4:E${3+tests.length}`).conditionalFormats.add('containsText',{text:'不通过',format:{fill:'#F4CCCC'}});
test.getRange(`E4:E${3+tests.length}`).conditionalFormats.add('containsText',{text:'待执行',format:{fill:'#FFF2CC'}});

for (const [sheet, range] of [[req,'A1:H15'],[cfg,'A1:F37'],[dt,'A1:F14'],[test,'A1:F14']]) {
 const rows=sheet.getRange(range); rows.format.autofitRows();
}

await fs.mkdir(outDir,{recursive:true});
for (const s of ['需求映射','内核配置清单','设备树清单','验证方案']) {
 const img=await wb.render({sheetName:s,autoCrop:'all',scale:1,format:'png'});
 await fs.writeFile(`${outDir}/${s}.png`,new Uint8Array(await img.arrayBuffer()));
}
console.log((await wb.inspect({kind:'table',range:'需求映射!A1:H15',include:'values,formulas',tableMaxRows:15,tableMaxCols:8,maxChars:6000})).ndjson);
console.log((await wb.inspect({kind:'match',searchTerm:'#REF!|#DIV/0!|#VALUE!|#NAME\\?|#N/A',options:{useRegex:true,maxResults:100},summary:'formula error scan'})).ndjson);
const xlsx=await SpreadsheetFile.exportXlsx(wb);
await xlsx.save(`${outDir}/T1022_操作系统11项需求_内核配置与设备树映射.xlsx`);
