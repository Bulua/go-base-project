// pages/MenuManagement.jsx
const MenuManagement = () => {
  const [modalOpen, setModalOpen] = React.useState(false);
  const [editMenu, setEditMenu] = React.useState(null);
  const [confirmOpen, setConfirmOpen] = React.useState(false);
  const [loading, setLoading] = React.useState(false);
  const [expandedRows, setExpandedRows] = React.useState([1, 2]);

  const menus = [
    { id:1, name:'系统管理', icon:'bi-gear', path:'/system', component:null, sort:1, visible:true, type:'dir', children:[
      { id:11, name:'用户管理', icon:'bi-people', path:'/system/user', component:'UserManagement', sort:1, visible:true, type:'menu', children:[] },
      { id:12, name:'角色管理', icon:'bi-shield-check', path:'/system/role', component:'RoleManagement', sort:2, visible:true, type:'menu', children:[] },
      { id:13, name:'菜单管理', icon:'bi-list-check', path:'/system/menu', component:'MenuManagement', sort:3, visible:true, type:'menu', children:[] },
      { id:14, name:'API管理', icon:'bi-cloud-upload', path:'/system/api', component:'ApiManagement', sort:4, visible:true, type:'menu', children:[] },
      { id:15, name:'字典管理', icon:'bi-book', path:'/system/dict', component:'DictManagement', sort:5, visible:true, type:'menu', children:[] },
    ]},
    { id:2, name:'日志审计', icon:'bi-journal-text', path:'/log', component:null, sort:2, visible:true, type:'dir', children:[
      { id:21, name:'操作历史', icon:'bi-clock-history', path:'/log/operation', component:'OperationHistory', sort:1, visible:true, type:'menu', children:[] },
      { id:22, name:'登录日志', icon:'bi-door-open', path:'/log/login', component:'LoginLogs', sort:2, visible:true, type:'menu', children:[] },
    ]},
    { id:3, name:'仪表盘', icon:'bi-speedometer2', path:'/dashboard', component:'Dashboard', sort:0, visible:true, type:'menu', children:[] },
  ];

  const typeMap = { dir:'目录', menu:'菜单', btn:'按钮' };
  const typeColorMap = { dir:'default', menu:'primary', btn:'info' };

  const toggleExpand = id => setExpandedRows(p => p.includes(id) ? p.filter(i=>i!==id) : [...p, id]);

  const handleSave = () => {
    setLoading(true);
    setTimeout(() => { setLoading(false); setModalOpen(false); window.toast.success(editMenu ? '菜单更新成功' : '菜单创建成功'); }, 700);
  };

  const renderRows = (items, depth = 0) => items.flatMap(item => {
    const isExpanded = expandedRows.includes(item.id);
    const hasChildren = item.children && item.children.length > 0;
    const row = (
      <tr key={item.id}
        onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-page)'}
        onMouseLeave={e => e.currentTarget.style.background = ''}
        style={{ cursor:'default' }}
      >
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)', width:'280px' }}>
          <div style={{ display:'flex', alignItems:'center', gap:'8px', paddingLeft: depth * 20 }}>
            {hasChildren ? (
              <span onClick={() => toggleExpand(item.id)} style={{ cursor:'pointer', color:'var(--text-tertiary)', fontSize:'12px', width:'14px', display:'inline-flex' }}>
                <i className={`bi bi-chevron-${isExpanded ? 'down' : 'right'}`} />
              </span>
            ) : <span style={{ width:'14px' }} />}
            <i className={`bi ${item.icon}`} style={{ color:'var(--primary)', fontSize:'14px', flexShrink:0 }} />
            <span style={{ fontSize:'13px', fontWeight: depth===0 ? '600' : '400' }}>{item.name}</span>
          </div>
        </td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)' }}><Tag color={typeColorMap[item.type]}>{typeMap[item.type]}</Tag></td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)', fontFamily:'var(--font-mono)', fontSize:'12px', color:'var(--text-secondary)' }}>{item.path || '—'}</td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)', fontSize:'12px', color:'var(--text-tertiary)' }}>{item.component || '—'}</td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)', textAlign:'center' }}>{item.sort}</td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)' }}>
          <div style={{ width:'36px', height:'20px', background: item.visible ? 'var(--primary)' : 'var(--border-strong)', borderRadius:'10px', position:'relative', cursor:'pointer', transition:'background 0.2s' }}>
            <div style={{ position:'absolute', top:'2px', left: item.visible ? '18px' : '2px', width:'16px', height:'16px', background:'#fff', borderRadius:'50%', transition:'left 0.2s', boxShadow:'0 1px 3px rgba(0,0,0,0.2)' }} />
          </div>
        </td>
        <td style={{ padding:'10px 14px', borderBottom:'1px solid var(--border)' }}>
          <div style={{ display:'flex', gap:'4px' }}>
            <Button size="xs" variant="ghost" icon={<i className="bi bi-pencil" />} onClick={() => { setEditMenu(item); setModalOpen(true); }}>编辑</Button>
            <Button size="xs" variant="secondary" icon={<i className="bi bi-plus" />} onClick={() => { setEditMenu(null); setModalOpen(true); }}>子菜单</Button>
            <Button size="xs" variant="text" danger icon={<i className="bi bi-trash" />} onClick={() => setConfirmOpen(true)} />
          </div>
        </td>
      </tr>
    );
    const childRows = (hasChildren && isExpanded) ? renderRows(item.children, depth + 1) : [];
    return [row, ...childRows];
  });

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">菜单管理</div>
          <div className="bp-page-subtitle">管理系统导航菜单与页面路由</div>
        </div>
        <div style={{ display:'flex', gap:'8px' }}>
          <Button variant="secondary" icon={<i className="bi bi-arrows-expand" />} onClick={() => setExpandedRows(menus.map(m=>m.id))}>展开全部</Button>
          <Button icon={<i className="bi bi-plus-lg" />} onClick={() => { setEditMenu(null); setModalOpen(true); }}>新增菜单</Button>
        </div>
      </div>

      <Card padding="0">
        <div style={{ overflowX:'auto' }}>
          <table style={{ width:'100%', borderCollapse:'collapse', fontSize:'13px' }}>
            <thead>
              <tr>
                {['菜单名称','类型','路由路径','组件','排序','显示','操作'].map((h,i) => (
                  <th key={h} style={{ padding:'10px 14px', textAlign: i===4 ? 'center' : 'left', fontWeight:'600', fontSize:'12px', color:'var(--text-secondary)', background:'var(--bg-page)', borderBottom:'1px solid var(--border)', whiteSpace:'nowrap', textTransform:'uppercase', letterSpacing:'0.04em', position:'sticky', top:0, zIndex:1 }}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>{renderRows(menus)}</tbody>
          </table>
        </div>
      </Card>

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editMenu ? '编辑菜单' : '新增菜单'} width={560}
        footer={<><Button variant="secondary" onClick={() => setModalOpen(false)}>取消</Button><Button onClick={handleSave} loading={loading}>保存</Button></>}>
        <Form style={{ gap:'16px' }}>
          <FormItem label="菜单类型" required>
            <RadioGroup value="menu" options={[{value:'dir',label:'目录'},{value:'menu',label:'菜单'},{value:'btn',label:'按钮'}]} />
          </FormItem>
          <FormRow cols={2}>
            <FormItem label="菜单名称" required><Input placeholder="菜单显示名称" defaultValue={editMenu?.name} /></FormItem>
            <FormItem label="图标" help="Bootstrap Icons 类名"><Input placeholder="bi-gear" defaultValue={editMenu?.icon} style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          </FormRow>
          <FormRow cols={2}>
            <FormItem label="路由路径"><Input placeholder="/module/page" defaultValue={editMenu?.path} style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
            <FormItem label="组件名称"><Input placeholder="ComponentName" defaultValue={editMenu?.component} style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          </FormRow>
          <FormRow cols={2}>
            <FormItem label="上级菜单"><Select options={['系统管理','日志审计']} placeholder="根目录" /></FormItem>
            <FormItem label="排序"><Input type="number" defaultValue={editMenu?.sort || 1} /></FormItem>
          </FormRow>
          <FormItem label="是否显示">
            <RadioGroup value="true" options={[{value:'true',label:'显示'},{value:'false',label:'隐藏'}]} />
          </FormItem>
        </Form>
      </Modal>

      <ConfirmModal open={confirmOpen} onClose={() => setConfirmOpen(false)}
        onConfirm={() => { setLoading(true); setTimeout(()=>{setLoading(false);setConfirmOpen(false);window.toast.success('菜单已删除');},600); }}
        title="删除菜单" type="danger" loading={loading}
        message="确认删除该菜单？子菜单将一并删除，已配置的路由将失效。" confirmText="确认删除" />
    </div>
  );
};
Object.assign(window, { MenuManagement });
