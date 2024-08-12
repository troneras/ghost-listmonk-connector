// src/lib/api-client.ts
import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8808';
const API_KEY = process.env.NEXT_PUBLIC_API_KEY || "your-api-key";

export const apiClient = axios.create({
    baseURL: API_BASE_URL.replace(/\/$/, ''), // Remove trailing slash if present
    headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${API_KEY}`,
    },
    withCredentials: true,
});