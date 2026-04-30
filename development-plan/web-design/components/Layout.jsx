// components/Layout.jsx — Sidebar + Navbar

const NAV_ITEMS = [
  { id: 'dashboard', label: '仪表盘', icon: 'bi-speedometer2' },
  {
    id: 'system', label: '系统管理', icon: 'bi-grid-3x3-gap',
    children: [
      { id: 'user',     label: '用户管理', icon: 'bi-people'         },
      { id: 'role',     label: '角色管理', icon: 'bi-shield-check'   },
      { id: 'menu',     label: '菜单管理', icon: 'bi-list-check'     },
      { id: 'api',      label: 'API 管理', icon: 'bi-cloud-upload'   },
      { id: 'dict',     label: '字典管理', icon: 'bi-book'           },
      { id: 'oplog',    label: '操作历史', icon: 'bi-clock-history'  },
      { id: 'loginlog', label: '登录日志', icon: 'bi-door-open'      },
      { id: 'config',   label: '系统配置', icon: 'bi-sliders'        },
    ],
  },
];

// ── Sidebar ────────────────────────────────────
const Sidebar = ({ currentPage, onNavigate, collapsed, onToggleCollapse }) => {
  const [openGroups, setOpenGroups] = React.useState(['system']);

  const toggleGroup = id => {
    setOpenGroups(prev => prev.includes(id) ? prev.filter(g => g !== id) : [...prev, id]);
  };

  const isChildActive = item => item.children && item.children.some(c => c.id === currentPage);

  return (
    <aside className="bp-sidebar" style={{
      width: collapsed ? 'var(--sidebar-col)' : 'var(--sidebar-w)',
      flexShrink: 0,
      background: 'var(--sidebar-bg)',
      display: 'flex', flexDirection: 'column',
      transition: 'width var(--t-slow)',
      overflow: 'hidden',
      position: 'relative', zIndex: 100,
      height: '100vh',
    }}>
      {/* Logo */}
      <div style={{
        height: 'var(--navbar-h)', display: 'flex', alignItems: 'center',
        padding: collapsed ? '0' : '0 20px',
        justifyContent: collapsed ? 'center' : 'flex-start',
        borderBottom: '1px solid var(--sidebar-divider)',
        flexShrink: 0, gap: '10px', overflow: 'hidden',
      }}>
        <div style={{
          width: '32px', height: '32px', borderRadius: '0px',
          background: 'var(--primary)', flexShrink: 0,
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          color: '#fff', fontWeight: '800', fontSize: '14px',
          clipPath: 'polygon(0 0,100% 0,100% 100%,0 100%)',
        }}>B</div>
        {!collapsed && (
          <span style={{
            color: 'var(--sidebar-logo-text)', fontWeight: '700', fontSize: '15px',
            whiteSpace: 'nowrap', letterSpacing: '-0.01em',
          }}>BaseProject</span>
        )}
      </div>

      {/* Navigation */}
      <nav style={{ flex: 1, overflowY: 'auto', overflowX: 'hidden', padding: '10px 0' }}>
        {NAV_ITEMS.map(item => (
          <NavItem
            key={item.id} item={item}
            currentPage={currentPage}
            onNavigate={onNavigate}
            collapsed={collapsed}
            isOpen={openGroups.includes(item.id)}
            onToggle={() => toggleGroup(item.id)}
            isChildActive={isChildActive(item)}
          />
        ))}
      </nav>

      {/* Collapse button */}
      <div style={{
        borderTop: '1px solid var(--sidebar-divider)',
        padding: '10px', display: 'flex',
        justifyContent: collapsed ? 'center' : 'flex-end',
      }}>
        <SidebarCollapseBtn collapsed={collapsed} onClick={onToggleCollapse} />
      </div>
    </aside>
  );
};

const SidebarCollapseBtn = ({ collapsed, onClick }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <button
      onClick={onClick}
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{
        width: '32px', height: '32px', borderRadius: 'var(--radius-md)',
        background: hov ? 'rgba(255,255,255,0.1)' : 'transparent',
        border: '1px solid var(--sidebar-divider)',
        color: 'var(--sidebar-item-text)', cursor: 'pointer',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        transition: 'background 0.15s', flexShrink: 0,
      }}
    >
      <i className={`bi bi-chevron-${collapsed ? 'right' : 'left'}`} style={{ fontSize: '12px' }} />
    </button>
  );
};

const NavItem = ({ item, currentPage, onNavigate, collapsed, isOpen, onToggle, isChildActive }) => {
  const isActive = item.id === currentPage;
  const hasChildren = item.children && item.children.length > 0;

  const itemStyle = (active, hov) => ({
    display: 'flex', alignItems: 'center',
    padding: collapsed ? '0' : '0 12px 0 14px',
    height: '38px',
    justifyContent: collapsed ? 'center' : 'flex-start',
    gap: '9px', cursor: 'pointer',
    borderRadius: '0px',
    margin: '1px 0',
    borderLeft: active ? '3px solid var(--primary)' : '3px solid transparent',
    background: active
      ? 'rgba(255,255,255,0.08)'
      : hov ? 'var(--sidebar-item-hover)' : 'transparent',
    color: active ? 'var(--sidebar-active-text)' : 'var(--sidebar-item-text)',
    transition: 'background 0.15s, color 0.15s, border-color 0.15s',
    userSelect: 'none',
    position: 'relative',
  });

  const [hov, setHov] = React.useState(false);
  const active = isActive || (hasChildren && isChildActive);

  return (
    <div>
      <div
        style={itemStyle(active && !hasChildren, hov)}
        onClick={() => hasChildren ? onToggle() : onNavigate(item.id)}
        onMouseEnter={() => setHov(true)}
        onMouseLeave={() => setHov(false)}
        title={collapsed ? item.label : undefined}
      >
        <i className={`bi ${item.icon}`} style={{ fontSize: '15px', flexShrink: 0, opacity: active ? 1 : 0.85 }} />
        {!collapsed && (
          <>
            <span style={{ flex: 1, fontSize: '13px', fontWeight: active ? '600' : '400', whiteSpace: 'nowrap' }}>
              {item.label}
            </span>
            {hasChildren && (
              <i className={`bi bi-chevron-${isOpen ? 'up' : 'down'}`}
                style={{ fontSize: '11px', opacity: 0.6, transition: 'transform 0.2s' }} />
            )}
          </>
        )}
      </div>
      {/* Sub-items */}
      {hasChildren && isOpen && !collapsed && (
        <div style={{ overflow: 'hidden', animation: 'bpSlideUp 0.18s ease' }}>
          {item.children.map(child => (
            <SubNavItem key={child.id} item={child} currentPage={currentPage} onNavigate={onNavigate} />
          ))}
        </div>
      )}
    </div>
  );
};

const SubNavItem = ({ item, currentPage, onNavigate }) => {
  const [hov, setHov] = React.useState(false);
  const isActive = item.id === currentPage;
  return (
    <div
      style={{
        display: 'flex', alignItems: 'center', height: '34px',
        padding: '0 12px 0 38px', gap: '8px',
        margin: '1px 0', borderRadius: '0px',
        borderLeft: isActive ? '3px solid var(--primary)' : '3px solid transparent',
        cursor: 'pointer', userSelect: 'none',
        background: isActive ? 'rgba(255,255,255,0.08)' : hov ? 'var(--sidebar-item-hover)' : 'transparent',
        color: isActive ? 'var(--sidebar-active-text)' : 'var(--sidebar-item-text)',
        transition: 'background 0.15s, color 0.15s, border-color 0.15s',
      }}
      onClick={() => onNavigate(item.id)}
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
    >
      <i className={`bi ${item.icon}`} style={{ fontSize: '13px', flexShrink: 0, opacity: 0.8 }} />
      <span style={{ fontSize: '13px', fontWeight: isActive ? '600' : '400' }}>{item.label}</span>
    </div>
  );
};

// ── Navbar / Top Header ────────────────────────
const Navbar = ({ currentPage, onLogout, username = 'Admin', onToggleDark, isDark }) => {
  const breadcrumbMap = {
    dashboard: ['仪表盘'],
    user:      ['系统管理', '用户管理'],
    role:      ['系统管理', '角色管理'],
    menu:      ['系统管理', '菜单管理'],
    api:       ['系统管理', 'API 管理'],
    dict:      ['系统管理', '字典管理'],
    oplog:     ['系统管理', '操作历史'],
    loginlog:  ['系统管理', '登录日志'],
    config:    ['系统管理', '系统配置'],
  };
  const crumbs = breadcrumbMap[currentPage] || [currentPage];

  const userMenuItems = [
    { icon: 'bi-person-circle', label: '个人中心', onClick: () => window.toast && window.toast.info('个人中心') },
    { icon: 'bi-gear',          label: '账户设置', onClick: () => window.toast && window.toast.info('账户设置') },
    { type: 'divider' },
    { icon: 'bi-box-arrow-right', label: '退出登录', danger: true, onClick: onLogout },
  ];

  return (
    <header className="bp-navbar" style={{
      height: 'var(--navbar-h)', background: 'var(--bg-surface)',
      borderBottom: '1px solid var(--border)',
      display: 'flex', alignItems: 'center',
      padding: '0 20px 0 24px', gap: '12px',
      flexShrink: 0, position: 'sticky', top: 0, zIndex: 99,
    }}>
      {/* Breadcrumb */}
      <div style={{ flex: 1, display: 'flex', alignItems: 'center', gap: '6px', fontSize: '13px', overflow: 'hidden' }}>
        <i className="bi bi-house" style={{ color: 'var(--text-tertiary)', fontSize: '12px' }} />
        {crumbs.map((c, i) => (
          <React.Fragment key={i}>
            {i > 0 && <i className="bi bi-chevron-right" style={{ fontSize: '10px', color: 'var(--text-tertiary)' }} />}
            <span style={{
              color: i === crumbs.length - 1 ? 'var(--text-primary)' : 'var(--text-secondary)',
              fontWeight: i === crumbs.length - 1 ? '600' : '400',
              whiteSpace: 'nowrap',
            }}>{c}</span>
          </React.Fragment>
        ))}
      </div>

      {/* Right actions */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
        {/* Dark mode toggle */}
        <NavIconBtn icon={isDark ? 'bi-sun' : 'bi-moon'} onClick={onToggleDark} title={isDark ? '切换浅色' : '切换深色'} />

        {/* Notifications */}
        <Badge count={3}>
          <NavIconBtn icon="bi-bell" onClick={() => window.toast && window.toast.info('暂无新消息')} title="通知" />
        </Badge>

        {/* User */}
        <Dropdown
          placement="bottom-right"
          trigger={
            <div style={{
              display: 'flex', alignItems: 'center', gap: '8px', padding: '5px 10px',
              borderRadius: 'var(--radius-md)', cursor: 'pointer', marginLeft: '4px',
              transition: 'background 0.15s',
            }}
              onMouseEnter={e => e.currentTarget.style.background = 'var(--bg-page)'}
              onMouseLeave={e => e.currentTarget.style.background = ''}
            >
              <div style={{
                width: '30px', height: '30px', borderRadius: 'var(--radius-full)',
                background: 'var(--primary)', display: 'flex', alignItems: 'center',
                justifyContent: 'center', color: '#fff', fontWeight: '700', fontSize: '12px',
                flexShrink: 0,
              }}>
                {username.charAt(0).toUpperCase()}
              </div>
              <span style={{ fontSize: '13px', fontWeight: '500', color: 'var(--text-primary)' }}>{username}</span>
              <i className="bi bi-chevron-down" style={{ fontSize: '11px', color: 'var(--text-tertiary)' }} />
            </div>
          }
          items={userMenuItems}
        />
      </div>
    </header>
  );
};

const NavIconBtn = ({ icon, onClick, title }) => {
  const [hov, setHov] = React.useState(false);
  return (
    <button
      onClick={onClick} title={title}
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{
        width: '36px', height: '36px', borderRadius: 'var(--radius-md)',
        background: hov ? 'var(--bg-page)' : 'transparent',
        border: 'none', cursor: 'pointer',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        color: hov ? 'var(--text-primary)' : 'var(--text-secondary)',
        transition: 'background 0.15s, color 0.15s', fontSize: '16px',
      }}
    >
      <i className={`bi ${icon}`} />
    </button>
  );
};

Object.assign(window, { Sidebar, Navbar, NAV_ITEMS });
