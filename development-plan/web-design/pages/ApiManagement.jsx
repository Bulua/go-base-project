// pages/ApiManagement.jsx
const ApiManagement = () => {
  const [search, setSearch] = React.useState('');
  const [methodFilter, setMethodFilter] = React.useState('');
  const [modalOpen, setModalOpen] = React.useState(false);
  const [editApi, setEditApi] = React.useState(null);
  const [confirmOpen, setConfirmOpen] = React.useState(false);
  const [loading, setLoading] = React.useState(false);
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(10);

  const apis = [
    { id:1, name:'获取用户列表', path:'/api/v1/users', method:'GET', module:'用户管理', auth:true, status:'active', description:'分页获取系统用户列表', callCount:18420 },
    { id:2, name:'创建用户', path:'/api/v1/users', method:'POST', module:'用户管理', auth:true, status:'active', description:'新增系统用户账号', callCount:342 },
    { id:3, name:'更新用户信息', path:'/api/v1/users/:id', method:'PUT', module:'用户管理', auth:true, status:'active', description:'修改指定用户的基本信息', callCount:1256 },
    { id:4, name:'删除用户', path:'/api/v1/users/:id', method:'DELETE', module:'用户管理', auth:true, status:'active', description:'软删除指定用户', callCount:89 },
    { id:5, name:'获取角色列表', path:'/api/v1/roles', method:'GET', module:'角色管理', auth:true, status:'active', description:'获取所有角色及权限信息', callCount:8930 },
    { id:6, name:'创建角色', path:'/api/v1/roles', method:'POST', module:'角色管理', auth:true, status:'active', description:'新建角色并配置权限', callCount:127 },
    { id:7, name:'用户登录', path:'/api/v1/auth/login', method:'POST', module:'认证授权', auth:false, status:'active', description:'账号密码登录，返回JWT Token', callCount:56780 },
    { id:8, name:'刷新Token', path:'/api/v1/auth/refresh', method:'POST', module:'认证授权', auth:false, status:'active', description:'使用RefreshToken换取新AccessToken', callCount:34210 },
    { id:9, name:'获取系统配置', path:'/api/v1/config', method:'GET', module:'系统配置', auth:true, status:'active', description:'获取当前系统配置项', callCount:4520 },
    { id:10, name:'更新系统配置', path:'/api/v1/config', method:'PATCH', module:'系统配置', auth:true, status:'inactive', description:'批量更新系统配置（已停用）', callCount:0 },
    { id:11, name:'获取操作日志', path:'/api/v1/logs/operation', method:'GET', module:'日志审计', auth:true, status:'active', description:'查询操作历史记录', callCount:2340 },
    { id:12, name:'获取登录日志', path:'/api/v1/logs/login', method:'GET', module:'日志审计', auth:true, status:'active', description:'查询用户登录记录', callCount:1890 },
  ];

  const filtered = apis.filter(a => {
    const q = search.toLowerCase();
    return (!q || a.name.includes(q) || a.path.includes(q)) && (!methodFilter || a.method === methodFilter);
  });
  const paginated = filtered.slice((page-1)*pageSize, page*pageSize);

  const handleSave = () => {
    setLoading(true);
    setTimeout(() => { setLoading(false); setModalOpen(false); window.toast.success(editApi ? 'API 更新成功' : 'API 创建成功'); }, 700);
  };

  const columns = [
    { key:'name', title:'接口名称', render:(_, r) => (
      <div>
        <div style={{ fontWeight:'500', fontSize:'13px', marginBottom:'2px' }}>{r.name}</div>
        <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
          <MethodTag method={r.method} />
          <code style={{ fontFamily:'var(--font-mono)', fontSize:'11px', color:'var(--text-secondary)' }}>{r.path}</code>
        </div>
      </div>
    )},
    { key:'module', title:'模块', dataIndex:'module', render: v => <Tag>{v}</Tag> },
    { key:'auth', title:'鉴权', dataIndex:'auth', render: v => v
      ? <Tag color="warning" dot>需要认证</Tag>
      : <Tag color="default" dot>公开接口</Tag>
    },
    { key:'callCount', title:'调用量', dataIndex:'callCount', sortable:true, render: v => (
      <span style={{ fontWeight:'600', fontSize:'13px', fontFamily:'var(--font-mono)', color: v > 10000 ? 'var(--primary)' : 'var(--text-primary)' }}>
        {v >= 10000 ? `${(v/10000).toFixed(1)}w` : v.toLocaleString()}
      </span>
    )},
    { key:'status', title:'状态', dataIndex:'status', render: v => <StatusBadge status={v} /> },
    { key:'actions', title:'操作', width:'130px', render:(_, r) => (
      <div style={{ display:'flex', gap:'4px' }}>
        <Button size="xs" variant="ghost" icon={<i className="bi bi-pencil" />} onClick={() => { setEditApi(r); setModalOpen(true); }}>编辑</Button>
        <Button size="xs" variant="secondary" icon={<i className="bi bi-eye" />} onClick={() => window.toast.info('查看接口详情')}>详情</Button>
        <Button size="xs" variant="text" danger icon={<i className="bi bi-trash" />} onClick={() => setConfirmOpen(true)} />
      </div>
    )},
  ];

  const methodOptions = [{value:'GET',label:'GET'},{value:'POST',label:'POST'},{value:'PUT',label:'PUT'},{value:'PATCH',label:'PATCH'},{value:'DELETE',label:'DELETE'}];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">API 管理</div>
          <div className="bp-page-subtitle">管理系统后端接口及调用鉴权配置</div>
        </div>
        <Button icon={<i className="bi bi-plus-lg" />} onClick={() => { setEditApi(null); setModalOpen(true); }}>新增接口</Button>
      </div>

      {/* Stats row */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(4, 1fr)', gap:'12px' }}>
        {[
          { label:'接口总数', value:apis.length, icon:'bi-cloud', color:'var(--primary)' },
          { label:'已启用', value:apis.filter(a=>a.status==='active').length, icon:'bi-check-circle', color:'var(--success)' },
          { label:'需鉴权', value:apis.filter(a=>a.auth).length, icon:'bi-lock', color:'var(--warning)' },
          { label:'今日调用', value:'12.4w', icon:'bi-lightning', color:'oklch(0.52 0.15 285)' },
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

      <Card padding="16px">
        <TableToolbar
          left={<>
            <SearchInput value={search} onChange={e => setSearch(e.target.value)} style={{ width:'220px' }} />
            <Select value={methodFilter} onChange={e => setMethodFilter(e.target.value)} options={methodOptions} placeholder="全部方法" style={{ width:'110px' }} />
          </>}
          right={<Button size="sm" variant="secondary" icon={<i className="bi bi-download" />} onClick={() => window.toast.info('导出 OpenAPI 文档')}>导出文档</Button>}
        />
        <Table columns={columns} dataSource={paginated} rowKey="id" />
        <div style={{ marginTop:'16px' }}>
          <Pagination total={filtered.length} page={page} pageSize={pageSize} onChange={setPage} onPageSizeChange={s => { setPageSize(s); setPage(1); }} />
        </div>
      </Card>

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editApi ? '编辑接口' : '新增接口'} size="md"
        footer={<><Button variant="secondary" onClick={() => setModalOpen(false)}>取消</Button><Button onClick={handleSave} loading={loading}>保存</Button></>}>
        <Form style={{ gap:'16px' }}>
          <FormItem label="接口名称" required><Input placeholder="请输入接口名称" defaultValue={editApi?.name} /></FormItem>
          <FormRow cols={2}>
            <FormItem label="请求方法" required><Select options={methodOptions} placeholder="请选择" defaultValue={editApi?.method} /></FormItem>
            <FormItem label="所属模块"><Select options={['用户管理','角色管理','认证授权','系统配置','日志审计']} placeholder="请选择" defaultValue={editApi?.module} /></FormItem>
          </FormRow>
          <FormItem label="接口路径" required help="支持路径参数，如 /users/:id"><Input placeholder="/api/v1/..." defaultValue={editApi?.path} style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          <FormItem label="接口描述"><Textarea placeholder="描述接口用途" rows={2} defaultValue={editApi?.description} /></FormItem>
          <FormRow cols={2}>
            <FormItem label="是否鉴权"><RadioGroup value={editApi?.auth ? 'true' : 'false'} options={[{value:'true',label:'需要认证'},{value:'false',label:'公开接口'}]} /></FormItem>
            <FormItem label="状态"><RadioGroup value={editApi?.status || 'active'} options={[{value:'active',label:'启用'},{value:'inactive',label:'停用'}]} /></FormItem>
          </FormRow>
        </Form>
      </Modal>

      <ConfirmModal open={confirmOpen} onClose={() => setConfirmOpen(false)}
        onConfirm={() => { setLoading(true); setTimeout(()=>{setLoading(false);setConfirmOpen(false);window.toast.success('接口已删除');},600); }}
        title="删除接口" type="danger" loading={loading}
        message="确认删除此 API 接口？相关调用将会失败，请确认已更新客户端配置。" confirmText="确认删除" />
    </div>
  );
};
Object.assign(window, { ApiManagement });
