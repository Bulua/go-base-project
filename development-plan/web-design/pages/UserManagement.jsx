// pages/UserManagement.jsx
const UserManagement = () => {
  const [search, setSearch] = React.useState('');
  const [statusFilter, setStatusFilter] = React.useState('');
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(10);
  const [selected, setSelected] = React.useState([]);
  const [modalOpen, setModalOpen] = React.useState(false);
  const [editUser, setEditUser] = React.useState(null);
  const [confirmOpen, setConfirmOpen] = React.useState(false);
  const [deleteTarget, setDeleteTarget] = React.useState(null);
  const [loading, setLoading] = React.useState(false);

  const allUsers = [
    { id:1, name:'张伟', username:'zhangwei', email:'zhangwei@corp.com', phone:'138****8001', role:'超级管理员', dept:'技术部', status:'active', created:'2023-06-01', lastLogin:'2024-01-15 09:23' },
    { id:2, name:'李娜', username:'lina', email:'lina@corp.com', phone:'139****8002', role:'运营管理员', dept:'运营部', status:'active', created:'2023-07-12', lastLogin:'2024-01-14 18:45' },
    { id:3, name:'王磊', username:'wanglei', email:'wanglei@corp.com', phone:'137****8003', role:'开发人员', dept:'技术部', status:'inactive', created:'2023-08-20', lastLogin:'2023-12-20 11:30' },
    { id:4, name:'陈静', username:'chenjing', email:'chenjing@corp.com', phone:'136****8004', role:'运营管理员', dept:'运营部', status:'active', created:'2023-09-05', lastLogin:'2024-01-13 14:22' },
    { id:5, name:'刘阳', username:'liuyang', email:'liuyang@corp.com', phone:'135****8005', role:'访客', dept:'市场部', status:'pending', created:'2024-01-10', lastLogin:'—' },
    { id:6, name:'赵敏', username:'zhaomin', email:'zhaomin@corp.com', phone:'134****8006', role:'财务人员', dept:'财务部', status:'active', created:'2023-10-18', lastLogin:'2024-01-12 08:55' },
    { id:7, name:'吴强', username:'wuqiang', email:'wuqiang@corp.com', phone:'133****8007', role:'开发人员', dept:'技术部', status:'active', created:'2023-11-22', lastLogin:'2024-01-15 16:40' },
    { id:8, name:'郑丽', username:'zhengli', email:'zhengli@corp.com', phone:'132****8008', role:'运营管理员', dept:'运营部', status:'inactive', created:'2023-05-30', lastLogin:'2023-11-05 10:15' },
    { id:9, name:'孙博', username:'sunbo', email:'sunbo@corp.com', phone:'131****8009', role:'开发人员', dept:'技术部', status:'active', created:'2024-01-05', lastLogin:'2024-01-15 11:20' },
    { id:10, name:'周芳', username:'zhoufang', email:'zhoufang@corp.com', phone:'130****8010', role:'访客', dept:'市场部', status:'active', created:'2024-01-08', lastLogin:'2024-01-14 13:50' },
    { id:11, name:'徐明', username:'xuming', email:'xuming@corp.com', phone:'158****8011', role:'财务人员', dept:'财务部', status:'active', created:'2023-12-01', lastLogin:'2024-01-10 09:30' },
    { id:12, name:'朱霞', username:'zhuxia', email:'zhuxia@corp.com', phone:'159****8012', role:'运营管理员', dept:'运营部', status:'pending', created:'2024-01-12', lastLogin:'—' },
  ];

  const filtered = allUsers.filter(u => {
    const q = search.toLowerCase();
    const matchSearch = !q || u.name.includes(q) || u.username.includes(q) || u.email.includes(q);
    const matchStatus = !statusFilter || u.status === statusFilter;
    return matchSearch && matchStatus;
  });
  const paginated = filtered.slice((page-1)*pageSize, page*pageSize);

  const handleSave = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      setModalOpen(false);
      window.toast.success(editUser ? '用户信息更新成功' : '用户创建成功');
    }, 800);
  };

  const handleDelete = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      setConfirmOpen(false);
      window.toast.success('用户已删除');
    }, 600);
  };

  const columns = [
    { key:'name', title:'用户', render:(_, r) => (
      <div style={{ display:'flex', alignItems:'center', gap:'9px' }}>
        <div style={{ width:'32px', height:'32px', borderRadius:'50%', background:'var(--primary-light)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:'13px', fontWeight:'700', color:'var(--primary)', flexShrink:0 }}>{r.name.charAt(0)}</div>
        <div>
          <div style={{ fontWeight:'500', fontSize:'13px' }}>{r.name}</div>
          <div style={{ fontSize:'11px', color:'var(--text-tertiary)' }}>{r.username}</div>
        </div>
      </div>
    )},
    { key:'email', title:'联系方式', render:(_, r) => (
      <div>
        <div style={{ fontSize:'13px' }}>{r.email}</div>
        <div style={{ fontSize:'11px', color:'var(--text-tertiary)' }}>{r.phone}</div>
      </div>
    )},
    { key:'role', title:'角色', dataIndex:'role', render: v => <Tag color="primary">{v}</Tag> },
    { key:'dept', title:'部门', dataIndex:'dept', render: v => <Tag>{v}</Tag> },
    { key:'status', title:'状态', dataIndex:'status', render: v => <StatusBadge status={v} /> },
    { key:'lastLogin', title:'最后登录', dataIndex:'lastLogin', render: v => <span style={{ fontSize:'12px', color:'var(--text-tertiary)' }}>{v}</span> },
    { key:'actions', title:'操作', width:'120px', render:(_, r) => (
      <div style={{ display:'flex', gap:'4px' }}>
        <Button size="xs" variant="ghost" icon={<i className="bi bi-pencil" />} onClick={() => { setEditUser(r); setModalOpen(true); }}>编辑</Button>
        <Button size="xs" variant="text" danger icon={<i className="bi bi-trash" />} onClick={() => { setDeleteTarget(r); setConfirmOpen(true); }} />
      </div>
    )},
  ];

  const [form, setForm] = React.useState({ name:'', username:'', email:'', phone:'', role:'', dept:'', status:'active' });
  React.useEffect(() => {
    if (editUser) setForm({ name:editUser.name, username:editUser.username, email:editUser.email, phone:editUser.phone, role:editUser.role, dept:editUser.dept, status:editUser.status });
    else setForm({ name:'', username:'', email:'', phone:'', role:'', dept:'', status:'active' });
  }, [editUser, modalOpen]);

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">用户管理</div>
          <div className="bp-page-subtitle">管理系统所有用户账号及权限</div>
        </div>
        <Button icon={<i className="bi bi-person-plus" />} onClick={() => { setEditUser(null); setModalOpen(true); }}>新增用户</Button>
      </div>

      <Card padding="16px">
        <TableToolbar
          left={<>
            <SearchInput value={search} onChange={e => setSearch(e.target.value)} onSearch={() => {}} style={{ width:'220px' }} />
            <Select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
              options={[{value:'active',label:'已启用'},{value:'inactive',label:'已禁用'},{value:'pending',label:'待审核'}]}
              placeholder="全部状态" style={{ width:'120px' }} />
          </>}
          right={<>
            {selected.length > 0 && <Tag color="warning">{selected.length} 项已选</Tag>}
            <Button size="sm" variant="secondary" icon={<i className="bi bi-download" />} onClick={() => window.toast.info('正在导出...')}>导出</Button>
          </>}
        />
        <Table columns={columns} dataSource={paginated} rowKey="id" selectable selectedKeys={selected} onSelectChange={setSelected} />
        <div style={{ marginTop:'16px' }}>
          <Pagination total={filtered.length} page={page} pageSize={pageSize} onChange={setPage} onPageSizeChange={s => { setPageSize(s); setPage(1); }} />
        </div>
      </Card>

      {/* Add/Edit Modal */}
      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editUser ? '编辑用户' : '新增用户'} width={560}
        footer={<><Button variant="secondary" onClick={() => setModalOpen(false)}>取消</Button><Button onClick={handleSave} loading={loading}>保存</Button></>}>
        <Form>
          <FormRow cols={2}>
            <FormItem label="姓名" required><Input value={form.name} onChange={e => setForm(p=>({...p,name:e.target.value}))} placeholder="请输入姓名" /></FormItem>
            <FormItem label="用户名" required><Input value={form.username} onChange={e => setForm(p=>({...p,username:e.target.value}))} placeholder="登录用户名" /></FormItem>
          </FormRow>
          <FormRow cols={2}>
            <FormItem label="邮箱"><Input value={form.email} onChange={e => setForm(p=>({...p,email:e.target.value}))} placeholder="email@example.com" /></FormItem>
            <FormItem label="手机号"><Input value={form.phone} onChange={e => setForm(p=>({...p,phone:e.target.value}))} placeholder="请输入手机号" /></FormItem>
          </FormRow>
          <FormRow cols={2}>
            <FormItem label="角色" required><Select value={form.role} onChange={e => setForm(p=>({...p,role:e.target.value}))} options={['超级管理员','运营管理员','开发人员','财务人员','访客']} placeholder="请选择角色" /></FormItem>
            <FormItem label="部门"><Select value={form.dept} onChange={e => setForm(p=>({...p,dept:e.target.value}))} options={['技术部','运营部','市场部','财务部','人事部']} placeholder="请选择部门" /></FormItem>
          </FormRow>
          <FormItem label="状态">
            <RadioGroup value={form.status} onChange={v => setForm(p=>({...p,status:v}))}
              options={[{value:'active',label:'启用'},{value:'inactive',label:'禁用'}]} />
          </FormItem>
          {!editUser && <FormItem label="初始密码" required help="用户首次登录后需修改密码"><Input type="password" placeholder="请设置初始密码" /></FormItem>}
        </Form>
      </Modal>

      <ConfirmModal open={confirmOpen} onClose={() => setConfirmOpen(false)} onConfirm={handleDelete}
        title="删除用户" type="danger" loading={loading}
        message={`确认删除用户「${deleteTarget?.name}」？此操作不可恢复，该用户的所有数据将被清除。`}
        confirmText="确认删除" />
    </div>
  );
};
Object.assign(window, { UserManagement });
