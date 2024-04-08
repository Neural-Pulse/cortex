// pages/api/search.ts
import type { NextApiRequest, NextApiResponse } from 'next';
import https from 'https';

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'POST') {
    const data = JSON.stringify(req.body);

    const options = {
      hostname: 'localhost',
      port: 9200,
      path: '/metadata_index3/_search?pretty',
      method: 'POST',
      headers: {
        'Authorization': 'ApiKey YkZkVHZvNEJoUHEwd1M2TmZfTTU6a3JscURnc19ReG1FUnFxZWJubGhDdw==',
        'Content-Type': 'application/json',
      },
      rejectUnauthorized: false, // Ignora a verificação de certificado SSL
    };

    const request = https.request(options, (response) => {
      let data = '';

      response.on('data', (chunk) => {
        data += chunk;
      });

      response.on('end', () => {
        res.status(200).json(JSON.parse(data));
      });
    });

    request.on('error', (error) => {
      console.error(error);
      res.status(500).json({ error: 'Erro ao se conectar com o Elasticsearch' });
    });

    request.write(data);
    request.end();
  } else {
    res.setHeader('Allow', ['POST']);
    res.status(405).end(`Method ${req.method} Not Allowed`);
  }
}
