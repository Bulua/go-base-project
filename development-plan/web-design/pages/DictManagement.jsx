// pages/DictManagement.jsx
const DictManagement = () => {
  const [activeDict, setActiveDict] = React.useState(1);
  const [modalOpen, setModalOpen] = React.useState(false);
  const [itemModalOpen, setItemModalOpen] = React.useState(false);
  const [loading, setLoading] = React.useState(false);

  const dicts = [
    { id:1, name:'用户状态', code:'user_status', description:'用户账号启用状态', itemCount:3, system:true },
    { id:2, name:'性别类型', code:'gender', description:'用户性别选项', itemCount:3, system:true },
    { id:3, name:'部门类型', code:'dept_type', description:'组织架构部门分类', itemCount:5, system:false },
    { id:4, name:'优先级', code:'priority', description:'任务优先级枚举', itemCount:4, system:false },
    { id:5, name:'操作类型', code:'op_type', description:'系统操作日志类型', itemCount:6, system:true },
    { id:6, name:'通知类型', code:'notice_type', description:'系统通知消息分类', itemCount:4, system:false },
  ];

  const dictItems = {
    1: [
      { id:1, label:'已启用', value:'active', sort:1, remark:'正常使用状态', color:'success' },
      { id:2, label:'已禁用', value:'inactive', sort:2, remark:'账号被禁用', color:'danger' },
      { id:3, label:'待审核', value:'pending', sort:3, remark:'等待管理员审核', color:'warning' },
    ],
    2: [
      { id:1, label:'男', value:'male', sort:1, remark:'', color:'primary' },
      { id:2, label:'女', value:'female', sort:2, remark:'', color:'danger' },
      { id:3, label:'未知', value:'unknown', sort:3, remark:'', color:'default' },
    ],
    3: [
      { id:1, label:'技术部', value:'tech', sort:1, remark:'', color:'primary' },
      { id:2, label:'运营部', value:'ops', sort:2, remark:'', color:'info' },
      { id:3, label:'市场部', value:'market', sort:3, remark:'', color:'success' },
      { id:4, label:'财务部', value:'finance', sort:4, remark:'', color:'warning' },
      { id:5, label:'人事部', value:'hr', sort:5, remark:'', color:'default' },
    ],
    4: [
      { id:1, label:'紧急', value:'urgent', sort:1, remark:'', color:'danger' },
      { id:2, label:'高', value:'high', sort:2, remark:'', color:'warning' },
      { id:3, label:'中', value:'medium', sort:3, remark:'', color:'primary' },
      { id:4, label:'低', value:'low', sort:4, remark:'', color:'default' },
    ],
  };

  const currentItems = dictItems[activeDict] || [];
  const currentDict = dicts.find(d => d.id === activeDict);

  const itemColumns = [
    { key:'label', title:'标签', dataIndex:'label', render: (v, r) => <Tag color={r.color}>{v}</Tag> },
    { key:'value', title:'键值', dataIndex:'value', render: v => <code style={{ fontFamily:'var(--font-mono)', fontSize:'12px', background:'var(--bg-page)', padding:'2px 7px', borderRadius:4, color:'var(--text-secondary)' }}>{v}</code> },
    { key:'sort', title:'排序', dataIndex:'sort', render: v => <span style={{ fontFamily:'var(--font-mono)', color:'var(--text-tertiary)' }}>{v}</span> },
    { key:'remark', title:'备注', dataIndex:'remark', render: v => <span style={{ fontSize:'12px', color:'var(--text-tertiary)' }}>{v || '—'}</span> },
    { key:'actions', title:'操作', width:'100px', render:(_, r) => (
      <div style={{ display:'flex', gap:'4px' }}>
        <Button size="xs" variant="ghost" icon={<i className="bi bi-pencil" />} onClick={() => setItemModalOpen(true)}>编辑</Button>
        <Button size="xs" variant="text" danger icon={<i className="bi bi-trash" />} onClick={() => window.toast.success('字典项已删除')} />
      </div>
    )},
  ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:'16px', animation:'bpFadeIn 0.25s ease' }}>
      <div className="bp-page-header">
        <div>
          <div className="bp-page-title">字典管理</div>
          <div className="bp-page-subtitle">管理系统枚举数据与下拉选项</div>
        </div>
        <Button icon={<i className="bi bi-plus-lg" />} onClick={() => setModalOpen(true)}>新增字典</Button>
      </div>

      <div style={{ display:'grid', gridTemplateColumns:'260px 1fr', gap:'16px', alignItems:'start' }}>
        {/* Dict list */}
        <Card padding="8px" title="字典列表">
          <div style={{ display:'flex', flexDirection:'column', gap:'2px' }}>
            {dicts.map(d => (
              <DictListItem key={d.id} dict={d} active={d.id === activeDict} onClick={() => setActiveDict(d.id)} />
            ))}
          </div>
        </Card>

        {/* Dict items */}
        <Card padding="16px"
          title={currentDict ? `${currentDict.name} · 字典项` : '字典项'}
          extra={
            <div style={{ display:'flex', gap:'8px', alignItems:'center' }}>
              {currentDict && <code style={{ fontFamily:'var(--font-mono)', fontSize:'11px', color:'var(--text-tertiary)', background:'var(--bg-page)', padding:'2px 8px', borderRadius:4 }}>{currentDict.code}</code>}
              <Button size="sm" icon={<i className="bi bi-plus" />} onClick={() => setItemModalOpen(true)}>新增项</Button>
            </div>
          }
        >
          <Table columns={itemColumns} dataSource={currentItems} rowKey="id"
            emptyText="该字典暂无数据项，点击「新增项」添加" />
        </Card>
      </div>

      {/* Add Dict Modal */}
      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="新增字典" width={480}
        footer={<><Button variant="secondary" onClick={() => setModalOpen(false)}>取消</Button><Button onClick={() => { setLoading(true); setTimeout(()=>{setLoading(false);setModalOpen(false);window.toast.success('字典创建成功');},700); }} loading={loading}>保存</Button></>}>
        <Form style={{ gap:'16px' }}>
          <FormRow cols={2}>
            <FormItem label="字典名称" required><Input placeholder="如：用户状态" /></FormItem>
            <FormItem label="字典编码" required help="唯一标识，英文+下划线"><Input placeholder="user_status" style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          </FormRow>
          <FormItem label="描述"><Textarea placeholder="字典用途说明" rows={2} /></FormItem>
        </Form>
      </Modal>

      {/* Add/Edit item modal */}
      <Modal open={itemModalOpen} onClose={() => setItemModalOpen(false)} title="编辑字典项" width={440}
        footer={<><Button variant="secondary" onClick={() => setItemModalOpen(false)}>取消</Button><Button onClick={() => { setItemModalOpen(false); window.toast.success('字典项已保存'); }}>保存</Button></>}>
        <Form style={{ gap:'16px' }}>
          <FormRow cols={2}>
            <FormItem label="标签" required><Input placeholder="显示文本" /></FormItem>
            <FormItem label="键值" required><Input placeholder="存储值" style={{ fontFamily:'var(--font-mono)' }} /></FormItem>
          </FormRow>
          <FormRow cols={2}>
            <FormItem label="颜色标签">
              <Select options={[{value:'default',label:'默认'},{value:'primary',label:'主色'},{value:'success',label:'成功'},{value:'warning',label:'警告'},{value:'danger',label:'危险'},{value:'info',label:'信息'}]} placeholder="选择颜色" />
            </FormItem>
            <FormItem label="排序"><Input type="number" defaultValue="1" /></FormItem>
          </FormRow>
          <FormItem label="备注"><Input placeholder="可选说明" /></FormItem>
        </Form>
      </Modal>
    </div>
  );
};

const DictListItem = ({ dict, active, onClick }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <div onClick={onClick}
      onMouseEnter={() => setHov(true)} onMouseLeave={() => setHov(false)}
      style={{
        padding:'10px 12px', borderRadius:'var(--radius-md)', cursor:'pointer', transition:'all 0.15s',
        background: active ? 'var(--primary-light)' : hov ? 'var(--bg-page)' : 'transparent',
        border: active ? '1px solid var(--primary-dim)' : '1px solid transparent',
      }}>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <span style={{ fontSize:'13px', fontWeight: active ? '600' : '400', color: active ? 'var(--primary)' : 'var(--text-primary)' }}>{dict.name}</span>
        <div style={{ display:'flex', alignItems:'center', gap:'5px' }}>
          {dict.system && <Tag color="info" style={{ fontSize:'10px' }}>内置</Tag>}
          <Tag style={{ fontSize:'10px' }}>{dict.itemCount}</Tag>
        </div>
      </div>
      <div style={{ fontSize:'11px', color:'var(--text-tertiary)', fontFamily:'var(--font-mono)', marginTop:'3px' }}>{dict.code}</div>
    </div>
  );
};

Object.assign(window, { DictManagement });
