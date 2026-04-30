// pages/RoleManagement.jsx
const RoleManagement = () => {
  const [modalOpen, setModalOpen] = React.useState(false);
  const [editRole, setEditRole] = React.useState(null);
  const [confirmOpen, setConfirmOpen] = React.useState(false);
  const [deleteTarget, setDeleteTarget] = React.useState(null);
  const [loading, setLoading] = React.useState(false);
  const [search, setSearch] = React.useState('');

  const roles = [
    { id:1, name:'超级管理员', code:'super_admin', description:'系统最高权限，可操作所有功能', userCount:2, permCount:128, status:'active', builtin:true, created:'2023-01-01' },
    { id:2, name:'运营管理员', code:'ops_admin', description:'负责内容运营及用户管理', userCount:5, permCount:64, status:'active', builtin:false, created:'2023-03-15' },
    { id:3, name:'开发人员', code:'developer', description:'系统开发及API接口管理', userCount:8, permCount:45, status:'active', builtin:false, created:'2023-03-20' },
    { id:4, name:'财务人员', code:'finance', description:'财务报表及账单管理', userCount:3, permCount:22, status:'active', builtin:false, created:'2023-04-10' },
    { id:5, name:'访客', code:'guest', description:'只读权限，仅可查看公开内容', userCount:12, permCount:8, status:'active', builtin:true, created:'2023-01-01' },
    { id:6, name:'审计员', code:'auditor', description:'审计日志查看，操作记录只读', userCount:1, permCount:15, status:'inactive', builtin:false, created:'2023-08-05' },
  ];

  const filtered = roles.filter(r => !search || r.name.includes(search) || r.code.includes(search));

  const PERMISSION_GROUPS = [
    { label:'用户管理', perms:['查看用户','新增用户','编辑用户','删除用户','重置密码'] },
    { label:'角色管理', perms:['查看角色','新增角色','编辑角色','删除角色','分配权限'] },
    { label:'菜单管理', perms:['查看菜单','新增菜单','编辑菜单','删除菜单'] },
    { label:'API管理',  perms:['查看API','新增API','编辑API','删除API'] },
    { label:'系统配置', perms:['查看配置','修改配置','导入导出'] },
  ];

  const [checkedPerms, setCheckedPerms] = React.useState({});
  const togglePerm = p => setCheckedPerms(prev => ({ ...prev, [p]: !prev[p] }));

  const handleSave = () => {
    setLoading(true);
    setTimeout(() => { setLoading(false); setModalOpen(false); window.toast.success(editRole ? '角色更新成功' : '角色创建成功'); }, 700);
  };

  const columns = [
    { key:'name', title:'角色名称', render:(_, r) => (
      <div style={{ display:'flex', alignItems:'center', gap:'10px' }}>
        <div style={{ width:'34px', height:'34px', borderRadius:'var(--radius-lg)', background:'var(--primary-light)', display:'flex', alignItems:'center', justifyContent:'center', color:'var(--primary)', fontSize:'15px', flexShrink:0 }}>
          <i className="bi bi-shield-check" />
        </div>
        <div>
          <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
            <span style={{ fontWeight:'500', fontSize:'13px' }}>{r.name}</span>
            {r.builtin && <Tag color="info" style={{ fontSize:'10px' }}>内置</Tag>}
          </div>
          <div style={{ fontSize:'11px', color:'var(--text-tertiary)', fontFamily:'var(--font-mono)' }}>{r.code}</div>
        </div>
      </div>
    )},
    { key:'description', title:'描述', dataIndex:'description', render: v => <span style={{ fontSize:'13px', color:'var(--text-secondary)' }}>{v}</span> },
    { key:'userCount', title:'用户数', dataIndex:'userCount', render: v => <span style={{ fontWeight:'600', color:'var(--primary)' }}>{v}</span> },
    { key:'permCount', title:'权限数', dataIndex:'permCount', render: v => <Tag color="primary">{v} 项</Tag> },
    { key:'status', title:'状态', dataIndex:'status', render: v => <StatusBadge status={v} /> },
    { key:'actions', title:'操作', width:'150px', render:(_, r) => (
      <div style={{ display:'flex', gap:'4px' }}>
        <Button size="xs" variant="ghost" icon={<i className="bi bi-pencil" />} onClick={() => { setEditRole(r); setModalOpen(true); }}>编辑</Button>
        <Button size="xs" variant="secondary" icon={<i className="bi bi-people" />} onClick={() => window.toast.info('查看角色用户')}>用户</Button>
        {!r.builtin && <Button size="xs" variant="text" danger icon={<i className="bi bi-trash" />} onClick={() => { setDeleteTarget(r); setConfirmOpen(true); }} />}
      </div>
    )},
  ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">角色管理</div>
          <div className="bp-page-subtitle">定义系统角色及其功能权限</div>
        </div>
        <Button icon={<i className="bi bi-shield-plus" />} onClick={() => { setEditRole(null); setCheckedPerms({}); setModalOpen(true); }}>新增角色</Button>
      </div>

      <Card padding="16px">
        <TableToolbar left={<SearchInput value={search} onChange={e => setSearch(e.target.value)} style={{ width:'220px' }} />} right={null} />
        <Table columns={columns} dataSource={filtered} rowKey="id" />
      </Card>

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editRole ? '编辑角色' : '新增角色'} size="lg"
        footer={<><Button variant="secondary" onClick={() => setModalOpen(false)}>取消</Button><Button onClick={handleSave} loading={loading}>保存</Button></>}>
        <Form style={{ gap:'20px' }}>
          <FormRow cols={2}>
            <FormItem label="角色名称" required><Input placeholder="请输入角色名称" defaultValue={editRole?.name} /></FormItem>
            <FormItem label="角色编码" required><Input placeholder="英文字母+下划线" defaultValue={editRole?.code} style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          </FormRow>
          <FormItem label="描述"><Textarea placeholder="角色功能描述" rows={2} defaultValue={editRole?.description} /></FormItem>
          <FormDivider label="权限配置" />
          <div style={{ display:'flex', flexDirection:'column', gap:'14px', maxHeight:'300px', overflowY:'auto', padding:'4px 0' }}>
            {PERMISSION_GROUPS.map(g => (
              <div key={g.label}>
                <div style={{ fontSize:'12px', fontWeight:'600', color:'var(--text-secondary)', marginBottom:'8px', textTransform:'uppercase', letterSpacing:'0.04em' }}>{g.label}</div>
                <div style={{ display:'flex', gap:'10px', flexWrap:'wrap' }}>
                  {g.perms.map(p => (
                    <Checkbox key={p} checked={!!checkedPerms[p]} onChange={() => togglePerm(p)}>{p}</Checkbox>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </Form>
      </Modal>

      <ConfirmModal open={confirmOpen} onClose={() => setConfirmOpen(false)} onConfirm={() => { setLoading(true); setTimeout(()=>{ setLoading(false); setConfirmOpen(false); window.toast.success('角色已删除'); },600); }}
        title="删除角色" type="danger" loading={loading}
        message={`确认删除角色「${deleteTarget?.name}」？该角色下的用户将失去对应权限。`} confirmText="确认删除" />
    </div>
  );
};
Object.assign(window, { RoleManagement });
