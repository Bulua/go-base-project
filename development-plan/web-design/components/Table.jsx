// components/Table.jsx
const Table = ({
  columns = [], dataSource = [], rowKey = 'id',
  loading = false, selectable = false,
  selectedKeys = [], onSelectChange,
  onSort, emptyText = '暂无数据',
  style: xStyle,
}) => {
  const [sortCol, setSortCol] = React.useState(null);
  const [sortDir, setSortDir] = React.useState('asc');

  const handleSort = col => {
    const newDir = sortCol === col.key && sortDir === 'asc' ? 'desc' : 'asc';
    setSortCol(col.key);
    setSortDir(newDir);
    onSort && onSort(col.key, newDir);
  };

  const allKeys = dataSource.map(r => r[rowKey]);
  const allSelected = allKeys.length > 0 && allKeys.every(k => selectedKeys.includes(k));
  const someSelected = allKeys.some(k => selectedKeys.includes(k)) && !allSelected;

  const toggleAll = () => {
    if (!onSelectChange) return;
    onSelectChange(allSelected ? [] : allKeys);
  };
  const toggleRow = key => {
    if (!onSelectChange) return;
    onSelectChange(selectedKeys.includes(key)
      ? selectedKeys.filter(k => k !== key)
      : [...selectedKeys, key]);
  };

  const thStyle = {
    padding: '10px 14px', textAlign: 'left', fontWeight: '600',
    fontSize: '12px', color: 'var(--text-secondary)',
    background: 'var(--bg-page)', borderBottom: '1px solid var(--border)',
    whiteSpace: 'nowrap', userSelect: 'none',
    textTransform: 'uppercase', letterSpacing: '0.04em',
    position: 'sticky', top: 0, zIndex: 1,
  };

  const SkeletonRows = () => (
    <>
      {[...Array(5)].map((_, i) => (
        <tr key={i}>
          {selectable && <td style={{ padding: '12px 14px', borderBottom: '1px solid var(--border)' }}><div className="bp-skeleton" style={{ width: 15, height: 15, borderRadius: 3 }} /></td>}
          {columns.map((c, j) => (
            <td key={j} style={{ padding: '12px 14px', borderBottom: '1px solid var(--border)' }}>
              <div className="bp-skeleton" style={{ height: 14, width: `${55 + Math.random() * 35}%`, borderRadius: 4 }} />
            </td>
          ))}
        </tr>
      ))}
    </>
  );

  return (
    <div style={{ overflowX: 'auto', ...xStyle }}>
      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '13px' }}>
        <thead>
          <tr>
            {selectable && (
              <th style={{ ...thStyle, width: '44px' }}>
                <Checkbox checked={allSelected} indeterminate={someSelected} onChange={toggleAll} />
              </th>
            )}
            {columns.map(col => (
              <th
                key={col.key || col.dataIndex}
                style={{
                  ...thStyle,
                  width: col.width,
                  textAlign: col.align || 'left',
                  cursor: col.sortable ? 'pointer' : 'default',
                }}
                onClick={() => col.sortable && handleSort(col)}
              >
                <span style={{ display: 'inline-flex', alignItems: 'center', gap: '4px' }}>
                  {col.title}
                  {col.sortable && (
                    <span style={{ display: 'flex', flexDirection: 'column', gap: '1px', marginLeft: '2px' }}>
                      <i className="bi bi-caret-up-fill" style={{ fontSize: '8px', color: sortCol === col.key && sortDir === 'asc' ? 'var(--primary)' : 'var(--border-strong)' }} />
                      <i className="bi bi-caret-down-fill" style={{ fontSize: '8px', color: sortCol === col.key && sortDir === 'desc' ? 'var(--primary)' : 'var(--border-strong)' }} />
                    </span>
                  )}
                </span>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {loading ? (
            <SkeletonRows />
          ) : dataSource.length === 0 ? (
            <tr>
              <td colSpan={columns.length + (selectable ? 1 : 0)}>
                <Empty description={emptyText} />
              </td>
            </tr>
          ) : (
            dataSource.map((row, ri) => {
              const key = row[rowKey] || ri;
              const isSelected = selectedKeys.includes(key);
              return (
                <TableRow
                  key={key} row={row} columns={columns} rowKey={rowKey}
                  selectable={selectable} selected={isSelected}
                  onToggle={() => toggleRow(key)}
                  isEven={ri % 2 === 0}
                />
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
};

const TableRow = ({ row, columns, rowKey, selectable, selected, onToggle, isEven }) => {
  const [hov, setHov] = React.useState(false);
  const tdStyle = {
    padding: '11px 14px',
    borderBottom: '1px solid var(--border)',
    verticalAlign: 'middle',
    color: 'var(--text-primary)',
    transition: 'background 0.1s',
  };
  const bg = selected ? 'var(--primary-light)' : hov ? 'var(--bg-page)' : 'transparent';

  return (
    <tr
      onMouseEnter={() => setHov(true)}
      onMouseLeave={() => setHov(false)}
      style={{ background: bg }}
    >
      {selectable && (
        <td style={tdStyle}>
          <Checkbox checked={selected} onChange={onToggle} />
        </td>
      )}
      {columns.map(col => {
        const val = col.dataIndex ? row[col.dataIndex] : row[col.key];
        return (
          <td key={col.key || col.dataIndex} style={{ ...tdStyle, textAlign: col.align || 'left', width: col.width }}>
            {col.render ? col.render(val, row) : (
              <span style={{ fontSize: '13px' }}>{val ?? '—'}</span>
            )}
          </td>
        );
      })}
    </tr>
  );
};

// ── TableToolbar — above table ─────────────────
const TableToolbar = ({ left, right, style: xStyle }) => (
  <div style={{
    display: 'flex', alignItems: 'center', justifyContent: 'space-between',
    gap: '12px', flexWrap: 'wrap', marginBottom: '14px', ...xStyle,
  }}>
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>{left}</div>
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>{right}</div>
  </div>
);

Object.assign(window, { Table, TableToolbar });
