// components/Pagination.jsx
const Pagination = ({
  total = 0, page = 1, pageSize = 10,
  onChange, onPageSizeChange,
  showTotal = true, showSizeChanger = true,
  pageSizeOptions = [10, 20, 50, 100],
}) => {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  const getPages = () => {
    if (totalPages <= 7) return Array.from({ length: totalPages }, (_, i) => i + 1);
    const pages = [1];
    if (page > 3) pages.push('...');
    for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) pages.push(i);
    if (page < totalPages - 2) pages.push('...');
    pages.push(totalPages);
    return pages;
  };

  const PageBtn = ({ p, current }) => {
    const [hov, setHov] = React.useState(false);
    const isActive = p === current;
    const isEllipsis = p === '...';
    return (
      <button
        onClick={() => !isEllipsis && onChange && onChange(p)}
        disabled={isEllipsis}
        onMouseEnter={() => setHov(true)}
        onMouseLeave={() => setHov(false)}
        style={{
          minWidth: '32px', height: '32px', padding: '0 4px',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          borderRadius: 'var(--radius-md)', border: '1px solid',
          fontSize: '13px', fontWeight: isActive ? '600' : '400',
          cursor: isEllipsis ? 'default' : 'pointer',
          fontFamily: 'var(--font-sans)',
          transition: 'all 0.12s',
          background: isActive ? 'var(--primary)' : (hov && !isEllipsis) ? 'var(--bg-page)' : 'var(--bg-surface)',
          color: isActive ? '#fff' : 'var(--text-primary)',
          borderColor: isActive ? 'var(--primary)' : 'var(--border)',
        }}
      >{p}</button>
    );
  };

  const NavBtn = ({ dir, disabled: dis }) => {
    const [hov, setHov] = React.useState(false);
    return (
      <button
        onClick={() => !dis && onChange && onChange(dir === 'prev' ? page - 1 : page + 1)}
        disabled={dis}
        onMouseEnter={() => setHov(true)}
        onMouseLeave={() => setHov(false)}
        style={{
          width: '32px', height: '32px', display: 'flex', alignItems: 'center', justifyContent: 'center',
          borderRadius: 'var(--radius-md)', border: '1px solid var(--border)',
          background: hov && !dis ? 'var(--bg-page)' : 'var(--bg-surface)',
          color: dis ? 'var(--text-disabled)' : 'var(--text-secondary)',
          cursor: dis ? 'not-allowed' : 'pointer',
          transition: 'all 0.12s', opacity: dis ? 0.5 : 1,
          fontFamily: 'var(--font-sans)',
        }}
      >
        <i className={`bi bi-chevron-${dir === 'prev' ? 'left' : 'right'}`} style={{ fontSize: '13px' }} />
      </button>
    );
  };

  const start = total === 0 ? 0 : (page - 1) * pageSize + 1;
  const end = Math.min(page * pageSize, total);

  return (
    <div style={{
      display: 'flex', alignItems: 'center', justifyContent: 'space-between',
      flexWrap: 'wrap', gap: '12px',
    }}>
      {showTotal && (
        <span style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
          共 <strong style={{ color: 'var(--text-primary)' }}>{total}</strong> 条，第 {start}–{end} 条
        </span>
      )}
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', flexWrap: 'wrap' }}>
        {showSizeChanger && (
          <select
            value={pageSize}
            onChange={e => onPageSizeChange && onPageSizeChange(Number(e.target.value))}
            style={{
              height: '32px', padding: '0 24px 0 10px', fontSize: '13px',
              border: '1px solid var(--border)', borderRadius: 'var(--radius-md)',
              background: 'var(--bg-surface)', color: 'var(--text-primary)',
              cursor: 'pointer', outline: 'none', fontFamily: 'var(--font-sans)',
              appearance: 'none',
            }}
          >
            {pageSizeOptions.map(s => <option key={s} value={s}>{s} 条/页</option>)}
          </select>
        )}
        <NavBtn dir="prev" disabled={page <= 1} />
        {getPages().map((p, i) => <PageBtn key={i} p={p} current={page} />)}
        <NavBtn dir="next" disabled={page >= totalPages} />
      </div>
    </div>
  );
};

Object.assign(window, { Pagination });
