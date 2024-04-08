'use client';

import https from 'https';
// Certifique-se de que 'use client' seja removido, pois não é necessário e pode causar erros.
import { useState } from "react";
import axios from 'axios';
import DataCard from './components/DataCard'; // Ajuste o caminho conforme necessário

export default function Home() {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);

  const handleSearch = async () => {
    try {
      const response = await axios.post('/api/search', {
        query: {
          query_string: {
            query: searchQuery,
          },
        },
      });
      // Ajuste para extrair os dados de _source
      const results = response.data.hits.hits.map(hit => hit._source);
      console.log(results);
      setSearchResults(results);
    } catch (error) {
      console.error('Error searching:', error);
    }
  };

  return (
    <main className="bg-white min-h-screen">
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold mb-4">Search Page</h1>
        <div className="mb-4">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search..."
            className="border border-gray-300 rounded-lg px-4 py-2 w-96"
          />
          <button
            onClick={handleSearch}
            className="ml-4 bg-blue-500 text-white rounded-lg px-4 py-2"
          >
            Search
          </button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {searchResults.map((result, index) => (
            <DataCard key={index} data={result} />
          ))}
        </div>
      </div>
    </main>
  );
}

