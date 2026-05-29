export default function TimeSlider({ value, onChange }) {
  return (
    <div style={{
      position: 'absolute',
      top: '20px',
      right: '20px',
      backgroundColor: 'white',
      padding: '15px',
      borderRadius: '8px',
      boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
      zIndex: 1000,
      maxWidth: '250px'
    }}>
      <label style={{ display: 'block', marginBottom: '10px', fontWeight: 'bold', fontSize: '14px' }}>
        Max Walking Time: {value} minutes
      </label>
      <input
        type="range"
        min="5"
        max="120"
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        style={{ width: '100%', cursor: 'pointer' }}
      />
      <div style={{ fontSize: '12px', color: '#666', marginTop: '8px' }}>
        Adjust to explore different route options
      </div>
    </div>
  );
}
