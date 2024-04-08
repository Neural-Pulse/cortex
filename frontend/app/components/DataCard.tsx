// components/ElasticsearchCard.tsx
import React from 'react';

function DataCard({ data }) {
  return (
    <div className="bg-white p-4 rounded shadow">
      <h2 className="text-lg font-bold mb-2">{data?.table_name}</h2>
      <p>{data?.description}</p>
      <div className="text-sm">
        <p><strong>Database:</strong> {data?.database_name}</p>
        <p><strong>Column:</strong> {data?.column_name}</p>
        <p><strong>Health:</strong> {data?.health}</p>
        <p><strong>Data Classification:</strong> {data?.data_classification}</p>
      </div>
    </div>
  );
}

export default DataCard;
