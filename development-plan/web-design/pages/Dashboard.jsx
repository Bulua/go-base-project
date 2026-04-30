// pages/Dashboard.jsx
const Dashboard = () => {
  const stats = [
    { title: '注册用户', value: '12,847', unit: '人', trend: 8.2, trendLabel: '较上月', icon: 'bi-people-fill', color: 'var(--primary)', iconBg: 'var(--primary-light)' },
    { title: '今日访问', value: '3,291', unit: '次', trend: 14.5, trendLabel: '较昨日', icon: 'bi-eye-fill', color: 'oklch(0.52 0.15 145)', iconBg: 'oklch(0.96 0.04 145)' },
    { title: 'API 调用', value: '186.4', unit: '万次', trend: -2.1, trendLabel: '较昨日', icon: 'bi-cloud-fill', color: 'oklch(0.52 0.15 285)', iconBg: 'oklch(0.96 0.04 285)' },
    { title: '系统角色', value: '18', unit: '个', trend: 0, trendLabel: '无变化', icon: 'bi-shield-fill-check', color: 'oklch(0.68 0.16 72)', iconBg: 'oklch(0.97 0.04 72)' },
  ];

  const recentUsers = [
    { id: 1, name: '张伟', email: 'zhangwei@example.com', role: '管理员', status: 'active', created: '2024-01-15' },
    { id: 2, name: '李娜', email: 'lina@example.com', role: '运营', status: 'active', created: '2024-01-14' },
    { id: 3, name: '王磊', email: 'wanglei@example.com', role: '开发', status: 'inactive', created: '2024-01-13' },
    { id: 4, name: '陈静', email: 'chenjing@example.com', role: '运营', status: 'active', created: '2024-01-12' },
    { id: 5, name: '刘阳', email: 'liuyang@example.com', role: '访客', status: 'pending', created: '2024-01-11' },
  ];

  const recentLogs = [
    { id: 1, user: 'admin', action: '修改用户权限', module: '用户管理', time: '2分钟前', result: 'success' },
    { id: 2, user: 'zhangwei', action: '导出数据报表', module: '报表中心', time: '15分钟前', result: 'success' },
    { id: 3, user: 'lina', action: '删除字典项', module: '字典管理', time: '1小时前', result: 'failed' },
    { id: 4, user: 'admin', action: '新增API接口', module: 'API管理', time: '2小时前', result: 'success' },
    { id: 5, user: 'wanglei', action: '重置用户密码', module: '用户管理', time: '3小时前', result: 'success' },
  ];

  const quickActions = [
    { icon: 'bi-person-plus', label: '新增用户', color: 'var(--primary)', bg: 'var(--primary-light)', action: () => window.toast.success('跳转到新增用户') },
    { icon: 'bi-shield-plus', label: '新增角色', color: 'oklch(0.52 0.15 145)', bg: 'oklch(0.96 0.04 145)', action: () => window.toast.success('跳转到新增角色') },
    { icon: 'bi-cloud-plus', label: '新增接口', color: 'oklch(0.52 0.15 285)', bg: 'oklch(0.96 0.04 285)', action: () => window.toast.success('跳转到新增接口') },
    { icon: 'bi-journal-plus', label: '新增字典', color: 'oklch(0.68 0.16 72)', bg: 'oklch(0.97 0.04 72)', action: () => window.toast.success('跳转到新增字典') },
  ];

  const userColumns = [
    { key: 'name', title: '姓名', dataIndex: 'name', render: (v, r) => (
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <div style={{ width: '28px', height: '28px', borderRadius: '50%', background: 'var(--primary-light)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '12px', fontWeight: '700', color: 'var(--primary)', flexShrink: 0 }}>
          {v.charAt(0)}
        </div>
        <div>
          <div style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>{v}</div>
          <div style={{ fontSize: '11px', color: 'var(--text-tertiary)' }}>{r.email}</div>
        </div>
      </div>
    )},
    { key: 'role', title: '角色', dataIndex: 'role', render: v => <Tag color="primary">{v}</Tag> },
    { key: 'status', title: '状态', dataIndex: 'status', render: v => <StatusBadge status={v} /> },
    { key: 'created', title: '注册日期', dataIndex: 'created', render: v => <span style={{ color: 'var(--text-secondary)', fontSize: '12px' }}>{v}</span> },
  ];

  const logColumns = [
    { key: 'user', title: '操作人', dataIndex: 'user', render: v => <code style={{ fontFamily: 'var(--font-mono)', fontSize: '12px', background: 'var(--bg-page)', padding: '2px 6px', borderRadius: 4 }}>{v}</code> },
    { key: 'action', title: '操作内容', dataIndex: 'action' },
    { key: 'module', title: '模块', dataIndex: 'module', render: v => <Tag>{v}</Tag> },
    { key: 'time', title: '时间', dataIndex: 'time', render: v => <span style={{ color: 'var(--text-tertiary)', fontSize: '12px' }}>{v}</span> },
    { key: 'result', title: '结果', dataIndex: 'result', render: v => <StatusBadge status={v} /> },
  ];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', animation: 'bpFadeIn 0.25s ease' }}>
      {/* Greeting */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: '10px' }}>
        <div>
          <h1 style={{ fontSize: '20px', fontWeight: '700', color: 'var(--text-primary)' }}>仪表盘</h1>
          <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginTop: '2px' }}>
            {new Date().toLocaleDateString('zh-CN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })} · 数据概览
          </p>
        </div>
        <Button
          icon={<i className="bi bi-arrow-clockwise" />}
          variant="secondary" size="sm"
          onClick={() => window.toast.success('数据已刷新')}
        >刷新数据</Button>
      </div>

      {/* Stat cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '16px' }}>
        {stats.map(s => <StatCard key={s.title} {...s} />)}
      </div>

      {/* Quick actions */}
      <SectionCard title="快捷操作" subtitle="常用功能入口">
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: '12px' }}>
          {quickActions.map(a => (
            <QuickActionBtn key={a.label} {...a} />
          ))}
        </div>
      </SectionCard>

      {/* Two-column row */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        {/* Trend chart placeholder */}
        <SectionCard title="访问趋势" subtitle="近 7 天日均访问量" action={<Tag color="primary">本周</Tag>}>
          <MiniChart />
        </SectionCard>

        {/* System overview */}
        <SectionCard title="系统状态" subtitle="资源使用率">
          <SystemStatus />
        </SectionCard>
      </div>

      {/* Tables row */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        <SectionCard title="最新用户" subtitle="最近注册的用户">
          <Table columns={userColumns} dataSource={recentUsers} rowKey="id" />
        </SectionCard>
        <SectionCard title="操作记录" subtitle="最近的系统操作">
          <Table columns={logColumns} dataSource={recentLogs} rowKey="id" />
        </SectionCard>
      </div>
    </div>
  );
};

const QuickActionBtn = ({ icon, label, color, bg, action }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <div
      onClick={action}
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{
        display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px',
        padding: '16px 12px', borderRadius: 'var(--radius-lg)',
        borderLeft: `${hov ? '3px' : '1px'} solid ${hov ? color : 'var(--border)'}`,
        borderTop: '1px solid var(--border)',
        borderRight: '1px solid var(--border)',
        borderBottom: '1px solid var(--border)',
        background: hov ? bg : 'var(--bg-surface)',
        cursor: 'pointer', transition: 'all 0.15s',
        boxShadow: hov ? 'var(--shadow-sm)' : 'none',
      }}
    >
      <div style={{ fontSize: '22px', color }}>
        <i className={`bi ${icon}`} />
      </div>
      <span style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>{label}</span>
    </div>
  );
};

const MiniChart = () => {
  const data = [42, 67, 58, 89, 73, 94, 81];
  const days = ['周一','周二','周三','周四','周五','周六','周日'];
  const max = Math.max(...data);
  return (
    <div style={{ display: 'flex', alignItems: 'flex-end', gap: '8px', height: '120px', padding: '0 4px' }}>
      {data.map((v, i) => (
        <div key={i} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '5px', height: '100%', justifyContent: 'flex-end' }}>
          <span style={{ fontSize: '11px', color: 'var(--text-tertiary)' }}>{v}</span>
          <BarItem height={`${(v / max) * 72}px`} active={i === 6} />
          <span style={{ fontSize: '10px', color: 'var(--text-tertiary)' }}>{days[i]}</span>
        </div>
      ))}
    </div>
  );
};

const BarItem = ({ height, active }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <div
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{
        width: '100%', height, borderRadius: 'var(--radius-sm) var(--radius-sm) 0 0',
        background: active || hov ? 'var(--primary)' : 'var(--primary-dim)',
        transition: 'background 0.15s, height 0.3s',
      }}
    />
  );
};

const SystemStatus = () => {
  const items = [
    { label: 'CPU 使用率', value: 34, color: 'var(--success)' },
    { label: '内存占用', value: 67, color: 'var(--warning)' },
    { label: '磁盘使用', value: 45, color: 'var(--primary)' },
    { label: '网络带宽', value: 22, color: 'oklch(0.52 0.15 285)' },
  ];
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
      {items.map(item => (
        <div key={item.label} style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>{item.label}</span>
            <span style={{ fontSize: '13px', fontWeight: '600', color: 'var(--text-primary)' }}>{item.value}%</span>
          </div>
          <div style={{ height: '6px', background: 'var(--bg-page)', borderRadius: 'var(--radius-full)', overflow: 'hidden' }}>
            <div style={{ height: '100%', width: `${item.value}%`, background: item.color, borderRadius: 'var(--radius-full)', transition: 'width 0.6s ease' }} />
          </div>
        </div>
      ))}
    </div>
  );
};

Object.assign(window, { Dashboard });
