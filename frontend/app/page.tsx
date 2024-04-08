'use client';

import https from 'https';
import { useState } from "react";

export default function Home() {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);

  const handleSearch = async () => {
    try {
      const response = await fetch("/api/search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          query: {
            query_string: {
              query: searchQuery,
            },
          },
        }),
      });
  
      if (!response.ok) {
        throw new Error(`Erro na requisição: ${response.statusText}`);
      }
  
      const data = await response.json();
  
      // Verifica se a propriedade 'hits' existe na resposta
      if (data.hits && data.hits.hits) {
        setSearchResults(data.hits.hits);
      } else {
        console.error('A resposta não contém a propriedade esperada "hits".', data);
        // Trate o caso de não haver 'hits' como achar melhor
      }
    } catch (error) {
      console.error('Erro ao buscar dados:', error);
      // Trate o erro como achar melhor
    }
  };

  return (
    <main className="flex flex-col items-center p-8">
      <div className="mb-8">
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

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
        {searchResults.map((result) => (
          <div key={result._id} className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-2">{result._source.title}</h2>
            <p className="text-gray-600">{result._source.description}</p>
          </div>
        ))}
      </div>
    </main>
  );
}
