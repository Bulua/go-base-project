// pages/LoginLogs.jsx
const LoginLogs = () => {
  const [search, setSearch] = React.useState('');
  const [resultFilter, setResultFilter] = React.useState('');
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(10);

  const logs = [
    { id:1,  user:'admin',    realName:'系统管理员', ip:'192.168.1.100', location:'内网',    browser:'Chrome 120', os:'Windows 11', result:'success', time:'2024-01-15 14:30:00', remark:'' },
    { id:2,  user:'zhangwei', realName:'张伟',     ip:'10.0.0.42',     location:'内网',    browser:'Edge 120',   os:'Windows 10', result:'success', time:'2024-01-15 09:22:11', remark:'' },
    { id:3,  user:'lina',     realName:'李娜',     ip:'10.0.0.88',     location:'内网',    browser:'Chrome 120', os:'macOS 14',   result:'success', time:'2024-01-15 08:55:40', remark:'' },
    { id:4,  user:'unknown',  realName:'—',        ip:'103.45.20.188',  location:'北京市',  browser:'Chrome 119', os:'Windows',    result:'failed',  time:'2024-01-15 03:12:55', remark:'密码错误 (3次)' },
    { id:5,  user:'wanglei',  realName:'王磊',     ip:'172.16.0.20',   location:'内网',    browser:'Firefox 121',os:'Linux',      result:'success', time:'2024-01-14 17:40:28', remark:'' },
    { id:6,  user:'admin',    realName:'系统管理员', ip:'192.168.1.100', location:'内网',    browser:'Chrome 120', os:'Windows 11', result:'success', time:'2024-01-14 09:05:12', remark:'' },
    { id:7,  user:'unknown',  realName:'—',        ip:'45.88.104.22',   location:'上海市',  browser:'Python/3.9', os:'Linux',      result:'failed',  time:'2024-01-14 02:08:44', remark:'账号不存在' },
    { id:8,  user:'chenjing', realName:'陈静',     ip:'10.0.0.55',     location:'内网',    browser:'Safari 17',  os:'macOS 14',   result:'success', time:'2024-01-14 08:30:05', remark:'' },
    { id:9,  user:'liuyang',  realName:'刘阳',     ip:'10.0.0.60',     location:'内网',    browser:'Chrome 120', os:'Windows 10', result:'failed',  time:'2024-01-13 16:22:38', remark:'密码错误 (1次)' },
    { id:10, user:'liuyang',  realName:'刘阳',     ip:'10.0.0.60',     location:'内网',    browser:'Chrome 120', os:'Windows 10', result:'success', time:'2024-01-13 16:23:15', remark:'' },
    { id:11, user:'zhaomin',  realName:'赵敏',     ip:'10.0.0.77',     location:'内网',    browser:'Edge 120',   os:'Windows 11', result:'success', time:'2024-01-13 09:12:00', remark:'' },
    { id:12, user:'unknown',  realName:'—',        ip:'91.134.168.22',  location:'法国',    browser:'curl/7.88',  os:'Linux',      result:'failed',  time:'2024-01-12 11:44:29', remark:'IP 已封禁' },
    { id:13, user:'wuqiang',  realName:'吴强',     ip:'10.0.0.91',     location:'内网',    browser:'Chrome 120', os:'Ubuntu',     result:'success', time:'2024-01-12 08:40:55', remark:'' },
    { id:14, user:'admin',    realName:'系统管理员', ip:'58.20.44.100',  location:'湖南省',  browser:'Chrome 120', os:'Windows 11', result:'success', time:'2024-01-11 21:15:33', remark:'异地登录提醒' },
  ];

  const filtered = logs.filter(l => {
    const q = search.toLowerCase();
    return (!q || l.user.includes(q) || l.realName.includes(q) || l.ip.includes(q))
      && (!resultFilter || l.result === resultFilter);
  });
  const paginated = filtered.slice((page-1)*pageSize, page*pageSize);

  const successCount = logs.filter(l=>l.result==='success').length;
  const failCount = logs.filter(l=>l.result==='failed').length;
  const uniqueIPs = new Set(logs.filter(l=>l.result==='failed').map(l=>l.ip)).size;

  const columns = [
    { key:'user', title:'登录账号', render:(_, r) => (
      <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
        <div style={{ width:'28px', height:'28px', borderRadius:'50%', background: r.result==='failed' ? 'var(--danger-light)' : 'var(--primary-light)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'11px', fontWeight:'700', color: r.result==='failed' ? 'var(--danger)' : 'var(--primary)', flexShrink:0 }}>
          {r.realName !== '—' ? r.realName.charAt(0) : '?'}
        </div>
        <div>
          <div style={{ fontSize:'13px', fontWeight:'500' }}>{r.realName !== '—' ? r.realName : <span style={{color:'var(--text-tertiary)'}}>未知用户</span>}</div>
          <code style={{ fontSize:'11px', color:'var(--text-tertiary)', fontFamily:'var(--font-mono)' }}>{r.user}</code>
        </div>
      </div>
    )},
    { key:'ip', title:'IP 地址', render:(_, r) => (
      <div>
        <code style={{ fontFamily:'var(--font-mono)', fontSize:'12px', color:'var(--text-secondary)', display:'block' }}>{r.ip}</code>
        <span style={{ fontSize:'11px', color:'var(--text-tertiary)' }}>{r.location}</span>
      </div>
    )},
    { key:'device', title:'设备', render:(_, r) => (
      <div style={{ fontSize:'12px' }}>
        <div style={{ color:'var(--text-secondary)' }}>{r.browser}</div>
        <div style={{ color:'var(--text-tertiary)' }}>{r.os}</div>
      </div>
    )},
    { key:'result', title:'状态', dataIndex:'result', render: v => <StatusBadge status={v} /> },
    { key:'remark', title:'备注', dataIndex:'remark', render: v => v
      ? <span style={{ fontSize:'12px', color:'var(--danger-text)', display:'flex', alignItems:'center', gap:'4px' }}><i className="bi bi-exclamation-circle" />{v}</span>
      : <span style={{ color:'var(--text-tertiary)', fontSize:'12px' }}>—</span>
    },
    { key:'time', title:'登录时间', dataIndex:'time', render: v => <span style={{ fontSize:'12px', color:'var(--text-tertiary)' }}>{v}</span> },
    { key:'actions', title:'', width:'60px', render:(_, r) => r.result==='failed' && (
      <Button size="xs" variant="text" danger onClick={() => window.toast.warning(`IP ${r.ip} 已加入黑名单`)}>封禁IP</Button>
    )},
  ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">登录日志</div>
          <div className="bp-page-subtitle">监控用户登录行为与异常访问</div>
        </div>
        <div style={{ display:'flex', gap:'8px' }}>
          <Button variant="secondary" icon={<i className="bi bi-shield-x" />} onClick={() => window.toast.warning('已清理 3 个异常IP')}>清理异常</Button>
          <Button variant="secondary" icon={<i className="bi bi-download" />} onClick={() => window.toast.info('正在导出...')}>导出</Button>
        </div>
      </div>

      {/* Stats */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(4, 1fr)', gap:'12px' }}>
        {[
          { label:'总登录次数', value:logs.length, icon:'bi-door-open', color:'var(--primary)' },
          { label:'成功登录', value:successCount, icon:'bi-check-circle', color:'var(--success)' },
          { label:'登录失败', value:failCount, icon:'bi-x-circle', color:'var(--danger)' },
          { label:'异常 IP', value:uniqueIPs, icon:'bi-shield-exclamation', color:'var(--warning)' },
        ].map(s => (
          <div key={s.label} className="bp-card-el" style={{ background:'var(--bg-surface)', border:'1px solid var(--border)', borderRadius:'var(--radius-xl)', padding:'14px 16px', display:'flex', alignItems:'center', gap:'12px' }}>
            <div style={{ width:'36px', height:'36px', borderRadius:'var(--radius-lg)', background:'var(--primary-light)', display:'flex', alignItems:'center', justifyContent:'center', color:s.color, fontSize:'16px', flexShrink:0 }}>
              <i className={`bi ${s.icon}`} />
            </div>
            <div>
              <div style={{ fontSize:'20px', fontWeight:'700', color:'var(--text-primary)', lineHeight:1 }}>{s.value}</div>
              <div style={{ fontSize:'12px', color:'var(--text-secondary)', marginTop:'2px' }}>{s.label}</div>
            </div>
          </div>
        ))}
      </div>

      {failCount > 0 && (
        <Alert type="warning" title="检测到异常登录尝试" message={`近期有 ${failCount} 次登录失败记录，其中 ${uniqueIPs} 个异常 IP，建议检查并封禁可疑来源。`} closable />
      )}

      <Card padding="16px">
        <TableToolbar
          left={<>
            <SearchInput value={search} onChange={e => setSearch(e.target.value)} placeholder="搜索用户或IP" style={{ width:'200px' }} />
            <Select value={resultFilter} onChange={e => setResultFilter(e.target.value)}
              options={[{value:'success',label:'登录成功'},{value:'failed',label:'登录失败'}]}
              placeholder="全部状态" style={{ width:'120px' }} />
          </>}
          right={<Tag color={resultFilter==='failed' ? 'danger' : 'primary'}>{filtered.length} 条记录</Tag>}
        />
        <Table columns={columns} dataSource={paginated} rowKey="id" />
        <div style={{ marginTop:'16px' }}>
          <Pagination total={filtered.length} page={page} pageSize={pageSize} onChange={setPage} onPageSizeChange={s=>{setPageSize(s);setPage(1);}} />
        </div>
      </Card>
    </div>
  );
};
Object.assign(window, { LoginLogs });
