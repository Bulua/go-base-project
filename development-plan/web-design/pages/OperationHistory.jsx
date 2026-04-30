// pages/OperationHistory.jsx
const OperationHistory = () => {
  const [search, setSearch] = React.useState('');
  const [moduleFilter, setModuleFilter] = React.useState('');
  const [resultFilter, setResultFilter] = React.useState('');
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(10);
  const [detailOpen, setDetailOpen] = React.useState(false);
  const [detailRecord, setDetailRecord] = React.useState(null);

  const logs = [
    { id:1, user:'admin', realName:'系统管理员', module:'用户管理', action:'修改用户权限', method:'PUT', path:'/api/v1/users/12/roles', ip:'192.168.1.100', location:'内网', browser:'Chrome 120', os:'Windows 11', result:'success', duration:234, time:'2024-01-15 14:32:18', requestBody:'{"roles":["ops_admin"]}', responseCode:200 },
    { id:2, user:'zhangwei', realName:'张伟', module:'报表中心', action:'导出数据报表', method:'GET', path:'/api/v1/reports/export', ip:'10.0.0.42', location:'内网', browser:'Edge 120', os:'Windows 10', result:'success', duration:1820, time:'2024-01-15 14:15:06', requestBody:'', responseCode:200 },
    { id:3, user:'lina', realName:'李娜', module:'字典管理', action:'删除字典项', method:'DELETE', path:'/api/v1/dict/items/8', ip:'10.0.0.88', location:'内网', browser:'Chrome 120', os:'macOS 14', result:'failed', duration:89, time:'2024-01-15 13:58:44', requestBody:'', responseCode:403 },
    { id:4, user:'admin', realName:'系统管理员', module:'API管理', action:'新增API接口', method:'POST', path:'/api/v1/apis', ip:'192.168.1.100', location:'内网', browser:'Chrome 120', os:'Windows 11', result:'success', duration:156, time:'2024-01-15 11:24:30', requestBody:'{"name":"新接口","path":"/test"}', responseCode:201 },
    { id:5, user:'wanglei', realName:'王磊', module:'用户管理', action:'重置用户密码', method:'POST', path:'/api/v1/users/5/reset-pwd', ip:'172.16.0.20', location:'内网', browser:'Firefox 121', os:'Linux', result:'success', duration:312, time:'2024-01-15 10:45:12', requestBody:'', responseCode:200 },
    { id:6, user:'admin', realName:'系统管理员', module:'系统配置', action:'修改系统参数', method:'PATCH', path:'/api/v1/config', ip:'192.168.1.100', location:'内网', browser:'Chrome 120', os:'Windows 11', result:'success', duration:198, time:'2024-01-15 09:30:55', requestBody:'{"uploadMaxSize":20}', responseCode:200 },
    { id:7, user:'chenjing', realName:'陈静', module:'角色管理', action:'新增角色', method:'POST', path:'/api/v1/roles', ip:'10.0.0.55', location:'内网', browser:'Safari 17', os:'macOS 14', result:'success', duration:267, time:'2024-01-14 17:22:38', requestBody:'{"name":"审计员"}', responseCode:201 },
    { id:8, user:'liuyang', realName:'刘阳', module:'用户管理', action:'查看用户列表', method:'GET', path:'/api/v1/users', ip:'10.0.0.60', location:'内网', browser:'Chrome 120', os:'Windows 10', result:'success', duration:45, time:'2024-01-14 16:18:20', requestBody:'', responseCode:200 },
    { id:9, user:'admin', realName:'系统管理员', module:'菜单管理', action:'删除菜单', method:'DELETE', path:'/api/v1/menus/99', ip:'192.168.1.100', location:'内网', browser:'Chrome 120', os:'Windows 11', result:'failed', duration:67, time:'2024-01-14 15:40:05', requestBody:'', responseCode:404 },
    { id:10, user:'zhaomin', realName:'赵敏', module:'系统配置', action:'导出备份', method:'POST', path:'/api/v1/config/backup', ip:'10.0.0.77', location:'内网', browser:'Edge 120', os:'Windows 11', result:'success', duration:4520, time:'2024-01-14 14:05:33', requestBody:'', responseCode:200 },
    { id:11, user:'wuqiang', realName:'吴强', module:'API管理', action:'更新接口权限', method:'PUT', path:'/api/v1/apis/3/perms', ip:'10.0.0.91', location:'内网', browser:'Chrome 120', os:'Ubuntu', result:'success', duration:188, time:'2024-01-14 11:28:19', requestBody:'{"auth":false}', responseCode:200 },
    { id:12, user:'lina', realName:'李娜', module:'用户管理', action:'批量导入用户', method:'POST', path:'/api/v1/users/import', ip:'10.0.0.88', location:'内网', browser:'Chrome 120', os:'macOS 14', result:'success', duration:2340, time:'2024-01-13 10:12:44', requestBody:'[multipart/form-data]', responseCode:200 },
  ];

  const modules = [...new Set(logs.map(l => l.module))];

  const filtered = logs.filter(l => {
    const q = search.toLowerCase();
    return (!q || l.user.includes(q) || l.realName.includes(q) || l.action.includes(q))
      && (!moduleFilter || l.module === moduleFilter)
      && (!resultFilter || l.result === resultFilter);
  });
  const paginated = filtered.slice((page-1)*pageSize, page*pageSize);

  const columns = [
    { key:'user', title:'操作用户', render:(_, r) => (
      <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
        <div style={{ width:'28px', height:'28px', borderRadius:'50%', background:'var(--primary-light)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'11px', fontWeight:'700', color:'var(--primary)', flexShrink:0 }}>
          {r.realName.charAt(0)}
        </div>
        <div>
          <div style={{ fontSize:'13px', fontWeight:'500' }}>{r.realName}</div>
          <code style={{ fontSize:'11px', color:'var(--text-tertiary)', fontFamily:'var(--font-mono)' }}>{r.user}</code>
        </div>
      </div>
    )},
    { key:'action', title:'操作内容', render:(_, r) => (
      <div>
        <div style={{ fontSize:'13px', marginBottom:'3px' }}>{r.action}</div>
        <div style={{ display:'flex', alignItems:'center', gap:'5px' }}>
          <MethodTag method={r.method} />
          <code style={{ fontSize:'11px', color:'var(--text-tertiary)', fontFamily:'var(--font-mono)' }}>{r.path}</code>
        </div>
      </div>
    )},
    { key:'module', title:'模块', dataIndex:'module', render: v => <Tag>{v}</Tag> },
    { key:'ip', title:'IP 地址', dataIndex:'ip', render: v => <code style={{ fontFamily:'var(--font-mono)', fontSize:'12px', color:'var(--text-secondary)' }}>{v}</code> },
    { key:'duration', title:'耗时', dataIndex:'duration', sortable:true, render: v => (
      <span style={{ fontSize:'12px', fontFamily:'var(--font-mono)', color: v > 1000 ? 'var(--warning-text)' : v > 2000 ? 'var(--danger-text)' : 'var(--text-secondary)' }}>
        {v >= 1000 ? `${(v/1000).toFixed(1)}s` : `${v}ms`}
      </span>
    )},
    { key:'result', title:'结果', dataIndex:'result', render: v => <StatusBadge status={v} /> },
    { key:'time', title:'时间', dataIndex:'time', render: v => <span style={{ fontSize:'12px', color:'var(--text-tertiary)' }}>{v}</span> },
    { key:'actions', title:'', width:'50px', render:(_, r) => (
      <Button size="xs" variant="ghost" onClick={() => { setDetailRecord(r); setDetailOpen(true); }}>详情</Button>
    )},
  ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">操作历史</div>
          <div className="bp-page-subtitle">记录用户对系统数据的所有变更操作</div>
        </div>
        <Button variant="secondary" icon={<i className="bi bi-download" />} onClick={() => window.toast.info('正在导出日志...')}>导出日志</Button>
      </div>

      <Card padding="16px">
        <TableToolbar
          left={<>
            <SearchInput value={search} onChange={e => setSearch(e.target.value)} placeholder="搜索用户或操作" style={{ width:'200px' }} />
            <Select value={moduleFilter} onChange={e => setModuleFilter(e.target.value)}
              options={modules.map(m => ({value:m, label:m}))} placeholder="全部模块" style={{ width:'120px' }} />
            <Select value={resultFilter} onChange={e => setResultFilter(e.target.value)}
              options={[{value:'success',label:'成功'},{value:'failed',label:'失败'}]} placeholder="全部结果" style={{ width:'100px' }} />
          </>}
          right={<Tag color="primary">{filtered.length} 条记录</Tag>}
        />
        <Table columns={columns} dataSource={paginated} rowKey="id" />
        <div style={{ marginTop:'16px' }}>
          <Pagination total={filtered.length} page={page} pageSize={pageSize} onChange={setPage} onPageSizeChange={s=>{setPageSize(s);setPage(1);}} />
        </div>
      </Card>

      {/* Detail modal */}
      <Modal open={detailOpen} onClose={() => setDetailOpen(false)} title="操作详情" size="md" footer={<Button onClick={() => setDetailOpen(false)}>关闭</Button>}>
        {detailRecord && (
          <div style={{ display:'flex', flexDirection:'column', gap:'16px' }}>
            <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'12px' }}>
              {[
                ['操作用户', `${detailRecord.realName} (${detailRecord.user})`],
                ['操作模块', detailRecord.module],
                ['操作内容', detailRecord.action],
                ['操作结果', detailRecord.result === 'success' ? '成功' : '失败'],
                ['请求方法', detailRecord.method],
                ['响应状态', detailRecord.responseCode],
                ['IP 地址', detailRecord.ip],
                ['地理位置', detailRecord.location],
                ['浏览器', detailRecord.browser],
                ['操作系统', detailRecord.os],
                ['耗时', `${detailRecord.duration}ms`],
                ['操作时间', detailRecord.time],
              ].map(([k,v]) => (
                <div key={k} style={{ display:'flex', flexDirection:'column', gap:'3px' }}>
                  <span style={{ fontSize:'11px', color:'var(--text-tertiary)', fontWeight:'600', textTransform:'uppercase', letterSpacing:'0.04em' }}>{k}</span>
                  <span style={{ fontSize:'13px', color:'var(--text-primary)' }}>{v}</span>
                </div>
              ))}
            </div>
            <FormDivider label="请求路径" />
            <code style={{ fontFamily:'var(--font-mono)', fontSize:'12px', background:'var(--bg-page)', padding:'8px 12px', borderRadius:'var(--radius-md)', color:'var(--text-secondary)', display:'block' }}>{detailRecord.method} {detailRecord.path}</code>
            {detailRecord.requestBody && <>
              <FormDivider label="请求数据" />
              <code style={{ fontFamily:'var(--font-mono)', fontSize:'12px', background:'var(--bg-page)', padding:'8px 12px', borderRadius:'var(--radius-md)', color:'var(--text-secondary)', display:'block', wordBreak:'break-all' }}>{detailRecord.requestBody}</code>
            </>}
          </div>
        )}
      </Modal>
    </div>
  );
};
Object.assign(window, { OperationHistory });
