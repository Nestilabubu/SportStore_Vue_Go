import axios from "axios";

const getBaseURL = () => {
  if (import.meta.env.PROD) {
    return import.meta.env.VITE_API_URL || "https://ваш-railway-проект.up.railway.app/api";
  }
  return "http://localhost:8080/api";
};

const api = axios.create({
  baseURL: getBaseURL(),
  withCredentials: true,
});

export default api;